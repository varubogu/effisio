-- トリガーの削除
DROP TRIGGER IF EXISTS update_permissions_updated_at ON permissions;
DROP TRIGGER IF EXISTS update_roles_updated_at ON roles;

-- トリガー関数の削除
DROP FUNCTION IF EXISTS update_updated_at_column();

-- テーブルの削除（外部キー制約のため逆順）
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
