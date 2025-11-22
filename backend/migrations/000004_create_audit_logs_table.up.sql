-- Create audit_logs table
BEGIN;

CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    action VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id VARCHAR(50) NOT NULL,
    changes JSONB NOT NULL DEFAULT '{}',
    ip_address VARCHAR(45),
    user_agent TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'success',
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- インデックスの作成
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_resource_type ON audit_logs(resource_type);
CREATE INDEX idx_audit_logs_resource_id ON audit_logs(resource_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_status ON audit_logs(status);
CREATE INDEX idx_audit_logs_user_id_created_at ON audit_logs(user_id, created_at DESC);

-- コメントの追加
COMMENT ON TABLE audit_logs IS '監査ログを記録するテーブル';
COMMENT ON COLUMN audit_logs.id IS '監査ログID';
COMMENT ON COLUMN audit_logs.user_id IS 'アクションを実行したユーザーID';
COMMENT ON COLUMN audit_logs.action IS 'アクション（create, read, update, delete, login, logout）';
COMMENT ON COLUMN audit_logs.resource_type IS 'リソースタイプ（user, role, organization等）';
COMMENT ON COLUMN audit_logs.resource_id IS 'リソースID';
COMMENT ON COLUMN audit_logs.changes IS 'JSONB形式の変更内容（変更前後の値を記録）';
COMMENT ON COLUMN audit_logs.ip_address IS 'クライアントIPアドレス';
COMMENT ON COLUMN audit_logs.user_agent IS 'ユーザーエージェント';
COMMENT ON COLUMN audit_logs.status IS 'ステータス（success, failed）';
COMMENT ON COLUMN audit_logs.error_message IS 'エラー発生時のメッセージ';
COMMENT ON COLUMN audit_logs.created_at IS '作成日時';

COMMIT;
