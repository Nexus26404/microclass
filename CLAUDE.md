# 微课脚本生成器 - 开发规范

## 修改后端代码后（自动执行）

修改任何 `.go` 文件后，必须执行以下步骤（不可跳过）：

```bash
fuser -k 8080/tcp 2>/dev/null
cd /home/hzhou/workspace/microclass/backend
GOPROXY=https://goproxy.cn,direct go build -o microclass-backend .
./microclass-backend
```

**关键规则**：
- 后端是编译好的二进制（`microclass-backend`），不是 `go run main.go`
- 必须先 kill 旧进程再重新编译启动，否则跑的是旧代码
- 日志写入 `./logs/server.log`，每次启动自动清空

## 启动命令

- **后端**：`/home/hzhou/workspace/microclass/backend/microclass-backend`
- **前端**：`cd /home/hzhou/workspace/microclass/frontend && npm run dev`（已包含 `--host`，局域网可访问）

## 项目结构

- 前端：`/home/hzhou/workspace/microclass/frontend/`
- 后端：`/home/hzhou/workspace/microclass/backend/`
- 日志：`/home/hzhou/workspace/microclass/backend/logs/server.log`

## 技术栈

- 前端：Vue3 + TS + TailwindCSS + Pinia + Vue Router
- 后端：Go + Gin + SQLite + go-openai SDK
- LLM：MiniMax（OpenAI 兼容格式）
