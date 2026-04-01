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

**Version**: 1.1.0 | **Ratified**: 2026-03-30 | **Last Amended**: 2026-04-01
