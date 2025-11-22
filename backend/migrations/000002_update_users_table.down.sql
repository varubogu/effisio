-- ロールバック用（逆の操作）

-- インデックスを削除
DROP INDEX IF EXISTS idx_users_department;
DROP INDEX IF EXISTS idx_users_status;

-- CHECK制約を削除
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_status_check;

-- is_activeカラムを再追加
ALTER TABLE users ADD COLUMN is_active BOOLEAN DEFAULT true;

-- statusからis_activeへ変換
UPDATE users SET is_active = CASE
    WHEN status = 'active' THEN true
    ELSE false
END;

-- statusカラムを削除
ALTER TABLE users DROP COLUMN status;

-- パスワードカラムをリネーム
ALTER TABLE users RENAME COLUMN password_hash TO password;

-- 新しいカラムを削除
ALTER TABLE users DROP COLUMN last_login;
ALTER TABLE users DROP COLUMN department;
ALTER TABLE users DROP COLUMN full_name;
