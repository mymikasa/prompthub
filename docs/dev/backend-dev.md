# 后端开发文档

## 1. 目标

后端负责 prompthub 的核心业务能力：认证、工作空间隔离、提示词管理、变量解析、版本历史、测试用例、运行记录和模型提供方配置。

MVP 阶段优先保证领域边界清晰、数据模型稳定、接口清楚、权限边界可靠。后端采用 DDD 风格的模块化单体：代码按领域拆分，但运行和部署仍然是一个 Gin 服务。

## 2. 技术栈建议

- 语言：Go
- Web 框架：Gin
- 数据库：MySQL
- API：REST + JSON
- 迁移：Goose、Atlas 或 golang-migrate 三选一
- 配置：环境变量
- 日志：结构化日志
- 测试：Go 标准测试框架

推荐先使用轻量组合：

- HTTP router：`gin`
- 数据访问：`database/sql` + MySQL driver，或 GORM；MVP 推荐先手写 repository，保持 SQL 明确可控
- 密码哈希：`bcrypt` 或 `argon2id`
- Mock 生成：`mockgen`
- 依赖注入：手动 DI 起步，复杂后可引入 Wire

架构约束：

- 单体服务：只有一个 API 进程、一个数据库、一个部署单元。
- DDD 分层：按领域拆 package，避免把所有业务堆在 handler 或 repository。
- 领域内聚：提示词、工作空间、认证、模型提供方分别作为边界明确的领域模块。
- 不做微服务：领域模块之间通过 Go 接口和应用服务协作，不通过 RPC 或消息队列强行拆分。
- 面向接口编程：service 层依赖接口，不直接依赖 MySQL、Gin、外部 SDK 的具体实现。

## 3. 推荐目录结构

```text
backend/
  cmd/
    main.go
  internal/
    domain/
    service/
    repo/
      dao/
    web/
  ioc/
  pkg/
  migrations/
  go.mod
  go.sum
```

轻量 DDD 分层职责：

- `cmd`：服务启动入口，只做配置读取、调用 `ioc` 初始化、启动 Gin。
- `ioc`：依赖注入和组件组装，创建 DB、DAO、repo、service、handler、Gin engine。
- `internal/web`：接口层，包含 Gin handler、middleware、DTO、响应封装和路由注册。
- `internal/service`：应用服务层，编排用例、事务、权限校验和跨领域协作。
- `internal/domain`：领域层，纯业务实体，无任何 ORM tag、无外部依赖。
- `internal/repo`：仓储层，面向 service 暴露领域模型操作，内部调用 DAO 并负责 model ↔ domain 转换。
- `internal/repo/dao`：数据访问层，只操作 `dao/model` 持久化模型，不依赖 domain。
- `internal/repo/dao/model`：持久化模型，带 GORM tag，与数据库表一一对应。
- `pkg`：可复用小工具，只放与业务无关的稳定能力，例如错误码、校验、加密辅助。

单个领域模块建议结构：

```text
internal/
  domain/
    prompt.go
    prompt_version.go
    prompt_variable.go
    errors.go
  service/
    prompt.go
    prompt_test.go
  repo/
    prompt.go
    dao/
      prompt.go
      prompt_version.go
  web/
    prompt.go
    middleware.go
ioc/
  app.go
  router.go
pkg/
  errors/
  validator/
```

目录职责：

- Handler 只做 HTTP 请求解析、DTO 校验、调用 service、返回响应。
- Service 负责具体用例，例如创建提示词、恢复版本、渲染提示词。
- Domain entity/value object 承载业务规则，纯 Go struct，不带 GORM tag。
- Repo（仓储层）调用 DAO 获取 model，转换为 domain entity 返回给 service。
- DAO 只处理 MySQL 读写，操作 `dao/model` 类型，不承载业务规则，不 import domain。
- `dao/model` 定义带 GORM tag 的持久化模型，用于数据库映射和迁移。
- 数据库事务由 service 层控制，具体事务实现由 `internal/repo/dao` 或 `ioc` 注入。

