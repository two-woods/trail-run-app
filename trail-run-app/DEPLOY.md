# 越野跑App 部署指南

## 环境要求

| 组件 | 版本 | 说明 |
|------|------|------|
| Docker | 24+ | 容器化 |
| Docker Compose | 2+ | 服务编排 |
| MySQL | 8.0 | 数据库 |
| Redis | 7+ | 缓存 |

---

## 快速部署

### 1. 启动所有服务

```bash
cd /home/ubuntu/claude_test/trail-run-app

# 启动MySQL和Redis
docker-compose up -d mysql redis

# 等待数据库就绪
sleep 10

# 启动后端和图片服务
docker-compose up -d --build
```

### 2. 验证服务

```bash
# 检查服务状态
docker-compose ps

# 测试后端API
curl http://localhost:8081/health

# 测试图片服务
curl http://localhost:8082/health
```

### 3. 服务地址

| 服务 | 地址 |
|------|------|
| 后端API | http://localhost:8081 |
| 图片服务 | http://localhost:8082 |
| MySQL | localhost:3306 |
| Redis | localhost:6379 |

---

## 详细配置

### 环境变量

后端默认配置：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `DB_HOST` | localhost | 数据库地址 |
| `DB_PORT` | 3306 | 数据库端口 |
| `DB_USER` | root | 数据库用户 |
| `DB_PASSWORD` | trail_run_password | 数据库密码 |
| `DB_NAME` | trail_run | 数据库名 |
| `REDIS_HOST` | localhost | Redis地址 |
| `REDIS_PORT` | 6379 | Redis端口 |
| `JWT_SECRET` | - | JWT密钥（生产环境必改） |
| `PORT` | 8081 | 后端监听端口 |
| `AMAP_KEY` | - | 高德地图API Key |

### 修改环境变量

创建 `.env` 文件：

```bash
cat > /home/ubuntu/claude_test/trail-run-app/.env << EOF
DB_HOST=mysql
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_secure_password
DB_NAME=trail_run
REDIS_HOST=redis
REDIS_PORT=6379
JWT_SECRET=your-very-secure-jwt-secret-change-this
PORT=8081
AMAP_KEY=your_amap_key
AMAP_KEY_SECRET=your_amap_key_secret
IMAGE_SERVICE_URL=http://localhost:8082
EOF
```

然后修改 `docker-compose.yml` 中的密码为一致。

---

## 数据库初始化

数据库表由 GORM AutoMigrate 自动创建。

### 手动初始化

```bash
# 进入MySQL容器
docker exec -it trail-run-mysql mysql -u root -p

# 创建数据库
CREATE DATABASE trail_run CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 种子数据

插入示例赛事数据：

```bash
docker exec trail-run-mysql mysql -u root -ptrail_run_password trail_run -e "
INSERT INTO races (id, name, date, location, province, city, distance_km, elevation_m, difficulty, start_lat, start_lng, end_lat, end_lng, status)
VALUES 
('00000000-0000-0000-0000-000000000001', '武汉马拉松', '2026-04-15', '武汉市', '湖北省', '武汉市', 42.195, 156, 'easy', 30.5728, 114.2525, 30.5728, 114.2525, 'published'),
('00000000-0000-0000-0000-000000000002', '北京TNF100', '2026-05-01', '北京市', '北京市', '北京市', 100, 4500, 'hard', 39.9042, 116.4074, 39.9042, 116.4074, 'published'),
('00000000-0000-0000-0000-000000000003', '杭州越野赛', '2026-04-20', '杭州市', '浙江省', '杭州市', 50, 2800, 'medium', 30.2741, 120.1551, 30.2741, 120.1551, 'published')
ON DUPLICATE KEY UPDATE name=name;
"
```

---

## 常用操作

### 查看日志

```bash
# 所有服务日志
docker-compose logs -f

# 指定服务日志
docker-compose logs -f backend
docker-compose logs -f python-image-service

# 最近100行
docker-compose logs --tail=100
```

### 重启服务

```bash
# 重启所有服务
docker-compose restart

# 重启指定服务
docker-compose restart backend
docker-compose restart python-image-service
```

### 停止服务

```bash
# 停止所有服务（保留数据）
docker-compose stop

# 停止并删除容器
docker-compose down

