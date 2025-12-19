# Docker 部署指南

本文档介绍如何使用 Docker 将博客后台系统部署到远程服务器。

## 准备工作

### 本地环境要求
- Docker 20.10+
- Docker Compose 2.0+

### 远程服务器要求
- Linux 服务器 (Ubuntu 20.04+, CentOS 7+, Debian 10+ 等)
- Docker 和 Docker Compose 已安装
- 开放端口: 8080 (应用), 3306 (MySQL, 可选), 6379 (Redis, 可选)

## 一、本地测试部署

在部署到远程服务器前,建议先在本地测试。

### 1. 修改配置文件

编辑 `config/config.yaml`,修改以下配置以适配 Docker 环境:

```yaml
app:
  mode: "release"  # 改为 release
  port: 8080       # 改为 8080

database:
  host: "mysql"        # 改为 mysql (Docker 服务名)
  username: "blog_user"    # 改为 blog_user
  password: "blog123456"   # 与 docker-compose.yml 中的密码一致

redis:
  host: "redis"    # 改为 redis (Docker 服务名)

jwt:
  secret: "CHANGE_THIS_TO_A_STRONG_SECRET"  # 修改为强密码
```

### 2. 修改敏感信息

编辑 `docker-compose.yml`,修改以下密码:
```yaml
MYSQL_ROOT_PASSWORD: your_strong_password
MYSQL_PASSWORD: your_strong_password
```

同时更新 `config/config.yaml` 中的 `database.password`。

### 3. 构建并启动

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 查看服务状态
docker-compose ps
```

### 4. 验证部署

```bash
# 检查应用是否正常运行
curl http://localhost:8080/api/v1/ping

# 或在浏览器访问
# http://localhost:8080
```

## 二、远程服务器部署

### 方法 1: 上传源码构建(推荐)

#### 1. 安装 Docker 和 Docker Compose

Ubuntu/Debian:
```bash
# 更新包索引
sudo apt-get update

# 安装 Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 启动 Docker
sudo systemctl start docker
sudo systemctl enable docker

# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 验证安装
docker --version
docker-compose --version
```

CentOS/RHEL:
```bash
# 安装 Docker
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo yum install -y docker-ce docker-ce-cli containerd.io

# 启动 Docker
sudo systemctl start docker
sudo systemctl enable docker

# 安装 Docker Compose (同上)
```

#### 2. 上传项目到服务器

```bash
# 在本地打包项目
tar -czf blog-backend.tar.gz --exclude='.git' --exclude='logs' --exclude='uploads' .

# 上传到服务器
scp blog-backend.tar.gz user@your-server-ip:/home/user/

# 或使用 rsync (推荐)
rsync -avz --exclude='.git' --exclude='logs' --exclude='uploads' ./ user@your-server-ip:/home/user/blog-backend/
```

#### 3. 在服务器上部署

```bash
# SSH 登录服务器
ssh user@your-server-ip

# 解压项目 (如果使用 tar)
tar -xzf blog-backend.tar.gz -C /home/user/blog-backend
cd /home/user/blog-backend

# 修改配置
cp config/config.docker.yaml config/config.yaml
vim config/config.yaml  # 修改 JWT secret 等敏感信息
vim docker-compose.yml  # 修改数据库密码

# 构建并启动
docker-compose up -d

# 查看日志
docker-compose logs -f app
```

### 方法 2: 使用预构建镜像

如果你有 Docker 镜像仓库(如 Docker Hub, 阿里云容器镜像服务):

#### 1. 本地构建并推送镜像

```bash
# 构建镜像
docker build -t your-registry/blog-backend:latest .

# 推送到仓库
docker push your-registry/blog-backend:latest
```

#### 2. 修改 docker-compose.yml

```yaml
services:
  app:
    image: your-registry/blog-backend:latest  # 使用远程镜像
    # build:
    #   context: .
    #   dockerfile: Dockerfile
```

#### 3. 在服务器上拉取并运行

```bash
# 只需要 docker-compose.yml 和 config 目录
ssh user@your-server-ip
mkdir -p blog-backend && cd blog-backend

# 上传 docker-compose.yml 和 config 目录
# 然后启动
docker-compose pull
docker-compose up -d
```

## 三、配置说明

### 端口映射

如果服务器端口被占用,可以修改 `docker-compose.yml` 中的端口映射:

```yaml
services:
  app:
    ports:
      - "9000:8080"  # 将 8080 改为 9000

  mysql:
    ports:
      - "13306:3306"  # 避免与宿主机 MySQL 冲突
