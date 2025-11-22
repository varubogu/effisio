-- タスクアクティビティログテーブルの作成
CREATE TABLE IF NOT EXISTS task_activities (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    activity_type VARCHAR(50) NOT NULL,
    field_name VARCHAR(100),
    old_value TEXT,
    new_value TEXT,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- インデックスの作成
CREATE INDEX idx_task_activities_task_id ON task_activities(task_id);
CREATE INDEX idx_task_activities_user_id ON task_activities(user_id);
CREATE INDEX idx_task_activities_activity_type ON task_activities(activity_type);
CREATE INDEX idx_task_activities_created_at ON task_activities(created_at);

-- コメント
COMMENT ON TABLE task_activities IS 'タスクアクティビティログテーブル';
COMMENT ON COLUMN task_activities.activity_type IS 'アクティビティタイプ: created, updated, status_changed, assigned, commented, tag_added, tag_removed等';
COMMENT ON COLUMN task_activities.field_name IS '変更されたフィールド名';
COMMENT ON COLUMN task_activities.old_value IS '変更前の値';
COMMENT ON COLUMN task_activities.new_value IS '変更後の値';
COMMENT ON COLUMN task_activities.metadata IS '追加のメタデータ（JSONB）';
