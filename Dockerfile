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
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-s -w' -o /blog-app ./

### Final image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /
COPY --from=builder /blog-app /blog-app

ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/blog-app"]
