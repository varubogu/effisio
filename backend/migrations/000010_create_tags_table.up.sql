-- タグマスタテーブルの作成
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    color VARCHAR(7) NOT NULL DEFAULT '#6B7280',
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- タスクタグ中間テーブルの作成（多対多リレーション）
CREATE TABLE IF NOT EXISTS task_tags (
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (task_id, tag_id)
);

-- インデックスの作成
CREATE INDEX idx_tags_name ON tags(name);
CREATE INDEX idx_task_tags_task_id ON task_tags(task_id);
CREATE INDEX idx_task_tags_tag_id ON task_tags(tag_id);

-- updated_at自動更新トリガー
CREATE OR REPLACE FUNCTION update_tags_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_tags_updated_at
    BEFORE UPDATE ON tags
    FOR EACH ROW
    EXECUTE FUNCTION update_tags_updated_at();

-- サンプルタグデータ
INSERT INTO tags (name, color, description) VALUES
    ('重要', '#EF4444', '重要度の高いタスク'),
    ('緊急', '#DC2626', '緊急対応が必要なタスク'),
    ('バグ', '#F59E0B', 'バグ修正タスク'),
    ('機能追加', '#3B82F6', '新機能追加タスク'),
    ('改善', '#10B981', '既存機能の改善タスク'),
    ('ドキュメント', '#8B5CF6', 'ドキュメント作成・更新タスク'),
    ('テスト', '#EC4899', 'テスト関連タスク'),
    ('レビュー', '#6366F1', 'レビュー待ちタスク')
ON CONFLICT (name) DO NOTHING;

-- コメント
COMMENT ON TABLE tags IS 'タグマスタテーブル';
COMMENT ON TABLE task_tags IS 'タスクタグ中間テーブル（多対多）';
COMMENT ON COLUMN tags.name IS 'タグ名（ユニーク）';
COMMENT ON COLUMN tags.color IS 'タグ色（16進数カラーコード）';