## 4. DDD 分层约定

### 4.1 依赖方向

依赖只能从外向内：

```text
web → service → repo → dao → dao/model
              ↘ domain ↗
ioc → web/service/repo/dao
```

严格的 import 约束：

| 层 | 可以 import | 禁止 import |
|---|---|---|
| `domain` | 标准库 | Gin、GORM、dao、repo、service、web |
| `dao/model` | 标准库、GORM | domain、repo、service、web |
| `dao` | 标准库、GORM、`dao/model` | domain、repo、service、web |
| `repo` | 标准库、domain、`dao`、`dao/model` | Gin、service、web |
| `service` | 标准库、domain、repo | Gin、dao、dao/model、web |
| `web` | Gin、service、domain（仅类型） | dao、dao/model、repo |
| `ioc` | 所有层 | — |

关键规则：

- `dao` 绝对不能 import `domain`。DAO 只操作 `dao/model` 持久化模型。
- `repo` 是 model ↔ domain 的转换边界。所有 `toModel*` / `toDomain*` 转换函数放在 `repo` 层。
- `domain` 是纯业务 struct，不带任何 GORM tag。
- `dao/model` 是持久化 struct，带 GORM tag，与表结构一一对应。
- `ioc` 负责把具体实现注入到 service 和 handler。

### 4.2 领域边界

MVP 领域对象先放在 `internal/domain` 下，保持少量文件即可；当某个领域明显变大时，再拆成子目录。

初始领域：

- `auth`：用户、密码、会话、登录注册。
- `workspaces`：工作空间、成员、角色、当前工作空间上下文。
- `prompts`：提示词、变量、标签、版本、测试用例、运行记录。
- `providers`：模型提供方配置、密钥保护、模型调用。

跨领域规则：

- `prompts` 不直接查询 auth 表，只接收 `UserID`、`WorkspaceID` 等身份上下文。
- `providers` 不直接修改 prompt，只提供模型执行能力。
- `workspaces` 负责判断成员关系，业务用例通过 service 层调用它。

### 4.3 Gin 使用边界

- Gin 只出现在 `internal/web` 和 `ioc`。
- 不要把 `*gin.Context` 传入 service 或 domain。
- Handler 应把请求转换成明确的 command/query DTO。
- Middleware 可把认证结果写入 request context，但 service 层接收的是显式 `Actor` 或 `RequestContext` 结构体。

Context 传递规则：

- Handler 从 `c.Request.Context()` 取出 `context.Context`，作为第一个参数传入 service。
- Service 透传 ctx 给 repo。
- Repo 透传 ctx 给 dao。
- DAO 用 `db.WithContext(ctx)` 执行所有数据库操作。
- 全链路 context 传递确保超时取消、链路追踪等能力不丢失。

示例：

```go
type Actor struct {
    UserID      uuid.UUID
    WorkspaceID uuid.UUID
    Role        workspace.Role
}
```

### 4.4 面向接口编程约定

原则：

- service 层依赖接口，repo/dao 层提供实现。
- 接口以用例需要为边界，不为了“抽象”而抽象。
- 接口尽量小，避免一个 repository 接口包含所有读写方法。
- 外部能力都通过接口进入 service 层，例如密码哈希、session、模型调用、密钥加密。

推荐接口放置规则：

- Service 需要的 repository 接口放在 `internal/service` 对应文件里，便于按用例定义最小接口。
- Domain 不关心持久化时，不强行在 domain 放 repository 接口。
- 具体 repository 实现放在 `internal/repo`。
- MySQL 访问细节放在 `internal/repo/dao`。

示例：

