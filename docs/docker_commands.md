# Docker 常用命令

## 查看所有容器

```bash
docker ps -a
```

### 作用说明

`docker ps -a` 用于列出所有容器，包括正在运行的和已经停止的容器。通过该命令可以查看容器的状态、名称、ID、创建时间等信息，方便管理和排查问题。

- `docker ps` 只显示正在运行的容器。
- `docker ps -a` 显示所有容器（包括已停止的）。

### 示例输出

```bash
CONTAINER ID   IMAGE         COMMAND                  CREATED          STATUS                      PORTS     NAMES
123456789abc   nginx:latest  "nginx -g 'daemon of…"   2 minutes ago    Exited (0) 1 minute ago               my-nginx
abcdef123456   redis:alpine  "docker-entrypoint.s…"   5 minutes ago    Up 5 minutes                6379/tcp  my-redis
```
