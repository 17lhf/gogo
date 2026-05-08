# 终端管理系统 · 需求文档

> 状态：已确认  
> 更新：2026-05-08  
> 所有待确认事项已逐一确认完毕

---

## 一、项目概述

基于 Go + Gin + GORM + PostgreSQL + Redis 构建的后端 API 服务，为线下**门店/业务终端**提供统一管理平台。

系统包含五个核心模块：用户管理、角色管理、菜单管理、终端管理、日志管理。

---

## 二、技术选型（已确认）

| 分类 | 选型 |
|------|------|
| Web 框架 | `gin-gonic/gin` |
| ORM | `gorm.io/gorm` + `gorm.io/driver/postgres` |
| 缓存 | `redis/go-redis/v9` |
| Session | `gin-contrib/sessions/redis`，Cookie 存 session_id，数据存 Redis |
| 权限 | `casbin/casbin/v2` + `casbin/gorm-adapter/v3`，策略存 PostgreSQL |
| API 文档 | `swaggo/swag` + `swaggo/gin-swagger` |
| 参数校验 | `go-playground/validator/v10`（Gin 内置） |
| 结构化日志 | `uber-go/zap` |
| 开发规范 | Google Go Style Guide |

---

## 三、已确认决策

| # | 问题 | 决策 |
|---|------|------|
| Q1 | 用户字段 | 保留 `real_name` + `phone` |
| Q2 | 部门/组织架构 | 不需要 |
| Q3 | 主键类型 | `BIGSERIAL` 自增 |
| Q4 | 删除策略 | **全部物理删除**，无软删除字段 |
| Q5 | 超级管理员 | 内置 `SUPER_ADMIN` 角色，绕过 Casbin 校验 |
| Q6 | 数据隔离 | 用户归属门店，只能操作自己门店的终端，新增 `user_stores` 关联表 |
| Q7 | Session 有效期 | **8 小时** |
| Q8 | 密码强度 | 最小 8 位，必须含大小写字母 + 数字 |
| Q9 | 终端 SN | **手动填写**，系统不自动生成 |
| Q10 | 心跳认证 | 设备专属 Token（注册时生成，存于 `terminals.device_token`） |
| Q11 | Config Push | **不需要**，终端主动拉取即可 |
| Q12 | DB 初始化 | 新项目直接建表，使用 GORM `AutoMigrate` |
| Q13 | 响应格式 | 标准 JSON 封装，见下方规范 |
| Q14 | 日志保留 | **保留 180 天**，定期清理过期数据 |

---

## 四、统一响应格式

```json
// 成功（单条）
{ "code": 0, "msg": "success", "data": { ... } }

// 成功（列表）
{ "code": 0, "msg": "success", "data": { "list": [...], "total": 100, "page": 1, "page_size": 20 } }

// 失败
{ "code": 40001, "msg": "参数错误：username 不能为空", "data": null }
```

错误码约定：

| 范围 | 含义 |
|------|------|
| 0 | 成功 |
| 40001–40099 | 参数错误 |
| 40101–40199 | 认证错误（未登录、Token 无效） |
| 40301–40399 | 权限不足 |
| 50001–50099 | 服务器内部错误 |

---

## 五、模块功能需求

### 5.1 用户管理（User）

| 接口 | 方法 | 路径 | 权限码 |
|------|------|------|--------|
| 用户列表 | GET | `/api/v1/users` | `sys:user:list` |
| 用户详情 | GET | `/api/v1/users/:id` | `sys:user:list` |
| 创建用户 | POST | `/api/v1/users` | `sys:user:add` |
| 更新用户 | PUT | `/api/v1/users/:id` | `sys:user:edit` |
| 删除用户 | DELETE | `/api/v1/users/:id` | `sys:user:delete` |
| 重置密码 | PUT | `/api/v1/users/:id/password` | `sys:user:reset-pwd` |
| 分配角色 | PUT | `/api/v1/users/:id/roles` | `sys:user:assign-role` |
| 分配门店 | PUT | `/api/v1/users/:id/stores` | `sys:user:assign-store` |
| 登录 | POST | `/api/v1/auth/login` | 公开 |
| 登出 | POST | `/api/v1/auth/logout` | 已登录 |
| 当前用户信息 + 菜单 | GET | `/api/v1/auth/me` | 已登录 |

**业务规则**

- 密码 bcrypt 存储（cost=12），最小 8 位，必须含大小写字母 + 数字
- 登录失败连续 5 次，账户锁定 30 分钟（Redis 计数器 `login_fail:{username}`）
- Session 有效期 8 小时，存于 Redis
- `SUPER_ADMIN` 角色在认证中间件中直接放行，不进入 Casbin

---

### 5.2 角色管理（Role）

| 接口 | 方法 | 路径 | 权限码 |
|------|------|------|--------|
| 角色列表 | GET | `/api/v1/roles` | `sys:role:list` |
| 创建角色 | POST | `/api/v1/roles` | `sys:role:add` |
| 更新角色 | PUT | `/api/v1/roles/:id` | `sys:role:edit` |
| 删除角色 | DELETE | `/api/v1/roles/:id` | `sys:role:delete` |
| 分配菜单权限 | PUT | `/api/v1/roles/:id/menus` | `sys:role:assign-menu` |
| 查询角色菜单 | GET | `/api/v1/roles/:id/menus` | `sys:role:list` |

