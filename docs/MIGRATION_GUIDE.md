# データベースマイグレーション運用ガイド

このガイドでは、Effisioプロジェクトにおけるデータベースマイグレーションの運用方法を説明します。

## マイグレーションツール

**golang-migrate** を使用します。

## インストール

```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Go install
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## ディレクトリ構造

```
backend/migrations/
├── 000001_initial_schema.up.sql
├── 000001_initial_schema.down.sql
├── 000002_add_users_table.up.sql
├── 000002_add_users_table.down.sql
├── 000003_add_roles_table.up.sql
└── 000003_add_roles_table.down.sql
```

## 命名規則

```
{version}_{description}.{up|down}.sql

例:
000001_initial_schema.up.sql
000001_initial_schema.down.sql
000002_add_users_table.up.sql
000002_add_users_table.down.sql
```

- **version**: 6桁の連番（`000001`, `000002`, ...）
- **description**: スネークケースで簡潔な説明
- **up**: マイグレーション適用用
- **down**: マイグレーションロールバック用

## マイグレーションファイルの作成

### 新規作成

```bash
cd backend/migrations

# 新しいマイグレーションファイルを作成
migrate create -ext sql -dir . -seq add_users_table
```

これにより以下のファイルが生成されます：

```
000001_add_users_table.up.sql
000001_add_users_table.down.sql
```

### upファイル（適用）の例

```sql
-- 000001_create_users_table.up.sql
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(255) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  full_name VARCHAR(255),
  department VARCHAR(255),
  role VARCHAR(50) NOT NULL DEFAULT 'user',
  status VARCHAR(50) NOT NULL DEFAULT 'active',
  last_login TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

### downファイル（ロールバック）の例

```sql
-- 000001_create_users_table.down.sql
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
```

## マイグレーションの実行

### 環境変数の設定

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/effisio_dev?sslmode=disable"
```

### マイグレーション適用

```bash
# すべてのマイグレーションを適用
migrate -path ./migrations -database "${DATABASE_URL}" up

# 1ステップだけ適用
migrate -path ./migrations -database "${DATABASE_URL}" up 1

# 特定のバージョンまで適用
migrate -path ./migrations -database "${DATABASE_URL}" goto 3
```

### マイグレーションロールバック

```bash
# 1ステップだけロールバック
migrate -path ./migrations -database "${DATABASE_URL}" down 1

# すべてロールバック
migrate -path ./migrations -database "${DATABASE_URL}" down
```

### バージョン確認

```bash
# 現在のマイグレーションバージョンを確認
migrate -path ./migrations -database "${DATABASE_URL}" version
```

### 強制的にバージョンを設定

```bash
# エラー状態から復旧する場合
migrate -path ./migrations -database "${DATABASE_URL}" force 2
```

## Makefileの活用

`backend/Makefile` を作成して便利にします：

```makefile
DB_URL=postgres://postgres:postgres@localhost:5432/effisio_dev?sslmode=disable

migrate-up:
	migrate -path ./migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DB_URL)" down 1

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ./migrations -seq $$name

migrate-version:
	migrate -path ./migrations -database "$(DB_URL)" version

migrate-force:
	@read -p "Enter version: " version; \
	migrate -path ./migrations -database "$(DB_URL)" force $$version
```

使用例：

```bash
cd backend

# マイグレーション適用
make migrate-up

# ロールバック
make migrate-down

# 新規作成
make migrate-create
# Enter migration name: add_audit_logs_table

# バージョン確認
make migrate-version
```

## ベストプラクティス

### 1. トランザクションを使用

```sql
-- up ファイル
BEGIN;

CREATE TABLE users (...);
CREATE INDEX ...;

COMMIT;
```

### 2. IF EXISTS / IF NOT EXISTS を使用

```sql
-- down ファイル
DROP TABLE IF EXISTS users;
DROP INDEX IF EXISTS idx_users_email;
```

### 3. 破壊的変更は慎重に

カラム削除やテーブル削除は本番環境で慎重に行う：

```sql
-- カラム削除（本番では要注意）
ALTER TABLE users DROP COLUMN IF EXISTS old_column;

-- データ保持しつつカラム名変更
ALTER TABLE users RENAME COLUMN old_name TO new_name;
```

### 4. データマイグレーションの分離

スキーマ変更とデータ変更は別のマイグレーションファイルに分ける：

```
000001_create_users_table.up.sql        # スキーマ
000002_migrate_user_data.up.sql         # データ
```

### 5. ロールバックの検証

マイグレーション作成後、必ず up → down → up のサイクルでテスト：

```bash
migrate -database "${DB_URL}" -path ./migrations up 1
migrate -database "${DB_URL}" -path ./migrations down 1
migrate -database "${DB_URL}" -path ./migrations up 1
```

## トラブルシューティング

### エラー: Dirty database version

**原因**: マイグレーション実行中にエラーが発生した場合

**解決方法**:

```bash
# 現在のバージョンを確認
migrate -path ./migrations -database "${DATABASE_URL}" version

# 強制的にバージョンをリセット（データロスの可能性あり）
migrate -path ./migrations -database "${DATABASE_URL}" force VERSION
```

### エラー: no change

**原因**: すでに最新のマイグレーションが適用済み

**解決方法**: 問題なし。新しいマイグレーションがない状態です。

### エラー: file does not exist

**原因**: マイグレーションファイルのパスが間違っている

**解決方法**:

```bash
# パスを確認
ls -la ./migrations/

# 正しいパスを指定
migrate -path ./migrations -database "${DATABASE_URL}" up
```

## 本番環境での運用

### 1. バックアップ

```bash
# マイグレーション前に必ずバックアップ
pg_dump -U postgres effisio_prod > backup_$(date +%Y%m%d_%H%M%S).sql
```

### 2. ステージング環境で検証

```bash
# ステージングで先にテスト
migrate -path ./migrations -database "${STAGING_DATABASE_URL}" up
```

### 3. メンテナンスウィンドウ

大規模なマイグレーションは、メンテナンスウィンドウ内で実行。

### 4. ロールバックプラン

必ずロールバック手順を準備してから実行。

---

## CI/CDとの統合

GitHub Actions で自動実行する例：

```yaml
# .github/workflows/migrate.yml
name: Database Migration

on:
  push:
    branches: [main, develop]
    paths:
      - 'backend/migrations/**'

jobs:
  migrate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Install migrate
        run: |
          curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
          sudo mv migrate /usr/local/bin/

      - name: Run migrations
        env:
          DATABASE_URL: ${{ secrets.DATABASE_URL }}
        run: |
          cd backend
          migrate -path ./migrations -database "${DATABASE_URL}" up
```

---

## まとめ

- **マイグレーションファイルは必ず up と down の両方を作成**
- **本番環境では必ずバックアップを取得**
- **ステージング環境で事前にテスト**
- **ロールバックプランを準備**

詳細は [golang-migrate のドキュメント](https://github.com/golang-migrate/migrate) を参照してください。
