# Docker 部署指南

本文档介绍如何使用 Docker 将博客后台系统部署到远程服务器(适用于宝塔面板)。

## 部署架构

- **Go 应用**: Docker 容器运行
- **MySQL**: 宝塔面板直接安装
- **Redis**: 宝塔面板直接安装

## 准备工作

### 服务器要求
- Linux 服务器 (Ubuntu 20.04+, CentOS 7+, Debian 10+ 等)
- 已安装宝塔面板
- 开放端口: 8080 (应用端口)

## 一、宝塔面板部署步骤

### 1. 安装必要软件

在宝塔面板的软件商店安装:
- **MySQL 8.0**
- **Redis**
- **Docker**

### 2. 创建数据库

在宝塔面板 → 数据库:
```sql
数据库名: blog_db
用户名: blog_user
密码: 设置一个强密码
字符集: utf8mb4
```

### 3. 上传项目

方式1: 使用宝塔文件管理器上传

方式2: 使用命令行
```bash
# 本地打包
tar -czf blog-backend.tar.gz --exclude='.git' --exclude='logs' --exclude='uploads' .

# 上传到服务器 (假设路径为 /www/wwwroot/blog-backend)
scp blog-backend.tar.gz root@your-server-ip:/www/wwwroot/

# 服务器上解压
ssh root@your-server-ip
cd /www/wwwroot
tar -xzf blog-backend.tar.gz
mv blog-backend blog-backend-dir  # 重命名
```

### 4. 修改配置文件

编辑 `config/config.yaml`:
```yaml
app:
  mode: "release"
  port: 8080

database:
  host: "127.0.0.1"      # 使用宿主机的 MySQL
  username: "blog_user"
  password: "你的数据库密码"
  dbname: "blog_db"

redis:
  host: "127.0.0.1"      # 使用宿主机的 Redis
  port: 6379

jwt:
  secret: "修改为强密码"
```

### 5. 构建 Docker 镜像

在宝塔终端或 SSH 中执行:
```bash
cd /www/wwwroot/blog-backend

# 构建镜像
docker build -t blog-backend:latest .

# 查看镜像
docker images | grep blog-backend
```

### 6. 运行容器

使用 `--network host` 模式(推荐,容器直接使用宿主机网络):
```bash
docker run -d \
  --name blog-backend \
  --network host \
  --restart always \
  -v /www/wwwroot/blog-backend/uploads:/app/uploads \
  -v /www/wwwroot/blog-backend/logs:/app/logs \
  -v /www/wwwroot/blog-backend/config:/app/config:ro \
  blog-backend:latest
```

或使用端口映射:
```bash
docker run -d \
  --name blog-backend \
  -p 8080:8080 \
  --restart always \
  --add-host host.docker.internal:host-gateway \
  -v /www/wwwroot/blog-backend/uploads:/app/uploads \
  -v /www/wwwroot/blog-backend/logs:/app/logs \
  -v /www/wwwroot/blog-backend/config:/app/config:ro \
  blog-backend:latest
```

### 7. 配置防火墙

在宝塔面板 → 安全 → 防火墙 → 添加规则:
- 端口: `8080`
- 协议: `TCP`
- 说明: `博客后台`

### 8. 验证部署

```bash
# 查看容器状态
docker ps | grep blog-backend

# 查看日志
docker logs -f blog-backend

# 测试接口
curl http://localhost:8080/api/v1/ping
```

或在浏览器访问: `http://your-server-ip:8080`

## 二、常用命令

### 容器管理
```bash
# 查看运行中的容器
docker ps

# 启动容器
docker start blog-backend

# 停止容器
docker stop blog-backend

# 重启容器
docker restart blog-backend

# 删除容器
docker rm -f blog-backend

# 查看日志
docker logs -f blog-backend
docker logs --tail 100 blog-backend  # 查看最后100行

# 进入容器
docker exec -it blog-backend sh
```

### 更新应用
```bash
# 停止并删除旧容器
docker stop blog-backend
docker rm blog-backend

# 重新构建镜像
cd /www/wwwroot/blog-backend
docker build -t blog-backend:latest .

# 启动新容器
docker run -d \
  --name blog-backend \
  --network host \
  --restart always \
  -v /www/wwwroot/blog-backend/uploads:/app/uploads \
  -v /www/wwwroot/blog-backend/logs:/app/logs \
  -v /www/wwwroot/blog-backend/config:/app/config:ro \
  blog-backend:latest
```

### 清理资源
```bash
# 清理未使用的镜像
docker image prune -a

# 清理所有未使用的资源
docker system prune -a

# 查看磁盘使用
docker system df
```

## 三、SSL/HTTPS 配置(使用宝塔)

### 1. 在宝塔创建网站

宝塔面板 → 网站 → 添加站点:
- 域名: `your-domain.com`
- 根目录: 任意(不使用)
- PHP版本: 纯静态

### 2. 配置反向代理

进入网站设置 → 反向代理 → 添加反向代理:
- 代理名称: `博客后台`
- 目标URL: `http://127.0.0.1:8080`
- 发送域名: `$host`
- 配置内容:
```nginx
location / {
    proxy_pass http://127.0.0.1:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

### 3. 申请 SSL 证书

网站设置 → SSL → Let's Encrypt:
- 点击"申请"
- 勾选"强制HTTPS"

现在可以通过 `https://your-domain.com` 访问了!

