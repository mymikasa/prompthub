# prompthub 开发文档

本目录用于沉淀 prompthub 的工程实现约定。MVP 产品范围见 [../mvp/prd-mvp.md](../mvp/prd-mvp.md)。

当前仓库仍处于初始阶段，`backend/` 和 `frontend/` 尚未脚手架化。以下文档定义第一版建议采用的工程结构、接口约定、开发流程和验收标准：

- [后端开发文档](./backend-dev.md)
- [后端 MVP 阶段计划](./backend-mvp-plan.md)
- [前端开发文档](./frontend-dev.md)

## 技术默认决策

- 后端：Go + Gin + MySQL + REST API，采用 DDD 风格的模块化单体，面向接口编程，使用 mockgen 生成 mock，并通过依赖注入组装组件。
- 前端：Vite + React + TypeScript。
- 通信：前端通过 `/api/*` 调用后端 REST 接口。
- 认证：MVP 使用邮箱密码登录，服务端签发会话。
- 数据隔离：所有业务数据必须按 `workspace_id` 隔离。
- Prompt 模板变量：统一使用 `{{variable_name}}` 语法。

## MVP 优先级

1. 认证、默认工作空间、基础布局。
2. 提示词 CRUD、标签、列表搜索。
3. 变量解析、渲染预览、复制。
4. 版本历史、恢复版本。
5. 测试用例、模型提供方配置、运行记录。

## 跨端约定

- API 响应统一使用 JSON。
- 时间字段统一返回 RFC 3339 字符串。
- ID 建议使用 UUID。
- 枚举值使用小写 snake_case，例如 `single_text`、`chat_messages`。
- 前后端共享的概念以 PRD 为准；如需变更，先更新 PRD 或开发文档。
