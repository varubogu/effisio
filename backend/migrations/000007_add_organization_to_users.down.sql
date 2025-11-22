-- インデックスの削除
DROP INDEX IF EXISTS idx_users_organization_id;

-- カラムの削除
ALTER TABLE users DROP COLUMN IF EXISTS organization_id;
