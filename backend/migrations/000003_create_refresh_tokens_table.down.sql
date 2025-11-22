-- インデックスを削除
DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;
DROP INDEX IF EXISTS idx_refresh_tokens_token_id;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;

-- テーブルを削除
DROP TABLE IF EXISTS refresh_tokens;