---

### 5.3 菜单管理（Menu）

**菜单类型**

| type | 含义 | perms 字段 |
|------|------|-----------|
| 1 | 目录 | 空 |
| 2 | 页面 | 空（页面可见性由菜单关联控制） |
| 3 | 按钮 | 必填，如 `sys:user:add` |

| 接口 | 方法 | 路径 | 权限码 |
|------|------|------|--------|
| 菜单树（全量） | GET | `/api/v1/menus` | `sys:menu:list` |
| 创建菜单 | POST | `/api/v1/menus` | `sys:menu:add` |
| 更新菜单 | PUT | `/api/v1/menus/:id` | `sys:menu:edit` |
| 删除菜单 | DELETE | `/api/v1/menus/:id` | `sys:menu:delete` |

`GET /api/v1/auth/me` 返回中包含当前用户的菜单树（根据角色过滤）。

---

### 5.4 终端管理（Terminal）

#### 门店（Store）

| 接口 | 方法 | 路径 | 权限码 |
|------|------|------|--------|
| 门店列表 | GET | `/api/v1/stores` | `sys:store:list` |
| 门店详情 | GET | `/api/v1/stores/:id` | `sys:store:list` |
| 创建门店 | POST | `/api/v1/stores` | `sys:store:add` |
| 更新门店 | PUT | `/api/v1/stores/:id` | `sys:store:edit` |
| 删除门店 | DELETE | `/api/v1/stores/:id` | `sys:store:delete` |

#### 终端（Terminal）

| 接口 | 方法 | 路径 | 权限码 |
|------|------|------|--------|
| 终端列表 | GET | `/api/v1/terminals` | `sys:terminal:list` |
| 终端详情 | GET | `/api/v1/terminals/:id` | `sys:terminal:list` |
| 创建终端 | POST | `/api/v1/terminals` | `sys:terminal:add` |
| 更新终端 | PUT | `/api/v1/terminals/:id` | `sys:terminal:edit` |
| 删除终端 | DELETE | `/api/v1/terminals/:id` | `sys:terminal:delete` |
| 心跳上报 | POST | `/api/v1/terminals/:sn/heartbeat` | 设备 Token（Header: `X-Device-Token`） |

**终端状态**：`online` / `offline` / `disabled`

**数据隔离规则**：

- 普通用户查询终端列表时，只返回其归属门店下的终端
- `SUPER_ADMIN` 可查看全量数据
- 中间件从 Session 中取用户 ID，查 `user_stores` 表得到可见门店列表，注入 Gin Context

---

### 5.5 日志管理（Log）

**① 操作审计日志** — Gin middleware 自动记录所有写操作（POST / PUT / DELETE）

| 接口 | 方法 | 路径 | 权限码 |
|------|------|------|--------|
| 操作日志列表 | GET | `/api/v1/logs/operations` | `sys:log:list` |

查询过滤：`user_id`、`action`、`status`、时间范围（`start_time` / `end_time`）

**② 终端设备日志** — 终端心跳超时、上下线事件

事件类型：`online` / `offline` / `heartbeat_timeout`

| 接口 | 方法 | 路径 | 权限码 |
|------|------|------|--------|
| 终端日志列表 | GET | `/api/v1/logs/terminals` | `sys:log:list` |

查询过滤：`sn`、`event_type`、时间范围

**日志保留策略**：保留 180 天，通过定时任务（Go `time.Ticker` 或外部 cron）定期删除 `created_at < NOW() - INTERVAL '180 days'` 的记录。

---

## 六、数据库设计

> 全部物理删除，无 `deleted_at` 字段。主键均为 `BIGSERIAL`。

### 用户相关

```sql
CREATE TABLE users (
    id            BIGSERIAL    PRIMARY KEY,
    username      VARCHAR(64)  NOT NULL UNIQUE,
    email         VARCHAR(128) NOT NULL UNIQUE,
    password      VARCHAR(256) NOT NULL,           -- bcrypt hash
    real_name     VARCHAR(64),
    phone         VARCHAR(20),
    status        SMALLINT     NOT NULL DEFAULT 1, -- 1=active 2=disabled 3=locked
    last_login_at TIMESTAMPTZ,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE roles (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(64)  NOT NULL,
    code        VARCHAR(64)  NOT NULL UNIQUE,     -- 如 SUPER_ADMIN, OPERATOR
    description VARCHAR(256),
    status      SMALLINT     NOT NULL DEFAULT 1,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE user_roles (
    user_id    BIGINT NOT NULL REFERENCES users(id),
    role_id    BIGINT NOT NULL REFERENCES roles(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);
```

### 菜单相关

