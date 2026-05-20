# Dashboard Statistics API

## Overview

Add two read-only statistics endpoints under `/api/v1/stats/*` that return aggregated dashboard data for terminals and users.

- `GET /api/v1/stats/terminals` — terminal status distribution, per-store breakdown, recent additions
- `GET /api/v1/stats/users` — user status distribution, per-role breakdown, recent additions

## Routes

```
GET /api/v1/stats/terminals
GET /api/v1/stats/users
```

Both routes sit in the existing `protected` middleware group (Auth → PasswordExpiry → Permission → DataScope). No Audit middleware (read-only endpoints).

Terminal stats respect DataScope (store-level isolation for non-SUPER_ADMIN users). User stats are global — no store-level scoping applies.

## Architecture

Follows the existing three-layer pattern:

```
handler → service → repository → GORM
```

### New files

| File | Purpose |
|------|---------|
| `internal/repository/stats.go` | `StatsRepository` interface + GORM implementation |
| `internal/service/stats.go` | `StatsService` — composes repo calls, accepts data-scope params |
| `internal/handler/stats.go` | `StatsHandler` — two handler methods |

### Modified files

| File | Change |
|------|--------|
| `internal/dto/response.go` | Add `TerminalStatsResp` and `UserStatsResp` DTOs |
| `internal/router/router.go` | Register `/stats/terminals` and `/stats/users` routes; add `StatsHandler` to `Dependencies` |
| `main.go` | Wire `StatsRepository` → `StatsService` → `StatsHandler` |

## Data Model

### Request

No request body. No query params (for now).

### Response: `GET /api/v1/stats/terminals`

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

`by_store` is scoped by DataScope: non-admin users only see stores they are assigned to. `status_distribution` and `recent_added` share the same scope.

### Response: `GET /api/v1/stats/users`

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

User stats are not store-scoped — role distribution and status distribution are global views.

## Repository

### StatsRepository interface

```go
type StatsRepository interface {
    // Terminal stats (storeIDs filters by data scope; empty = all)
    TerminalStatusDistribution(ctx context.Context, storeIDs []int64) (map[string]int64, error)
    TerminalByStore(ctx context.Context, storeIDs []int64) ([]TerminalStoreStat, error)
    TerminalRecentAdded(ctx context.Context, storeIDs []int64) (last7Days int64, last30Days int64, err error)

    // User stats (global, no data scope)
    UserStatusDistribution(ctx context.Context) (map[int16]int64, error)
    UserByRole(ctx context.Context) ([]UserRoleStat, error)
    UserRecentAdded(ctx context.Context) (last7Days int64, last30Days int64, err error)
}
```

### Queries (GORM)

- **TerminalStatusDistribution**: `SELECT status, COUNT(*) FROM terminals WHERE ... GROUP BY status`
- **TerminalByStore**: `SELECT store_id, status, COUNT(*) FROM terminals JOIN stores ... WHERE ... GROUP BY store_id, status` — service layer pivots rows into `[]TerminalStoreStat`
- **TerminalRecentAdded**: `SELECT COUNT(*) FROM terminals WHERE created_at >= ?`
- **UserStatusDistribution**: `SELECT status, COUNT(*) FROM users GROUP BY status`
- **UserByRole**: `SELECT role_id, COUNT(*) FROM user_roles JOIN roles ... GROUP BY role_id`
- **UserRecentAdded**: `SELECT COUNT(*) FROM users WHERE created_at >= ?`

## DTOs

```go
// TerminalStatsResp is the top-level response for terminal stats.
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

// UserStatsResp is the top-level response for user stats.
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

## Permissions

SUPER_ADMIN bypasses Casbin and always has access. For other roles, policies must be added via the existing role API (`PUT /api/v1/roles/:id/menus` or direct Casbin policy insertion) with the resource paths `/api/v1/stats/terminals` and `/api/v1/stats/users`.

The seed/default setup does not auto-create Casbin policies for stats endpoints — this is consistent with how all other API permissions are managed.

## Error Handling

- Database errors return `50001` (internal server error)
- Empty data (e.g., no terminals in a store) returns zero counts, not null / 404
- Permission denied from Casbin returns `40301` (standard middleware behavior)

## Testing

| Layer | What to test |
|-------|-------------|
| Repository | Mock GORM or use SQLite in-memory; verify aggregate queries return correct counts |
| Service | Mock repository; verify data transformation (e.g., row pivoting for by_store) |
| Handler | Mock service; verify HTTP 200, response structure, data-scope context reading |

## Scope considerations

- Only two endpoints (terminals + users). Stores and logs stats are out of scope.
- No time-range filtering (yet) — recent_added is always 7d/30d from now.
- No caching in this iteration — queries hit the database directly.
- DataScope for terminal stats reads store IDs from gin context (set by existing DataScope middleware).
