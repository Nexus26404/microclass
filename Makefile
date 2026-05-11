.PHONY: all build build-frontend build-backend run clean

# 项目根目录
ROOT_DIR := $(shell pwd)
FRONTEND_DIR := $(ROOT_DIR)/frontend
BACKEND_DIR := $(ROOT_DIR)/backend
BUILD_DIR := $(ROOT_DIR)/build

all: build

# 构建前端 + 后端
build: build-frontend build-backend

# 构建前端
build-frontend:
	cd $(FRONTEND_DIR) && npm run build

# 复制前端 dist 到 build/dist
	mkdir -p $(BUILD_DIR)
	cp -r $(FRONTEND_DIR)/dist $(BUILD_DIR)/dist

# 构建后端
build-backend:
	@# 跨平台停止旧进程（兼容 macOS/Linux）
	@if lsof -ti:8080 > /dev/null 2>&1; then \
		lsof -ti:8080 | xargs kill 2>/dev/null || true; \
		sleep 1; \
	elif pkill -f microclass-backend 2>/dev/null; then \
		sleep 1; \
	fi
	cd $(BACKEND_DIR) && GOPROXY=https://goproxy.cn,direct go build -o $(BUILD_DIR)/microclass-backend .

# 运行
run: build
	cd $(BUILD_DIR) && ./microclass-backend

# 清理
clean:
	rm -rf $(BUILD_DIR)
	rm -rf $(FRONTEND_DIR)/dist