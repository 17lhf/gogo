# 终端管理系统 · 需求文档

> 状态：已确认  
> 更新：2026-05-08  
> 所有待确认事项已逐一确认完毕

---

## 一、项目概述

后端 API 服务，为线下**门店/业务终端**提供统一管理平台。

系统包含五个核心模块：用户管理、角色管理、菜单管理、终端管理、日志管理。

---

## 二、已确认决策

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
| Q12 | DB 初始化 | 表结构由人工管理，项目不执行 AutoMigrate |
| Q13 | 响应格式 | 标准 JSON 封装，见下方规范 |
| Q14 | 日志保留 | **保留 180 天**，定期清理过期数据 |
| Q15 | 心跳超时判定 | Redis TTL + keyspace notification 过期回调，阈值 **60 秒** |
| Q16 | 终端禁用策略 | 禁用即断连：清理 Redis key，下次心跳返回特定错误码 |
| Q17 | 终端创建方式 | **管理员预录入**，手动填写 SN，系统生成 UUID device_token |
| Q18 | device_token | UUID 格式，终端通过 `POST /api/v1/terminals/:sn/rotate-token` 主动轮换；管理员不可手动重置；泄露时运维从 Redis 移除 |
| Q19 | 终端删除 | 清理 Redis heartbeat key；若终端继续心跳返回 404"终端不存在"；历史操作日志和终端日志中数据保留 |
| Q20 | 审计日志范围 | **选择性记录**读操作，白名单：`GET /api/v1/auth/me`，其余 GET 不记录 |
| Q21 | 日志清理调度 | 后续引入 gocron 等定时任务框架内部执行 |
| Q22 | 门店删除 | 门店下有终端时**禁止删除** |
| Q23 | store.code | 门店编号（如 `SH001`），纯内部标识，用户不可见 |
| Q24 | 数据权限粒度 | **门店级隔离**，门店内所有终端对同一管理员全可见 |
| Q25 | 超级管理员初始化 | **数据库迁移脚本**：INSERT `SUPER_ADMIN` 角色 + `admin` 用户 |
| Q26 | 用户改密 | 提供独立改密接口 `PUT /api/v1/auth/password`，需提交旧密码 |
| Q27 | 密码策略 | bcrypt cost=12，最小 8 位含大小写+数字，**365 天过期**，不禁止密码历史重复 |
| Q28 | 首次登录 | 管理员重置密码后首次登录**强制修改密码** |
| Q29 | 密码过期处理 | 密码过期后**允许登录**，但中间件拦截非改密接口（`PUT /api/v1/auth/password`），前端强制跳转改密页 |

---

## 三、统一响应格式

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

## 四、模块功能需求

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
| 修改密码 | PUT | `/api/v1/auth/password` | 已登录 |

**业务规则**

密码策略：
- bcrypt 存储（cost=12），最小 8 位，必须含大小写字母 + 数字
- 密码有效期 **365 天**（从 `password_updated_at` 起算）
- 管理员重置密码后，`must_change_password` 置为 `true`，用户首次登录**强制修改密码**
- 密码过期后**允许登录**，但中间件拦截非改密接口（仅放行 `PUT /api/v1/auth/password`），前端强制跳转改密页
- 不禁止密码历史重复使用

登录与锁定：
- 登录失败连续 5 次，账户锁定 30 分钟（Redis 计数器 `login_fail:{username}`）
- Session 有效期 8 小时，存于 Redis

管理员重置密码：
- 管理员通过 `PUT /api/v1/users/:id/password` 重置任意用户密码，重置后 `must_change_password = true`
- 用户通过 `PUT /api/v1/auth/password` 自行修改密码（需提交旧密码验证），修改后 `password_updated_at` 更新，`must_change_password = false`

超级管理员：
- `SUPER_ADMIN` 角色在认证中间件中直接放行，不进入 Casbin 校验

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

**业务规则**：

- `code` 为门店编号（如 `SH001`），纯内部标识，用户不可见
- 删除门店时，若该门店下存在终端（`terminals.store_id` 引用该门店），**禁止删除**

#### 终端（Terminal）

| 接口 | 方法 | 路径 | 权限码 |
|------|------|------|--------|
| 终端列表 | GET | `/api/v1/terminals` | `sys:terminal:list` |
| 终端详情 | GET | `/api/v1/terminals/:id` | `sys:terminal:list` |
| 创建终端 | POST | `/api/v1/terminals` | `sys:terminal:add` |
| 更新终端 | PUT | `/api/v1/terminals/:id` | `sys:terminal:edit` |
| 删除终端 | DELETE | `/api/v1/terminals/:id` | `sys:terminal:delete` |
| 心跳上报 | POST | `/api/v1/terminals/:sn/heartbeat` | 设备 Token（Header: `X-Device-Token`） |
| Token 轮换 | POST | `/api/v1/terminals/:sn/rotate-token` | 设备 Token（Header: `X-Device-Token`） |

**终端状态**：`online` / `offline` / `disabled`

**状态流转**：

