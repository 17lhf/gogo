# gogo - 终端管理系统

基于 Go 的后端 API 服务，为线下门店/业务终端提供统一管理平台，包含用户管理、角色管理、菜单管理、终端管理和日志管理五大核心模块。

## 技术栈

| 类别 | 方案 |
|------|------|
| 语言 | Go 1.26 |
| HTTP 框架 | Gin |
| ORM | GORM + pgx 混用 |
| 数据库 | PostgreSQL 16, pgxpool/v5 |
| 缓存 | Redis 7, go-redis/v9 |
| 认证 | golang-jwt/jwt/v5, bcrypt (cost=12) |
| 权限 | Casbin + GORM Adapter |
| 校验 | go-playground/validator |
| 日志 | slog (结构化日志) |
| 测试 | testify |
| Mock | uber-go/mock |
| UUID | google/uuid |
| 配置管理 | viper |
| 国际化 | go-i18n/v2 + BurntSushi/toml |

## 项目结构

```
gogo/
├── main.go                       # 入口，依赖注入
├── internal/
│   ├── config/                   # viper 配置管理 (env + 默认值)
│   │   ├── config.go             # 配置聚合
│   │   ├── postgres.go           # PostgreSQL 连接参数
│   │   ├── redis.go              # Redis 连接参数
│   │   ├── server.go             # HTTP 端口
│   │   └── auth.go               # JWT 密钥、Session TTL、锁定策略
│   ├── db/
│   │   ├── postgres.go           # pgxpool 连接池
│   │   ├── gorm.go               # GORM 初始化
│   │   └── seed.go               # 初始数据播种
│   ├── cache/
│   │   ├── redis.go              # Redis 客户端
│   │   ├── session.go            # JWT Session 存储 (8h TTL)
│   │   ├── lockout.go            # 登录失败计数器 (5次/30min)
│   │   └── heartbeat.go          # 终端心跳 TTL + keyspace 过期监听
│   ├── model/                    # GORM 模型
│   │   ├── user.go               # 用户
│   │   ├── role.go               # 角色
│   │   ├── menu.go               # 菜单 (目录/页面/按钮)
│   │   ├── store.go              # 门店
│   │   ├── terminal.go           # 终端
│   │   ├── operation_log.go      # 操作审计日志
│   │   └── terminal_log.go       # 终端事件日志
│   ├── repository/               # 数据访问层 (接口+实现)
│   │   ├── user.go               # 用户 CRUD + 角色/门店关联
│   │   ├── role.go               # 角色 CRUD + 菜单关联
│   │   ├── menu.go               # 菜单 CRUD + 树形查询
│   │   ├── store.go              # 门店 CRUD + 终端检测
│   │   ├── terminal.go           # 终端 CRUD + 心跳/状态
│   │   └── log.go                # 操作日志 + 终端日志 + 过期清理
│   ├── service/                  # 业务逻辑层
│   │   ├── auth.go               # 登录/登出/个人信息/改密
│   │   ├── user.go               # 用户管理
│   │   ├── role.go               # 角色管理
│   │   ├── menu.go               # 菜单管理 + 树构建
│   │   ├── store.go              # 门店管理
│   │   └── terminal.go           # 终端管理 + 心跳 + Token 轮换 + 状态机
│   ├── handler/                  # HTTP 处理层 (薄层)
│   │   ├── auth.go               # /api/v1/auth/*
│   │   ├── user.go               # /api/v1/users/*
│   │   ├── role.go               # /api/v1/roles/*
│   │   ├── menu.go               # /api/v1/menus/*
│   │   ├── store.go              # /api/v1/stores/*
│   │   ├── terminal.go           # /api/v1/terminals/*
│   │   └── log.go                # /api/v1/logs/*
│   ├── middleware/               # 中间件链
│   │   ├── auth.go               # JWT 解析 + Redis Session 验证
│   │   ├── permission.go         # Casbin RBAC (SUPER_ADMIN 绕过)
│   │   ├── datascope.go          # 门店级数据隔离
│   │   ├── password_expiry.go    # 密码过期/强制改密检查
│   │   └── audit.go              # 操作审计日志 (写操作+白名单读)
│   ├── router/
│   │   └── router.go             # 路由注册 + 中间件编排
│   ├── i18n/
│   │   ├── i18n.go               # Bundle 初始化 + 语言检测中间件
│   │   ├── messages.go           # 消息 ID 常量
│   │   └── locales/
│   │       ├── active.zh-CN.toml # 中文翻译 (默认回退)
│   │       └── active.en-US.toml # 英文翻译
│   ├── casbin/
│   │   └── enforcer.go           # Casbin 初始化 + 策略管理
│   ├── dto/
│   │   ├── request.go            # 所有请求 DTO (含 validator 标签)
│   │   └── response.go           # 响应 DTO
│   └── pkg/                      # 工具包
│       ├── jwt.go                # JWT 生成/验证 (HS256, 8h)
│       ├── hash.go               # bcrypt 哈希/验证 + 密码强度校验
│       └── response.go           # 统一 JSON 响应 (code/msg/data)
├── test/
│   └── simple.http               # HTTP 测试文件
├── docker-compose.yml            # PostgreSQL 16 + Redis 7
├── .env.example                  # 环境变量模板
├── go.mod
└── go.sum
```

