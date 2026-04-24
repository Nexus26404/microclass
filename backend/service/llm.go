package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
	"microclass-backend/config"
	"microclass-backend/model"
)

type LLMService struct {
	cfg *config.Config
}

func NewLLMService(cfg *config.Config) *LLMService {
	return &LLMService{cfg: cfg}
}

func (s *LLMService) GenerateScript(req model.GenerateRequest) (*model.Script, error) {
	cfg := config.Get()
	if cfg.MockMode {
		return s.mockGenerate(req)
	}

	var content strings.Builder
	script, err := s.generateStream(req, func(chunk string) {
		content.WriteString(chunk)
	})
	if err != nil {
		return nil, err
	}
	_ = content // silence unused warning
	return script, nil
}

func (s *LLMService) GenerateScriptStream(req model.GenerateRequest, onChunk func(string)) (*model.Script, error) {
	cfg := config.Get()
	if cfg.MockMode {
		script, err := s.mockGenerate(req)
		if err != nil {
			return nil, err
		}
		// 模拟流式输出
		data, _ := json.Marshal(script)
		for i := 0; i < len(data); i += 20 {
			end := i + 20
			if end > len(data) {
				end = len(data)
			}
			onChunk(string(data[i:end]))
		}
		return script, nil
	}
	return s.generateStream(req, onChunk)
}

