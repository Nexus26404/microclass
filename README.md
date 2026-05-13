# 微课脚本生成器

基于 LLM 自动生成微课教学脚本的工具，支持流式输出、脚本历史管理、Markdown 导出。

## 技术栈

- **前端**：Vue 3 + TypeScript + TailwindCSS + Pinia + Vue Router
- **后端**：Go + Gin + SQLite + go-openai SDK
- **LLM**：MiniMax Text-01（OpenAI 兼容格式）

## 项目结构

```
microclass/
├── frontend/          # Vue3 前端
│   ├── src/
│   │   ├── views/     # 页面组件（Home、History、Settings）
│   │   ├── stores/    # Pinia 状态管理
│   │   └── components/# 可复用组件
│   └── dist/          # 构建产物（前端资源）
├── backend/           # Go 后端
│   ├── handler/       # HTTP handlers
│   ├── service/       # LLM 服务
│   ├── store/         # SQLite 操作
│   └── main.go        # 入口
├── build/             # 构建输出目录
└── Makefile
```

## 快速开始

### 开发模式

```bash
# 前端（热重载，局域网可访问）
cd frontend && npm run dev

# 后端（另一个终端）
cd backend && ./microclass-backend
```

### 构建

```bash
make build    # 构建前端 + 后端
make run      # 构建并运行
make clean    # 清理
```

构建产物输出到 `build/` 目录：
```
build/
├── microclass-backend  # 可执行文件
├── dist/              # 前端静态资源
├── data/              # SQLite 数据库
└── logs/              # 日志
```

### 环境变量（后端）

| 变量 | 默认值 | 说明 |
|------|--------|------|
| PORT | 8080 | 服务端口 |
| MOCK_MODE | true | Mock 模式（默认开启） |
| OPENAI_API_KEY | - | API Key |
| OPENAI_BASE_URL | https://api.minimax.chat | API 地址 |
| OPENAI_MODEL | MiniMax-Text-01 | 模型名称 |
| DB_PATH | ./data/microclass.db | 数据库路径 |

## 功能

- [x] 微课脚本生成（流式/非流式）
- [x] 脚本历史管理
- [x] 设置管理（API 配置、Mock 模式）
- [x] 参考文件上传
- [x] Markdown 导出

## API 端点

| Method | Endpoint | 说明 |
|--------|----------|------|
| POST | `/api/generate-stream` | 流式生成脚本 |
| POST | `/api/generate` | 非流式生成 |
| GET | `/api/scripts` | 列出所有脚本 |
| GET | `/api/scripts/:id` | 获取单个脚本 |
| DELETE | `/api/scripts/:id` | 删除脚本 |
| GET/PUT | `/api/settings` | 获取/更新设置 |
| POST | `/api/upload` | 上传参考文件 |

## 未来计划

- [ ] 桌面端（考虑 Wails + Vue）
- [ ] 移动端（考虑 Capacitor）

## License

MIT