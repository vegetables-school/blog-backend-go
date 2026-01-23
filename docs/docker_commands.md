
# Docker 常用命令与容器互联说明

## 查看所有容器

```bash
docker ps -a

---

## docker-compose 连接另一个容器内 MongoDB 的配置说明

### 场景说明

在使用 docker-compose 部署时，通常会将应用（如 blog-backend-go）和数据库（如 MongoDB）分别作为独立的服务运行在不同的容器中。此时，应用需要通过服务名来连接 MongoDB 容器，而不是使用 localhost 或 127.0.0.1。

### 步骤与注意事项

1. **docker-compose.yml 配置示例**

```yaml
services:
  app:
    # ...其他配置...
    depends_on:
      - mongo
    environment:
      - MONGO_URI=mongodb://mongo:27017/blog
  mongo:
    image: mongo:latest
    restart: always
    ports:
      - "27017:27017"
    volumes:
      - mongo-data:/data/db

volumes:
  mongo-data:
```

2. **连接字符串写法**

	 - 连接 MongoDB 时，host 应写为 `mongo`（即 compose 文件中 mongo 服务的名称），而不是 `localhost`。
	 - 例如：`mongodb://mongo:27017/blog`

3. **depends_on 的作用**

	 - `depends_on` 保证 app 服务在 mongo 服务启动后再启动，但不能保证 MongoDB 已经完全可用。建议应用内实现重试机制。

4. **端口映射说明**

	 - `ports` 字段用于将容器端口映射到主机，便于本地调试。容器间通信无需依赖端口映射，直接通过服务名和容器端口即可。

5. **网络说明**

	 - docker-compose 默认会为同一项目的服务创建一个专用网络，服务间可通过服务名互相访问。

6. **常见问题排查**

	 - 连接报错 `ECONNREFUSED` 或 `No suitable servers found`，请检查：
		 - 连接字符串 host 是否为 mongo（服务名）
		 - MongoDB 服务是否正常启动
		 - 应用是否有重试机制

---

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

---

- `docker ps` 只显示正在运行的容器。
- `docker ps -a` 显示所有容器（包括已停止的）。

### 示例输出

```bash
CONTAINER ID   IMAGE         COMMAND                  CREATED          STATUS                      PORTS     NAMES
123456789abc   nginx:latest  "nginx -g 'daemon of…"   2 minutes ago    Exited (0) 1 minute ago               my-nginx
abcdef123456   redis:alpine  "docker-entrypoint.s…"   5 minutes ago    Up 5 minutes                6379/tcp  my-redis
```