```sql
CREATE TABLE menus (
    id         BIGSERIAL    PRIMARY KEY,
    parent_id  BIGINT       NOT NULL DEFAULT 0,
    name       VARCHAR(64)  NOT NULL,
    path       VARCHAR(256),
    component  VARCHAR(256),
    icon       VARCHAR(64),
    type       SMALLINT     NOT NULL,             -- 1=目录 2=页面 3=按钮
    perms      VARCHAR(128),                      -- 按钮权限码
    sort_order INT          NOT NULL DEFAULT 0,
    visible    BOOLEAN      NOT NULL DEFAULT TRUE,
    status     SMALLINT     NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE role_menus (
    role_id    BIGINT NOT NULL REFERENCES roles(id),
    menu_id    BIGINT NOT NULL REFERENCES menus(id),
    PRIMARY KEY (role_id, menu_id)
);
```

### 门店 & 终端相关

```sql
CREATE TABLE stores (
    id         BIGSERIAL    PRIMARY KEY,
    name       VARCHAR(128) NOT NULL,
    code       VARCHAR(64)  NOT NULL UNIQUE,
    address    VARCHAR(256),
    status     SMALLINT     NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- 用户-门店归属（数据隔离）
CREATE TABLE user_stores (
    user_id    BIGINT NOT NULL REFERENCES users(id),
    store_id   BIGINT NOT NULL REFERENCES stores(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, store_id)
);

CREATE TABLE terminals (
    id                BIGSERIAL   PRIMARY KEY,
    sn                VARCHAR(64) NOT NULL UNIQUE, -- 手动填入设备序列号
    name              VARCHAR(128) NOT NULL,
    type              VARCHAR(64),                 -- 终端类型，如 POS/KIOSK
    store_id          BIGINT      REFERENCES stores(id),
    status            VARCHAR(16) NOT NULL DEFAULT 'offline',
    ip_address        VARCHAR(45),
    mac_address       VARCHAR(17),
    device_token      VARCHAR(256) NOT NULL,       -- 心跳认证 Token，注册时生成
    last_heartbeat_at TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 日志相关

```sql
CREATE TABLE operation_logs (
    id            BIGSERIAL    PRIMARY KEY,
    user_id       BIGINT,
    username      VARCHAR(64),
    action        VARCHAR(128) NOT NULL,
    resource_type VARCHAR(64),
    resource_id   VARCHAR(64),
    detail        JSONB,
    ip            VARCHAR(45),
    user_agent    VARCHAR(512),
    status        SMALLINT     NOT NULL,           -- 1=成功 2=失败
    duration_ms   INT,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_operation_logs_created_at ON operation_logs(created_at);

CREATE TABLE terminal_logs (
    id          BIGSERIAL   PRIMARY KEY,
    terminal_id BIGINT,
    sn          VARCHAR(64) NOT NULL,
    event_type  VARCHAR(32) NOT NULL,              -- online/offline/heartbeat_timeout
    detail      JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_terminal_logs_created_at ON terminal_logs(created_at);
```

> 日志表均建 `created_at` 索引，用于 180 天定期清理查询加速。

### Casbin 策略表

由 `casbin/gorm-adapter` 自动创建，存储 `(role, resource, action)` 三元组规则。

---

## 七、项目目录结构

遵循 Google Go Style Guide：

```
gogo/
├── cmd/
│   └── server/
│       └── main.go              # 程序入口，Wire 依赖注入
├── internal/
│   ├── config/                  # 配置加载
│   ├── db/                      # GORM 初始化 + AutoMigrate
│   ├── cache/                   # Redis 客户端
│   ├── model/                   # GORM 模型（纯结构体）
│   │   ├── user.go
│   │   ├── role.go
│   │   ├── menu.go
│   │   ├── store.go
│   │   ├── terminal.go
│   │   └── log.go
│   ├── repository/              # 数据访问层（接口 + GORM 实现）
│   ├── service/                 # 业务逻辑层（依赖 repository 接口）
│   ├── handler/                 # Gin Handler（依赖 service 接口）
│   ├── middleware/
│   │   ├── auth.go              # Session 认证，注入当前用户到 Context
│   │   ├── casbin.go            # RBAC 权限校验（SUPER_ADMIN 跳过）
│   │   ├── store_scope.go       # 数据隔离：查询时限制可见门店范围
│   │   └── audit.go             # 写操作自动记录审计日志
│   ├── router/
│   │   └── router.go            # 路由注册
│   └── pkg/
│       ├── response/            # 统一响应封装
│       ├── errcode/             # 错误码定义
│       └── password/            # bcrypt 工具
├── docs/                        # swaggo 生成的 Swagger 文件
├── go.mod
└── go.sum
```

---

## 八、权限校验流程

```
请求进入
  → auth.go 中间件：校验 Session，取出 user_id + roles，注入 Context
      ↓ 未登录 → 返回 401
  → casbin.go 中间件：判断 roles 是否含 SUPER_ADMIN
      ↓ 含 SUPER_ADMIN → 直接放行
      ↓ 普通角色 → Casbin Enforce(role, path, method) → 无权限返回 403
  → store_scope.go 中间件（终端/门店接口）：从 user_stores 查可见门店，注入 Context
  → Handler 执行业务逻辑
  → audit.go 中间件（After）：写操作记录审计日志
```
