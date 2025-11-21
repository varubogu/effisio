# データベース設計

## テーブル設計概要

### 1. users テーブル（ユーザー）

```sql
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
```

**カラム説明**:
- `id`: ユーザーID (主キー)
- `username`: ログインユーザー名
- `email`: メールアドレス
- `password_hash`: パスワードハッシュ (bcrypt推奨)
- `full_name`: ユーザーの実名
- `department`: 部門
- `role`: ロール (admin, user, viewer)
- `status`: ユーザーステータス (active, inactive, suspended)
- `last_login`: 最後のログイン時刻
- `created_at`: 作成日時
- `updated_at`: 更新日時
- `deleted_at`: 論理削除日時 (soft delete)

---

### 2. roles テーブル（ロール管理）

```sql
CREATE TABLE roles (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(100) UNIQUE NOT NULL,
  description TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO roles (name, description) VALUES
  ('admin', 'システム管理者'),
  ('manager', 'マネージャー'),
  ('user', '一般ユーザー'),
  ('viewer', '閲覧ユーザー');
```

---

### 3. permissions テーブル（権限）

```sql
CREATE TABLE permissions (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(100) UNIQUE NOT NULL,
  description TEXT,
  resource VARCHAR(100) NOT NULL,
  action VARCHAR(100) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 例
INSERT INTO permissions (name, description, resource, action) VALUES
  ('users:read', 'ユーザー一覧閲覧', 'users', 'read'),
  ('users:create', 'ユーザー作成', 'users', 'create'),
  ('users:update', 'ユーザー更新', 'users', 'update'),
  ('users:delete', 'ユーザー削除', 'users', 'delete');
```

---

### 4. role_permissions テーブル（ロール権限の関連付け）

```sql
CREATE TABLE role_permissions (
  id BIGSERIAL PRIMARY KEY,
  role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(role_id, permission_id)
);

CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
```

---

### 5. audit_logs テーブル（監査ログ）

```sql
CREATE TABLE audit_logs (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
  action VARCHAR(100) NOT NULL,
  resource_type VARCHAR(100) NOT NULL,
  resource_id BIGINT,
  old_values JSONB,
  new_values JSONB,
  ip_address VARCHAR(45),
  user_agent TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
```

---

### 6. sessions テーブル（セッション管理）

```sql
CREATE TABLE sessions (
  id VARCHAR(255) PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token VARCHAR(1000) NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  ip_address VARCHAR(45),
  user_agent TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
```

---

### 7. organizations テーブル（組織・部門）

```sql
CREATE TABLE organizations (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  parent_id BIGINT REFERENCES organizations(id),
  status VARCHAR(50) NOT NULL DEFAULT 'active',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_organizations_parent_id ON organizations(parent_id);
```

---

## スキーマ初期化スクリプト

### migrations/001_initial_schema.up.sql

マイグレーションツール (`migrate`, `sql-migrate` など) を使用して管理：

```
migrations/
├── 001_initial_schema.up.sql
├── 001_initial_schema.down.sql
├── 002_add_audit_logs.up.sql
├── 002_add_audit_logs.down.sql
...
```

---

## インデックス戦略

### 優先度の高いインデックス

1. `users.email`: ユーザーのログイン時に頻繁に使用
2. `users.username`: ユーザー検索に使用
3. `audit_logs.user_id`: 監査ログ検索
4. `audit_logs.created_at`: ログ期間検索
5. `sessions.user_id`: セッション検索
6. `sessions.expires_at`: 期限切れセッション削除
7. `role_permissions.role_id`: ロール権限取得

---

## パフォーマンス考慮事項

### クエリ最適化
- SELECT で必要なカラムのみを指定
- 多数のジョインは避ける（キャッシュの活用を検討）
- LIMIT/OFFSET ではなくカーソルベースページネーションを使用

### キャッシュ戦略
- Redis でユーザー情報をキャッシュ
- ロール・権限情報はメモリにロード
- 24時間の TTL を設定

### バックアップ戦略
- 日次フルバックアップ
- 時間ごとの差分バックアップ
- 本番環境でのWALアーカイブ

---

## マイグレーション管理ツール

推奨: **golang-migrate** または **Flyway**

```bash
# インストール
brew install golang-migrate

# マイグレーション実行
migrate -path ./migrations -database "postgres://..." up

# ロールバック
migrate -path ./migrations -database "postgres://..." down
```

---

## テーブル選定ポイント

- [ ] テーブル設計の確定
- [ ] マイグレーションツールの決定
- [ ] キャッシュ戦略の詳細化
- [ ] バックアップリカバリ計画の策定
