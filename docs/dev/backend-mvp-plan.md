# 后端 MVP 阶段计划

本文档按 PRD 范围拆分后端 MVP 的完整实现阶段。后端采用 Go + Gin + MySQL 的模块化单体，目录遵循 `cmd`、`ioc`、`internal/domain`、`internal/service`、`internal/repo`、`internal/repo/dao`、`internal/web`、`pkg`。

## 阶段 0：工程初始化

- 初始化 Go module。
- 接入 Gin。
- 建立目录：`cmd`、`ioc`、`internal/domain`、`internal/service`、`internal/repo`、`internal/repo/dao`、`internal/web`、`pkg`。
- 配置加载：`APP_ENV`、`HTTP_ADDR`、`MYSQL_DSN`、`SESSION_SECRET`、`SECRET_ENCRYPTION_KEY`。
- MySQL 连接。
- 健康检查接口：`GET /healthz`。
- 统一响应、统一错误结构。
- recover middleware、request log middleware。

验收：服务可启动，健康检查可用。

## 阶段 1：数据库迁移和基础表

- 建立 migrations 机制。
- 创建基础表：
  - `users`
  - `workspaces`
  - `workspace_members`
  - `prompts`
  - `prompt_versions`
  - `prompt_variables`
  - `prompt_test_cases`
  - `prompt_runs`
  - `tags`
  - `prompt_tags`
  - `provider_configs`
- 加基础索引：
  - `workspace_id`
  - `user_id`
  - `prompt_id`
  - `prompt_version_id`
  - `status`
  - `created_at`
  - `updated_at`

验收：空 MySQL 库可以完整迁移成功。

## 阶段 2：认证与默认工作空间

- 邮箱密码注册：`POST /api/auth/signup`。
- 登录：`POST /api/auth/login`。
- 退出：`POST /api/auth/logout`。
- 当前用户：`GET /api/me`。
- 密码哈希。
- Session 或安全 Cookie。
- 注册后自动创建默认 workspace。
- 创建 `workspace_members` owner 记录。
- Auth middleware 注入当前用户和 workspace。

验收：用户能注册登录，业务接口能拿到 `user_id` 和 `workspace_id`。

## 阶段 3：权限和工作空间隔离

- 定义 `Actor` / `RequestContext`。
- 所有业务查询强制带 `workspace_id`。
- 私有 prompt 只允许创建者访问。
- workspace prompt 允许成员访问。
- 创建者和 workspace owner 可以更新、归档、恢复。
- 跨 workspace 资源返回 `404` 或统一无权限错误。

验收：不同 workspace 数据不能互相访问。

## 阶段 4：Prompt 基础 CRUD

- 创建 prompt：`POST /api/prompts`。
- 列表：`GET /api/prompts`。
- 详情：`GET /api/prompts/{id}`。
- 更新：`PATCH /api/prompts/{id}`。
- 归档：`POST /api/prompts/{id}/archive`。
- 恢复：`POST /api/prompts/{id}/restore`。
- 支持字段：
  - `title`
  - `slug`
  - `description`
  - `body`
  - `message_format`
  - `visibility`
  - `status`
  - `target_provider`
  - `target_model`
  - `default_temperature`
  - `default_max_tokens`
  - `usage_notes`

验收：提示词可以完整创建、查看、编辑、归档、恢复。

## 阶段 5：标签和搜索筛选

- 创建或复用 tag。
- 绑定 prompt 和 tag。
- `GET /api/tags`。
- prompt 列表支持：
  - `keyword`
  - `status`
  - `tag`
  - `target_provider`
  - `target_model`
  - `updated_time`
- 默认隐藏 `archived`。
- MySQL 初期可用 `LIKE`，后续再加 FULLTEXT。

验收：用户可以按关键词、标签、状态、模型筛选提示词。

## 阶段 6：变量解析

- 支持 `{{variable_name}}`。
- 变量名校验。
- 从 prompt body 或 chat messages 推断变量。
- 保存 `prompt_variables`。
- 支持变量字段：
  - `name`
  - `label`
  - `description`
  - `required`
  - `default_value`
  - `example_value`
- 更新 prompt 时同步变量集合。
- 保留用户编辑过的变量元数据。

验收：保存 prompt 后能看到推断出的变量，并能编辑变量元数据。

## 阶段 7：Prompt 渲染和复制支持

- 渲染接口：`POST /api/prompts/{id}/render`。
- 校验必填变量。
- 替换所有变量占位符。
- 支持 `single_text`。
- 支持 `chat_messages`。
- 返回 plain text。
- 返回 JSON messages。
- 返回未填变量错误。