```go
type PromptRepository interface {
    FindByID(ctx context.Context, workspaceID, promptID uuid.UUID) (*Prompt, error)
    Save(ctx context.Context, prompt *Prompt) error
}

type VersionRepository interface {
    NextVersionNumber(ctx context.Context, promptID uuid.UUID) (int, error)
    Save(ctx context.Context, version *PromptVersion) error
}
```

### 4.5 Repository 约定

Repository 是具体实现层（非接口），放在 `internal/repo`：

```text
internal/repo/user.go         # 调用 DAO，转换 model → domain
internal/repo/workspace.go
internal/repo/dao/user.go     # 原始 DB 操作，只返回 model 类型
internal/repo/dao/model/      # GORM 持久化模型
```

三层协作示例：

```go
// dao/model/user.go — 持久化模型（GORM tag）
type User struct {
    ID           int64  `gorm:"primaryKey;autoIncrement"`
    Email        string `gorm:"type:varchar(255);uniqueIndex;not null"`
    PasswordHash string `gorm:"type:varchar(255);not null"`
}

// dao/user.go — 原始 DB 操作，返回 model 类型
func (d *UserDAO) FindByID(ctx context.Context, id int64) (*model.User, error) { ... }

// repo/user.go — 仓储层，转换 model → domain
func (r *UserRepo) FindByID(ctx context.Context, id int64) (*domain.User, error) {
    m, err := r.dao.FindByID(ctx, id)
    return toDomainUser(m), nil  // 转换边界
}
```

规则：

- Repository 方法必须显式接收 `context.Context` 作为第一个参数。
- Repository 方法必须显式接收 `workspaceID`，避免漏掉租户隔离。
- Repository 返回 domain entity，不直接返回数据库 model。
- 数据库 model 到 domain entity 的转换放在 repo 层。
- DAO model 不向 service/web 暴露。
- Repository 实现不处理 HTTP 错误，只返回领域错误或基础设施错误。

### 4.6 事务约定

涉及多个写入的用例必须在 service 层开启事务：

- 创建提示词和初始版本。
- 更新提示词和创建新版本。
- 恢复历史版本和创建新版本。
- 运行测试用例和保存运行记录。

建议定义事务接口：

```go
type TxManager interface {
    WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}
```

具体实现可以放在 `internal/repo/dao` 或 `internal/repo`，由 `ioc` 注入到 service。

### 4.7 mockgen 约定

所有 service 依赖的接口都应能用 `mockgen` 生成 mock，用于单元测试。

推荐生成路径：

```text
internal/service/mocks/
  mock_prompt_repository.go
  mock_workspace_checker.go
  mock_model_runner.go
```

推荐命令：

```text
mockgen -source=internal/service/prompt.go -destination=internal/service/mocks/mock_prompt_repository.go -package=servicemocks
mockgen -source=internal/service/provider.go -destination=internal/service/mocks/mock_provider.go -package=servicemocks
```

也可以在接口文件旁放 `go:generate` 注释，统一通过 `go generate ./...` 生成：

```go
//go:generate mockgen -source=prompt.go -destination=./mocks/mock_prompt_repository.go -package=servicemocks
```

规则：

- mock 文件不手写，只通过 `mockgen` 生成。
- 生成的 mock 可以提交到仓库，保证 CI 不依赖开发者本地工具状态。
- 修改接口后必须重新运行 `go generate ./...` 或对应 `mockgen` 命令。
- service 层测试优先使用 mock 验证用例编排、权限判断、事务行为。
- domain 层测试不需要 mock 数据库，应直接测试实体和值对象规则。

### 4.8 依赖注入约定

MVP 先使用显式构造函数和 `ioc` 手动注入：

```go
type PromptService struct {
    prompts PromptRepository
    versions VersionRepository
    tx       TxManager
}

func NewPromptService(
    prompts PromptRepository,
    versions VersionRepository,
    tx TxManager,
) *PromptService {
    return &PromptService{
        prompts: prompts,
        versions: versions,
        tx: tx,
    }
}
```

