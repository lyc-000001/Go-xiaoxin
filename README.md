# Blog Backend

一个基于 Go + Gin + MySQL + Redis 的企业级博客后台系统。

## 技术栈

- **Web框架**: Gin
- **数据库**: MySQL
- **缓存**: Redis
- **ORM**: GORM
- **日志**: Zap
- **认证**: JWT
- **配置管理**: Viper

## 项目结构

```
blog-backend/
├── cmd/
│   └── server/          # 应用入口
│       └── main.go
├── config/              # 配置文件
│   └── config.yaml
├── internal/            # 私有应用代码
│   ├── controllers/     # 控制器层
│   ├── models/          # 数据模型
│   ├── services/        # 业务逻辑层
│   ├── middleware/      # 中间件
│   ├── routes/          # 路由配置
│   └── utils/           # 工具函数
├── pkg/                 # 公共库代码
│   ├── config/          # 配置管理
│   ├── database/        # 数据库连接
│   ├── redis/           # Redis连接
│   ├── logger/          # 日志系统
│   └── jwt/             # JWT工具
├── uploads/             # 上传文件目录
├── logs/                # 日志文件目录
├── go.mod
├── go.sum
└── README.md
```

## 功能特性

- ✅ 用户注册、登录、认证
- ✅ JWT Token 认证
- ✅ 文章CRUD操作
- ✅ 文章分类管理
- ✅ 文章标签系统
- ✅ 文件上传功能
- ✅ Redis缓存支持
- ✅ 结构化日志记录
- ✅ 跨域支持（CORS）
- ✅ 角色权限控制

## 环境要求

- Go 1.21+
- MySQL 5.7+
- Redis 5.0+

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd blog-backend
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置数据库

创建MySQL数据库：

```sql
CREATE DATABASE blog_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. 修改配置

编辑 `config/config.yaml`，修改数据库和Redis连接信息：

```yaml
database:
  host: "127.0.0.1"
  port: 3306
  username: "root"
  password: "your_password"
  dbname: "blog_db"

redis:
  host: "127.0.0.1"
  port: 6379
  password: ""
```

### 5. 运行项目

```bash
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动。

## API 文档

### 公开接口

#### 用户注册
```
POST /api/v1/register
Content-Type: application/json

{
  "username": "user123",
  "password": "password123",
  "email": "user@example.com",
  "nickname": "昵称"
}
```

#### 用户登录
```
POST /api/v1/login
Content-Type: application/json

{
  "username": "user123",
  "password": "password123"
}

Response:
{
  "code": 200,
  "msg": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

#### 获取文章列表
```
GET /api/v1/articles?page=1&page_size=10&status=1&category_id=1
```

#### 获取文章详情
```
GET /api/v1/articles/:id
```

#### 获取分类列表
```
GET /api/v1/categories
```

### 需要认证的接口

需要在请求头中添加：
```
Authorization: Bearer <token>
```

#### 获取当前用户信息
```
GET /api/v1/user/profile
```

#### 更新用户信息
```
PUT /api/v1/user/profile
Content-Type: application/json

{
  "nickname": "新昵称",
  "email": "new@example.com",
  "avatar": "avatar_url"
}
```

#### 修改密码
```
PUT /api/v1/user/password
Content-Type: application/json

{
  "old_password": "old_password",
  "new_password": "new_password"
}
```

#### 创建文章
```
POST /api/v1/articles
Content-Type: application/json

{
  "title": "文章标题",
  "description": "文章描述",
  "content": "文章内容",
  "cover": "封面图片URL",
  "category_id": 1,
  "tag_ids": [1, 2, 3],
  "status": 1,
  "is_top": false
}
```

#### 更新文章
```
PUT /api/v1/articles/:id
Content-Type: application/json

{
  "title": "更新后的标题",
  "content": "更新后的内容",
  ...
}
```

#### 删除文章
```
DELETE /api/v1/articles/:id
```

#### 点赞文章
```
POST /api/v1/articles/:id/like
```

#### 文件上传
```
POST /api/v1/upload
Content-Type: multipart/form-data

file: <file>

Response:
{
  "code": 200,
  "msg": "success",
  "data": {
    "path": "2024-01-01/uuid.jpg",
    "url": "http://localhost:8080/uploads/2024-01-01/uuid.jpg"
  }
}
```

### 管理员接口

需要管理员角色权限。

#### 创建分类
```
POST /api/v1/admin/categories
Content-Type: application/json

{
  "name": "分类名称",
  "description": "分类描述",
  "sort": 0
}
```

#### 更新分类
```
PUT /api/v1/admin/categories/:id
Content-Type: application/json

{
  "name": "更新后的名称",
  "description": "更新后的描述",
  "sort": 1
}
```

#### 删除分类
```
DELETE /api/v1/admin/categories/:id
```

## 开发说明

### 数据模型

- **User**: 用户模型
- **Article**: 文章模型
- **Category**: 分类模型
- **Tag**: 标签模型
- **Comment**: 评论模型

### 添加新功能

1. 在 `internal/models/` 中定义数据模型
2. 在 `internal/services/` 中实现业务逻辑
3. 在 `internal/controllers/` 中创建控制器
4. 在 `internal/routes/routes.go` 中注册路由

## 配置说明

主要配置项说明：

- `app.mode`: 运行模式（debug, release, test）
- `app.port`: 服务端口
- `database`: 数据库配置
- `redis`: Redis配置
- `jwt.secret`: JWT密钥（生产环境请务必修改）
- `jwt.expire_hours`: Token过期时间（小时）
- `log`: 日志配置
- `upload`: 文件上传配置

## 注意事项

1. 生产环境请修改 `jwt.secret` 为强密码
2. 建议使用环境变量管理敏感配置
3. 定期备份数据库
4. 根据实际情况调整连接池参数

## 许可证

MIT License
