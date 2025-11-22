-- roles テーブル: システムのロールを定義
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- permissions テーブル: システムの権限を定義
CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- role_permissions テーブル: ロールと権限の多対多の関連
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, permission_id)
);

-- インデックス作成
CREATE INDEX idx_permissions_resource ON permissions(resource);
CREATE INDEX idx_permissions_action ON permissions(action);
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);

-- 初期ロールの投入
INSERT INTO roles (name, description) VALUES
    ('admin', 'システム管理者 - 全ての操作が可能'),
    ('manager', 'マネージャー - ユーザーとタスクの管理が可能'),
    ('user', '一般ユーザー - タスクの作成・編集が可能'),
    ('viewer', '閲覧者 - 読み取り専用');

-- 初期権限の投入
INSERT INTO permissions (name, resource, action, description) VALUES
    -- ユーザー管理権限
    ('users:read', 'users', 'read', 'ユーザー情報の閲覧'),
    ('users:write', 'users', 'write', 'ユーザー情報の作成・編集'),
    ('users:delete', 'users', 'delete', 'ユーザーの削除'),

    -- タスク管理権限
    ('tasks:read', 'tasks', 'read', 'タスクの閲覧'),
    ('tasks:write', 'tasks', 'write', 'タスクの作成・編集'),
    ('tasks:delete', 'tasks', 'delete', 'タスクの削除'),

    -- 設定権限
    ('settings:read', 'settings', 'read', 'システム設定の閲覧'),
    ('settings:write', 'settings', 'write', 'システム設定の変更'),

    -- ロール・権限管理権限
    ('roles:read', 'roles', 'read', 'ロール情報の閲覧'),
    ('roles:write', 'roles', 'write', 'ロール情報の作成・編集'),
    ('permissions:read', 'permissions', 'read', '権限情報の閲覧'),
    ('permissions:write', 'permissions', 'write', '権限情報の作成・編集'),

    -- 監査ログ権限
    ('audit:read', 'audit', 'read', '監査ログの閲覧');

-- ロールと権限の関連付け
-- admin: 全ての権限
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin';

-- manager: ユーザー閲覧、タスク管理、監査ログ閲覧
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'manager'
  AND p.name IN (
    'users:read',
    'tasks:read',
    'tasks:write',
    'tasks:delete',
    'audit:read'
  );

-- user: タスクの読み書き
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'user'
  AND p.name IN (
    'tasks:read',
    'tasks:write'
  );

-- viewer: 読み取り専用
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'viewer'
  AND p.name IN (
    'tasks:read'
  );

-- updated_at を自動更新するトリガー関数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- roles テーブルのトリガー
CREATE TRIGGER update_roles_updated_at
    BEFORE UPDATE ON roles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- permissions テーブルのトリガー
CREATE TRIGGER update_permissions_updated_at
    BEFORE UPDATE ON permissions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
