-- users テーブルに organization_id カラムを追加
ALTER TABLE users ADD COLUMN organization_id INTEGER REFERENCES organizations(id) ON DELETE SET NULL;

-- インデックス作成
CREATE INDEX idx_users_organization_id ON users(organization_id);

COMMENT ON COLUMN users.organization_id IS 'ユーザーが所属する組織のID';