func (s *LLMService) generateStream(req model.GenerateRequest, onChunk func(string)) (*model.Script, error) {
	styleDesc := map[string]string{
		"formal":       "正式严谨，逻辑清晰",
		"casual":       "轻松口语化，像朋友聊天",
		"storytelling": "讲故事的方式，引人入胜",
	}

	systemPrompt := `你是一个专业的微课脚本生成助手。请根据用户输入的信息，生成一个结构化的微课脚本。

要求：
1. 脚本内容要口语化，适合录制短视频课程
2. 每个段落要包含：标题、内容、预估时长（秒）
3. 内容要生动有趣，避免照本宣科
4. 直接返回 JSON，不要包含其他文字

返回格式：
{
  "sections": [
    {"type": "opening", "title": "话题导入", "content": "...", "estimatedDuration": 30},
    {"type": "knowledge", "title": "知识点讲解", "content": "...", "estimatedDuration": 90},
    {"type": "example", "title": "实践案例", "content": "...", "estimatedDuration": 60},
    {"type": "summary", "title": "总结回顾", "content": "...", "estimatedDuration": 30}
  ]
}`

	userPrompt := fmt.Sprintf(`微课信息：
- 主题：%s
- 学科/行业：%s
- 目标受众：%s
- 预计时长：%d秒
- 风格：%s
- 自定义结构：%s
%s

请生成完整的微课脚本。`,
		req.Topic,
		req.Discipline,
		req.TargetAudience,
		req.TargetDuration,
		styleDesc[req.Style],
		req.CustomStructure,
		formatReference(req.ReferenceText, req.ReferenceFileURL),
	)

	cfg := config.Get()

	var client *openai.Client
	if cfg.OpenAIURL != "" {
		c := openai.DefaultConfig(cfg.OpenAIKey)
		c.BaseURL = cfg.OpenAIURL + "/v1"
		client = openai.NewClientWithConfig(c)
	} else {
		client = openai.NewClient(cfg.OpenAIKey)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	stream, err := client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model: cfg.OpenAIModel,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: userPrompt},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("LLM API error: %v", err)
	}
	defer stream.Close()

	var content strings.Builder
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("stream error: %v", err)
		}
		if len(chunk.Choices) == 0 {
			continue
		}
		// 思考内容走单独字段，不混入正文
		if chunk.Choices[0].Delta.ReasoningContent != "" {
			continue
		}
		text := chunk.Choices[0].Delta.Content
		content.WriteString(text)
		onChunk(text)
	}

	contentStr := content.String()
	log.Printf("[llm] 收到响应（%d字）: %.200s", len(contentStr), contentStr)

	jsonStr := extractJSON(contentStr)

	var raw struct {
		Sections []model.Section `json:"sections"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		log.Printf("[llm] JSON 解析失败: %v", err)
		return nil, fmt.Errorf("failed to parse LLM response: %v", err)
	}

	log.Printf("[llm] 解析成功，共 %d 个章节", len(raw.Sections))

	now := time.Now()
	script := &model.Script{
		ID:             fmt.Sprintf("%d", now.UnixNano()),
		Topic:          req.Topic,
		Discipline:     req.Discipline,
		TargetAudience: req.TargetAudience,
		Style:          req.Style,
		Sections:       raw.Sections,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	total := 0
	for i := range script.Sections {
		if script.Sections[i].ID == "" {
			script.Sections[i].ID = fmt.Sprintf("%d-%d", now.UnixNano(), i+1)
		}
		total += script.Sections[i].EstimatedDuration
	}
	script.EstimatedTotalDuration = total

	log.Printf("[llm] 脚本生成成功，ID: %s", script.ID)
	return script, nil
}

func (s *LLMService) mockGenerate(req model.GenerateRequest) (*model.Script, error) {
	now := time.Now()
	nano := now.UnixNano()

	styleDesc := map[string]string{
		"formal":       "正式严谨，逻辑清晰",
		"casual":       "轻松口语化，像朋友聊天",
		"storytelling": "讲故事的方式，引人入胜",
	}

	discipline := req.Discipline
	if discipline == "" {
		discipline = "通用"
	}

	audience := req.TargetAudience
	if audience == "" {
		audience = "通用受众"
	}

	totalDuration := req.TargetDuration
	if totalDuration == 0 {
		totalDuration = 180
	}

	openingDur := int(float64(totalDuration) * 0.12)
	knowledgeDur := int(float64(totalDuration) * 0.40)
	exampleDur := int(float64(totalDuration) * 0.28)
	summaryDur := totalDuration - openingDur - knowledgeDur - exampleDur

	var sections []model.Section
	if req.CustomStructure != "" {
		parts := parseCustomStructure(req.CustomStructure)
		for i, part := range parts {
			sectionType := inferSectionType(part)
			partDuration := totalDuration / len(parts)
			sections = append(sections, model.Section{
				ID:                fmt.Sprintf("%d-%d", nano, i+1),
				Type:              sectionType,
				Title:             part,
				Content:           buildSectionContent(part, req.Topic, discipline, audience, styleDesc[req.Style]),
				EstimatedDuration: partDuration,
			})
		}
	} else {
		sections = []model.Section{
			{
				ID:                fmt.Sprintf("%d-1", nano),
				Type:              "opening",
				Title:             "话题导入",
				Content:           fmt.Sprintf("【%s · %s】\n\n大家好，欢迎来到今天的微课！\n\n我是你们的主讲老师，今天我们来一起学习「%s」。\n\n在开始之前，我想问大家一个问题：%s\n\n带着这个问题，让我们进入今天的内容。", discipline, audience, req.Topic, generateHookQuestion(req.Topic, discipline)),
				EstimatedDuration: openingDur,
			},
			{
				ID:                fmt.Sprintf("%d-2", nano),
				Type:              "knowledge",
				Title:             "知识点讲解",
				Content:           fmt.Sprintf("【%s】%s\n\n%s\n\n%s\n\n这里有一个关键点需要记住：%s", discipline, req.Topic, explainConcept(req.Topic, discipline), giveExample(req.Topic, discipline), highlightKey(req.Topic)),
				EstimatedDuration: knowledgeDur,
			},
			{
				ID:                fmt.Sprintf("%d-3", nano),
				Type:              "example",
				Title:             "实践案例",
				Content:           fmt.Sprintf("理论讲完了，让我们来看一个实际案例。\n\n场景：%s\n\n步骤一：%s\n步骤二：%s\n步骤三：%s\n\n%s", buildScenario(req.Topic, discipline), buildStep1(req.Topic), buildStep2(req.Topic), buildStep3(req.Topic), showResult(req.Topic)),
				EstimatedDuration: exampleDur,
			},
			{
				ID:                fmt.Sprintf("%d-4", nano),
				Type:              "summary",
				Title:             "总结回顾",
				Content:           fmt.Sprintf("今天我们学习了「%s」，%s\n\n%s\n\n%s\n\n今天的课就到这里，如果觉得有帮助，欢迎收藏和分享！", req.Topic, discipline, summarizePoints(req.Topic), callToAction(req.Topic)),
				EstimatedDuration: summaryDur,
			},
		}
	}

	return &model.Script{
		ID:                      fmt.Sprintf("%d", nano),
		Topic:                   req.Topic,
		Discipline:             discipline,
		TargetAudience:         audience,
		Style:                  req.Style,
		EstimatedTotalDuration: totalDuration,
		Sections:               sections,
		CreatedAt:              now,
		UpdatedAt:              now,
	}, nil
}

func parseCustomStructure(raw string) []string {
	var parts []string
	for _, p := range splitWithTrim(raw, "----") {
		trimmed := trimBrackets(p)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	if len(parts) == 0 {
		return []string{"话题导入", "知识点讲解", "实践案例", "总结回顾"}
	}
	return parts
}

func splitWithTrim(s, sep string) []string {
	result := []string{}
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trimBrackets(s string) string {
	s = trim(s, "[")
	s = trim(s, "]")
	s = trim(s, "（")
	s = trim(s, "）")
	s = trim(s, "(")
	s = trim(s, ")")
	return trimSpace(s)
}

func trim(s, cutset string) string {
	for len(s) > 0 && contains(cutset, s[0]) {
		s = s[1:]
	}
	for len(s) > 0 && contains(cutset, s[len(s)-1]) {
		s = s[:len(s)-1]
	}
	return s
}

func trimSpace(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[0] == '\n') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t' || s[len(s)-1] == '\n') {
		s = s[:len(s)-1]
	}
	return s
}

func contains(set string, c byte) bool {
	for i := 0; i < len(set); i++ {
		if set[i] == c {
			return true
		}
	}
	return false
}

func inferSectionType(title string) string {
	t := trimBrackets(title)
	for _, kw := range []string{"开场", "导入", "引入", "开头", "引入语"} {
		if len(t) >= len(kw) && (contains(t, kw[0]) || t[:len(kw)] == kw) {
			return "opening"
		}
	}
	for _, kw := range []string{"总结", "回顾", "收尾", "结尾"} {
		if len(t) >= len(kw) && contains(t, kw[0]) {
			return "summary"
		}
	}
	for _, kw := range []string{"案例", "实践", "练习", "习题", "操作", "应用"} {
		if len(t) >= len(kw) && contains(t, kw[0]) {
			return "example"
		}
	}
	return "knowledge"
}

func buildSectionContent(title, topic, discipline, audience, styleDesc string) string {
	return fmt.Sprintf("【%s · %s · %s】\n\n%s\n\n本节重点：%s", discipline, audience, title, topic, topic)
}

func generateHookQuestion(topic, discipline string) string {
	return fmt.Sprintf("你知道%s和我们的日常生活有什么联系吗？", topic)
}

func explainConcept(topic, discipline string) string {
	return fmt.Sprintf("首先，让我们来理解什么是%s。\n\n在%s领域，%s是一个非常基础但又非常重要的概念。\n\n简单来说，它的核心原理是...", topic, discipline, topic)
}

func giveExample(topic, discipline string) string {
	return fmt.Sprintf("举个例子，在%s的实际工作中，%s经常被用来解决...问题。", discipline, topic)
}

func highlightKey(topic string) string {
	return fmt.Sprintf("记住：理解%s的关键在于把握它的核心逻辑，而不是死记硬背。", topic)
}

func buildScenario(topic, discipline string) string {
	return fmt.Sprintf("小李是一名%s从业者，今天他需要用%s来解决一个实际问题。", discipline, topic)
}

func buildStep1(topic string) string { return "第一步：理解问题的本质，明确目标" }
func buildStep2(topic string) string { return "第二步：按照方法逐步操作，注意细节" }
func buildStep3(topic string) string { return "第三步：验证结果，检查是否有遗漏" }

func showResult(topic string) string {
	return fmt.Sprintf("最终，小李成功用%s解决了问题！这就是%s在实际中的应用。", topic, topic)
}

func summarizePoints(topic string) string {
	return fmt.Sprintf("通过今天的学习，我们需要掌握以下要点：\n\n1. %s的基本概念\n2. %s的核心原理\n3. %s的实际应用场景", topic, topic, topic)
}

func callToAction(topic string) string {
	return fmt.Sprintf("如果你对%s还有其他疑问，欢迎在评论区留言，我会一一解答。下节课见！", topic)
}

func formatReference(text, fileURL string) string {
	if text == "" && fileURL == "" {
		return ""
	}
	if text != "" {
		if len(text) > 3000 {
			text = text[:3000] + "\n[...参考资料已被截断...]"
		}
		return "---\n参考资料：\n" + text + "\n---"
	}
	return ""
}

func stripThinkTags(content string) string {
	for {
		open := strings.Index(content, "<think>")
		if open == -1 {
			break
		}
		close := strings.Index(content, "</think>")
		if close == -1 {
			break
		}
		content = content[:open] + content[close+len("</think>"):]
	}
	return content
}

func extractJSON(content string) string {
	idx := strings.Index(content, "```json")
	if idx != -1 {
		start := idx + len("```json")
		end := strings.Index(content[start:], "```")
		if end != -1 {
			return strings.TrimSpace(content[start : start+end])
		}
	}
	idx = strings.Index(content, "{")
	if idx == -1 {
		return content
	}
	depth := 0
	for i := idx; i < len(content); i++ {
		if content[i] == '{' {
			depth++
		} else if content[i] == '}' {
			depth--
			if depth == 0 {
				return content[idx : i+1]
			}
		}
	}
	return content
}

