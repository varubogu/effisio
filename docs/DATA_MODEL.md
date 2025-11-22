# データモデル詳細仕様書

このドキュメントでは、Effisioプロジェクトで使用する全てのデータモデルの詳細を定義します。

## 目次

- [基本方針](#基本方針)
- [共通カラム](#共通カラム)
- [テーブル定義](#テーブル定義)
- [リレーション図](#リレーション図)
- [マイグレーション順序](#マイグレーション順序)
- [バリデーションルール](#バリデーションルール)

---

## 基本方針

### 命名規則

- **テーブル名**: 小文字スネークケース、複数形 (例: `users`, `audit_logs`)
- **カラム名**: 小文字スネークケース (例: `user_id`, `created_at`)
- **主キー**: `id` (BIGSERIAL)
- **外部キー**: `{テーブル名}_id` (例: `user_id`, `role_id`)

### データ型の使用基準

- **ID**: `BIGSERIAL` (将来的なスケールに対応)
- **文字列**: `VARCHAR(n)` (n は最大長を明示)
- **長文**: `TEXT`
- **真偽値**: `BOOLEAN`
- **日時**: `TIMESTAMP` (常にUTCで保存)
- **JSON**: `JSONB` (インデックス可能、検索高速)

### ソフトデリート

全てのマスターテーブルに `deleted_at TIMESTAMP` を設けてソフトデリートを実装します。

```sql
-- 削除フラグの確認
WHERE deleted_at IS NULL  -- 有効なレコード
WHERE deleted_at IS NOT NULL  -- 削除済みレコード
```

---

## 共通カラム

全てのテーブルに以下のカラムを含めます：

```sql
id BIGSERIAL PRIMARY KEY,
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
deleted_at TIMESTAMP  -- ソフトデリート用（マスターテーブルのみ）
```

**updated_at の自動更新トリガー:**

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 各テーブルに適用
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

---

## テーブル定義

### 1. users テーブル

**目的**: システムユーザーの管理

**完全なスキーマ:**

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    department VARCHAR(100),
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    last_login TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- インデックス
CREATE UNIQUE INDEX idx_users_username ON users(username) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- 制約
ALTER TABLE users ADD CONSTRAINT check_user_role
    CHECK (role IN ('admin', 'manager', 'user', 'viewer'));

ALTER TABLE users ADD CONSTRAINT check_user_status
    CHECK (status IN ('active', 'inactive', 'suspended'));

-- コメント
COMMENT ON TABLE users IS 'システムユーザー情報';
COMMENT ON COLUMN users.id IS 'ユーザーID（主キー）';
COMMENT ON COLUMN users.username IS 'ログインユーザー名（一意、3-50文字）';
COMMENT ON COLUMN users.email IS 'メールアドレス（一意）';
COMMENT ON COLUMN users.password_hash IS 'bcryptハッシュ化されたパスワード';
COMMENT ON COLUMN users.full_name IS 'ユーザーの実名（表示用）';
COMMENT ON COLUMN users.department IS '所属部署';
COMMENT ON COLUMN users.role IS 'ユーザーロール（admin/manager/user/viewer）';
COMMENT ON COLUMN users.status IS 'アカウント状態（active/inactive/suspended）';
COMMENT ON COLUMN users.last_login IS '最終ログイン日時';
```

**カラム詳細:**

| カラム名 | 型 | NULL | デフォルト | 説明 |
|---------|-----|------|-----------|------|
| id | BIGSERIAL | NO | AUTO | ユーザーID |
| username | VARCHAR(50) | NO | - | ログイン用ユーザー名 |
| email | VARCHAR(255) | NO | - | メールアドレス |
| password_hash | VARCHAR(255) | NO | - | bcryptハッシュ（cost=10） |
| full_name | VARCHAR(255) | YES | NULL | 実名 |
| department | VARCHAR(100) | YES | NULL | 部署名 |
| role | VARCHAR(20) | NO | 'user' | ロール |
| status | VARCHAR(20) | NO | 'active' | ステータス |
| last_login | TIMESTAMP | YES | NULL | 最終ログイン |
| created_at | TIMESTAMP | NO | NOW() | 作成日時 |
| updated_at | TIMESTAMP | NO | NOW() | 更新日時 |
| deleted_at | TIMESTAMP | YES | NULL | 削除日時 |

**Goモデル定義:**

```go
type User struct {
    ID           uint           `gorm:"primarykey" json:"id"`
    Username     string         `gorm:"uniqueIndex;not null;size:50" json:"username"`
    Email        string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
    PasswordHash string         `gorm:"column:password_hash;not null;size:255" json:"-"`
    FullName     *string        `gorm:"size:255" json:"full_name"`
    Department   *string        `gorm:"size:100" json:"department"`
    Role         string         `gorm:"not null;size:20;default:user" json:"role"`
    Status       string         `gorm:"not null;size:20;default:active" json:"status"`
    LastLogin    *time.Time     `json:"last_login"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
```

---

### 2. roles テーブル

**目的**: ロール定義の管理

```sql
CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 初期データ
INSERT INTO roles (name, display_name, description) VALUES
    ('admin', 'システム管理者', '全ての操作が可能'),
    ('manager', 'マネージャー', 'ユーザー管理と閲覧が可能'),
    ('user', '一般ユーザー', '基本的な操作が可能'),
    ('viewer', '閲覧者', '読み取り専用');

COMMENT ON TABLE roles IS 'ロール定義';
```

**カラム詳細:**

| カラム名 | 型 | NULL | 説明 |
|---------|-----|------|------|
| id | BIGSERIAL | NO | ロールID |
| name | VARCHAR(100) | NO | システム内部名（英数字） |
| display_name | VARCHAR(255) | NO | 表示名（日本語可） |
| description | TEXT | YES | ロールの説明 |

---

### 3. permissions テーブル

**目的**: 権限の定義

```sql
CREATE TABLE permissions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 初期データ例
INSERT INTO permissions (name, display_name, description, resource, action) VALUES
    ('users:read', 'ユーザー閲覧', 'ユーザー一覧・詳細の閲覧', 'users', 'read'),
    ('users:create', 'ユーザー作成', '新規ユーザーの作成', 'users', 'create'),
    ('users:update', 'ユーザー更新', 'ユーザー情報の更新', 'users', 'update'),
    ('users:delete', 'ユーザー削除', 'ユーザーの削除', 'users', 'delete'),
    ('dashboard:read', 'ダッシュボード閲覧', 'ダッシュボードの閲覧', 'dashboard', 'read');

CREATE INDEX idx_permissions_resource ON permissions(resource);

COMMENT ON TABLE permissions IS '権限定義';
```

**権限命名規則:**

```
{resource}:{action}

resource: users, dashboard, settings, audit_logs, etc.
action: read, create, update, delete
```

---

### 4. role_permissions テーブル

**目的**: ロールと権限の関連付け

```sql
CREATE TABLE role_permissions (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(role_id, permission_id)
);

CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);

COMMENT ON TABLE role_permissions IS 'ロールと権限の関連';
```

**初期データ例:**

```sql
-- admin: 全ての権限
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p WHERE r.name = 'admin';

-- manager: ユーザー管理とダッシュボード閲覧
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'manager' AND p.name IN ('users:read', 'users:create', 'users:update', 'dashboard:read');

-- user: 基本権限
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'user' AND p.name IN ('dashboard:read');

-- viewer: 閲覧のみ
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'viewer' AND p.name LIKE '%:read';
```

---

### 5. sessions テーブル

**目的**: ユーザーセッション管理（リフレッシュトークン保存）

```sql
CREATE TABLE sessions (
    id VARCHAR(255) PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token VARCHAR(1000) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

COMMENT ON TABLE sessions IS 'ユーザーセッション情報';
COMMENT ON COLUMN sessions.id IS 'セッションID（UUID）';
COMMENT ON COLUMN sessions.refresh_token IS 'リフレッシュトークン（JWT）';
COMMENT ON COLUMN sessions.expires_at IS 'セッション有効期限';
```

**期限切れセッション削除（定期実行）:**

```sql
-- 期限切れセッションを削除
DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP;
```

---

### 6. audit_logs テーブル

**目的**: 監査ログの記録

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
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);

COMMENT ON TABLE audit_logs IS '監査ログ';
COMMENT ON COLUMN audit_logs.action IS '実行されたアクション（create/update/delete/login等）';
COMMENT ON COLUMN audit_logs.resource_type IS '対象リソースの種類（users/roles等）';
COMMENT ON COLUMN audit_logs.resource_id IS '対象リソースのID';
COMMENT ON COLUMN audit_logs.old_values IS '変更前の値（JSON）';
COMMENT ON COLUMN audit_logs.new_values IS '変更後の値（JSON）';
```

**ログ例:**

```json
{
  "user_id": 1,
  "action": "update",
  "resource_type": "users",
  "resource_id": 5,
  "old_values": {"role": "user", "status": "active"},
  "new_values": {"role": "manager", "status": "active"},
  "ip_address": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "created_at": "2024-01-16T10:30:00Z"
}
```

---

### 7. organizations テーブル

**目的**: 組織・部門の階層管理

```sql
CREATE TABLE organizations (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE,
    description TEXT,
    parent_id BIGINT REFERENCES organizations(id) ON DELETE SET NULL,
    level INT NOT NULL DEFAULT 0,
    path VARCHAR(1000),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_organizations_parent_id ON organizations(parent_id);
CREATE INDEX idx_organizations_path ON organizations USING gin(path gin_trgm_ops);
CREATE INDEX idx_organizations_status ON organizations(status);

COMMENT ON TABLE organizations IS '組織・部門マスタ';
COMMENT ON COLUMN organizations.code IS '組織コード（一意）';
COMMENT ON COLUMN organizations.parent_id IS '親組織ID（階層構造）';
COMMENT ON COLUMN organizations.level IS '階層レベル（0=トップ）';
COMMENT ON COLUMN organizations.path IS '組織パス（例: /1/3/5）';
```

**階層例:**

```
会社 (id=1, parent_id=NULL, level=0, path=/1/)
├── 開発部 (id=2, parent_id=1, level=1, path=/1/2/)
│   ├── フロントエンド (id=4, parent_id=2, level=2, path=/1/2/4/)
│   └── バックエンド (id=5, parent_id=2, level=2, path=/1/2/5/)
└── 営業部 (id=3, parent_id=1, level=1, path=/1/3/)
```

---

## リレーション図

```
┌─────────────┐
│    users    │
├─────────────┤
│ id (PK)     │
│ username    │
│ email       │
│ role        │────┐
│ ...         │    │
└─────────────┘    │
       │           │
       │ 1:N       │
       │           │
       ▼           ▼
┌─────────────┐   ┌──────────────┐
│  sessions   │   │ audit_logs   │
├─────────────┤   ├──────────────┤
│ id (PK)     │   │ id (PK)      │
│ user_id (FK)│   │ user_id (FK) │
│ ...         │   │ action       │
└─────────────┘   │ resource_type│
                  │ ...          │
                  └──────────────┘

┌─────────────┐       ┌──────────────────┐       ┌──────────────┐
│    roles    │       │ role_permissions │       │ permissions  │
├─────────────┤       ├──────────────────┤       ├──────────────┤
│ id (PK)     │◄──────│ role_id (FK)     │       │ id (PK)      │
│ name        │ 1:N   │ permission_id(FK)│───────► name         │
│ ...         │       └──────────────────┘ N:M   │ resource     │
└─────────────┘                                   │ action       │
                                                  └──────────────┘

┌──────────────┐
│organizations │
├──────────────┤
│ id (PK)      │◄────┐
│ parent_id(FK)│─────┘ (自己参照)
│ name         │
│ path         │
└──────────────┘
```

---

## マイグレーション順序

マイグレーションは依存関係に基づいて以下の順序で実行します：

```
000001_create_users_table.sql          # 依存なし
000002_create_roles_table.sql          # 依存なし
000003_create_permissions_table.sql    # 依存なし
000004_create_role_permissions.sql     # roles, permissions に依存
000005_create_sessions_table.sql       # users に依存
000006_create_audit_logs_table.sql     # users に依存
000007_create_organizations_table.sql  # 依存なし（自己参照）
000008_add_triggers.sql                # 全テーブルに依存
000009_insert_seed_data.sql            # roles, permissions に依存
```

**各マイグレーションには必ず `.up.sql` と `.down.sql` を用意:**

```bash
migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_roles_table.up.sql
├── 000002_create_roles_table.down.sql
...
```

---

## バリデーションルール

### users テーブル

| フィールド | ルール | エラーメッセージ |
|-----------|--------|-----------------|
| username | 必須、3-50文字、英数字とアンダースコアのみ | "ユーザー名は3-50文字の英数字で入力してください" |
| email | 必須、メール形式、最大255文字 | "有効なメールアドレスを入力してください" |
| password | 必須、最低8文字、英大小+数字+記号を含む | "パスワードは8文字以上で、英大小文字、数字、記号を含めてください" |
| full_name | 任意、最大255文字 | "氏名は255文字以内で入力してください" |
| department | 任意、最大100文字 | "部署名は100文字以内で入力してください" |
| role | 必須、列挙値 | "有効なロールを選択してください" |
| status | 必須、列挙値 | "有効なステータスを選択してください" |

### Goでのバリデーション実装例

```go
type CreateUserRequest struct {
    Username   string `json:"username" binding:"required,min=3,max=50,alphanum"`
    Email      string `json:"email" binding:"required,email,max=255"`
    Password   string `json:"password" binding:"required,min=8,password"`
    FullName   string `json:"full_name" binding:"omitempty,max=255"`
    Department string `json:"department" binding:"omitempty,max=100"`
    Role       string `json:"role" binding:"omitempty,oneof=admin manager user viewer"`
}

// カスタムバリデーション: パスワード強度
func ValidatePassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()

    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
    hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

    return hasUpper && hasLower && hasNumber && hasSpecial
}
```

---

## デフォルト値一覧

| テーブル | カラム | デフォルト値 | 説明 |
|---------|--------|------------|------|
| users | role | 'user' | 新規ユーザーは一般ユーザー |
| users | status | 'active' | 新規ユーザーはアクティブ |
| organizations | status | 'active' | 新規組織はアクティブ |
| organizations | level | 0 | トップレベル組織 |
| 全テーブル | created_at | CURRENT_TIMESTAMP | 作成時刻を自動記録 |
| 全テーブル | updated_at | CURRENT_TIMESTAMP | 更新時刻を自動記録 |

---

## NULL許可ポリシー

### NULL を許可するケース
- オプショナルな情報（full_name, department等）
- 後から設定される情報（last_login等）
- 削除フラグ（deleted_at）
- 外部キーで削除時にNULLになる場合（audit_logs.user_id等）

### NULL を許可しないケース
- 必須の識別情報（username, email等）
- システムで必須の値（role, status等）
- タイムスタンプ（created_at, updated_at）

---

## パフォーマンス考慮事項

### インデックス戦略

1. **主キー**: 自動的にインデックス作成
2. **外部キー**: 常にインデックス作成
3. **UNIQUE制約**: 自動的にインデックス作成
4. **検索に使用するカラム**: インデックス作成
5. **ソート条件に使用するカラム**: インデックス作成

### 複合インデックスの考慮

```sql
-- よくある検索パターンに対応
CREATE INDEX idx_users_role_status ON users(role, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_audit_logs_user_created ON audit_logs(user_id, created_at DESC);
```

### パーティショニング（将来的）

audit_logs テーブルは月次パーティショニングを検討：

```sql
-- 例: 2024年1月のパーティション
CREATE TABLE audit_logs_2024_01 PARTITION OF audit_logs
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');
```

---

## 開発時の注意事項

1. **マイグレーションは必ず可逆的に**: down.sql で確実にロールバック可能にする
2. **外部キー制約は慎重に**: ON DELETE の動作を明確にする
3. **インデックスの追加は慎重に**: 書き込み性能への影響を考慮
4. **JSONB は適度に使用**: 構造化データは通常のカラムとして定義
5. **トランザクション分離レベル**: デフォルトの READ COMMITTED を使用
6. **文字コード**: UTF-8 を使用
7. **タイムゾーン**: 全て UTC で保存、表示時にローカル時刻に変換

---

## 変更履歴

| 日付 | バージョン | 変更内容 |
|------|-----------|---------|
| 2024-01-16 | 1.0.0 | 初版作成 |