```
offline ──心跳恢复──→ online
online  ──60s超时──→ offline  (Redis key 过期回调)
online  ──管理员禁用─→ disabled (同时清理 Redis key，下次心跳返回错误码)
offline ──管理员禁用─→ disabled (同上)
disabled──管理员启用─→ offline  (等下次心跳自动切 online)
```

**心跳机制**：

- 心跳上报时在 Redis 中 SET key `heartbeat:{sn}`，TTL = 60 秒
- 通过 Redis keyspace notification 监听过期事件，超时回调：更新 `terminals.status` → `offline`，写入 `terminal_logs`（event_type = `heartbeat_timeout`）
- 心跳恢复（`offline` → `online`）时写入 `terminal_logs`（event_type = `online`）
- 管理员禁用/启用终端时也写入对应 `terminal_logs` 记录
- `disabled` 状态的终端心跳返回错误，终端应停止上报

**终端创建与凭证**：

- 管理员在后管页面预录入终端（填写 SN、名称、类型、归属门店），系统生成 UUID 格式 `device_token`
- 管理员将 SN + device_token 下发配置到终端设备
- 终端创建后初始状态为 `offline`，首次心跳成功后自动切 `online`

**device_token 轮换**：

- 终端通过 `POST /api/v1/terminals/:sn/rotate-token`（携带旧 token）换取新 token，旧 token 立即失效
- 管理员**不可**手动重置 token
- token 泄露时由运维从 Redis 移除对应 key

**终端删除**：

- 删除终端时清理 Redis heartbeat key
- 若该终端继续心跳，返回 404 "终端不存在"
- 操作日志和终端日志中历史引用该终端的数据保留

**数据隔离规则**：

- 普通用户查询终端列表时，只返回其归属门店下的终端
- `SUPER_ADMIN` 可查看全量数据
- 数据权限为**门店级**，门店内所有终端对同一管理员全可见
- 中间件从 Session 中取用户 ID，查 `user_stores` 表得到可见门店列表，注入请求上下文

---

### 5.5 日志管理（Log）

**① 操作审计日志** — 中间件自动记录所有写操作（POST / PUT / DELETE），以及白名单内的读操作

| 接口 | 方法 | 路径 | 权限码 |
|------|------|------|--------|
| 操作日志列表 | GET | `/api/v1/logs/operations` | `sys:log:list` |

查询过滤：`user_id`、`action`、`status`、时间范围（`start_time` / `end_time`）

审计白名单（读操作记录）：
- `GET /api/v1/auth/me` — 当前用户信息查询

其余 GET 接口不记录审计日志。

**② 终端设备日志** — 终端心跳超时、上下线、管理员禁用/启用事件

事件类型：`online` / `offline` / `heartbeat_timeout` / `disabled` / `enabled`

| 接口 | 方法 | 路径 | 权限码 |
|------|------|------|--------|
| 终端日志列表 | GET | `/api/v1/logs/terminals` | `sys:log:list` |

查询过滤：`sn`、`event_type`、时间范围

**日志保留策略**：保留 180 天，通过 gocron 等定时任务框架执行清理，定期删除 `created_at < NOW() - INTERVAL '180 days'` 的记录。

---

## 五、数据库设计

> 全部物理删除，无 `deleted_at` 字段。主键均为 `BIGSERIAL`。

### 用户相关

```sql
CREATE TABLE users (
    id                    BIGSERIAL    PRIMARY KEY,
    username              VARCHAR(64)  NOT NULL UNIQUE,
    email                 VARCHAR(128) NOT NULL UNIQUE,
    password              VARCHAR(256) NOT NULL,           -- bcrypt hash
    real_name             VARCHAR(64),
    phone                 VARCHAR(20),
    status                SMALLINT     NOT NULL DEFAULT 1, -- 1=active 2=disabled 3=locked
    must_change_password  BOOLEAN      NOT NULL DEFAULT FALSE, -- 首次登录/重置后强制改密
    password_updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),  -- 密码修改时间（用于365天过期判定）
    last_login_at         TIMESTAMPTZ,
    created_at            TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ  NOT NULL DEFAULT NOW()
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
    event_type  VARCHAR(32) NOT NULL,              -- online/offline/heartbeat_timeout/disabled/enabled
    detail      JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_terminal_logs_created_at ON terminal_logs(created_at);
```

> 日志表均建 `created_at` 索引，用于 180 天定期清理查询加速。

### 权限策略表

由权限框架自动管理，存储 `(role, resource, action)` 三元组规则。

---

## 七、权限校验流程

```
请求进入
  → auth 中间件：校验 Session，取出 user_id + roles，注入请求上下文
      ↓ 未登录 → 返回 401
  → 密码过期检查中间件：若 password_updated_at > 365 天 或 must_change_password = true
      ↓ 仅放行 PUT /api/v1/auth/password，其他接口返回 403（强制改密）
  → 权限中间件：判断 roles 是否含 SUPER_ADMIN
      ↓ 含 SUPER_ADMIN → 直接放行
      ↓ 普通角色 → RBAC Enforce(role, path, method) → 无权限返回 403
  → 数据隔离中间件（终端/门店接口）：从 user_stores 查可见门店，注入请求上下文
  → Handler 执行业务逻辑
  → 审计中间件（After）：写操作 + 白名单读操作记录审计日志
```
