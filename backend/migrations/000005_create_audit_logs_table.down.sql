-- インデックスの削除
DROP INDEX IF EXISTS idx_audit_logs_resource_action;
DROP INDEX IF EXISTS idx_audit_logs_user_action;
DROP INDEX IF EXISTS idx_audit_logs_resource_id;
DROP INDEX IF EXISTS idx_audit_logs_created_at;
DROP INDEX IF EXISTS idx_audit_logs_resource;
DROP INDEX IF EXISTS idx_audit_logs_action;
DROP INDEX IF EXISTS idx_audit_logs_user_id;

-- テーブルの削除
DROP TABLE IF EXISTS audit_logs;
