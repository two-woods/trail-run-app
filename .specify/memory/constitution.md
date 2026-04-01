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

### VII. Web Testing (NON-NEGOTIABLE)
每次功能改动必须使用Playwright进行页面测试验证：
- 使用Playwright MCP工具进行自动化Web测试
- 测试覆盖：页面加载、关键元素存在、数据展示正确
- 登录流程测试：注册→登录→Token存储→登出
- 列表页测试：数据加载、分页、筛选
- 详情页测试：数据完整展示、导航跳转
- 提交前必须完成Web测试，测试报告随commit记录

## 技术栈约束

### 技术选型
- 后端：Go 1.21 + Gin + GORM + MySQL 8.0 + Redis 7
- 前端：Uni-app (Vue3) + Vite + H5
- 图片服务：Python FastAPI + Pillow + cairosvg
- 地图：高德地图API
- 对象存储：阿里云OSS
- 测试：Playwright (Web Testing)

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
- Web Testing：每次功能改动必须使用Playwright测试

## Governance

宪法优先于其他实践规范。修订要求：
- 新增原则需明确说明原因
- 重大变更需要迁移计划
- 所有PR必须验证合规性

**Version**: 1.3.0 | **Ratified**: 2026-03-30 | **Last Amended**: 2026-04-01

## Sync Impact Report

### Version Change
- Old: 1.2.0
- New: 1.3.0 (MINOR bump - 新增Web Testing原则)

### Modified Principles
- None (new principle added)

### Added Sections
- VII. Web Testing (NON-NEGOTIABLE) - Playwright自动化测试要求

### Templates Requiring Updates
- ⚠ pending: `.specify/templates/plan-template.md` - 检查Constitution Check是否需要更新
- ⚠ pending: `.specify/templates/tasks-template.md` - 检查任务分类是否需要更新

### Follow-up TODOs
- 配置Playwright测试环境
- 编写核心页面的Playwright测试用例

### Rationale
用户明确要求"每次改动功能,都要使用web testing,凭借playwright来测试下页面"。这确保：
1. 功能改动后立即验证页面行为正确
2. 捕获回归问题在早期阶段
3. 提供可视化的测试报告记录
