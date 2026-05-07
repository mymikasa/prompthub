# 前端开发文档

## 1. 目标

前端负责提供清晰、高效的提示词管理体验。MVP 不做营销站，不做复杂视觉展示，第一屏应直接进入可操作的产品界面。

核心体验：

- 快速浏览和搜索提示词。
- 顺畅创建、编辑和归档提示词。
- 清楚填写变量并复制渲染结果。
- 查看版本历史、测试用例和运行记录。
- 安全配置模型提供方。

## 2. 技术栈建议

- 构建工具：Vite
- 框架：React
- 语言：TypeScript
- 路由：React Router
- 数据请求：TanStack Query
- 表单：React Hook Form
- 校验：Zod
- 样式：Tailwind CSS
- 图标：lucide-react

如果后续需要 SSR、公开分享页 SEO 或边缘部署，再评估 Next.js。

## 3. 推荐目录结构

```text
frontend/
  src/
    app/
      App.tsx
      router.tsx
      providers.tsx
    components/
      ui/
      layout/
    features/
      auth/
      prompts/
      settings/
      runs/
    lib/
      api.ts
      errors.ts
      format.ts
    styles/
      globals.css
    main.tsx
  index.html
  package.json
  tsconfig.json
  vite.config.ts
```

目录职责：

- `app`：应用入口、路由和全局 providers。
- `components/ui`：通用按钮、输入框、弹窗、菜单、标签、空状态。
- `components/layout`：侧边栏、顶部栏、页面容器。
- `features/auth`：登录、注册、当前用户。
- `features/prompts`：提示词列表、编辑器、变量、版本、测试用例。
- `features/runs`：运行记录页面和详情。
- `features/settings`：个人资料、工作空间、模型提供方配置。
- `lib/api.ts`：统一 API client。

## 4. 页面结构

主导航：

- 提示词
- 运行记录
- 设置

路由建议：

```text
/login
/signup
/prompts
/prompts/new
/prompts/:promptId
/prompts/:promptId/edit
/prompts/:promptId/versions/:versionId
/runs
/settings
/settings/provider
```

提示词详情页建议使用 tabs：

- 概览
- 编辑器
- 变量
- 测试用例
- 运行记录
- 版本

## 5. UI 设计原则

- 产品界面优先，避免营销式 landing page。
- 信息密度适中，列表和表单要适合反复使用。
- 操作按钮使用清晰图标和短文本，图标优先使用 `lucide-react`。
- 页面区块不要层层嵌套卡片。
- 空状态要给出下一步操作，例如“创建提示词”。
- 复制、保存、归档、恢复等操作必须有即时反馈。
- 危险操作必须确认。
- 长提示词和长输出必须可滚动，不能撑坏页面布局。
- 移动端至少保证可查看、搜索、复制；复杂编辑可优先优化桌面端。

## 6. 核心数据类型

前端类型应和后端 API 对齐：

```ts
type PromptMessageFormat = "single_text" | "chat_messages";
type PromptVisibility = "private" | "workspace";
type PromptStatus = "draft" | "active" | "deprecated" | "archived";
type MessageRole = "system" | "user" | "assistant";

type Prompt = {
  id: string;
  title: string;
  slug: string;
  description: string;
  messageFormat: PromptMessageFormat;
  visibility: PromptVisibility;
  status: PromptStatus;
  targetProvider: string;
  targetModel: string;
  defaultTemperature: number | null;
  defaultMaxTokens: number | null;
  usageNotes: string;
  tags: string[];
  createdAt: string;
  updatedAt: string;
};
```

命名约定：

- TypeScript 内部使用 `camelCase`。
- API JSON 如果后端返回 `snake_case`，在 API client 层做转换，或全项目统一接受 `snake_case`。二者选一个，不要混用。

## 7. API Client 约定

统一封装 `fetch`：

- 自动添加 `Content-Type: application/json`。
- 自动携带 cookie/session。
- 统一解析错误响应。
- 对 `401` 做登录态失效处理。

错误对象格式：

```ts
type ApiError = {
  code: string;
  message: string;
  details?: unknown;
};
```

前端不直接拼接复杂 API 逻辑，业务请求集中在各 feature 的 `api.ts` 中，例如：

```text
features/prompts/api.ts
features/auth/api.ts
features/settings/api.ts
```

## 8. 认证流程

页面行为：

- 未登录访问业务页面时跳转到 `/login`。
- 登录成功后跳转到 `/prompts`。
- 注册成功后自动进入 `/prompts`。
- `GET /api/me` 用于恢复登录态。

