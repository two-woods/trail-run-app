# Quickstart: 越野跑App MVP

**Date**: 2026-03-30

---

## 1. 项目结构

```
trail-run-app/
├── backend/            # Go后端 (cmd/server + internal/)
├── image-service/     # Python图片服务
├── frontend/          # Uni-app项目
└── docs/              # 设计文档
```

---

## 2. 环境准备

### 必需工具

| 工具 | 版本 | 用途 |
|------|------|------|
| Go | 1.21+ | 后端开发 |
| Python | 3.11+ | 图片服务 |
| Docker | 24+ | 容器化 |
| Docker Compose | 2+ | 本地编排 |
| HBuilderX | 最新 | Uni-app开发 |
| MySQL | 8.0 | 数据库 |
| Redis | 7+ | 缓存 |

### 环境变量

```bash
# backend/.env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=trail_run
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=your_jwt_secret
OSS_ACCESS_KEY=your_key
OSS_SECRET_KEY=your_secret
OSS_BUCKET=trail-run
IMAGE_SERVICE_URL=http://localhost:5000

# image-service/.env
OSS_ACCESS_KEY=your_key
OSS_SECRET_KEY=your_secret
OSS_BUCKET=trail-run
OSS_ENDPOINT=oss-cn-hangzhou.aliyuncs.com
```

---

## 3. 后端启动

```bash
cd backend

# 安装依赖
go mod download

# 数据库迁移 (待实现)
# migrate -database $DB_URL up

# 启动服务
go run cmd/server/main.go
```

服务运行在 `http://localhost:8080`

---

## 4. 图片服务启动

```bash
cd image-service

# 安装依赖
pip install -r requirements.txt

# 启动服务
python -m uvicorn app.main:app --host 0.0.0.0 --port 5000
```

服务运行在 `http://localhost:5000`

---

## 5. Docker Compose 启动 (推荐)

```bash
docker-compose up -d
```

启动:
- Go backend: `:8080`
- Python image service: `:5000`
- MySQL: `:3306`
- Redis: `:6379`

---

## 6. 前端开发

### HBuilderX 打开项目

```bash
cd frontend
# 用 HBuilderX 打开 frontend 目录
```

### 配置 API 地址

```javascript
// frontend/src/config.js
export const API_BASE_URL = 'http://localhost:8080/api/v1'
```

### 运行预览

- **浏览器**: 点击 "运行" → 选择 "运行到浏览器"
- **小程序**: 点击 "运行" → 选择 "微信小程序" (需安装开发者工具)
- **真机**: 连接手机后选择对应平台

---

## 7. 数据库初始化

```sql
CREATE DATABASE trail_run CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 核心表结构 (详见 data-model.md)
CREATE TABLE users (...);
CREATE TABLE races (...);
CREATE TABLE aid_stations (...);
CREATE TABLE equipment (...);
CREATE TABLE results (...);
CREATE TABLE equipment_checks (...);
CREATE TABLE favorite_races (...);
```

---

## 8. 测试API

```bash
# 注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123","nickname":"测试用户"}'

# 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'

# 获取赛事列表 (需token)
curl http://localhost:8080/api/v1/races \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 9. 常用命令

```bash
# 后端测试
cd backend && go test ./...

# 图片服务测试
cd image-service && python -m pytest

# Docker 重启服务
docker-compose restart backend

# 查看日志
docker-compose logs -f backend
```

---

## 10. 下一步

1. 实现用户注册/登录 (POST /auth/*)
2. 实现赛事CRUD (GET /races)
3. 实现当日规划 (GET /races/:id/plan)
4. 实现完赛记录 (GET/POST /results)
5. 实现图片生成 (POST /results/:id/image)
