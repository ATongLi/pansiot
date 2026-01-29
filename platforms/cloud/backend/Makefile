.PHONY: all build run test clean deps

# 变量定义
APP_NAME=pansiot-cloud
CMD_DIR=./cmd/server
BUILD_DIR=./build
CONFIG_FILE=./configs/config.yaml

# Go参数
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOGET=$(GOCMD) get

# 默认目标
all: deps build

# 下载依赖
deps:
	$(GOMOD) tidy
	$(GOMOD) download

# 运行
run:
	$(GORUN) $(CMD_DIR)/main.go

# 运行（指定配置文件）
run-dev:
	$(GORUN) $(CMD_DIR)/main.go -c $(CONFIG_FILE)

# 构建
build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

# 测试
test:
	$(GOTEST) -v ./...

# 测试（覆盖率）
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# 清理
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# 格式化代码
fmt:
	gofmt -s -w .

# 代码检查
lint:
	golangci-lint run

# Docker构建
docker-build:
	docker build -t $(APP_NAME):latest .

# Docker运行
docker-run:
	docker-compose up -d

# Docker停止
docker-down:
	docker-compose down

# 数据库迁移
migrate-up:
	@echo "Running database migrations..."
	# TODO: 实现数据库迁移命令

migrate-down:
	@echo "Rolling back database migrations..."
	# TODO: 实现数据库回滚命令

# 生成Swagger文档
swagger:
	swag init -g $(CMD_DIR)/main.go -o ./docs
