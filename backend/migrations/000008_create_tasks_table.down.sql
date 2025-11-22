-- トリガーの削除
DROP TRIGGER IF EXISTS trigger_update_tasks_completed_at ON tasks;
DROP TRIGGER IF EXISTS trigger_update_tasks_updated_at ON tasks;

-- トリガー関数の削除
DROP FUNCTION IF EXISTS update_tasks_completed_at();
DROP FUNCTION IF EXISTS update_tasks_updated_at();

-- テーブルの削除
DROP TABLE IF EXISTS tasks;
