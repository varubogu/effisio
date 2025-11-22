-- 新しいカラムを追加
ALTER TABLE users ADD COLUMN full_name VARCHAR(100);
ALTER TABLE users ADD COLUMN department VARCHAR(100);
ALTER TABLE users ADD COLUMN last_login TIMESTAMP;

-- パスワードカラムをリネーム
ALTER TABLE users RENAME COLUMN password TO password_hash;

-- ステータスカラムを追加（一時的にNULL許可）
ALTER TABLE users ADD COLUMN status VARCHAR(20);

-- 既存データのステータスを設定（is_activeからstatusへ変換）
UPDATE users SET status = CASE
    WHEN is_active = true THEN 'active'
    WHEN is_active = false THEN 'inactive'
    ELSE 'active'
END;

-- statusをNOT NULLに変更してデフォルト値を設定
ALTER TABLE users ALTER COLUMN status SET NOT NULL;
ALTER TABLE users ALTER COLUMN status SET DEFAULT 'active';

-- is_activeカラムを削除
ALTER TABLE users DROP COLUMN is_active;

-- statusにCHECK制約を追加
ALTER TABLE users ADD CONSTRAINT users_status_check
    CHECK (status IN ('active', 'inactive', 'suspended'));

-- インデックスを追加（パフォーマンス向上）
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_department ON users(department);