UI 状态：

- 登录中显示加载状态。
- 登录失败展示后端返回的错误信息。
- 退出登录后清空 query cache。

## 9. 提示词列表

功能：

- 关键词搜索。
- 状态筛选。
- 标签筛选。
- 模型筛选。
- 默认隐藏归档提示词。
- 创建提示词入口。

列表字段：

- 标题
- 描述
- 状态
- 标签
- 目标模型
- 更新时间
- 创建者，后续团队版需要

交互：

- 搜索输入应 debounce。
- 点击行进入详情页。
- 空列表展示创建入口。
- 加载和错误状态要明确。

## 10. 提示词编辑器

编辑字段：

- 标题
- 描述
- 正文或 chat messages
- 标签
- 可见性
- 状态
- 模型提供方
- 模型名称
- temperature
- max tokens
- 使用说明

编辑要求：

- 保留正文格式。
- 支持 `single_text` 和 `chat_messages` 两种模式。
- 变量占位符如 `{{topic}}` 应尽量高亮。
- 保存前展示推断出的变量列表。
- 保存成功后刷新提示词详情和版本信息。
- 离开未保存页面前应提醒。

## 11. 变量和渲染

变量页功能：

- 查看系统推断变量。
- 编辑变量标签、描述、必填、默认值、示例值。
- 填写变量值并预览渲染结果。
- 复制纯文本。
- 对 chat prompt 复制 JSON messages。

校验：

- 必填变量为空时不能渲染。
- 变量名不合法时展示错误。
- 渲染结果中如果仍有未填占位符，应提示用户。

## 12. 版本历史

版本页功能：

- 展示版本号、作者、时间、变更说明。
- 查看某个版本快照。
- 恢复历史版本。

MVP 可暂不实现 diff，但页面结构应预留 diff 区域。

恢复行为：

- 点击恢复前弹确认框。
- 恢复成功后跳转到最新版本或提示词详情。
- 恢复会创建新版本，前端文案要避免暗示“覆盖历史”。

## 13. 测试用例和运行记录

测试用例功能：

- 创建、编辑、删除测试用例。
- 填写变量值。
- 记录预期行为说明。
- 在 provider config 可用时运行测试。

运行记录展示：

- 运行时间
- Prompt 版本
- 模型
- 输入变量摘要
- 输出摘要
- 成功或失败状态
- 延迟和 token 使用量，如有

运行失败时：

- 展示可理解错误。
- 不展示 API Key 或敏感请求头。

## 14. 设置页

设置分区：

- 个人资料
- 工作空间
- 模型提供方

模型提供方配置：

- Base URL
- API Key
- Default Model

安全交互：

- 已保存 API Key 时只展示“已配置”，不展示原值。
- 更新 API Key 需要重新输入。
- 保存成功后清空输入框中的 API Key。

## 15. 状态管理

优先使用：

- TanStack Query 管理服务端状态。
- React local state 管理表单局部状态。
- 避免在 MVP 引入全局状态库，除非登录态和 UI 状态已经无法清晰维护。

Query key 建议：

```ts
["me"]
["prompts", filters]
["prompt", promptId]
["promptVersions", promptId]
["promptRuns", promptId]
["providerConfig"]
```

## 16. 测试策略

MVP 前端至少覆盖：

- 变量渲染表单校验。
- Prompt editor 的保存请求。
- 登录失败和登录成功流程。
- API error 展示。

建议工具：

- Vitest
- React Testing Library
- Playwright，后续用于关键端到端流程

建议命令：

```text
pnpm test
pnpm build
```

## 17. 本地开发流程

建议命令：

```text
cd frontend
pnpm install
pnpm dev
```

开发顺序：

1. 初始化 Vite React TypeScript。
2. 建立路由、布局和 API client。
3. 实现登录、注册和 `/api/me` 登录态恢复。
4. 实现提示词列表和详情。
5. 实现提示词创建、编辑和归档。
6. 实现变量渲染和复制。
7. 实现版本、测试用例、运行记录和设置页。

## 18. 前端验收清单

- 应用可本地启动。
- 未登录访问业务页会跳转登录。
- 登录后能进入提示词列表。
- 提示词列表支持搜索、筛选和空状态。
- 可以创建、编辑、归档、恢复提示词。
- 可以填写变量并复制渲染结果。
- chat prompt 可以复制 JSON messages。
- 可以查看版本历史并触发恢复。
- Provider API Key 保存后不在页面回显。
- 构建命令通过。
