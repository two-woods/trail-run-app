# 越野跑App Constitution

## Core Principles

### I. 全栈容器化 (NON-NEGOTIABLE)
所有服务必须容器化部署。架构要求：
- 后端服务（Go/Gin）、前端服务（Uni-app H5）、Python图片服务、MySQL、Redis等所有依赖服务必须通过Docker Compose管理
- 开发环境必须使用容器化，确保本地与生产环境一致
- 禁止在容器外直接运行服务（除了本地调试短暂运行）
- 每个服务必须有对应的Dockerfile和docker-compose.yml配置

### II. API优先设计
后端API驱动前端实现：
- 所有功能通过RESTful API暴露
- 前后端分离，各自独立部署
- 统一响应格式：{success, data, error}

### III. 数据完整性
数据模型规范：
- 使用GORM进行数据库操作
- 所有Model必须定义完整字段和验证
- 敏感数据不硬编码，使用环境变量

### IV. 错误处理与日志
错误处理规范：
- 统一错误响应格式
- 日志分级：Info、Warn、Error
- 关键操作必须记录日志

### V. 前端状态安全
前端安全要求：
- Token存储在Storage中
- 401响应自动跳转登录页
- 用户输入必须验证

### VI. 模拟数据优先 (NON-NEGOTIABLE)
所有模块必须提供模拟数据进行开发和测试：
- 每个列表类接口必须支持分页返回模拟数据
- 每个详情类接口必须返回完整的模拟数据对象
- 模拟数据必须包含所有字段的典型值（不能只是占位符）
- 模拟数据文件名规范：使用 `mock_` 前缀，如 `mock_races.json`
- 模拟数据存放位置：后端 `/internal/testdata/` 或前端 `/src/testdata/`

## 技术栈约束

### 技术选型
- 后端：Go 1.21 + Gin + GORM + MySQL 8.0 + Redis 7
- 前端：Uni-app (Vue3) + Vite + H5
- 图片服务：Python FastAPI + Pillow + cairosvg
- 地图：高德地图API
- 对象存储：阿里云OSS

### 部署要求
- Nginx作为反向代理和静态资源服务
- 所有服务通过docker-compose编排
- 环境变量通过.env文件管理（不提交到版本控制）

## 开发工作流

### 代码审查
- 所有代码变更必须经过审查
- 提交前运行构建验证
- 禁止直接推送主分支

### 测试要求
- 核心功能需要测试覆盖
- 集成测试验证API端点
- E2E测试关键用户流程

## Governance

宪法优先于其他实践规范。修订要求：
- 新增原则需明确说明原因
- 重大变更需要迁移计划
- 所有PR必须验证合规性

**Version**: 1.2.0 | **Ratified**: 2026-03-30 | **Last Amended**: 2026-04-01

## Sync Impact Report

### Version Change
- Old: 1.1.0
- New: 1.2.0 (MINOR bump - 新增模拟数据原则)

### Modified Principles
- None (new principle added)

### Added Sections
- VI. 模拟数据优先 (NON-NEGOTIABLE) - 所有模块必须提供模拟数据

### Templates Requiring Updates
- ⚠ pending: `.specify/templates/plan-template.md` - 检查Constitution Check是否需要更新
- ⚠ pending: `.specify/templates/tasks-template.md` - 检查任务分类是否需要更新

### Follow-up TODOs
- 为现有模块补充mock数据文件
- 在tasks.md中添加模拟数据相关任务（如有）

### Rationale
用户明确要求"所有模块都需要有模拟数据,比如列表展示"。这确保：
1. 前端开发可以在后端API未完成时独立进行
2. 测试可以在没有真实数据库的情况下运行
3. Demo和演示可以脱离实际环境运行
