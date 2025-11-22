-- トリガーの削除
DROP TRIGGER IF EXISTS trigger_update_tags_updated_at ON tags;

-- トリガー関数の削除
DROP FUNCTION IF EXISTS update_tags_updated_at();

-- テーブルの削除（task_tagsを先に削除）
DROP TABLE IF EXISTS task_tags;
DROP TABLE IF EXISTS tags;
