-- organizations テーブル: 階層的な組織構造を管理
CREATE TABLE IF NOT EXISTS organizations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE,
    parent_id INTEGER REFERENCES organizations(id) ON DELETE CASCADE,
    path VARCHAR(1000),
    level INTEGER NOT NULL DEFAULT 0,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- インデックス作成
CREATE INDEX idx_organizations_parent_id ON organizations(parent_id);
CREATE INDEX idx_organizations_path ON organizations(path);
CREATE INDEX idx_organizations_code ON organizations(code);
CREATE INDEX idx_organizations_level ON organizations(level);

-- パスを自動生成・更新するトリガー関数
CREATE OR REPLACE FUNCTION update_organization_path()
RETURNS TRIGGER AS $$
DECLARE
    parent_path VARCHAR(1000);
BEGIN
    -- 親が存在する場合
    IF NEW.parent_id IS NOT NULL THEN
        SELECT path, level INTO parent_path, NEW.level
        FROM organizations
        WHERE id = NEW.parent_id;

        NEW.path := parent_path || '/' || NEW.id::TEXT;
        NEW.level := NEW.level + 1;
    ELSE
        -- ルート組織の場合
        NEW.path := '/' || NEW.id::TEXT;
        NEW.level := 0;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- パス更新トリガー
CREATE TRIGGER trigger_update_organization_path
    BEFORE INSERT OR UPDATE OF parent_id ON organizations
    FOR EACH ROW
    EXECUTE FUNCTION update_organization_path();

-- updated_at を自動更新するトリガー（既存の関数を使用）
CREATE TRIGGER update_organizations_updated_at
    BEFORE UPDATE ON organizations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 初期データの投入（サンプル組織構造）
INSERT INTO organizations (name, code, description) VALUES
    ('本社', 'HQ', '本社組織');

-- 本社のIDを取得して部門を作成
DO $$
DECLARE
    hq_id INTEGER;
    dev_id INTEGER;
    sales_id INTEGER;
BEGIN
    SELECT id INTO hq_id FROM organizations WHERE code = 'HQ';

    INSERT INTO organizations (name, code, parent_id, description) VALUES
        ('開発部', 'DEV', hq_id, '開発部門'),
        ('営業部', 'SALES', hq_id, '営業部門'),
        ('管理部', 'ADMIN', hq_id, '管理部門');

    SELECT id INTO dev_id FROM organizations WHERE code = 'DEV';
    SELECT id INTO sales_id FROM organizations WHERE code = 'SALES';

    -- 開発部の下に課を作成
    INSERT INTO organizations (name, code, parent_id, description) VALUES
        ('フロントエンド課', 'DEV-FE', dev_id, 'フロントエンド開発'),
        ('バックエンド課', 'DEV-BE', dev_id, 'バックエンド開発');

    -- 営業部の下に課を作成
    INSERT INTO organizations (name, code, parent_id, description) VALUES
        ('国内営業課', 'SALES-DOM', sales_id, '国内営業'),
        ('海外営業課', 'SALES-INT', sales_id, '海外営業');
END $$;

COMMENT ON TABLE organizations IS '組織階層テーブル。自己参照外部キーで階層構造を表現';
COMMENT ON COLUMN organizations.path IS 'パス（例: /1/3/5）。検索高速化のため自動生成';
COMMENT ON COLUMN organizations.level IS '階層レベル（0=ルート）';