```

### 数据持久化

数据卷确保数据不会丢失:
- `mysql-data`: MySQL 数据
- `redis-data`: Redis 数据
- `./uploads`: 上传的文件
- `./logs`: 应用日志

### 环境变量

可以通过环境变量覆盖配置:

```yaml
services:
  app:
    environment:
      - GIN_MODE=release
      - DB_HOST=mysql
      - DB_PASSWORD=your_password
```

需要修改代码以支持环境变量配置。

## 四、常用命令

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 重启服务
docker-compose restart

# 重启单个服务
docker-compose restart app

# 查看日志
docker-compose logs -f
docker-compose logs -f app     # 只看应用日志
docker-compose logs -f mysql   # 只看 MySQL 日志

# 查看服务状态
docker-compose ps

# 进入容器
docker-compose exec app sh
docker-compose exec mysql bash

# 更新代码后重新构建
docker-compose down
docker-compose build --no-cache
docker-compose up -d

# 清理无用的镜像和容器
docker system prune -a
```

## 五、数据库初始化

首次部署后,需要创建表结构:

### 方法 1: 自动迁移 (推荐)

应用启动时会自动创建表(如果代码中有 AutoMigrate)。

### 方法 2: 手动执行 SQL

```bash
# 创建初始化 SQL 文件
vim init.sql

# 进入 MySQL 容器执行
docker-compose exec mysql mysql -ublog_user -pblog123456 blog_db < init.sql

# 或直接进入容器
docker-compose exec mysql mysql -ublog_user -pblog123456 blog_db
```

## 六、SSL/HTTPS 配置

### 使用 Nginx 反向代理 (推荐)

```bash
# 安装 Nginx
sudo apt-get install nginx

# 配置 Nginx
sudo vim /etc/nginx/sites-available/blog-backend
```

Nginx 配置示例:
```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

使用 Let's Encrypt 配置 HTTPS:
```bash
# 安装 Certbot
sudo apt-get install certbot python3-certbot-nginx

# 获取证书并自动配置 Nginx
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo certbot renew --dry-run
```

## 七、监控和维护

### 查看资源使用

```bash
# 查看容器资源使用
docker stats

# 查看磁盘使用
df -h
du -sh logs/
du -sh uploads/
```

### 日志管理

```bash
# 清理旧日志 (应用日志已配置自动清理)
# 清理 Docker 日志
sudo truncate -s 0 /var/lib/docker/containers/*/*-json.log
```

### 数据备份

```bash
# 备份 MySQL 数据
docker-compose exec mysql mysqldump -ublog_user -pblog123456 blog_db > backup_$(date +%Y%m%d).sql

# 备份上传文件
tar -czf uploads_backup_$(date +%Y%m%d).tar.gz uploads/

# 恢复数据库
docker-compose exec -T mysql mysql -ublog_user -pblog123456 blog_db < backup_20240101.sql
```

## 八、故障排查

### 应用无法启动

```bash
# 查看详细日志
docker-compose logs -f app

# 常见问题:
# 1. 数据库连接失败 -> 检查 config.yaml 配置
# 2. 端口被占用 -> 修改 docker-compose.yml 端口映射
# 3. 权限问题 -> 检查 uploads 和 logs 目录权限
```

### 数据库连接失败

```bash
# 检查 MySQL 是否正常
docker-compose ps
docker-compose logs mysql

# 测试连接
docker-compose exec mysql mysql -ublog_user -pblog123456 -e "SELECT 1"

# 检查网络
docker-compose exec app ping mysql
```

### Redis 连接失败

```bash
# 检查 Redis 状态
docker-compose exec redis redis-cli ping

# 查看 Redis 日志
docker-compose logs redis
```

## 九、安全建议

1. **修改默认密码**: 务必修改 docker-compose.yml 中的数据库密码
2. **修改 JWT Secret**: 修改 config.yaml 中的 jwt.secret
3. **配置防火墙**: 只开放必要的端口 (80, 443, 8080)
4. **使用非 root 用户**: 运行 Docker 容器
5. **定期更新**: 定期更新基础镜像和依赖
6. **备份数据**: 定期备份数据库和上传文件
7. **使用 HTTPS**: 生产环境务必配置 SSL 证书

## 十、性能优化

### 1. 数据库优化

```yaml
# 调整 MySQL 配置
services:
  mysql:
    command:
      - --max_connections=500
      - --innodb_buffer_pool_size=1G
```

### 2. 应用优化

在 `config/config.yaml` 中调整:
```yaml
database:
  max_idle_conns: 20
  max_open_conns: 200
```

### 3. Redis 优化

```yaml
services:
  redis:
    command: redis-server --maxmemory 256mb --maxmemory-policy allkeys-lru
```

## 支持

如有问题,请提交 Issue 或查看应用日志进行排查。
