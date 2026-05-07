# PromptHub

提示词管理/共享平台。Go + Gin + MySQL 模块化单体架构，React + TypeScript 前端。

## 功能

- 提示词 CRUD（支持 Single Text / Chat Messages 两种格式）
- 变量占位符 `{{variable}}` 自动提取与渲染
- 版本管理（快照 + 恢复）
- 测试用例管理
- 多模型 Provider 配置（OpenAI Compatible / DeepSeek / 智谱 / 通义千问 / Moonshot / MiniMax）
- API Key 加密存储（AES-256-GCM）
- 在线运行提示词，记录运行历史
- 工作空间隔离 + 权限控制

## 快速开始

### Docker Compose（推荐）

```bash
cp .env.example .env  # 按需修改密钥
docker compose up -d
```

访问 `http://localhost:8080`。

### 本地开发

**依赖**: Go 1.26+, Node.js 22+, pnpm, MySQL 8.0

```bash
# 启动 MySQL
docker compose up -d mysql

# 后端
cd backend
cp .env.example .env
go run cmd/main.go

# 前端（另一个终端）
cd frontend
pnpm install
pnpm dev
```

前端开发服务器 `http://localhost:5173`，自动代理 `/api` 到后端。

## 项目结构

```
├── backend/
│   ├── cmd/main.go              # 入口
│   ├── ioc/                     # 依赖初始化
│   ├── internal/
│   │   ├── config/              # 配置
│   │   ├── domain/              # 领域模型
│   │   ├── service/             # 业务逻辑
│   │   ├── repo/                # 仓储层
│   │   │   └── dao/             # GORM 实现
│   │   └── web/
│   │       ├── handler/         # HTTP handler
│   │       ├── middleware/      # 中间件
│   │       ├── result/          # 统一响应
│   │       └── router/          # 路由
│   └── pkg/                     # 公共工具
│       ├── crypto/              # AES 加密
│       ├── openai/              # OpenAI-compatible client
│       ├── session/             # Session 签名
│       └── vars/                # 变量提取
├── frontend/                    # React + Vite
├── docker-compose.yml
├── Dockerfile
└── Makefile
```

## 常用命令

```bash
make build      # 编译后端
make run        # 运行后端
make test       # 跑测试
make vet        # 静态检查
make tidy       # 整理依赖
make clean      # 清理 8080 端口
```

## 环境变量

| 变量 | 说明 | 必填 |
|------|------|------|
| `APP_ENV` | 环境（development/production） | 否 |
| `HTTP_ADDR` | 监听地址 | 否，默认 `:8080` |
| `MYSQL_DSN` | MySQL 连接串 | 是 |
| `SESSION_SECRET` | Session 密钥 | 否 |
| `SECRET_ENCRYPTION_KEY` | API Key 加密密钥（64 位 hex） | 否（生产必填） |

## API 概览

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/auth/signup` | 注册 |
| POST | `/api/auth/login` | 登录 |
| GET/POST/PATCH | `/api/prompts` | 提示词 CRUD |
| POST | `/api/prompts/:id/run` | 运行提示词 |
| GET | `/api/prompts/:id/runs` | 运行记录 |
| GET | `/api/runs` | 全局运行记录 |
| GET/PATCH | `/api/prompts/:id/variables` | 变量管理 |
| GET/POST/PATCH/DELETE | `/api/prompts/:id/test-cases` | 测试用例 |
| GET/POST | `/api/prompts/:id/versions` | 版本管理 |
| GET/POST/DELETE | `/api/settings/providers` | Provider 配置 |

## License

MIT