# 删除容器和网络（保留数据卷）
docker-compose down --remove-orphans

# 完全清理（删除数据）
docker-compose down -v
```

### 重建服务

```bash
# 重新构建并启动
docker-compose up -d --build
```

---

## 生产环境部署

### 1. 安全配置

```bash
# 修改JWT密钥
export JWT_SECRET=$(openssl rand -base64 32)

# 配置防火墙
sudo ufw allow 8081/tcp  # 后端API
sudo ufw allow 8082/tcp  # 图片服务
```

### 2. 使用Nginx反向代理

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location /api/ {
        proxy_pass http://127.0.0.1:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /img/ {
        proxy_pass http://127.0.0.1:8082;
        proxy_set_header Host $host;
    }
}
```

### 3. HTTPS配置

使用Let's Encrypt：

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

### 4. 开机自启

```bash
# 创建systemd服务
sudo cat > /etc/systemd/system/trail-run-app.service << EOF
[Unit]
Description=Trail Run App
Requires=docker-compose.target
After=docker-compose.target

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/home/ubuntu/claude_test/trail-run-app
ExecStart=/usr/bin/docker-compose up -d
ExecStop=/usr/bin/docker-compose stop
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
EOF

# 启用服务
sudo systemctl enable trail-run-app.service
sudo systemctl start trail-run-app.service
```

---

## 前端部署

### HBuilderX 构建

1. 在 Windows/Mac 用 HBuilderX 打开 `frontend/` 目录
2. 修改 `services/api.js` 中的 `BASE_URL` 为服务器地址
3. 点击 **发行** → **网站-H5**
4. 将构建产物上传到服务器的 nginx 目录

### Nginx 配置

```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /var/www/trail-run-app;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://127.0.0.1:8081;
    }

    location /img/ {
        proxy_pass http://127.0.0.1:8082;
    }
}
```

---

## 故障排查

### 后端无法启动

```bash
# 检查端口占用
lsof -i :8081
netstat -tlnp | grep 8081

# 查看后端日志
docker-compose logs backend
```

### 数据库连接失败

```bash
# 检查MySQL状态
docker-compose ps mysql

# 测试连接
docker exec -it trail-run-mysql mysql -u root -p
```

### 图片服务报错

```bash
# 查看日志
docker-compose logs python-image-service

# 常见错误：字体缺失
docker exec trail-run-image-service fc-list
```

---

## 完整部署脚本

```bash
#!/bin/bash
# deploy.sh - 越野跑App部署脚本

set -e

cd /home/ubuntu/claude_test/trail-run-app

echo "=== 1. 停止旧服务 ==="
docker-compose down --remove-orphans 2>/dev/null || true

echo "=== 2. 启动MySQL和Redis ==="
docker-compose up -d mysql redis

echo "=== 3. 等待数据库就绪 ==="
sleep 15

echo "=== 4. 启动所有服务 ==="
docker-compose up -d --build

echo "=== 5. 等待服务启动 ==="
sleep 10

echo "=== 6. 验证服务 ==="
echo -n "后端: "
curl -s http://localhost:8081/health | grep -q "ok" && echo "OK" || echo "FAIL"
echo -n "图片服务: "
curl -s http://localhost:8082/health | grep -q "ok" && echo "OK" || echo "FAIL"

echo ""
echo "=== 部署完成 ==="
echo "后端API: http://localhost:8081"
echo "图片服务: http://localhost:8082"
```

---

## API文档

后端API运行后可访问：

- Swagger文档: http://localhost:8081/swagger/index.html (如已配置)
- API测试: http://localhost:8081/health

### 主要接口

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/auth/register | 用户注册 |
| POST | /api/v1/auth/login | 用户登录 |
| GET | /api/v1/users/me | 获取当前用户 |
| GET | /api/v1/races | 赛事列表 |
| GET | /api/v1/races/:id | 赛事详情 |
| GET | /api/v1/races/:id/plan | 当日规划 |
| GET | /api/v1/races/:id/equipment | 装备清单 |
| POST | /api/v1/races/:id/equipment/check | 更新装备状态 |
| GET | /api/v1/results | 我的比赛 |
| POST | /api/v1/results | 创建完赛记录 |
| GET | /api/v1/results/:id | 比赛详情 |
| POST | /api/v1/results/:id/image | 生成完赛图片 |
