# Tasks: 越野跑App MVP

**Input**: Design documents from `/specs/002-trail-run-app/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 项目脚手架和基础设施搭建

- [x] T001 [P] 创建项目目录结构（backend/, image-service/, frontend/, docs/）
- [x] T002 [P] 初始化 Go(Gin) 后端项目，配置 go.mod 依赖
- [x] T003 [P] 初始化 Python FastAPI 图片服务项目
- [x] T004 [P] 初始化 Uni-app 前端项目
- [x] T005 配置 Docker Compose 本地开发环境（MySQL + Redis + 服务）
- [x] T006 配置环境变量管理（.env 示例文件）

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 核心基础设施，所有用户故事的前置依赖

**⚠️ CRITICAL**: 用户故事实现前必须完成此阶段

- [x] T007 初始化 MySQL 数据库和表结构（users, races, aid_stations, equipment, results, equipment_checks, favorite_races）
- [x] T008 [P] 实现 Go JWT 认证中间件
- [x] T009 [P] 配置 Gin 路由结构（api/v1/* 分组）
- [x] T010 创建 User/Race/AidStation/Equipment/Result 基础 Model（Go + GORM）
- [x] T011 实现统一错误处理和日志基础设施
- [x] T012 配置 Redis 连接和会话管理

**Checkpoint**: 基础就绪，用户故事可并行开始

---

## Phase 3: 用户认证模块

**Goal**: 用户注册、登录、个人信息管理

- [x] T013 [P] POST /api/v1/auth/register 注册接口
- [x] T014 [P] POST /api/v1/auth/login 登录接口
- [x] T015 GET /api/v1/users/me 获取当前用户信息
- [x] T016 Uni-app 端：实现登录/注册页面和 token 存储

---

## Phase 4: 赛事管理模块

**Goal**: 赛事列表查询、赛事详情

**依赖**: Phase 2 完成

- [ ] T017 [P] GET /api/v1/races 赛事列表（分页、筛选）
- [ ] T018 [P] GET /api/v1/races/:id 赛事详情（含补给站、装备清单）
- [ ] T019 [P] GET /api/v1/races/:id/plan 当日规划数据聚合接口
- [ ] T020 Uni-app 端：赛事列表页和赛事详情页

---

## Phase 5: 赛程当日规划（装备+通勤+补给+天气）

**Goal**: 比赛前的一站式信息聚合

**依赖**: Phase 4 完成

### 装备清单
- [x] T021 GET /api/v1/races/:id/equipment 获取装备清单
- [x] T022 POST /api/v1/races/:id/equipment/check 更新装备勾选状态
- [x] T023 实现装备查漏补缺逻辑（基于距离/天气/爬升）
- [x] T024 Uni-app 端：装备清单页面（分类展示、勾选交互）

### 通勤方式
- [x] T025 [P] 高德地图API对接（地理编码、POI搜索）
- [x] T026 实现起点导航、停车场指引接口
- [x] T027 Uni-app 端：一键跳转高德导航

### 补给站信息
- [x] T028 [P] 补给站标记点渲染（markers + callout）
- [x] T029 Uni-app 端：补给站列表和地图标记

### 天气预报
- [x] T030 [P] 高德天气API对接
- [x] T031 Uni-app 端：天气预报展示组件

---

## Phase 6: 我的比赛 + 完赛记录

**Goal**: 完赛记录管理、GPX轨迹解析、交互式地图展示

**依赖**: Phase 3 完成

### 完赛记录 CRUD
- [ ] T032 [P] GET /api/v1/results 我的比赛列表
- [ ] T033 [P] POST /api/v1/results 创建完赛记录
- [ ] T034 GET /api/v1/results/:id 完赛记录详情
- [ ] T035 Uni-app 端：我的比赛页面

### GPX轨迹解析
- [ ] T036 [P] Go gpxgo库集成，解析GPX文件
- [ ] T037 计算轨迹统计（距离、爬升、海拔变化）
- [ ] T038 转换GPX为GeoJSON供地图渲染

### 参赛地图展示
- [ ] T039 [P] Uni-app 高德地图组件集成
- [ ] T040 轨迹polyline渲染 + 渐变配色（海拔/速度）
- [ ] T041 起点/终点/补给站标记点
- [ ] T042 地图缩放、拖拽交互

---

## Phase 7: 轨迹图片生成

**Goal**: 生成可分享的炫酷完赛轨迹图片

**依赖**: Phase 6 完成

- [ ] T043 [P] Python FastAPI 图片生成服务搭建
- [ ] T044 [P] GPX轨迹解析 + 地图底图叠加
- [ ] T045 实现图片元素：赛事名称、日期、轨迹图、关键数据
- [ ] T046 POST /api/v1/results/:id/image 触发图片生成
- [ ] T047 OSS上传集成，生成可访问URL
- [ ] T048 Uni-app 端：图片预览和分享功能（保存相册/微信）

---

## Phase 8: 收尾与优化

**Purpose**: 跨功能 concern 和上线准备

- [ ] T049 [P] 离线缓存策略实现（Uni Storage）
- [ ] T050 [P] 代码清理和 README 完善
- [ ] T051 单元测试补全（Go testing）
- [ ] T052 Docker Compose 生产镜像构建配置
- [ ] T053 完善 quickstart.md

---

## Dependencies & Execution Order

### Phase 依赖

```
Phase 1 (Setup) ──────────────────────────────┐
    │                                            │
Phase 2 (Foundational) ────────────────────────┤
    │                                            │
    ├── Phase 3 (认证) ─────────────────────────┤
    ├── Phase 4 (赛事) ─────────────────────────┤
    │        │                                   │
    │   Phase 5 (当日规划) ◄─────────────────────┤
    │        │                                   │
    │   Phase 6 (完赛记录) ◄─────────────────────┤
    │        │                                   │
    │   Phase 7 (图片生成) ◄─────────────────────┘
    │
Phase 8 (收尾)
```

### 可并行任务（标记 [P]）

- T001, T002, T003, T004, T005, T006 可并行
- T008, T009, T010, T011, T012 可并行（Phase 2内）
- T013, T014, T015, T016 可并行
- T017, T018, T019, T020 可并行
- T021, T025, T028, T030, T036, T039, T043, T044 可并行

---

## MVP 优先实现路径

如果资源有限，按此顺序实现：

1. **第一波**：Phase 1 → Phase 2 → Phase 3（用户认证）
2. **第二波**：Phase 4（赛事基础）+ Phase 6（完赛记录 + 地图）
3. **第三波**：Phase 5（当日规划）
4. **第四波**：Phase 7（图片生成）
5. **第五波**：Phase 8（优化）
