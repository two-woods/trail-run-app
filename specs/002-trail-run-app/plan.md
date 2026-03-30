# Implementation Plan: 越野跑App MVP

**Branch**: `002-trail-run-app` | **Date**: 2026-03-30 | **Spec**: [spec.md](../spec.md)
**Input**: Feature specification from `/specs/002-trail-run-app/spec.md`

## Summary

为越野跑爱好者开发一款多端App（iOS/Android/小程序），核心功能包括赛程当日规划（装备清单、通勤、补给站、天气）和完赛记录管理（轨迹展示+可分享图片）。

技术方案：Uni-app (Vue) 前端 + Go后端 + Python图片生成服务 + MySQL + Redis + 高德地图

## Technical Context

**Language/Version**:
- 后端: Go 1.21+
- 前端: Vue 3 + Uni-app
- 图片服务: Python 3.11+

**Primary Dependencies**:
- 后端: Gorilla Mux, GORM, go-redis, go-oss-sdk
- 前端: uni-app, 高德地图小程序SDK
- 图片: Pillow, Matplotlib, cairosvg

**Storage**: MySQL (主数据) + Redis (缓存/会话) + OSS/COS (文件存储)

**Testing**: Go testing (backend), Uni-app test (frontend)

**Target Platform**: iOS 15+, Android 10+, 微信小程序, 支付宝小程序 + Linux server

**Project Type**: 移动App + REST API后端 + Python微服务

**Performance Goals**: MVP阶段无严格性能要求，关注功能完整性

**Scale/Scope**: MVP阶段小规模，后续可扩展

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

 constitution.md 为模板文件，尚未填充实际约束规则。
 本项目将遵循以下基本原则：
 1. MVP优先：先实现核心功能，不过度设计
 2. 前后端分离：API优先，接口先行
 3. 渐进式开发：先跑通流程，再优化体验

**Status**: ✅ PASS (constitution为模板，无冲突约束)

## Project Structure

### Documentation (this feature)

```text
specs/002-trail-run-app/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md            # Phase 2 output
```

### Source Code (repository root)

```text
trail-run-app/
├── backend/            # Go后端服务
│   ├── cmd/
│   │   └── server/     # 主入口
│   ├── internal/
│   │   ├── handler/    # HTTP handlers
│   │   ├── service/    # 业务逻辑
│   │   ├── repository/ # 数据访问
│   │   └── model/      # 数据模型
│   ├── pkg/
│   │   └── utils/      # 工具函数
│   └── go.mod
├── image-service/       # Python图片生成服务
│   ├── app/
│   │   └── generator.py
│   ├── requirements.txt
│   └── Dockerfile
├── frontend/           # Uni-app项目
│   ├── pages/
│   ├── components/
│   ├── services/       # API调用
│   └── utils/
└── docs/               # 设计文档
    └── specs/
```

**Structure Decision**:
- 前后端分离架构
- Go后端提供REST API
- Python作为独立图片生成微服务
- Uni-app管理多端页面

## Phase 0: Research

### Unknowns Identified

| Item | Status | Notes |
|------|--------|-------|
| 高德地图小程序SDK对接方式 | NEEDS CLARIFICATION | key申请、数据格式、API限制 |
| Python图片服务部署方式 | NEEDS CLARIFICATION | 独立服务 vs sidecar vs 内联 |
| GPX轨迹解析方案 | NEEDS CLARIFICATION | Go库选择或自解析 |
| 天气API数据源 | NEEDS CLARIFICATION | 高德天气 vs 第三方 |
| 离线缓存策略 | NEEDS CLARIFICATION | Uni-app Storage方案 |

### Research Tasks

1. **Uni-app + Go后端集成模式** — REST API设计规范
2. **高德地图小程序SDK** — 使用方法和对接细节
3. **Go + Python微服务通信** — gRPC vs HTTP
4. **GPX轨迹解析** — Go库或自解析方案
5. **天气API对接** — 高德天气API vs OpenWeatherMap

## Phase 1: Design & Contracts

### Data Model (from spec.md entities)

```
User
├── id (uuid)
├── nickname
├── avatar_url
├──itra_account (optional)
└── privacy_settings

Race
├── id (uuid)
├── name
├── date
├── location
├── distance_km
├── elevation_m
├── difficulty
├──itra_points
├── route_gpx_url
└── status (draft/published)

AidStation
├── id (uuid)
├── race_id (fk)
├── name
├── distance_km
├── elevation_m
├── supplies[]
└── close_time

Equipment
├── id (uuid)
├── race_id (fk)
├── name
├── is_mandatory
└── category

Result
├── id (uuid)
├── user_id (fk)
├── race_id (fk)
├── finish_time
├── ranking
├── ranking_age_group
├── gpx_url
├── photos[]
├── generated_image_url
└── created_at

EquipmentCheck
├── id (uuid)
├── user_id (fk)
├── race_id (fk)
├── equipment_id (fk)
└── checked (bool)
```

### Interface Contracts

**REST API Endpoints (MVP)**:

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/v1/auth/register | 用户注册 |
| POST | /api/v1/auth/login | 用户登录 |
| GET | /api/v1/users/me | 当前用户信息 |
| GET | /api/v1/races | 赛事列表 |
| GET | /api/v1/races/:id | 赛事详情 |
| GET | /api/v1/races/:id/plan | 当日规划数据 |
| GET | /api/v1/races/:id/equipment | 装备清单 |
| POST | /api/v1/races/:id/equipment/check | 更新装备勾选 |
| GET | /api/v1/results | 我的比赛列表 |
| POST | /api/v1/results | 创建完赛记录 |
| GET | /api/v1/results/:id | 完赛记录详情 |
| POST | /api/v1/results/:id/image | 生成轨迹图片 |

**Python Image Service**:
- POST /generate — 接收GPX数据+参数，返回生成的图片URL
- Input: {race_name, date, gpx_data, stats}
- Output: {image_url}

---

**Phase 1 artifacts**: research.md, data-model.md, contracts/, quickstart.md 待生成
