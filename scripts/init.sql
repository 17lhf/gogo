-- ============================================================================
-- gogo PostgreSQL Database Initialization Script
-- Database: gogo_dev
-- User:     gogo
-- Password: gogo123
-- ============================================================================

-- ----------------------------------------------------------------------------
-- 1. Create database (run as superuser / postgres)
-- ----------------------------------------------------------------------------
-- CREATE DATABASE gogo_dev
--     WITH ENCODING 'UTF8'
--     LC_COLLATE = 'en_US.UTF-8'
--     LC_CTYPE = 'en_US.UTF-8'
--     TEMPLATE template0;
--
-- COMMENT ON DATABASE gogo_dev IS 'gogo system database';

-- ----------------------------------------------------------------------------
-- 2. Tables
-- ----------------------------------------------------------------------------

BEGIN;

-- 2.1 users
CREATE TABLE IF NOT EXISTS users (
    id                  BIGSERIAL       PRIMARY KEY,
    username            VARCHAR(64)     NOT NULL,
    email               VARCHAR(128)    NOT NULL,
    password            VARCHAR(256)    NOT NULL,
    real_name           VARCHAR(64)     NOT NULL DEFAULT '',
    phone               VARCHAR(20)     NOT NULL DEFAULT '',
    status              SMALLINT        NOT NULL DEFAULT 1,   -- 1=enabled, 2=disabled, 3=locked
    must_change_password BOOLEAN        NOT NULL DEFAULT FALSE,
    password_updated_at TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    last_login_at       TIMESTAMPTZ,
    created_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users (username);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email    ON users (email);
COMMENT ON TABLE users IS 'System users';
COMMENT ON COLUMN users.status IS '1=enabled, 2=disabled, 3=locked';

-- 2.2 roles
CREATE TABLE IF NOT EXISTS roles (
    id          BIGSERIAL       PRIMARY KEY,
    name        VARCHAR(64)     NOT NULL,
    code        VARCHAR(64)     NOT NULL,
    description VARCHAR(256)    NOT NULL DEFAULT '',
    status      SMALLINT        NOT NULL DEFAULT 1,   -- 1=enabled, 2=disabled
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_code ON roles (code);
COMMENT ON TABLE roles IS 'RBAC roles';

-- 2.3 menus
CREATE TABLE IF NOT EXISTS menus (
    id          BIGSERIAL       PRIMARY KEY,
    parent_id   BIGINT          NOT NULL DEFAULT 0,
    name        VARCHAR(64)     NOT NULL,
    path        VARCHAR(256)    NOT NULL DEFAULT '',
    component   VARCHAR(256)    NOT NULL DEFAULT '',
    icon        VARCHAR(64)     NOT NULL DEFAULT '',
    type        SMALLINT        NOT NULL,              -- 1=directory, 2=page, 3=button
    perms       VARCHAR(128)    NOT NULL DEFAULT '',
    api_path    VARCHAR(256)    NOT NULL DEFAULT '',
    api_method  VARCHAR(10)     NOT NULL DEFAULT '',
    sort_order  INT             NOT NULL DEFAULT 0,
    visible     BOOLEAN         NOT NULL DEFAULT TRUE,
    status      SMALLINT        NOT NULL DEFAULT 1,    -- 1=enabled, 2=disabled
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_menus_parent_id ON menus (parent_id);
COMMENT ON TABLE menus IS 'Menu items (directory, page, or button)';
COMMENT ON COLUMN menus.type IS '1=directory, 2=page, 3=button';
COMMENT ON COLUMN menus.perms IS 'Permission identifier, e.g. sys:user:add';
COMMENT ON COLUMN menus.api_path IS 'API path for Casbin enforcement (button type), e.g. /api/v1/stores';
COMMENT ON COLUMN menus.api_method IS 'HTTP method for Casbin enforcement (button type), e.g. GET';

-- 2.4 stores
CREATE TABLE IF NOT EXISTS stores (
    id          BIGSERIAL       PRIMARY KEY,
    name        VARCHAR(128)    NOT NULL,
    code        VARCHAR(64)     NOT NULL,
    address     VARCHAR(256)    NOT NULL DEFAULT '',
    status      SMALLINT        NOT NULL DEFAULT 1,    -- 1=enabled, 2=disabled
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_stores_code ON stores (code);
COMMENT ON TABLE stores IS 'Physical store locations';

-- 2.5 terminals
CREATE TABLE IF NOT EXISTS terminals (
    id                BIGSERIAL       PRIMARY KEY,
    sn                VARCHAR(64)     NOT NULL,
    name              VARCHAR(128)    NOT NULL,
    type              VARCHAR(64)     NOT NULL DEFAULT '',
    store_id          BIGINT,
    status            VARCHAR(16)     NOT NULL DEFAULT 'offline',  -- offline, online, disabled, enabled
    ip_address        VARCHAR(45)     NOT NULL DEFAULT '',          -- supports IPv6
    mac_address       VARCHAR(17)     NOT NULL DEFAULT '',
    device_token      VARCHAR(256)    NOT NULL,
    last_heartbeat_at TIMESTAMPTZ,
    created_at        TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_terminals_sn       ON terminals (sn);
CREATE INDEX        IF NOT EXISTS idx_terminals_store_id ON terminals (store_id);
COMMENT ON TABLE terminals IS 'Edge devices registered in the system';
COMMENT ON COLUMN terminals.status IS 'offline, online, disabled, enabled';

-- 2.6 operation_logs
CREATE TABLE IF NOT EXISTS operation_logs (
    id            BIGSERIAL       PRIMARY KEY,
    user_id       BIGINT,
    username      VARCHAR(64)     NOT NULL DEFAULT '',
    action        VARCHAR(128)    NOT NULL,
    resource_type VARCHAR(64)     NOT NULL DEFAULT '',
    resource_id   VARCHAR(64)     NOT NULL DEFAULT '',
    detail        JSONB,
    ip            VARCHAR(45)     NOT NULL DEFAULT '',
    user_agent    VARCHAR(512)    NOT NULL DEFAULT '',
    status        SMALLINT        NOT NULL,              -- 1=success, 2=failure
    duration_ms   INT             NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_operation_logs_created_at ON operation_logs (created_at);
COMMENT ON TABLE operation_logs IS 'Admin operation audit logs';
COMMENT ON COLUMN operation_logs.status IS '1=success, 2=failure';

-- 2.7 terminal_logs
CREATE TABLE IF NOT EXISTS terminal_logs (
    id          BIGSERIAL       PRIMARY KEY,
    terminal_id BIGINT,
    sn          VARCHAR(64)     NOT NULL,
    event_type  VARCHAR(32)     NOT NULL,
    detail      JSONB,
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_terminal_logs_created_at ON terminal_logs (created_at);
COMMENT ON TABLE terminal_logs IS 'Terminal event logs (online/offline/heartbeat_timeout/etc.)';

-- 2.8 user_roles (many-to-many junction: users <-> roles)
CREATE TABLE IF NOT EXISTS user_roles (
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    PRIMARY KEY (user_id, role_id)
);
COMMENT ON TABLE user_roles IS 'Many-to-many: users <-> roles';

-- 2.9 user_stores (many-to-many junction: users <-> stores)
CREATE TABLE IF NOT EXISTS user_stores (
    user_id  BIGINT NOT NULL,
    store_id BIGINT NOT NULL,
    PRIMARY KEY (user_id, store_id)
);
COMMENT ON TABLE user_stores IS 'Many-to-many: users <-> stores';

-- 2.10 role_menus (many-to-many junction: roles <-> menus)
CREATE TABLE IF NOT EXISTS role_menus (
    role_id BIGINT NOT NULL,
    menu_id BIGINT NOT NULL,
    PRIMARY KEY (role_id, menu_id)
);
COMMENT ON TABLE role_menus IS 'Many-to-many: roles <-> menus';

-- 2.11 casbin_rule (RBAC policy rules managed by casbin gorm-adapter v3)
CREATE TABLE IF NOT EXISTS casbin_rule (
    id    BIGSERIAL    PRIMARY KEY,
    ptype VARCHAR(100) NOT NULL DEFAULT '',
    v0    VARCHAR(100) NOT NULL DEFAULT '',
    v1    VARCHAR(100) NOT NULL DEFAULT '',
    v2    VARCHAR(100) NOT NULL DEFAULT '',
    v3    VARCHAR(100) NOT NULL DEFAULT '',
    v4    VARCHAR(100) NOT NULL DEFAULT '',
    v5    VARCHAR(100) NOT NULL DEFAULT ''
);
COMMENT ON TABLE casbin_rule IS 'Casbin RBAC policy rules';

COMMIT;

-- ============================================================================
-- 3. Initial data (idempotent — safe to run multiple times)
-- ============================================================================

BEGIN;

-- 3.1 SUPER_ADMIN role
INSERT INTO roles (name, code, description, status, created_at, updated_at)
SELECT '超级管理员', 'SUPER_ADMIN', '系统内置超级管理员，拥有所有权限', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM roles WHERE code = 'SUPER_ADMIN');

-- 3.2 Admin user (password: Admin123!, bcrypt cost=12)
INSERT INTO users (username, email, password, real_name, status, created_at, updated_at)
SELECT 'admin', 'admin@system.local',
       '$2a$12$K/y2IE9Z3PXfa9R2TwTzzOMSIYm9ySmYc.eVFMbu1PlGyv2IPY142',
       '系统管理员', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'admin');

-- 3.3 user_roles: admin -> SUPER_ADMIN
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u CROSS JOIN roles r
WHERE u.username = 'admin' AND r.code = 'SUPER_ADMIN'
  AND NOT EXISTS (SELECT 1 FROM user_roles ur WHERE ur.user_id = u.id AND ur.role_id = r.id);

-- 3.4 Menu tree (directory + pages + button permissions)
-- 3.4.1 Directory: 系统管理
INSERT INTO menus (parent_id, name, path, component, icon, type, sort_order, visible, status, created_at, updated_at)
SELECT 0, '系统管理', '/system', 'Layout', 'system', 1, 1, true, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE name = '系统管理' AND parent_id = 0 AND type = 1);

-- 3.4.2 Pages (children of 系统管理)
INSERT INTO menus (parent_id, name, path, component, icon, type, sort_order, visible, status, created_at, updated_at)
SELECT d.id, p.name, p.path, p.component, p.icon, 2, p.sort_order, true, 1, NOW(), NOW()
FROM (SELECT id FROM menus WHERE name = '系统管理' AND parent_id = 0 AND type = 1) d
CROSS JOIN (VALUES
    ('用户管理', '/system/user',    'system/user/index',     'user',     1),
    ('角色管理', '/system/role',    'system/role/index',     'role',     2),
    ('菜单管理', '/system/menu',    'system/menu/index',     'menu',     3),
    ('门店管理', '/system/store',   'system/store/index',    'store',    4),
    ('终端管理', '/system/terminal','system/terminal/index', 'terminal', 5),
    ('日志管理', '/system/log',     'system/log/index',      'log',      6)
) AS p(name, path, component, icon, sort_order)
WHERE NOT EXISTS (
    SELECT 1 FROM menus c
    WHERE c.parent_id = d.id AND c.name = p.name AND c.type = 2
);

-- 3.4.3 Button permissions (children of their respective page menus)
-- Each button stores its api_path and api_method; Casbin policies are derived from these.
INSERT INTO menus (parent_id, name, perms, api_path, api_method, type, sort_order, visible, status, created_at, updated_at)
SELECT pg.id, b.name, b.perms, b.api_path, b.api_method, 3, b.sort_order, true, 1, NOW(), NOW()
FROM (VALUES
    ('用户管理', '创建用户',   'sys:user:add',          '/api/v1/users',              'POST',   1),
    ('用户管理', '编辑用户',   'sys:user:edit',         '/api/v1/users/:id',          'PUT',    2),
    ('用户管理', '删除用户',   'sys:user:delete',       '/api/v1/users/:id',          'DELETE', 3),
    ('用户管理', '重置密码',   'sys:user:reset-pwd',    '/api/v1/users/:id/password', 'PUT',    4),
    ('用户管理', '分配角色',   'sys:user:assign-role',  '/api/v1/users/:id/roles',    'PUT',    5),
    ('用户管理', '分配门店',   'sys:user:assign-store', '/api/v1/users/:id/stores',   'PUT',    6),
    ('用户管理', '用户列表',   'sys:user:list',         '/api/v1/users',              'GET',    7),
    ('角色管理', '创建角色',   'sys:role:add',          '/api/v1/roles',              'POST',   1),
    ('角色管理', '编辑角色',   'sys:role:edit',         '/api/v1/roles/:id',          'PUT',    2),
    ('角色管理', '删除角色',   'sys:role:delete',       '/api/v1/roles/:id',          'DELETE', 3),
    ('角色管理', '分配菜单',   'sys:role:assign-menu',  '/api/v1/roles/:id/menus',    'PUT',    4),
    ('角色管理', '角色列表',   'sys:role:list',         '/api/v1/roles',              'GET',    5),
    ('菜单管理', '创建菜单',   'sys:menu:add',          '/api/v1/menus',              'POST',   1),
    ('菜单管理', '编辑菜单',   'sys:menu:edit',         '/api/v1/menus/:id',          'PUT',    2),
    ('菜单管理', '删除菜单',   'sys:menu:delete',       '/api/v1/menus/:id',          'DELETE', 3),
    ('菜单管理', '菜单列表',   'sys:menu:list',         '/api/v1/menus',              'GET',    4),
    ('门店管理', '创建门店',   'sys:store:add',         '/api/v1/stores',             'POST',   1),
    ('门店管理', '编辑门店',   'sys:store:edit',        '/api/v1/stores/:id',         'PUT',    2),
    ('门店管理', '删除门店',   'sys:store:delete',      '/api/v1/stores/:id',         'DELETE', 3),
    ('门店管理', '门店列表',   'sys:store:list',        '/api/v1/stores',             'GET',    4),
    ('终端管理', '创建终端',   'sys:terminal:add',      '/api/v1/terminals',          'POST',   1),
    ('终端管理', '编辑终端',   'sys:terminal:edit',     '/api/v1/terminals/:id',      'PUT',    2),
    ('终端管理', '删除终端',   'sys:terminal:delete',   '/api/v1/terminals/:id',      'DELETE', 3),
    ('终端管理', '终端列表',   'sys:terminal:list',     '/api/v1/terminals',          'GET',    4),
    ('日志管理', '操作日志列表', 'sys:log:list',         '/api/v1/logs/operations',    'GET',    1),
    ('日志管理', '终端日志列表', 'sys:log:list',         '/api/v1/logs/terminals',     'GET',    2)
) AS b(page_name, name, perms, api_path, api_method, sort_order)
JOIN menus pg ON pg.name = b.page_name AND pg.type = 2
    AND pg.parent_id = (SELECT id FROM menus WHERE name = '系统管理' AND parent_id = 0 AND type = 1)
WHERE NOT EXISTS (
    SELECT 1 FROM menus c
    WHERE c.parent_id = pg.id AND c.name = b.name AND c.type = 3
);

-- 3.5 role_menus: SUPER_ADMIN -> all menus
INSERT INTO role_menus (role_id, menu_id)
SELECT r.id, m.id
FROM roles r CROSS JOIN menus m
WHERE r.code = 'SUPER_ADMIN'
  AND NOT EXISTS (SELECT 1 FROM role_menus rm WHERE rm.role_id = r.id AND rm.menu_id = m.id);

-- 3.6 casbin_rule: SUPER_ADMIN policies derived from PermsToPolicies mapping
-- Each INSERT uses WHERE NOT EXISTS for idempotency.
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'SUPER_ADMIN', p.path, p.method
FROM (VALUES
    ('/api/v1/users',              'POST'),
    ('/api/v1/users/:id',          'PUT'),
    ('/api/v1/users/:id',          'DELETE'),
    ('/api/v1/users/:id/password', 'PUT'),
    ('/api/v1/users/:id/roles',    'PUT'),
    ('/api/v1/users/:id/stores',   'PUT'),
    ('/api/v1/users',              'GET'),
    ('/api/v1/roles',              'POST'),
    ('/api/v1/roles/:id',          'PUT'),
    ('/api/v1/roles/:id',          'DELETE'),
    ('/api/v1/roles/:id/menus',    'PUT'),
    ('/api/v1/roles',              'GET'),
    ('/api/v1/menus',              'POST'),
    ('/api/v1/menus/:id',          'PUT'),
    ('/api/v1/menus/:id',          'DELETE'),
    ('/api/v1/menus',              'GET'),
    ('/api/v1/stores',             'POST'),
    ('/api/v1/stores/:id',         'PUT'),
    ('/api/v1/stores/:id',         'DELETE'),
    ('/api/v1/stores',             'GET'),
    ('/api/v1/terminals',          'POST'),
    ('/api/v1/terminals/:id',      'PUT'),
    ('/api/v1/terminals/:id',      'DELETE'),
    ('/api/v1/terminals',          'GET'),
    ('/api/v1/logs/operations',    'GET'),
    ('/api/v1/logs/terminals',     'GET')
) AS p(path, method)
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule cr
    WHERE cr.ptype = 'p' AND cr.v0 = 'SUPER_ADMIN' AND cr.v1 = p.path AND cr.v2 = p.method
);

COMMIT;

-- ============================================================================
-- Usage:
--   psql -U postgres -f scripts/init.sql
-- Or run the CREATE DATABASE part first, then:
--   psql -U gogo -d gogo_dev -f scripts/init.sql
-- ============================================================================
