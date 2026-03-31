# Quickstart: 越野跑App MVP

**Date**: 2026-03-30
**Status**: Phase 1-7 Completed ✅

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
| Python | 3.10+ | 图片服务 |
| Docker | 24+ | 容器化 |
| Docker Compose | 2+ | 本地编排 |
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
AMAP_KEY=your_amap_key
AMAP_KEY_SECRET=your_amap_secret
IMAGE_SERVICE_URL=http://localhost:8082
PORT=8081
```

---

## 3. 后端启动

```bash
cd backend

# 安装依赖
go mod download

# 启动服务（自动创建数据库表）
go run cmd/server/main.go
```

服务运行在 `http://localhost:8081`

---

## 4. 图片服务启动

```bash
cd image-service

# 安装依赖
pip install -r requirements.txt

# 启动服务
python -m uvicorn app.main:app --host 0.0.0.0 --port 8082 --reload
```

服务运行在 `http://localhost:8082`

---

## 5. Docker Compose 启动 (推荐)

```bash
docker-compose up -d
```

启动:
- Go backend: `:8081`
- Python image service: `:8082`
- MySQL: `:3306`
- Redis: `:6379`

---

## 6. 前端开发

### 使用 HBuilderX

```bash
cd frontend
# 用 HBuilderX 打开 frontend 目录
```

### 配置 API 地址

```javascript
// frontend/services/api.js
const BASE_URL = 'http://localhost:8081/api/v1'
```

### 运行预览

- **微信小程序**: 点击 "运行" → 选择 "微信小程序" (需安装开发者工具)
- **H5**: 点击 "运行" → 选择 "运行到浏览器"
- **真机**: 连接手机后选择对应平台

**注意**: Ubuntu服务器无HBuilderX，需在Windows/Mac开发

---

## 7. 数据库

数据库表由 GORM AutoMigrate 自动创建:

- `users` - 用户表
- `races` - 赛事表
- `aid_stations` - 补给站表
- `equipment` - 装备表
- `results` - 完赛记录表
- `equipment_checks` - 装备检查表
- `favorite_races` - 收藏赛事表

---

## 8. 测试API

```bash
# 健康检查
curl http://localhost:8081/health

# 注册
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123","nickname":"测试用户"}'

# 登录
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'

# 获取赛事列表 (需token)
curl http://localhost:8081/api/v1/races \
  -H "Authorization: Bearer YOUR_TOKEN"

# 获取我的比赛
curl http://localhost:8081/api/v1/results \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 9. 已实现功能

### Phase 1-2: 基础设施 ✅
- 项目结构搭建
- Go + Gin 后端
- MySQL + Redis 集成
- JWT 认证

### Phase 3: 用户认证 ✅
- 注册/登录
- Token 认证

### Phase 4: 赛事管理 ✅
- 赛事列表/详情
- 当日规划聚合

### Phase 5: 赛程当日规划 ✅
- 装备清单 + 查漏
- 高德地图集成
- 通勤导航
- 补给站标记
- 天气预报

### Phase 6: 完赛记录 + GPX ✅
- 成绩 CRUD
- GPX 轨迹解析
- GeoJSON 转换
- 交互式地图

### Phase 7: 图片生成 ✅
- 完赛证书生成
- 渐变轨迹可视化
- 海拔剖面图
- 保存相册

### Phase 8: 收尾优化 🔄
- 离线缓存策略
- README 完善
- Docker 配置

---

## 10. 常用命令

```bash
# 后端编译
cd backend && go build ./...

# 后端测试
cd backend && go test ./...

# Docker 重启
docker-compose restart backend

# 查看日志
docker-compose logs -f backend

# 清理Docker
docker-compose down -v
```

---

## 11. 下一步

- [ ] Phase 8: 单元测试
- [ ] Phase 8: 生产环境Docker配置
- [ ] 微信小程序发布配置
- [ ] 高德地图AppKey配置
