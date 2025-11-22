-- audit_logs テーブル: システムの全ての重要操作を記録
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    resource VARCHAR(50) NOT NULL,
    resource_id VARCHAR(255),
    ip_address VARCHAR(45),
    user_agent TEXT,
    changes JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- インデックス作成（クエリ最適化）
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_resource_id ON audit_logs(resource, resource_id);

-- 複合インデックス（よく使われるクエリパターン用）
CREATE INDEX idx_audit_logs_user_action ON audit_logs(user_id, action);
CREATE INDEX idx_audit_logs_resource_action ON audit_logs(resource, action);

-- パーティショニング用のコメント（将来的な最適化のため）
COMMENT ON TABLE audit_logs IS '監査ログテーブル。将来的にパーティショニング検討';
COMMENT ON COLUMN audit_logs.changes IS 'JSONB型で変更前後の値を保存。例: {"before": {...}, "after": {...}}';
