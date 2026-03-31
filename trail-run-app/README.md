# 越野跑 App (Trail Run App)

一款为越野跑爱好者设计的跨平台移动应用，提供赛事信息、装备清单、通勤规划、完赛记录和轨迹分享等功能。

## 快速部署

详见 [DEPLOY.md](./DEPLOY.md)

## 功能特性

### 核心功能
- **用户认证** - 注册、登录、JWT认证
- **赛事管理** - 赛事列表、详情查看、赛事规划
- **赛程当日规划**
  - 装备清单（智能推荐、勾选管理）
  - 通勤方式（起点导航、停车场指引）
  - 补给站信息（地图标记）
  - 天气预报
- **完赛记录** - 成绩管理、GPX轨迹解析
- **交互式地图** - 轨迹展示、起点终点标记
- **图片生成** - 完赛证书自动生成

## 技术栈

### 后端
- Go 1.21+
- Gin Web Framework
- GORM (MySQL)
- Redis (缓存/会话)
- JWT (认证)

### 前端
- Uni-app (Vue 3)
- 高德地图 SDK

### 服务
- Python FastAPI (图片生成)

### 基础设施
- MySQL 8.0
- Redis 7
- Docker Compose

## 项目结构

```
trail-run-app/
├── backend/              # Go 后端
│   ├── cmd/server/       # 入口
│   ├── internal/
│   │   ├── handler/      # HTTP处理器
│   │   ├── model/        # 数据模型
│   │   └── service/      # 业务逻辑
│   └── pkg/
│       ├── middleware/   # 中间件
│       └── amap/        # 高德地图SDK
├── frontend/            # Uni-app 前端
│   ├── pages/           # 页面
│   ├── components/      # 组件
│   ├── services/        # API服务
│   └── utils/           # 工具
├── image-service/       # Python 图片生成服务
└── docker-compose.yml   # Docker配置
```

## 快速开始

### 环境要求
- Docker & Docker Compose
- Go 1.21+
- Node.js 18+
- Python 3.10+ (本地开发)

### 1. 启动基础设施

```bash
docker-compose up -d mysql redis
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 配置数据库密码等
```

### 3. 启动后端

```bash
cd backend
go run cmd/server/main.go
```

后端服务运行在 http://localhost:8080

### 4. 启动图片服务（可选）

```bash
cd image-service
pip install -r requirements.txt
uvicorn app.main:app --reload --port 8082
```

### 5. 前端开发

```bash
cd frontend
npm install
npm run dev
```

## API 文档

启动后端后访问: http://localhost:8080/swagger/index.html

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

## 配置说明

### 高德地图配置

在 `backend/pkg/amap/client.go` 中配置您的API Key:

```go
AmapKey := os.Getenv("AMAP_KEY")
AmapKeySecret := os.Getenv("AMAP_KEY_SECRET")
```

### 数据库配置

```bash
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=trail_run
```

### Redis配置

```bash
REDIS_HOST=localhost
REDIS_PORT=6379
```

## 开发说明

### 数据库迁移

首次启动时，GORM 会自动创建所有表结构。

### 前端路由

页面配置在 `pages.json` 中定义，使用 Uni-app 的路由系统。

### 地图功能

前端使用高德地图 Uni-app SDK，需在 manifest.json 中配置 AppKey。

## 许可证

MIT
