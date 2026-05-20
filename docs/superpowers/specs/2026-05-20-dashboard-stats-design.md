# Dashboard 统计 API 设计

## 概述

在 `/api/v1/stats/*` 下新增两个只读统计接口，返回终端和用户的聚合数据，用于首页看板展示。

- `GET /api/v1/stats/terminals` — 终端状态分布、门店维度统计、近期新增
- `GET /api/v1/stats/users` — 用户状态分布、角色维度统计、近期新增

## 路由

```
GET /api/v1/stats/terminals
GET /api/v1/stats/users
```

两个路由挂载在已有的 `protected` 中间件组下（Auth → PasswordExpiry → Permission → DataScope）。不加 Audit 中间件（只读接口无需审计）。

终端统计受 DataScope 约束（非 SUPER_ADMIN 仅能看到归属门店下的终端数据）。用户统计为全局视角，不做门店级隔离。

## 架构

沿用现有三层模式：

```
handler → service → repository → GORM
```

### 新增文件

| 文件 | 用途 |
|------|------|
| `internal/repository/stats.go` | `StatsRepository` 接口 + GORM 实现 |
| `internal/service/stats.go` | `StatsService`，组合 repo 调用，接收数据权限参数 |
| `internal/handler/stats.go` | `StatsHandler`，两个 handler 方法 |

### 修改文件

| 文件 | 变更 |
|------|------|
| `internal/dto/response.go` | 新增 `TerminalStatsResp` 和 `UserStatsResp` 响应 DTO |
| `internal/router/router.go` | 注册路由；`Dependencies` 新增 `StatsHandler` |
| `main.go` | 依赖注入：`StatsRepository` → `StatsService` → `StatsHandler` |

## 数据结构

### Request

无请求体，无查询参数（当前版本）。

### 响应：`GET /api/v1/stats/terminals`

```json
{
  "code": 0,
  "data": {
    "status_distribution": {
      "online": 10,
      "offline": 5,
      "disabled": 2
    },
    "by_store": [
      {
        "store_id": 1,
        "store_name": "王府井店",
        "total": 8,
        "online": 3,
        "offline": 4,
        "disabled": 1
      }
    ],
    "recent_added": {
      "last_7_days": 3,
      "last_30_days": 8
    }
  }
}
```

`by_store`、`status_distribution` 和 `recent_added` 均受 DataScope 约束：非管理员只看到归属门店下的数据。

### 响应：`GET /api/v1/stats/users`

```json
{
  "code": 0,
  "data": {
    "status_distribution": {
      "enabled": 20,
      "disabled": 3,
      "locked": 1
    },
    "by_role": [
      {
        "role_id": 1,
        "role_name": "SUPER_ADMIN",
        "count": 2
      }
    ],
    "recent_added": {
      "last_7_days": 2,
      "last_30_days": 5
    }
  }
}
```

用户统计不做门店隔离，角色分布和状态分布为全局视图。

## Repository

### StatsRepository 接口

```go
type StatsRepository interface {
    // 终端统计（storeIDs 用于数据范围过滤，空切片 = 全量）
    TerminalStatusDistribution(ctx context.Context, storeIDs []int64) (map[string]int64, error)
    TerminalByStore(ctx context.Context, storeIDs []int64) ([]TerminalStoreStat, error)
    TerminalRecentAdded(ctx context.Context, storeIDs []int64) (last7Days int64, last30Days int64, err error)

    // 用户统计（全局，不做数据隔离）
    UserStatusDistribution(ctx context.Context) (map[int16]int64, error)
    UserByRole(ctx context.Context) ([]UserRoleStat, error)
    UserRecentAdded(ctx context.Context) (last7Days int64, last30Days int64, err error)
}
```

### 查询说明 (GORM)

- **TerminalStatusDistribution**: `SELECT status, COUNT(*) FROM terminals WHERE ... GROUP BY status`，返回 map
- **TerminalByStore**: `SELECT store_id, status, COUNT(*) FROM terminals JOIN stores ... WHERE ... GROUP BY store_id, status` — service 层将行转为 `[]TerminalStoreStat`
- **TerminalRecentAdded**: `SELECT COUNT(*) FROM terminals WHERE created_at >= ?`，分别统计 7 天和 30 天
- **UserStatusDistribution**: `SELECT status, COUNT(*) FROM users GROUP BY status`，service 层将 int16 status 转为字符串 key
- **UserByRole**: `SELECT role_id, COUNT(*) FROM user_roles JOIN roles ... GROUP BY role_id`
- **UserRecentAdded**: `SELECT COUNT(*) FROM users WHERE created_at >= ?`

## DTO

```go
// TerminalStatsResp 终端统计响应
type TerminalStatsResp struct {
    StatusDistribution map[string]int64        `json:"status_distribution"`
    ByStore            []TerminalStoreStatItem `json:"by_store"`
    RecentAdded        RecentAddedStat         `json:"recent_added"`
}

type TerminalStoreStatItem struct {
    StoreID   int64  `json:"store_id"`
    StoreName string `json:"store_name"`
    Total     int64  `json:"total"`
    Online    int64  `json:"online"`
    Offline   int64  `json:"offline"`
    Disabled  int64  `json:"disabled"`
}

type RecentAddedStat struct {
    Last7Days  int64 `json:"last_7_days"`
    Last30Days int64 `json:"last_30_days"`
}

// UserStatsResp 用户统计响应
type UserStatsResp struct {
    StatusDistribution map[string]int64    `json:"status_distribution"`
    ByRole             []UserRoleStatItem `json:"by_role"`
    RecentAdded        RecentAddedStat    `json:"recent_added"`
}

type UserRoleStatItem struct {
    RoleID   int64  `json:"role_id"`
    RoleName string `json:"role_name"`
    Count    int64  `json:"count"`
}
```

## 权限

SUPER_ADMIN 绕过 Casbin，始终有访问权限。其他角色需要通过现有的角色管理 API 添加策略（资源路径为 `/api/v1/stats/terminals` 和 `/api/v1/stats/users`）。

统计数据接口的 Casbin 策略不会在 seed 阶段自动创建，与其他 API 权限的管理方式保持一致。

## 错误处理

- 数据库异常返回 `50001`（服务器内部错误）
- 空数据（如某门店下无终端）返回零值计数，不返回 null 或 404
- 无权限访问由 Permission 中间件统一返回 `40301`

## 测试

| 层级 | 测试内容 |
|------|----------|
| Repository | Mock GORM 或 SQLite 内存数据库；验证聚合查询返回正确的计数 |
| Service | Mock repository；验证数据转换逻辑（如 by_store 行转列） |
| Handler | Mock service；验证 HTTP 200、响应结构、DataScope 上下文读取 |

## 范围约束

- 仅实现终端统计和用户统计两个接口，门店和日志统计不在当前范围
- 暂不支持自定义时间范围过滤，recent_added 固定为当前时间的 7 天 / 30 天
- 不引入缓存，每次请求直接查询数据库
- 终端统计的 DataScope 从 gin context 读取 storeID 列表（由已有 DataScope 中间件注入）
