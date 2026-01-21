# Docker 常用命令

## 查看所有容器

```bash
docker ps -a

---

## 1Panel 卷挂载说明

### 1. 在 1Panel 创建卷

- 卷名称：`blog-backend-go-dev`
- 挂载点（宿主机目录）：`/var/lib/docker/volumes/blog-backend-go-dev/_data`

### 2. docker-compose.yml 配置

```yaml
services:
	app:
		# ...其他配置...
		volumes:
			- blog-backend-go-dev:/blog-backend-go-dev

volumes:
	blog-backend-go-dev:
		external: true
```

### 3. 1Panel 容器目录配置

- 服务器目录：`/var/lib/docker/volumes/blog-backend-go-dev/_data`
- 容器目录：`/blog-backend-go-dev`

这样配置后，容器内 `/blog-backend-go-dev` 目录的数据会自动同步到宿主机卷。

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
