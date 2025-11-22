-- トリガーの削除
DROP TRIGGER IF EXISTS update_organizations_updated_at ON organizations;
DROP TRIGGER IF EXISTS trigger_update_organization_path ON organizations;

-- トリガー関数の削除
DROP FUNCTION IF EXISTS update_organization_path();

-- テーブルの削除
DROP TABLE IF EXISTS organizations;
