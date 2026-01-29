# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -o pansiot-cloud cmd/server/main.go

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装必要工具
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 从构建阶段复制二进制文件
COPY --from=builder /app/pansiot-cloud .
COPY --from=builder /app/configs ./configs

# 创建日志目录
RUN mkdir -p logs

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./pansiot-cloud"]
