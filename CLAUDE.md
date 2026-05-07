# CLAUDE.md

Always reply in Chinese.
This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**prompthub** — 提示词管理/共享平台。Go + Gin + MySQL 模块化单体架构。

- **Remote**: `git@github.com:mymikasa/prompthub.git`
- **License**: MIT (Copyright 2026 Mikasa)

## Common Commands

```bash
cd backend

# 编译
go build ./...

# 静态检查
go vet ./...

# 整理依赖
go mod tidy

# 运行（需要先配置 .env）
cp .env.example .env  # 编辑 MySQL DSN 等配置
go run cmd/main.go

# 测试
go test ./...
```

## Architecture

```
backend/
├── cmd/main.go              # 入口，组装启动
├── ioc/                     # 依赖初始化（DB 等）
├── internal/
│   ├── config/              # 配置加载（环境变量 + .env）
│   ├── domain/              # 领域模型
│   ├── service/             # 业务逻辑层
│   ├── repo/                # 仓储接口
│   │   └── dao/             # GORM 数据访问实现
│   └── web/
│       ├── handler/         # HTTP handler
│       ├── middleware/       # 中间件（recover、request log）
│       ├── result/          # 统一响应结构
│       └── router/          # 路由注册
├── pkg/                     # 公共工具包
└── frontend/                # 前端应用（独立）
```

## Key Patterns

- **配置**: 通过 `godotenv` 加载 `.env`，`os.Getenv` 读取，必填字段在启动时校验
- **响应**: 统一 `result.Response{code, message, data}`，通过 `result.OK` / `result.BadRequest` 等方法返回
- **错误码**: 五位数字，前三位对应 HTTP 状态码（如 `40000` = 400，`50000` = 500）
- **中间件**: `middleware.Recovery()` 捕获 panic，`middleware.RequestLog()` 记录请求日志
- **ORM**: GORM + MySQL，开发模式开启 Info 日志

## Environment Variables

| 变量 | 说明 | 必填 |
|------|------|------|
| `APP_ENV` | 环境（development/production） | 否，默认 development |
| `HTTP_ADDR` | 监听地址 | 否，默认 :8080 |
| `MYSQL_DSN` | MySQL 连接串 | 是 |
| `SESSION_SECRET` | Session 密钥 | 否（阶段 2 需要） |
| `SECRET_ENCRYPTION_KEY` | 加密密钥 | 否（阶段 10 需要） |