## 中间件链

请求处理流程如下：

```
Recovery → i18n (语言检测)
  → Auth (JWT 验证 + Redis Session 校验)
    → PasswordExpiry (密码过期检查，仅放行 /auth/password 和 /auth/me)
      → Permission (Casbin RBAC，SUPER_ADMIN 直接放行)
        → DataScope (注入用户可访问的门店 ID 列表)
          → Handler (执行业务逻辑)
            → Audit (事后记录写操作 + 白名单 GET)
```

## API 接口

### 认证 (Auth)
| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | `/api/v1/auth/login` | 登录 | 公开 |
| POST | `/api/v1/auth/logout` | 登出 | 已登录 |
| GET | `/api/v1/auth/me` | 当前用户信息+菜单树 | 已登录 |
| PUT | `/api/v1/auth/password` | 修改密码 | 已登录 |

### 用户管理 (User)
| 方法 | 路径 | 权限码 |
|------|------|--------|
| GET | `/api/v1/users` | `sys:user:list` |
| GET | `/api/v1/users/:id` | `sys:user:list` |
| POST | `/api/v1/users` | `sys:user:add` |
| PUT | `/api/v1/users/:id` | `sys:user:edit` |
| DELETE | `/api/v1/users/:id` | `sys:user:delete` |
| PUT | `/api/v1/users/:id/password` | `sys:user:reset-pwd` |
| PUT | `/api/v1/users/:id/roles` | `sys:user:assign-role` |
| PUT | `/api/v1/users/:id/stores` | `sys:user:assign-store` |

### 角色管理 (Role)
| 方法 | 路径 | 权限码 |
|------|------|--------|
| GET | `/api/v1/roles` | `sys:role:list` |
| GET | `/api/v1/roles/:id` | `sys:role:list` |
| POST | `/api/v1/roles` | `sys:role:add` |
| PUT | `/api/v1/roles/:id` | `sys:role:edit` |
| DELETE | `/api/v1/roles/:id` | `sys:role:delete` |
| GET | `/api/v1/roles/:id/menus` | `sys:role:list` |
| PUT | `/api/v1/roles/:id/menus` | `sys:role:assign-menu` |

### 菜单管理 (Menu)
| 方法 | 路径 | 权限码 |
|------|------|--------|
| GET | `/api/v1/menus` | `sys:menu:list` |
| GET | `/api/v1/menus/:id` | `sys:menu:list` |
| POST | `/api/v1/menus` | `sys:menu:add` |
| PUT | `/api/v1/menus/:id` | `sys:menu:edit` |
| DELETE | `/api/v1/menus/:id` | `sys:menu:delete` |

### 门店管理 (Store)
| 方法 | 路径 | 权限码 |
|------|------|--------|
| GET | `/api/v1/stores` | `sys:store:list` |
| GET | `/api/v1/stores/:id` | `sys:store:list` |
| POST | `/api/v1/stores` | `sys:store:add` |
| PUT | `/api/v1/stores/:id` | `sys:store:edit` |
| DELETE | `/api/v1/stores/:id` | `sys:store:delete` |

### 终端管理 (Terminal)
| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/v1/terminals` | 终端列表 | `sys:terminal:list` |
| GET | `/api/v1/terminals/:id` | 终端详情 | `sys:terminal:list` |
| POST | `/api/v1/terminals` | 创建终端 | `sys:terminal:add` |
| PUT | `/api/v1/terminals/:id` | 更新终端 | `sys:terminal:edit` |
| DELETE | `/api/v1/terminals/:id` | 删除终端 | `sys:terminal:delete` |
| POST | `/api/v1/terminals/:sn/heartbeat` | 心跳上报 | X-Device-Token |
| POST | `/api/v1/terminals/:sn/rotate-token` | Token 轮换 | X-Device-Token |

### 日志管理 (Log)
| 方法 | 路径 | 权限码 |
|------|------|--------|
| GET | `/api/v1/logs/operations` | `sys:log:list` |
| GET | `/api/v1/logs/terminals` | `sys:log:list` |

## 统一响应格式

```json
// 成功（单条）
{ "code": 0, "msg": "success", "data": { ... } }