验收：前端可以填写变量并获得可复制的最终 prompt。

## 阶段 8：版本管理

- 创建 prompt 时生成版本 `1`。
- 更新以下内容时生成新版本：
  - body/message blocks
  - variables
  - tags
  - model settings
  - status
  - usage notes
- 版本 snapshot 使用 JSON。
- 版本列表：`GET /api/prompts/{id}/versions`。
- 版本详情：`GET /api/prompts/{id}/versions/{version_id}`。
- 恢复版本：`POST /api/prompts/{id}/versions/{version_id}/restore`。
- 恢复历史版本时创建新的最新版本，不覆盖历史。

验收：提示词能查看历史版本，并能恢复历史版本。

## 阶段 9：测试用例

- 测试用例列表：`GET /api/prompts/{id}/test-cases`。
- 创建：`POST /api/prompts/{id}/test-cases`。
- 更新：`PATCH /api/prompts/{id}/test-cases/{test_case_id}`。
- 删除：`DELETE /api/prompts/{id}/test-cases/{test_case_id}`。
- 支持字段：
  - `name`
  - `variable_values`
  - `expected_behavior`
  - `expected_output`
- 校验 `variable_values` 和 prompt variables 的关系。

验收：每个 prompt 可以维护多个测试用例。

## 阶段 10：模型 Provider 配置

- 查看 provider config：`GET /api/settings/provider`。
- 保存 provider config：`PUT /api/settings/provider`。
- 支持 OpenAI-compatible：
  - `base_url`
  - `api_key`
  - `default_model`
- API Key 加密保存。
- API 响应只返回 `has_api_key`，不回显原始 key。
- 日志中不能出现 API Key。

验收：用户可以保存模型提供方配置，密钥不会泄露。

## 阶段 11：Prompt 运行

- 运行接口：`POST /api/prompts/{id}/run`。
- 支持直接变量输入运行。
- 支持基于 test case 运行。
- 调用 OpenAI-compatible API。
- 记录：
  - `prompt_id`
  - `prompt_version_id`
  - `test_case_id`
  - `provider`
  - `model`
  - `input_variables`
  - `rendered_prompt_snapshot`
  - `output_text`
  - `latency`
  - `token_usage`
  - `error_message`
  - `created_by`
  - `created_at`
- Provider 错误转成安全错误信息。

验收：可以运行 prompt，并保存成功或失败结果。

## 阶段 12：运行记录

- prompt 维度运行记录：`GET /api/prompts/{id}/runs`。
- 可选全局运行记录：`GET /api/runs`。
- 支持分页。
- 支持按成功/失败、模型、时间筛选。
- 返回输入摘要、输出摘要、版本号、模型信息。

验收：用户可以查看某个 prompt 的历史运行结果。

## 阶段 13：Mock 和单元测试

- service 层接口补 `go:generate`。
- 使用 `mockgen` 生成 mocks。
- 覆盖：
  - 注册和默认 workspace 创建。
  - prompt CRUD 权限。
  - workspace 隔离。
  - 变量解析和渲染。
  - 版本创建和恢复。
  - provider config 不泄露 API Key。
  - run 失败时不泄露敏感信息。
- 命令：
  - `go generate ./...`
  - `go test ./...`
  - `go vet ./...`

验收：核心业务逻辑有稳定单元测试。

## 阶段 14：前后端联调支持

- CORS 或同源代理策略。
- Cookie/session 配置适配本地开发。
- API 错误码稳定。
- API 返回字段和前端类型对齐。
- 提供基础 seed 数据，可选。
- 补充 README 本地启动说明。

验收：前端可以完整接入后端 API，跑通 MVP 主流程。

## 阶段 15：MVP 验收闭环

完整流程必须跑通：

1. 用户注册。
2. 自动创建默认 workspace。
3. 创建 prompt。
4. 自动解析变量。
5. 编辑变量元数据。
6. 渲染 prompt。
7. 复制 plain text 或 chat JSON。
8. 更新 prompt。
9. 自动生成版本。
10. 查看版本历史。
11. 恢复历史版本。
12. 搜索和筛选 prompt。
13. 创建测试用例。
14. 配置 provider。
15. 运行测试用例。
16. 查看运行记录。
17. 归档和恢复 prompt。

验收：以上主流程全部通过，后端 MVP 完成。


// sk-c0010ccd9b06467a8e676f733441d891