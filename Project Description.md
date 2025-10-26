# Go 博客管理后端

一个使用 Go 语言编写的简单博客管理 RESTful API 服务。

## 功能特性

- ✅ 获取所有博客文章
- ✅ 获取单篇博客文章
- ✅ 创建新博客文章
- ✅ 更新博客文章
- ✅ 删除博客文章
- ✅ 健康检查端点

## 技术栈

- **Go** 1.21+
- **Gorilla Mux** - HTTP 路由器

## 快速开始

### 前置要求

- 安装 Go 1.21 或更高版本

### 安装依赖

```bash
go mod download
```

### 运行服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

## API 文档

### 基础 URL

```
http://localhost:8080/api
```

### 端点列表

#### 1. 获取所有博客文章

```
GET /api/posts
```

**响应示例:**
```json
[
  {
    "id": 1,
    "title": "欢迎使用 Go 博客系统",
    "content": "这是第一篇博客文章。Go 是一门很棒的语言！",
    "author": "管理员",
    "created_at": "2025-10-27T10:00:00Z",
    "updated_at": "2025-10-27T10:00:00Z"
  }
]
```

#### 2. 获取单篇博客文章

```
GET /api/posts/{id}
```

**响应示例:**
```json
{
  "id": 1,
  "title": "欢迎使用 Go 博客系统",
  "content": "这是第一篇博客文章。Go 是一门很棒的语言！",
  "author": "管理员",
  "created_at": "2025-10-27T10:00:00Z",
  "updated_at": "2025-10-27T10:00:00Z"
}
```

#### 3. 创建新博客文章

```
POST /api/posts
```

**请求体:**
```json
{
  "title": "我的新文章",
  "content": "文章内容...",
  "author": "作者姓名"
}
```

**响应示例:**
```json
{
  "id": 2,
  "title": "我的新文章",
  "content": "文章内容...",
  "author": "作者姓名",
  "created_at": "2025-10-27T10:30:00Z",
  "updated_at": "2025-10-27T10:30:00Z"
}
```

#### 4. 更新博客文章

```
PUT /api/posts/{id}
```

**请求体:**
```json
{
  "title": "更新后的标题",
  "content": "更新后的内容...",
  "author": "作者姓名"
}
```

#### 5. 删除博客文章

```
DELETE /api/posts/{id}
```

**响应:** 204 No Content

#### 6. 健康检查

```
GET /health
```

**响应:** "服务运行正常"

## 使用示例

### 使用 curl 测试

#### 获取所有文章
```bash
curl http://localhost:8080/api/posts
```

#### 创建新文章
```bash
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -d '{"title":"测试文章","content":"这是测试内容","author":"张三"}'
```

#### 获取单篇文章
```bash
curl http://localhost:8080/api/posts/1
```

#### 更新文章
```bash
curl -X PUT http://localhost:8080/api/posts/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"更新的标题","content":"更新的内容","author":"张三"}'
```

#### 删除文章
```bash
curl -X DELETE http://localhost:8080/api/posts/1
```

## 项目结构

```
blog-backend-go/
├── main.go       # 主程序文件
├── go.mod        # Go 模块依赖
└── README.md     # 项目说明文档
```

## 注意事项

⚠️ **这是一个最小示例项目**，使用内存存储数据。生产环境中应该：

1. 使用真实的数据库（如 PostgreSQL、MySQL、MongoDB）
2. 添加身份认证和授权
3. 添加输入验证
4. 添加日志记录
5. 添加错误处理
6. 添加 CORS 支持
7. 使用配置文件管理环境变量
8. 添加单元测试

## 下一步改进

- [ ] 集成数据库（PostgreSQL/MySQL）
- [ ] 添加用户认证（JWT）
- [ ] 添加分页功能
- [ ] 添加文章搜索功能
- [ ] 添加文章分类和标签
- [ ] 添加评论功能
- [ ] 添加日志中间件
- [ ] 编写单元测试

## 许可证

MIT
