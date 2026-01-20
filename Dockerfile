## ----------------------------
# Multi-stage build for Go app
## ----------------------------

## ----------------------------
# 多阶段构建 Go 应用
## ----------------------------

FROM golang:1.24-alpine AS builder
WORKDIR /src  # 设置工作目录为/src

# 仅拷贝依赖文件用于缓存依赖，加速构建
COPY go.mod go.sum ./
# 设置 Go 模块代理，安装依赖工具，下载依赖
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    apk add --no-cache ca-certificates git && \
    go mod download

# 拷贝项目所有文件到构建环境
COPY . .

# 编译 Go 程序为静态二进制，输出到/app目录
RUN mkdir -p /app \
    && CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w' -o /app/blog-backend-go-dev ./

### 构建最终运行镜像
FROM alpine:3.18
# 安装 CA 证书
RUN apk add --no-cache ca-certificates
# 设置工作目录为/app
WORKDIR /app
# 从构建阶段复制二进制文件到运行环境
COPY --from=builder /app/blog-backend-go-dev /app/blog-backend-go-dev

# 设置环境变量
ENV PORT=8080
# 暴露端口
EXPOSE 8080
# ENTRYPOINT 指定容器启动时默认执行的主命令。
# 这里设置为 /app/blog-backend-go-dev，容器启动后会自动运行该可执行文件。
ENTRYPOINT ["/app/blog-backend-go-dev"]
