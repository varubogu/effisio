-- タスクテーブルの作成
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'TODO',
    priority VARCHAR(10) NOT NULL DEFAULT 'MEDIUM',
    assigned_to INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_by INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id INTEGER REFERENCES organizations(id) ON DELETE SET NULL,
    due_date TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- インデックスの作成
CREATE INDEX idx_tasks_assigned_to ON tasks(assigned_to);
CREATE INDEX idx_tasks_created_by ON tasks(created_by);
CREATE INDEX idx_tasks_organization_id ON tasks(organization_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);
CREATE INDEX idx_tasks_created_at ON tasks(created_at);

-- ステータスと優先度の制約
ALTER TABLE tasks ADD CONSTRAINT check_task_status
    CHECK (status IN ('TODO', 'IN_PROGRESS', 'IN_REVIEW', 'DONE', 'CANCELLED'));

ALTER TABLE tasks ADD CONSTRAINT check_task_priority
    CHECK (priority IN ('LOW', 'MEDIUM', 'HIGH', 'URGENT'));

-- updated_at自動更新トリガー
CREATE OR REPLACE FUNCTION update_tasks_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_tasks_updated_at
    BEFORE UPDATE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_tasks_updated_at();

-- ステータスがDONEまたはCANCELLEDになった時にcompleted_atを自動設定
CREATE OR REPLACE FUNCTION update_tasks_completed_at()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status IN ('DONE', 'CANCELLED') AND OLD.status NOT IN ('DONE', 'CANCELLED') THEN
        NEW.completed_at = CURRENT_TIMESTAMP;
    ELSIF NEW.status NOT IN ('DONE', 'CANCELLED') THEN
        NEW.completed_at = NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_tasks_completed_at
    BEFORE UPDATE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_tasks_completed_at();

-- コメント
COMMENT ON TABLE tasks IS 'タスク管理テーブル';
COMMENT ON COLUMN tasks.status IS 'タスクステータス: TODO, IN_PROGRESS, IN_REVIEW, DONE, CANCELLED';
COMMENT ON COLUMN tasks.priority IS 'タスク優先度: LOW, MEDIUM, HIGH, URGENT';
COMMENT ON COLUMN tasks.assigned_to IS '担当者（ユーザーID）';
COMMENT ON COLUMN tasks.created_by IS '作成者（ユーザーID）';
COMMENT ON COLUMN tasks.organization_id IS '所属組織ID';
COMMENT ON COLUMN tasks.completed_at IS '完了日時（ステータスがDONE/CANCELLEDの時に自動設定）';
