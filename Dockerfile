## ----------------------------
# Multi-stage build for Go app
## ----------------------------

FROM golang:1.24-alpine AS builder
WORKDIR /src

# Ensure go modules are downloaded separately for build cache
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    apk add --no-cache ca-certificates git && \
    go mod download

COPY . .

# Build a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w' -o /blog-backend-go-dev ./

### Final image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /
COPY --from=builder /blog-backend-go-dev /blog-backend-go-dev

ENV PORT=8080
EXPOSE 8080
# ENTRYPOINT 指定容器启动时默认执行的主命令。
# 这里设置为 /blog-backend-go-dev，容器启动后会自动运行该可执行文件。
ENTRYPOINT ["/blog-backend-go-dev"]
