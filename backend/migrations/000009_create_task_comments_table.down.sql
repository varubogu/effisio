-- トリガーの削除
DROP TRIGGER IF EXISTS trigger_update_task_comments_updated_at ON task_comments;

-- トリガー関数の削除
DROP FUNCTION IF EXISTS update_task_comments_updated_at();

-- テーブルの削除
DROP TABLE IF EXISTS task_comments;
