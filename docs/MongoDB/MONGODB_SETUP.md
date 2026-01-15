## Mongo Express 配置说明

以下为 `mongo-express` 的 docker-compose 配置示例，并对每一项进行了详细注释：

```yaml
networks:
 1panel-network:                # 指定使用的外部网络，便于与其他服务通信
  external: true
services:
 mongo-express:
  container_name: ${CONTAINER_NAME}   # 容器名称，可通过环境变量自定义
  deploy:
   resources:
    limits:
     cpus: ${CPUS}          # 限制容器可用的 CPU 数量
     memory: ${MEMORY_LIMIT}# 限制容器可用的内存大小
  environment:
   ME_CONFIG_BASICAUTH: true      # 启用基本认证，防止未授权访问
   ME_CONFIG_BASICAUTH_PASSWORD: ${BASICAUTH_PASSWORD} # 访问 mongo-express 的密码
   ME_CONFIG_BASICAUTH_USERNAME: ${BASICAUTH_USERNAME} # 访问 mongo-express 的用户名
   ME_CONFIG_MONGODB_URL: mongodb://${PANEL_DB_ROOT_USER}:${PANEL_DB_ROOT_PASSWORD}@${MONGO_HOST}:27017 # 连接 MongoDB 的 URL，包含用户名、密码和主机(数据库服务)
  image: mongo-express:1.0.2-20      # 使用的 mongo-express 镜像及版本
  labels:
   createdBy: Apps                # 自定义标签，可用于标识容器用途
  networks:
   - 1panel-network               # 连接到指定的网络
  ports:
   - ${HOST_IP}:${PANEL_APP_PORT_HTTP}:8081 # 将主机端口映射到容器 8081 端口
  restart: on-failure:5              # 失败时自动重启，最多重启 5 次
```

> **说明**：
>
> - 请根据实际环境设置各个 `${}` 变量的值。
> - `mongo-express` 是一个可视化 MongoDB 管理工具，适合开发和测试环境使用。
>
# MongoDB 快速指南

## 1. env 配置

**作用：** 存储 MongoDB 连接信息，避免硬编码敏感数据

```dotenv
# MongoDB 连接字符串（包含用户名、密码、主机地址）
# host.docker.internal 用于 Docker 容器访问宿主机 MongoDB
MONGO_URI=mongodb://username:password@host.docker.internal:27017

# 数据库名称 - 博客项目使用的数据库
DB_NAME=blogs-db

# 集合名称 - 数据库中的表，存储博客文档
COLLECTION_NAME=blogs-db-dev
```

## 2. 初始化指令

**作用：** 在 MongoDB 中创建数据库和集合，初始化项目

```javascript
// 连接到 MongoDB 服务（使用实际的用户名和密码）
mongosh "mongodb://localhost:27017" --username <user> --password <pwd>

// 切换到 blogs-db 数据库（不存在则自动创建）
use blogs-db;

// 插入初始文档，同时创建 blogs-db-dev 集合
db["blogs-db-dev"].insertOne({test: "init"});

// 验证集合已创建
show collections;

// 退出 mongosh
.exit
```

## 3. Go 代码使用

**作用：** 在后端代码中连接和操作 MongoDB 数据库

```go
// 第一步：连接 MongoDB 服务
// 使用 .env 中的 MONGO_URI，包含认证信息
client, _ := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))

// 第二步：获取数据库对象
// DB_NAME = "blogs-db"
db := client.Database(os.Getenv("DB_NAME"))

// 第三步：获取集合对象
// COLLECTION_NAME = "blogs-db-dev"，相当于传统数据库中的表
collection := db.Collection(os.Getenv("COLLECTION_NAME"))

// 第四步：执行数据库操作
collection.Find(ctx, filter)           // 查询数据
collection.InsertOne(ctx, doc)         // 插入单条数据
collection.UpdateOne(ctx, filter, update)  // 更新数据
collection.DeleteOne(ctx, filter)      // 删除数据
```