组装位置：

- `ioc/` 创建配置、数据库连接、DAO、repository、service、handler。
- 组装顺序：`DB → DAO → Repo → Service → Handler → Router`

```go
// cmd/main.go 组装示例
userDAO := dao.NewUserDAO(db)
workspaceDAO := dao.NewWorkspaceDAO(db)

userRepo := repo.NewUserRepo(userDAO)
workspaceRepo := repo.NewWorkspaceRepo(workspaceDAO)

authService := service.NewAuthService(userRepo, workspaceRepo)
authHandler := handler.NewAuthHandler(authService, cfg.SessionSecret, cfg.IsDev())
```

规则：

- 不在业务代码里直接 `new` 基础设施实现。
- 不使用全局单例保存 repository、service 或数据库连接。
- Handler 通过构造函数接收 service。
- Service 通过构造函数接收 repo 具体类型。
- 如果手动 DI 变得冗长，再引入 Google Wire；Wire 配置仍放在 `ioc` 附近。

## 5. 配置

MVP 后端至少需要以下环境变量：

```text
APP_ENV=development
HTTP_ADDR=:8080
MYSQL_DSN=prompthub:prompthub@tcp(127.0.0.1:3306)/prompthub?parseTime=true&loc=Local&charset=utf8mb4
SESSION_SECRET=change-me
SECRET_ENCRYPTION_KEY=32-byte-base64-or-hex-key
```

规则：

- `.env` 不提交到仓库。
- 缺少关键配置时服务启动必须失败。
- 生产环境不能使用默认 `SESSION_SECRET`。
- API Key 加密密钥必须和数据库分离管理。

## 6. 数据模型

MVP 初始表：

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

通用字段建议：

```text
id char(36) primary key
created_at datetime(3) not null
updated_at datetime(3) not null
deleted_at datetime(3) null
```

关键约束：

- 业务表必须包含 `workspace_id`，并在查询时强制过滤。
- `prompt_versions` 是不可变历史，不允许更新历史快照。
- `prompt_runs` 绑定 `prompt_version_id`，不能只绑定当前 prompt。
- `provider_configs` 不能存明文 API Key。

## 7. 核心枚举

```text
prompt_message_format: single_text | chat_messages
prompt_visibility: private | workspace
prompt_status: draft | active | deprecated | archived
message_role: system | user | assistant
workspace_role: owner | member
```

枚举值进入数据库和 API 时都使用小写 snake_case。

## 8. API 约定

基础路径：

```text
/api
```

统一响应：

```json
{
  "data": {}
}
```

错误响应：

```json
{
  "error": {
    "code": "validation_error",
    "message": "请求参数无效",
    "details": {}
  }
}
```

分页响应：

```json
{
  "data": [],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100
  }
}
```

HTTP 状态码：

- `200`：查询或更新成功。
- `201`：创建成功。
- `204`：删除或无内容成功。
- `400`：参数错误。
- `401`：未登录。
- `403`：无权限。
- `404`：资源不存在或不在当前工作空间。
- `409`：状态冲突。
- `500`：服务端错误。

## 9. REST 接口清单

```text
POST   /api/auth/signup
POST   /api/auth/login
POST   /api/auth/logout
GET    /api/me

GET    /api/prompts
POST   /api/prompts
GET    /api/prompts/{id}
PATCH  /api/prompts/{id}
POST   /api/prompts/{id}/archive
POST   /api/prompts/{id}/restore

GET    /api/prompts/{id}/versions
GET    /api/prompts/{id}/versions/{version_id}
POST   /api/prompts/{id}/versions/{version_id}/restore

GET    /api/prompts/{id}/test-cases
POST   /api/prompts/{id}/test-cases
PATCH  /api/prompts/{id}/test-cases/{test_case_id}
DELETE /api/prompts/{id}/test-cases/{test_case_id}

POST   /api/prompts/{id}/render
POST   /api/prompts/{id}/run
GET    /api/prompts/{id}/runs

GET    /api/tags

GET    /api/settings/provider
PUT    /api/settings/provider
```

