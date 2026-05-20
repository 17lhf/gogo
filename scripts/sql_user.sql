-- 1. 创建用户 gogo，密码为 gogo123
-- 注意：这里先不赋予 CREATEDB 等全局权限，只赋予 LOGIN 权限
CREATE USER gogo WITH PASSWORD 'gogo123' LOGIN;

-- 2. 赋予 gogo 用户对 gogo_dev 数据库的连接权限
GRANT CONNECT ON DATABASE gogo_dev TO gogo;

-- 3. 【关键步骤】切换到 gogo_dev 数据库
-- 在 pgAdmin 中，你可以点击工具栏的下拉菜单切换当前连接的数据库为 gogo_dev
-- 或者在同一个脚本窗口继续执行下面的命令（确保上下文已切换）

-- 4. 赋予 Schema 的使用权限 (默认通常是 public schema)
-- 这一步是必须的，否则用户连上库也看不到表
GRANT USAGE ON SCHEMA public TO gogo;

-- 5. 赋予该 Schema 下所有表的增删改查权限
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO gogo;

-- 6. 赋予序列(自增主键)的权限
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO gogo;

-- 7. 【可选但推荐】设置默认权限
-- 这样以后你在 gogo_dev 库里新建表时，gogo 用户自动拥有权限，不用每次都重新授权
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO gogo;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO gogo;

-- 将 public 模式的所有权转让给 app_user
ALTER SCHEMA public OWNER TO app_user;

-- 一次性把 public schema 下所有表的所有权转给 gogo
DO $$
  DECLARE
      r RECORD;
  BEGIN
      FOR r IN SELECT tablename FROM pg_tables WHERE schemaname = 'public'
      LOOP
          EXECUTE 'ALTER TABLE ' || quote_ident(r.tablename) || ' OWNER TO gogo';
      END LOOP;
  END $$;

  DO $$
  DECLARE
      r RECORD;
  BEGIN
      FOR r IN SELECT sequencename FROM pg_sequences WHERE schemaname = 'public'
      LOOP
          EXECUTE 'ALTER SEQUENCE ' || quote_ident(r.sequencename) || ' OWNER TO gogo';
      END LOOP;
  END $$;