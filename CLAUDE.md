# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

微课脚本生成器 - 基于 LLM 生成微课教学脚本的 Web 应用。前端 Vue3，后端 Go + Gin，SQLite 存储，使用 MiniMax GPT（OpenAI 兼容格式）。

## 开发命令

### 前端
```bash
cd frontend && npm run dev        # 启动开发服务器（已含 --host，局域网可访问）
cd frontend && npm run build      # 构建生产版本（需先运行 vue-tsc -b）
```

### 后端
```bash
# 修改 .go 文件后必须重新编译
fuser -k 8080/tcp 2>/dev/null
cd backend
GOPROXY=https://goproxy.cn,direct go build -o microclass-backend .
./microclass-backend
```

**关键**：后端是编译后的二进制 `microclass-backend`，不是 `go run main.go`。日志写入 `./logs/server.log`，每次启动自动清空。

## 架构

### 前端 (Vue3 + TS)
- `src/views/` - 三个页面：HomePage（生成）、HistoryPage（历史）、SettingsPage（设置）
- `src/stores/script.ts` - Pinia store，管理脚本状态
- `src/components/ScriptDisplay.vue` - 脚本渲染与内联编辑
- 路由：`/`、`/history`、`/settings`
- Vite 代理 `/api/*` 到 `http://localhost:8080`

### 后端 (Go + Gin)
- `handler/` - HTTP handlers（script、settings、upload）
- `service/llm.go` - LLM 服务（mock 模式 + 真实 API 流式生成）
- `store/sqlite.go` - SQLite 数据库操作
- API 全部在 `/api/*` 前缀下

### API 端点
| Method | Endpoint | 说明 |
|--------|----------|------|
| POST | `/api/generate-stream` | 流式生成脚本（SSE） |
| POST | `/api/generate` | 非流式生成 |
| GET | `/api/scripts` | 列出所有脚本 |
| DELETE | `/api/scripts/:id` | 删除脚本 |
| GET/PUT | `/api/settings` | 获取/更新设置 |
| POST | `/api/upload` | 上传参考文件 |

### 数据模型
- `Script` - 微课脚本（主题、受众、风格、时长、章节列表）
- `Settings` - API 配置（mock 模式、API Key、URL、模型）

## 环境变量（后端）

| 变量 | 默认值 | 说明 |
|------|--------|------|
| PORT | 8080 | 服务端口 |
| MOCK_MODE | true | 使用 mock 生成（默认开启） |
| OPENAI_API_KEY | - | 真实 API Key |
| OPENAI_BASE_URL | https://api.minimax.chat | API 地址 |
| OPENAI_MODEL | MiniMax-Text-01 | 模型名称 |
| DB_PATH | ./data/microclass.db | SQLite 数据库路径 |

## 技术栈

- 前端：Vue 3.5 + TypeScript 5.6 + TailwindCSS 3.4 + Pinia 2.2 + Vue Router 4.6 + Vite 6.1
- 后端：Go 1.26 + Gin 1.10 + SQLite (mattn/go-sqlite3) + go-openai SDK
- LLM：MiniMax Text-01（OpenAI 兼容格式）