## 10. 认证和权限

MVP 使用服务端会话或安全 cookie 均可。需要满足：

- 登录后可以通过 `GET /api/me` 获取当前用户和默认工作空间。
- 所有业务接口必须要求登录。
- 所有业务查询必须带当前 `workspace_id` 范围。
- 私有提示词只允许创建者访问。
- 工作空间成员可以查看和运行 `workspace` 可见提示词。
- 创建者和 workspace owner 可以更新、归档、恢复提示词。

实现注意：

- 不要通过前端传入的 `workspace_id` 直接决定数据范围，应以服务端会话中的当前工作空间为准。
- 对外统一返回 `404` 可以避免泄露跨工作空间资源是否存在。

## 11. Prompt 变量解析

变量语法：

```text
{{variable_name}}
```

变量名规则：

- 只能包含字母、数字和下划线。
- 必须以字母开头。
- 区分大小写不推荐，建议保存时统一转小写。

后端职责：

- 保存提示词时解析占位符。
- 返回推断出的变量列表。
- 渲染前校验必填变量。
- 渲染后不能残留未处理的必填占位符。

渲染接口应保存输入快照，避免后续 prompt 修改影响历史运行记录。

## 12. 版本策略

以下字段变化时创建新版本：

- 提示词正文或消息块
- 变量定义
- 标签
- 模型设置
- 状态
- 使用说明

版本规则：

- 创建 prompt 时生成版本 `1`。
- 每次更新生成递增版本号。
- 恢复历史版本时，不覆盖历史版本，而是创建新的最新版本。
- `prompt_versions.snapshot` 使用 JSON 保存完整快照。

## 13. 模型提供方和密钥

MVP 支持 OpenAI-compatible 配置：

```text
base_url
api_key
default_model
```

安全规则：

- `api_key` 只接收写入，不允许读出。
- 返回 provider config 时只返回 `has_api_key: true/false`。
- 日志、错误、运行记录中不能包含 API Key。
- 运行失败时返回可诊断的错误摘要，不直接透传包含敏感信息的 provider 原始错误。

## 14. 测试策略

最少测试覆盖：

- 变量解析和渲染。
- Prompt CRUD 的权限隔离。
- 版本创建和恢复。
- 私有提示词访问控制。
- Provider config 不返回 API Key。
- Service 的事务行为。
- Domain 层核心规则不依赖数据库即可测试。

建议测试命令：

```text
go test ./...
go vet ./...
```

## 15. 本地开发流程

建议命令：

```text
cd backend
go mod tidy
go generate ./...
go run ./cmd
```

数据库建议通过 Docker 启动 MySQL，后续可以在仓库根目录补 `docker-compose.yml`。

开发顺序：

1. 初始化 Go module、Gin 服务和健康检查。
2. 接入 MySQL 和 migrations。
3. 建立 DDD 基础目录、错误响应、事务接口和 repository 接口。
4. 实现认证、用户和默认工作空间。
5. 实现提示词 CRUD 和 workspace scope。
6. 实现变量解析、渲染和版本历史。
7. 实现测试用例、provider config 和运行记录。

## 16. 后端验收清单

- 服务可本地启动。
- 数据库迁移可从空库成功执行。
- Gin 只出现在 `internal/web` 和 `ioc`。
- Domain 层不依赖数据库、Gin 或外部 SDK。
- 注册后自动创建默认工作空间。
- 所有业务接口未登录返回 `401`。
- 跨工作空间访问被拒绝。
- 创建、更新、归档、恢复提示词可用。
- Prompt 更新生成版本记录。
- 历史版本可恢复为新版本。
- 渲染接口能校验必填变量。
- Provider API Key 不会出现在响应和日志中。