## 四、监控和维护

### 查看资源使用

```bash
# 查看容器资源使用
docker stats blog-backend

# 查看磁盘使用
df -h
du -sh /www/wwwroot/blog-backend/logs
du -sh /www/wwwroot/blog-backend/uploads
```

### 日志管理

```bash
# 应用日志(自动管理,配置在 config.yaml 中)
tail -f /www/wwwroot/blog-backend/logs/app.log

# Docker 容器日志
docker logs --tail 200 blog-backend

# 清理 Docker 日志(如果过大)
truncate -s 0 $(docker inspect --format='{{.LogPath}}' blog-backend)
```

### 数据备份

```bash
# 备份 MySQL 数据
mysqldump -ublog_user -p blog_db > backup_$(date +%Y%m%d).sql

# 或使用宝塔面板:数据库 → 备份

# 备份上传文件
tar -czf uploads_backup_$(date +%Y%m%d).tar.gz /www/wwwroot/blog-backend/uploads/

# 恢复数据库
mysql -ublog_user -p blog_db < backup_20240101.sql
```

## 五、故障排查

### 应用无法启动

```bash
# 查看容器日志
docker logs blog-backend

# 常见问题:
# 1. 数据库连接失败 -> 检查 config.yaml 和数据库配置
# 2. 端口被占用 -> 检查 8080 端口: netstat -tulnp | grep 8080
# 3. 权限问题 -> 检查 uploads 和 logs 目录权限: chmod 755
```

### 数据库连接失败

```bash
# 检查 MySQL 是否运行
systemctl status mysql
# 或在宝塔面板检查 MySQL 状态

# 测试连接
mysql -ublog_user -p -e "SELECT 1"

# 检查防火墙(如果容器用端口映射模式)
iptables -L | grep 3306
```

### Redis 连接失败

```bash
# 检查 Redis 状态
redis-cli ping
# 或在宝塔面板检查 Redis 状态

# 如果使用密码
redis-cli -a your_password ping
```

### 容器网络问题

如果使用 `--network host` 模式但无法连接数据库:
```bash
# 检查 MySQL 是否监听 127.0.0.1
netstat -tulnp | grep 3306

# 如果只监听 localhost,需要配置 MySQL 允许本地连接
# 编辑 /etc/mysql/mysql.conf.d/mysqld.cnf
bind-address = 127.0.0.1
```

如果使用端口映射模式:
```bash
# 确保容器可以访问宿主机
docker exec -it blog-backend ping host.docker.internal

# 如果不通,尝试使用宿主机 IP
ip addr show docker0  # 查看 docker0 的 IP
# 在 config.yaml 中使用这个 IP (通常是 172.17.0.1)
```

## 六、安全建议

1. **修改默认密码**: 数据库密码和 JWT Secret
2. **配置防火墙**: 只开放必要的端口 (80, 443)
3. **使用 HTTPS**: 配置 SSL 证书
4. **定期更新**: 更新系统和 Docker 镜像
5. **定期备份**: 备份数据库和上传文件
6. **限制访问**: 数据库只允许本地访问
7. **日志监控**: 定期检查应用日志

## 七、性能优化

### 数据库优化

在宝塔面板 → 数据库 → 性能调整:
- 调整 `max_connections`
- 调整 `innodb_buffer_pool_size`

或编辑 MySQL 配置文件。

### 应用优化

在 `config/config.yaml` 中调整:
```yaml
database:
  max_idle_conns: 20
  max_open_conns: 200
  conn_max_lifetime: 3600
```

### Redis 优化

在宝塔面板 → Redis → 配置修改:
```
maxmemory 256mb
maxmemory-policy allkeys-lru
```

## 常见问题 FAQ

### 1. 如何更新代码?

```bash
# 上传新代码到服务器
# 然后执行:
cd /www/wwwroot/blog-backend
docker stop blog-backend
docker rm blog-backend
docker build -t blog-backend:latest .
docker run -d --name blog-backend --network host --restart always \
  -v /www/wwwroot/blog-backend/uploads:/app/uploads \
  -v /www/wwwroot/blog-backend/logs:/app/logs \
  -v /www/wwwroot/blog-backend/config:/app/config:ro \
  blog-backend:latest
```

### 2. 如何查看实时日志?

```bash
docker logs -f blog-backend
# 或
tail -f /www/wwwroot/blog-backend/logs/app.log
```

### 3. 如何修改配置?

```bash
# 修改配置文件
vim /www/wwwroot/blog-backend/config/config.yaml

# 重启容器生效
docker restart blog-backend
```

### 4. 容器自动重启吗?

使用了 `--restart always` 参数,容器会在:
- Docker 重启后自动启动
- 容器异常退出后自动重启
- 服务器重启后自动启动

### 5. 如何绑定多个域名?

在宝塔面板网站设置中添加多个域名即可,反向代理配置保持不变。

## 支持

如有问题,请查看:
- 应用日志: `/www/wwwroot/blog-backend/logs/app.log`
- 容器日志: `docker logs blog-backend`
- 宝塔面板的错误日志

或提交 Issue 获取帮助。