// 成功（列表）
{ "code": 0, "msg": "success", "data": { "list": [...], "total": 100, "page": 1, "page_size": 20 } }

// 失败
{ "code": 40001, "msg": "参数错误：username 不能为空", "data": null }
```

### 错误码
| 范围 | 含义 |
|------|------|
| 0 | 成功 |
| 40001–40099 | 参数错误 |
| 40101–40199 | 认证错误 |
| 40301–40399 | 权限不足 |
| 50001–50099 | 服务器内部错误 |

## 国际化 (i18n)

API 响应消息支持多语言，通过 `Accept-Language` 请求头切换。

| 语言 | Accept-Language | 说明 |
|------|-----------------|------|
| 简体中文 | `zh-CN` (或不传) | 默认语言 |
| English | `en-US` | |

错误码 `code` 不受语言影响，始终为数字，可用于前端编程判断；`msg` 字段根据语言返回对应文本。

### 新增语言

1. 在 `internal/i18n/locales/` 下创建 `active.{lang}.toml`
2. 按现有 TOML 格式翻译所有消息 ID
3. 在 `internal/i18n/i18n.go` 的 `init()` 中 `LoadMessageFileFS` 加载新文件

## 核心设计

### 终端状态机
```
offline ──心跳恢复──→ online
online  ──60s超时──→ offline  (Redis key 过期回调)
online  ──管理员禁用─→ disabled (清理 Redis key)
offline ──管理员禁用─→ disabled
disabled──管理员启用─→ offline  (等下次心跳切 online)
```

### 密码策略
- bcrypt 存储 (cost=12)
- 最小 8 位，必须含大写字母 + 小写字母 + 数字
- 有效期 365 天 (从 `password_updated_at` 起算)
- 连续 5 次登录失败锁定 30 分钟
- 管理员重置密码后首次登录强制修改密码

### 数据隔离
- `SUPER_ADMIN` 可查看全量数据
- 普通用户只能访问其归属门店下的终端
- 数据权限为门店级，门店内所有终端对同一管理员全可见

### 心跳机制
- 终端心跳时：Redis 设置 `heartbeat:{sn}` key，TTL = 60 秒
- Redis keyspace notification 监听过期事件，超时回调更新 `terminals.status` → `offline`
- `disabled` 状态终端心跳返回 `40305` 错误码
- 心跳认证通过 `X-Device-Token` 请求头

### 审计日志
- 自动记录所有 POST/PUT/DELETE 操作
- 白名单读操作仅记录 `GET /api/v1/auth/me`
- 日志保留 180 天，支持定时清理

## 快速开始

### 前置要求
- Go 1.26+
- Docker Compose

### 启动
```bash
# 启动 PostgreSQL 和 Redis
docker-compose up -d

# 安装依赖
go mod tidy

# 运行服务（开发环境自动加载 .env.example）
go run main.go
```

`APP_ENV` 默认为 `development`，启动时会自动加载项目根目录下的 `.env.example`。
生产环境部署时设置 `APP_ENV=production`，应用将直接读取系统环境变量，不依赖 `.env.example` 文件。

服务默认监听 `:8080`。

### 默认账户
| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | Admin123! | SUPER_ADMIN |

### 运行测试
```bash
# 单元测试 (无需数据库)
go test ./internal/... -v

# 集成测试 (需要 PostgreSQL 和 Redis)
go test ./internal/... -v -tags=integration
```

## 测试覆盖

| 包 | 测试文件 | 覆盖内容 |
|----|----------|----------|
| `pkg` | hash_test.go, jwt_test.go, response_test.go | bcrypt 哈希/校验、JWT 签发/解析、响应格式 |
| `service` | auth_test.go, user_test.go, menu_test.go, terminal_test.go | 认证逻辑、用户 CRUD、菜单树构建、终端操作 |
| `handler` | auth_test.go | 中间件上下文、Auth Bearer Token、handler 输入验证 |
| `middleware` | permission_test.go | SUPER_ADMIN 绕过、角色检查 |
| `router` | router_test.go | 路由注册、404 处理、响应格式 |
