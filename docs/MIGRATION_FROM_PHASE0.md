# Phase 0からの移行ガイド

このドキュメントでは、Phase 0で作成された基本実装から、設計仕様に完全準拠した実装への移行手順を説明します。

## 目次

- [移行が必要な理由](#移行が必要な理由)
- [変更点の概要](#変更点の概要)
- [移行前の準備](#移行前の準備)
- [ステップバイステップ移行手順](#ステップバイステップ移行手順)
- [テストと検証](#テストと検証)
- [ロールバック手順](#ロールバック手順)
- [チェックリスト](#チェックリスト)

---

## 移行が必要な理由

Phase 0では、開発環境を迅速に立ち上げるため、簡略化されたスキーマとモデルを使用していました。しかし、本番運用や今後の機能追加を考慮すると、設計書通りの完全なスキーマに移行する必要があります。

**主な問題点:**

1. **セキュリティ**: `password` カラムは `password_hash` とすべき（ハッシュ化されていることを明示）
2. **ユーザー情報不足**: `full_name`（フルネーム）と `department`（部署）が欠けている
3. **ステータス管理**: `is_active` (boolean) では不十分。`status` (enum) で詳細な状態管理が必要
4. **監査証跡**: `last_login` がなく、ユーザーの最終ログイン時刻を追跡できない

---

## 変更点の概要

### データベーススキーマの変更

| 項目 | Phase 0 (現在) | 設計仕様 (移行後) | 変更種別 |
|------|---------------|------------------|---------|
| **追加カラム** | - | `full_name VARCHAR(100)` | ➕ 追加 |
| | - | `department VARCHAR(100)` | ➕ 追加 |
| | - | `last_login TIMESTAMP` | ➕ 追加 |
| **カラム名変更** | `password VARCHAR(255)` | `password_hash VARCHAR(255)` | 🔄 リネーム |
| | `is_active BOOLEAN` | `status VARCHAR(20)` | 🔄 型変更 |
| **デフォルト値** | `is_active DEFAULT true` | `status DEFAULT 'active'` | 🔄 変更 |

### Goモデルの変更

**Phase 0 (現在):**
```go
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Email     string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Password  string         `gorm:"not null;size:255" json:"-"`
	Role      string         `gorm:"not null;size:20;default:'user'" json:"role"`
	IsActive  bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

**設計仕様 (移行後):**
```go
type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Username     string         `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Email        string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	FullName     string         `gorm:"size:100" json:"full_name"`                    // ➕ 追加
	Department   string         `gorm:"size:100" json:"department"`                   // ➕ 追加
	PasswordHash string         `gorm:"not null;size:255;column:password_hash" json:"-"` // 🔄 リネーム
	Role         string         `gorm:"not null;size:20;default:'user'" json:"role"`
	Status       string         `gorm:"not null;size:20;default:'active'" json:"status"` // 🔄 型変更
	LastLogin    *time.Time     `json:"last_login"`                                   // ➕ 追加
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
```

### APIレスポンスの変更

**Phase 0 (現在):**
```json
{
  "id": 1,
  "username": "admin",
  "email": "admin@example.com",
  "role": "admin",
  "is_active": true,
  "created_at": "2025-01-15T10:00:00Z",
  "updated_at": "2025-01-15T10:00:00Z"
}
```

**設計仕様 (移行後):**
```json
{
  "id": 1,
  "username": "admin",
  "email": "admin@example.com",
  "full_name": "管理者 太郎",
  "department": "IT部",
  "role": "admin",
  "status": "active",
  "last_login": "2025-01-20T14:30:00Z",
  "created_at": "2025-01-15T10:00:00Z",
  "updated_at": "2025-01-15T10:00:00Z"
}
```

---

## 移行前の準備

### 1. データバックアップ

**必須作業:** 移行前に必ずデータベース全体をバックアップしてください。

```bash
# PostgreSQLのバックアップ
docker-compose exec postgres pg_dump -U postgres effisio_dev > backup_phase0_$(date +%Y%m%d_%H%M%S).sql

# バックアップの確認
ls -lh backup_phase0_*.sql
```

### 2. 現在のデータ確認

```bash
# 既存のユーザー数を確認
docker-compose exec postgres psql -U postgres -d effisio_dev -c "SELECT COUNT(*) FROM users;"

# 既存のユーザーデータをエクスポート
docker-compose exec postgres psql -U postgres -d effisio_dev -c "COPY users TO STDOUT CSV HEADER;" > users_backup.csv
```

### 3. 開発環境の停止

```bash
# 全てのコンテナを停止
docker-compose down
```

---

## ステップバイステップ移行手順

### ステップ1: 新しいマイグレーションファイルを作成

**ファイル: `backend/migrations/000002_update_users_table.up.sql`**

```sql
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
```

**ファイル: `backend/migrations/000002_update_users_table.down.sql`**

```sql
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
```

### ステップ2: Goモデルを更新

**ファイル: `backend/internal/model/user.go`**

```go
package model

import (
	"time"
	"gorm.io/gorm"
)

// User ユーザーモデル
type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Username     string         `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Email        string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	FullName     string         `gorm:"size:100" json:"full_name"`
	Department   string         `gorm:"size:100" json:"department"`
	PasswordHash string         `gorm:"not null;size:255;column:password_hash" json:"-"`
	Role         string         `gorm:"not null;size:20;default:'user'" json:"role"`
	Status       string         `gorm:"not null;size:20;default:'active'" json:"status"`
	LastLogin    *time.Time     `json:"last_login"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName テーブル名を指定
func (User) TableName() string {
	return "users"
}

// UserResponse APIレスポンス用
type UserResponse struct {
	ID         uint       `json:"id"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	FullName   string     `json:"full_name"`
	Department string     `json:"department"`
	Role       string     `json:"role"`
	Status     string     `json:"status"`
	LastLogin  *time.Time `json:"last_login"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// ToResponse UserからUserResponseへ変換
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:         u.ID,
		Username:   u.Username,
		Email:      u.Email,
		FullName:   u.FullName,
		Department: u.Department,
		Role:       u.Role,
		Status:     u.Status,
		LastLogin:  u.LastLogin,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

// CreateUserRequest ユーザー作成リクエスト
type CreateUserRequest struct {
	Username   string `json:"username" binding:"required,min=3,max=50"`
	Email      string `json:"email" binding:"required,email"`
	FullName   string `json:"full_name" binding:"max=100"`
	Department string `json:"department" binding:"max=100"`
	Password   string `json:"password" binding:"required,min=8"`
	Role       string `json:"role" binding:"required,oneof=admin manager user viewer"`
}

// UpdateUserRequest ユーザー更新リクエスト
type UpdateUserRequest struct {
	Email      *string `json:"email" binding:"omitempty,email"`
	FullName   *string `json:"full_name" binding:"omitempty,max=100"`
	Department *string `json:"department" binding:"omitempty,max=100"`
	Role       *string `json:"role" binding:"omitempty,oneof=admin manager user viewer"`
	Status     *string `json:"status" binding:"omitempty,oneof=active inactive suspended"`
}

// UserStatus ステータス定数
const (
	UserStatusActive    = "active"
	UserStatusInactive  = "inactive"
	UserStatusSuspended = "suspended"
)

// UserRole ロール定数
const (
	RoleAdmin   = "admin"
	RoleManager = "manager"
	RoleUser    = "user"
	RoleViewer  = "viewer"
)

// IsValidStatus ステータスの妥当性チェック
func IsValidStatus(status string) bool {
	return status == UserStatusActive || status == UserStatusInactive || status == UserStatusSuspended
}

// IsValidRole ロールの妥当性チェック
func IsValidRole(role string) bool {
	return role == RoleAdmin || role == RoleManager || role == RoleUser || role == RoleViewer
}
```

### ステップ3: サービス層を更新

**ファイル: `backend/internal/service/user.go`**

主な変更点:
- `Password` → `PasswordHash` に変更
- `IsActive` → `Status` に変更
- `FullName`, `Department` の追加

```go
// Create メソッドを更新
func (s *UserService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodePasswordHashError, err)
	}

	// ユーザーモデルを作成
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		FullName:     req.FullName,     // ➕ 追加
		Department:   req.Department,   // ➕ 追加
		PasswordHash: string(hashedPassword), // 🔄 変更
		Role:         req.Role,
		Status:       model.UserStatusActive, // 🔄 変更
	}

	// データベースに保存
	if err := s.repo.Create(ctx, user); err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return user.ToResponse(), nil
}

// Update メソッドを更新
func (s *UserService) Update(ctx context.Context, id uint, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	// 既存ユーザーを取得
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, util.NewNotFoundError(util.ErrCodeUserNotFound, err)
	}

	// 更新データを適用
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.FullName != nil {
		user.FullName = *req.FullName // ➕ 追加
	}
	if req.Department != nil {
		user.Department = *req.Department // ➕ 追加
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.Status != nil {
		user.Status = *req.Status // 🔄 変更
	}

	// データベースを更新
	if err := s.repo.Update(ctx, user); err != nil {
		s.logger.Error("Failed to update user", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	return user.ToResponse(), nil
}
```

### ステップ4: シードデータを更新

**ファイル: `backend/scripts/seed.sh`**

```bash
#!/bin/bash

set -e

# データベース接続情報
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-effisio_dev}

echo "🌱 シードデータを投入しています..."

# bcryptハッシュ（cost=10）
# admin123 -> $2a$10$X8yI6qZ...
# manager123 -> $2a$10$Y9zJ7rA...
# user123 -> $2a$10$Z0aK8sB...
# viewer123 -> $2a$10$A1bL9tC...

PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME <<EOF

-- 既存のシードデータをクリア
DELETE FROM users WHERE username IN ('admin', 'manager', 'testuser', 'viewer');

-- 管理者
INSERT INTO users (username, email, full_name, department, password_hash, role, status, created_at, updated_at)
VALUES (
  'admin',
  'admin@example.com',
  '管理者 太郎',
  'IT部',
  '\$2a\$10\$X8yI6qZvKZH5mP3nR4tVH.YqJ5mN6oP7qR8sT9uV0wX1yZ2aB3cD4',
  'admin',
  'active',
  NOW(),
  NOW()
);

-- マネージャー
INSERT INTO users (username, email, full_name, department, password_hash, role, status, created_at, updated_at)
VALUES (
  'manager',
  'manager@example.com',
  '管理 次郎',
  '営業部',
  '\$2a\$10\$Y9zJ7rAvLAI6nQ4oS5uWI.ZrK6nO7pQ8rS9tU0vW1xY2zA3bC4dE5',
  'manager',
  'active',
  NOW(),
  NOW()
);

-- 一般ユーザー
INSERT INTO users (username, email, full_name, department, password_hash, role, status, created_at, updated_at)
VALUES (
  'testuser',
  'testuser@example.com',
  'テスト 三郎',
  '開発部',
  '\$2a\$10\$Z0aK8sBwMBJ7oR5pT6vXJ.AsL7oP8qR9sT0uV1wX2yZ3aB4cD5eF6',
  'user',
  'active',
  NOW(),
  NOW()
);

-- 閲覧者
INSERT INTO users (username, email, full_name, department, password_hash, role, status, created_at, updated_at)
VALUES (
  'viewer',
  'viewer@example.com',
  '閲覧 四郎',
  '総務部',
  '\$2a\$10\$A1bL9tCxNCK8pS6qU7wYK.BtM8pQ9rS0tU1vW2xY3zA4bC5dD6eG7',
  'viewer',
  'active',
  NOW(),
  NOW()
);

-- 停止中のユーザー（テスト用）
INSERT INTO users (username, email, full_name, department, password_hash, role, status, created_at, updated_at)
VALUES (
  'suspended_user',
  'suspended@example.com',
  '停止 五郎',
  'なし',
  '\$2a\$10\$B2cM0uDyODL9qT7rV8xZL.CuN9qR0sT1uV2wX3yZ4aB5cD6eE7fH8',
  'user',
  'suspended',
  NOW(),
  NOW()
);

EOF

echo "✅ シードデータの投入が完了しました"
echo ""
echo "作成されたユーザー:"
echo "  admin     / admin123     (管理者 太郎 - IT部)"
echo "  manager   / manager123   (管理 次郎 - 営業部)"
echo "  testuser  / user123      (テスト 三郎 - 開発部)"
echo "  viewer    / viewer123    (閲覧 四郎 - 総務部)"
echo "  suspended_user / suspended123 (停止 五郎 - 停止中)"
```

### ステップ5: マイグレーション実行

```bash
# 1. 開発環境を起動
make dev

# 2. 別のターミナルで、新しいマイグレーションを実行
make migrate-up

# 期待される出力:
# 📊 マイグレーションを実行しています...
# 000002_update_users_table.up.sql
# ✅ マイグレーション完了

# 3. マイグレーション結果を確認
docker-compose exec postgres psql -U postgres -d effisio_dev -c "\d users"

# 期待される出力（新しいカラムが追加されている）:
#                                          Table "public.users"
#     Column     |            Type             | Nullable |              Default
# ---------------+-----------------------------+----------+-----------------------------------
#  id            | integer                     | not null | nextval('users_id_seq'::regclass)
#  username      | character varying(50)       | not null |
#  email         | character varying(255)      | not null |
#  full_name     | character varying(100)      |          |
#  department    | character varying(100)      |          |
#  password_hash | character varying(255)      | not null |
#  role          | character varying(20)       | not null | 'user'::character varying
#  status        | character varying(20)       | not null | 'active'::character varying
#  last_login    | timestamp without time zone |          |
#  created_at    | timestamp without time zone | not null | CURRENT_TIMESTAMP
#  updated_at    | timestamp without time zone | not null | CURRENT_TIMESTAMP
#  deleted_at    | timestamp without time zone |          |
```

### ステップ6: シードデータ再投入

```bash
# 新しいスキーマに対応したシードデータを投入
make seed

# 期待される出力:
# 🌱 シードデータを投入しています...
# ✅ シードデータの投入が完了しました
#
# 作成されたユーザー:
#   admin     / admin123     (管理者 太郎 - IT部)
#   manager   / manager123   (管理 次郎 - 営業部)
#   testuser  / user123      (テスト 三郎 - 開発部)
#   viewer    / viewer123    (閲覧 四郎 - 総務部)
```

### ステップ7: Go依存関係の更新とビルド

```bash
# backendディレクトリに移動
cd backend

# Go modulesを整理
go mod tidy

# ビルドして構文エラーがないか確認
make build

# 期待される出力:
# 🔨 ビルドしています...
# ✅ ビルド完了: bin/server
```

### ステップ8: サーバー再起動

```bash
# プロジェクトルートに戻る
cd ..

# Dockerコンテナを再起動
docker-compose restart backend

# ログを確認
docker-compose logs -f backend

# 期待される出力:
# backend_1  | {"level":"info","ts":1642345678,"msg":"Server starting","port":"8080"}
```

---

## テストと検証

### 1. APIエンドポイントの動作確認

```bash
# Pingエンドポイント
curl http://localhost:8080/api/v1/ping

# 期待される出力:
# {"message":"pong"}

# ユーザー一覧取得
curl http://localhost:8080/api/v1/users | jq

# 期待される出力（新しいフィールドが含まれている）:
# {
#   "users": [
#     {
#       "id": 1,
#       "username": "admin",
#       "email": "admin@example.com",
#       "full_name": "管理者 太郎",
#       "department": "IT部",
#       "role": "admin",
#       "status": "active",
#       "last_login": null,
#       "created_at": "2025-01-15T10:00:00Z",
#       "updated_at": "2025-01-15T10:00:00Z"
#     }
#   ]
# }
```

### 2. データベースの直接確認

```bash
# PostgreSQLに接続
docker-compose exec postgres psql -U postgres -d effisio_dev

# 全ユーザーを確認
SELECT id, username, email, full_name, department, role, status FROM users;

# 期待される出力:
#  id |    username     |        email         |   full_name  | department |  role   | status
# ----+-----------------+----------------------+--------------+------------+---------+--------
#   1 | admin           | admin@example.com    | 管理者 太郎   | IT部       | admin   | active
#   2 | manager         | manager@example.com  | 管理 次郎     | 営業部     | manager | active
#   3 | testuser        | testuser@example.com | テスト 三郎   | 開発部     | user    | active
#   4 | viewer          | viewer@example.com   | 閲覧 四郎     | 総務部     | viewer  | active
#   5 | suspended_user  | suspended@example.com| 停止 五郎     | なし       | user    | suspended

# 終了
\q
```

### 3. ユーザー作成のテスト

```bash
# 新しいユーザーを作成
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "email": "newuser@example.com",
    "full_name": "新規 太郎",
    "department": "人事部",
    "password": "newuser123",
    "role": "user"
  }' | jq

# 期待される出力:
# {
#   "user": {
#     "id": 6,
#     "username": "newuser",
#     "email": "newuser@example.com",
#     "full_name": "新規 太郎",
#     "department": "人事部",
#     "role": "user",
#     "status": "active",
#     "last_login": null,
#     "created_at": "2025-01-20T15:00:00Z",
#     "updated_at": "2025-01-20T15:00:00Z"
#   }
# }
```

### 4. ユーザー更新のテスト

```bash
# ユーザー情報を更新
curl -X PUT http://localhost:8080/api/v1/users/6 \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "新規 次郎",
    "department": "総務部",
    "status": "inactive"
  }' | jq

# 期待される出力:
# {
#   "user": {
#     "id": 6,
#     "username": "newuser",
#     "email": "newuser@example.com",
#     "full_name": "新規 次郎",
#     "department": "総務部",
#     "role": "user",
#     "status": "inactive",
#     "last_login": null,
#     "created_at": "2025-01-20T15:00:00Z",
#     "updated_at": "2025-01-20T15:01:00Z"
#   }
# }
```

### 5. フロントエンドの確認

```bash
# ブラウザで確認
# http://localhost:3000/users にアクセス

# 期待される表示:
# - 全ユーザーの一覧が表示される
# - 新しいフィールド（full_name, department, status）が表示される
```

**注意:** フロントエンドのTypeScript型定義も更新する必要があります（後述）。

### 6. テストスイートの実行

```bash
# バックエンドテストを実行
cd backend
make test

# 期待される出力:
# === RUN   TestUserRepository_Create
# --- PASS: TestUserRepository_Create (0.05s)
# ...
# PASS
# coverage: 75.0% of statements
```

---

## フロントエンドの更新

### TypeScript型定義の更新

**ファイル: `frontend/src/types/user.ts`**

```typescript
export interface User {
  id: number;
  username: string;
  email: string;
  full_name: string;      // ➕ 追加
  department: string;     // ➕ 追加
  role: 'admin' | 'manager' | 'user' | 'viewer';
  status: 'active' | 'inactive' | 'suspended'; // 🔄 変更（is_active から status へ）
  last_login: string | null; // ➕ 追加
  created_at: string;
  updated_at: string;
}

export interface CreateUserRequest {
  username: string;
  email: string;
  full_name?: string;     // ➕ 追加
  department?: string;    // ➕ 追加
  password: string;
  role: 'admin' | 'manager' | 'user' | 'viewer';
}

export interface UpdateUserRequest {
  email?: string;
  full_name?: string;     // ➕ 追加
  department?: string;    // ➕ 追加
  role?: 'admin' | 'manager' | 'user' | 'viewer';
  status?: 'active' | 'inactive' | 'suspended'; // 🔄 変更
}

export type UserStatus = 'active' | 'inactive' | 'suspended'; // ➕ 追加
export type UserRole = 'admin' | 'manager' | 'user' | 'viewer';
```

### ユーザーリストコンポーネントの更新

**ファイル: `frontend/src/components/users/UserList.tsx`**

```typescript
export function UserList({ users }: UserListProps) {
  if (users.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500">
        ユーザーがいません
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ユーザー名</th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">メール</th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">氏名</th> {/* ➕ 追加 */}
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">部署</th> {/* ➕ 追加 */}
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ロール</th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ステータス</th> {/* 🔄 変更 */}
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">最終ログイン</th> {/* ➕ 追加 */}
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {users.map((user) => (
            <tr key={user.id} className="hover:bg-gray-50">
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{user.id}</td>
              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{user.username}</td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{user.email}</td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{user.full_name || '-'}</td> {/* ➕ 追加 */}
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{user.department || '-'}</td> {/* ➕ 追加 */}
              <td className="px-6 py-4 whitespace-nowrap">
                <RoleBadge role={user.role} />
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <StatusBadge status={user.status} /> {/* 🔄 変更 */}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {user.last_login ? new Date(user.last_login).toLocaleString('ja-JP') : '未ログイン'} {/* ➕ 追加 */}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
```

### ステータスバッジコンポーネントの作成

**ファイル: `frontend/src/components/users/StatusBadge.tsx`**

```typescript
import { UserStatus } from '@/types/user';

interface StatusBadgeProps {
  status: UserStatus;
}

export function StatusBadge({ status }: StatusBadgeProps) {
  const styles = {
    active: 'bg-green-100 text-green-800',
    inactive: 'bg-gray-100 text-gray-800',
    suspended: 'bg-red-100 text-red-800',
  };

  const labels = {
    active: 'アクティブ',
    inactive: '非アクティブ',
    suspended: '停止中',
  };

  return (
    <span className={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${styles[status]}`}>
      {labels[status]}
    </span>
  );
}
```

### フロントエンドのビルドとテスト

```bash
cd frontend

# TypeScriptの型チェック
npm run type-check

# 期待される出力（エラーなし）:
# ✓ Type checking completed successfully

# リンターを実行
npm run lint

# ビルド
npm run build

# 期待される出力:
# ✓ Compiled successfully

# 開発サーバー起動
npm run dev
```

---

## ロールバック手順

移行に問題が発生した場合、以下の手順でPhase 0に戻すことができます。

### 1. マイグレーションのロールバック

```bash
# マイグレーションを1つ戻す
make migrate-down

# または、特定のバージョンまで戻す
migrate -path backend/migrations -database "postgresql://postgres:postgres@localhost:5432/effisio_dev?sslmode=disable" down 1
```

### 2. データベースバックアップからリストア

```bash
# 完全にリストアする場合
docker-compose down -v
docker-compose up -d postgres

# バックアップからリストア
cat backup_phase0_20250120_100000.sql | docker-compose exec -T postgres psql -U postgres -d effisio_dev
```

### 3. Goモデルを旧版に戻す

```bash
# Gitで変更を戻す
git checkout HEAD~1 backend/internal/model/user.go
git checkout HEAD~1 backend/internal/service/user.go
git checkout HEAD~1 backend/scripts/seed.sh
```

### 4. フロントエンドの型定義を戻す

```bash
git checkout HEAD~1 frontend/src/types/user.ts
git checkout HEAD~1 frontend/src/components/users/UserList.tsx
```

### 5. 動作確認

```bash
# サーバー再起動
docker-compose restart backend

# APIテスト
curl http://localhost:8080/api/v1/users | jq

# 旧フォーマットが返ることを確認
# {
#   "users": [
#     {
#       "id": 1,
#       "username": "admin",
#       "email": "admin@example.com",
#       "role": "admin",
#       "is_active": true,
#       ...
#     }
#   ]
# }
```

---

## チェックリスト

移行作業を確実に完了させるためのチェックリストです。

### 移行前

- [ ] データベース全体のバックアップを取得した
- [ ] 既存のユーザーデータをCSVでエクスポートした
- [ ] 開発環境を停止した
- [ ] 移行手順書を読んで理解した

### マイグレーション

- [ ] 新しいマイグレーションファイル（up/down）を作成した
- [ ] マイグレーションファイルの構文が正しいことを確認した
- [ ] `make migrate-up` を実行してエラーなく完了した
- [ ] `\d users` でスキーマ変更を確認した

### コード更新

- [ ] `backend/internal/model/user.go` を更新した
- [ ] `backend/internal/service/user.go` を更新した
- [ ] `backend/scripts/seed.sh` を更新した
- [ ] `frontend/src/types/user.ts` を更新した
- [ ] `frontend/src/components/users/UserList.tsx` を更新した
- [ ] `frontend/src/components/users/StatusBadge.tsx` を作成した

### テストと検証

- [ ] `make seed` でシードデータを投入した
- [ ] `make build` でビルドエラーがないことを確認した
- [ ] `curl http://localhost:8080/api/v1/ping` が成功した
- [ ] `curl http://localhost:8080/api/v1/users` で新しいフィールドが返ることを確認した
- [ ] PostgreSQLに直接接続してデータを確認した
- [ ] 新しいユーザー作成APIが動作することを確認した
- [ ] ユーザー更新APIで新しいフィールドを更新できることを確認した
- [ ] ユーザー削除APIが動作することを確認した
- [ ] フロントエンドで新しいフィールドが表示されることを確認した
- [ ] `npm run type-check` でTypeScriptエラーがないことを確認した
- [ ] `make test` でバックエンドテストがパスした

### ドキュメント

- [ ] CHANGELOG.md に移行内容を記載した
- [ ] README.md の必要箇所を更新した
- [ ] API仕様書を最新化した

### 完了

- [ ] 全ての変更をコミットした
- [ ] Git tagを作成した（例: `v0.2.0-phase1-start`）
- [ ] チームメンバーに移行完了を通知した

---

## よくある質問

### Q1. 既存のユーザーデータは失われませんか？

A1. いいえ、失われません。マイグレーションは既存データを保持したまま、新しいカラムを追加します。`is_active` から `status` への変換も自動で行われます。

### Q2. full_name や department が NULL のままでも大丈夫ですか？

A2. はい、これらのフィールドは NULL を許可しています。後から管理画面などで入力してもらうことを想定しています。

### Q3. マイグレーションが途中で失敗したらどうすればいいですか？

A3. PostgreSQL はトランザクション内でマイグレーションを実行するため、失敗した場合は自動的にロールバックされます。エラーメッセージを確認し、マイグレーションファイルを修正してから再実行してください。

### Q4. Phase 0 のコードに戻したい場合は？

A4. [ロールバック手順](#ロールバック手順) に従って、マイグレーションを戻し、Git で変更をリバートしてください。

### Q5. フロントエンドの型エラーが出ます

A5. `frontend/src/types/user.ts` が正しく更新されているか確認してください。また、`npm run type-check` で詳細なエラーメッセージを確認できます。

### Q6. シードデータが投入されません

A6. マイグレーションが完了していることを確認してください。また、PostgreSQL コンテナが起動していることを `docker ps` で確認してください。

---

## まとめ

この移行ガイドに従うことで、Phase 0 の簡略化された実装から、設計書通りの完全な実装に移行できます。

**移行後の利点:**

1. **セキュリティ向上**: `password_hash` という明示的な命名
2. **ユーザー情報の充実**: フルネーム、部署の追加
3. **詳細なステータス管理**: active/inactive/suspended の3状態
4. **監査証跡**: 最終ログイン時刻の記録

**次のステップ:**

移行が完了したら、**[PHASE1_IMPLEMENTATION_STEPS.md](PHASE1_IMPLEMENTATION_STEPS.md)** に進んで、Phase 1 の残りのタスクを実装してください。

---

## サポート

問題が発生した場合:

1. **[TROUBLESHOOTING.md](TROUBLESHOOTING.md)** を確認
2. GitHub Issues で検索
3. 新しい Issue を作成

---

**最終更新**: 2025-01-20
