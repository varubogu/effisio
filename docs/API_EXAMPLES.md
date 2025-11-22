# API仕様の具体例

このドキュメントでは、Effisio APIの各エンドポイントに対する具体的なリクエスト/レスポンス例とcURLコマンドを提供します。

## 目次

- [共通仕様](#共通仕様)
- [認証API](#認証api)
- [ユーザーAPI](#ユーザーapi)
- [ロール・権限API](#ロール権限api)
- [組織API](#組織api)
- [ダッシュボードAPI](#ダッシュボードapi)
- [監査ログAPI](#監査ログapi)

---

## 共通仕様

### ベースURL

```
開発環境: http://localhost:8080/api/v1
本番環境: https://api.effisio.example.com/api/v1
```

### 共通ヘッダー

```
Content-Type: application/json
Accept: application/json
Authorization: Bearer {access_token}  # 認証が必要なエンドポイント
```

### レスポンス形式

**成功:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    // データ
  }
}
```

**エラー:**
```json
{
  "code": 400,
  "message": "error",
  "error": {
    "code": "USER_002",
    "message": "ユーザー名は既に使用されています"
  }
}
```

---

## 認証API

### POST /auth/login - ログイン

**リクエスト:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

**リクエストボディ:**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**レスポンス (200 OK):**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFkbWluIiwiZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsInJvbGUiOiJhZG1pbiIsInBlcm1pc3Npb25zIjpbInVzZXJzOnJlYWQiLCJ1c2VyczpjcmVhdGUiLCJ1c2Vyczp1cGRhdGUiLCJ1c2VyczpkZWxldGUiXSwiaXNzIjoiZWZmaXNpby1hcGkiLCJhdWQiOiJlZmZpc2lvLWNsaWVudCIsImV4cCI6MTcwNTQwNDAwMCwiaWF0IjoxNzA1NDAwNDAwfQ.abcdefg...",
    "expires_in": 900,
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "full_name": "システム管理者",
      "department": "IT部",
      "role": "admin",
      "status": "active",
      "last_login": "2024-01-16T10:30:00Z",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-16T10:30:00Z"
    }
  }
}
```

**エラーレスポンス (401 Unauthorized):**
```json
{
  "code": 401,
  "message": "error",
  "error": {
    "code": "AUTH_001",
    "message": "認証に失敗しました"
  }
}
```

---

### POST /auth/refresh - トークンリフレッシュ

**リクエスト:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -H "Cookie: refresh_token=eyJhbGc..." \
  -d '{}'
```

**レスポンス (200 OK):**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900
  }
}
```

---

### POST /auth/logout - ログアウト

**リクエスト:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer {access_token}" \
  -H "Cookie: refresh_token=eyJhbGc..."
```

**レスポンス (200 OK):**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "message": "Logged out successfully"
  }
}
```

---

## ユーザーAPI

### GET /users - ユーザー一覧取得

**リクエスト:**
```bash
curl -X GET "http://localhost:8080/api/v1/users?page=1&per_page=20&role=admin&status=active&search=john" \
  -H "Authorization: Bearer {access_token}"
```

**クエリパラメータ:**
- `page` (int, optional): ページ番号（デフォルト: 1）
- `per_page` (int, optional): 1ページあたりのアイテム数（デフォルト: 20、最大: 100）
- `role` (string, optional): ロールでフィルタ (admin, manager, user, viewer)
- `status` (string, optional): ステータスでフィルタ (active, inactive, suspended)
- `search` (string, optional): 名前またはメールで検索

**レスポンス (200 OK):**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "full_name": "John Doe",
      "department": "Engineering",
      "role": "admin",
      "status": "active",
      "last_login": "2024-01-16T10:30:00Z",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-15T14:20:00Z"
    },
    {
      "id": 2,
      "username": "jane_smith",
      "email": "jane@example.com",
      "full_name": "Jane Smith",
      "department": "HR",
      "role": "manager",
      "status": "active",
      "last_login": "2024-01-15T09:15:00Z",
      "created_at": "2024-01-02T00:00:00Z",
      "updated_at": "2024-01-15T09:15:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 50,
    "total_pages": 3,
    "has_next": true,
    "has_prev": false
  }
}
```

---

### GET /users/:id - ユーザー詳細取得

**リクエスト:**
```bash
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer {access_token}"
```

**レスポンス (200 OK):**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "full_name": "John Doe",
      "department": "Engineering",
      "role": "admin",
      "status": "active",
      "last_login": "2024-01-16T10:30:00Z",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-15T14:20:00Z"
    }
  }
}
```

**エラーレスポンス (404 Not Found):**
```json
{
  "code": 404,
  "message": "error",
  "error": {
    "code": "USER_001",
    "message": "ユーザーが見つかりません"
  }
}
```

---

### POST /users - ユーザー作成

**リクエスト:**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "new_user",
    "email": "newuser@example.com",
    "password": "SecurePass123!",
    "full_name": "New User",
    "department": "Sales",
    "role": "user"
  }'
```

**リクエストボディ:**
```json
{
  "username": "new_user",
  "email": "newuser@example.com",
  "password": "SecurePass123!",
  "full_name": "New User",
  "department": "Sales",
  "role": "user"
}
```

**レスポンス (201 Created):**
```json
{
  "code": 201,
  "message": "created",
  "data": {
    "user": {
      "id": 3,
      "username": "new_user",
      "email": "newuser@example.com",
      "full_name": "New User",
      "department": "Sales",
      "role": "user",
      "status": "active",
      "last_login": null,
      "created_at": "2024-01-16T11:00:00Z",
      "updated_at": "2024-01-16T11:00:00Z"
    }
  }
}
```

**バリデーションエラー (422 Unprocessable Entity):**
```json
{
  "code": 422,
  "message": "validation error",
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "リクエストの検証に失敗しました",
    "details": [
      {
        "field": "username",
        "message": "usernameは3文字以上である必要があります"
      },
      {
        "field": "password",
        "message": "パスワードは8文字以上で、英大小文字、数字、記号を含めてください"
      }
    ]
  }
}
```

**重複エラー (409 Conflict):**
```json
{
  "code": 409,
  "message": "error",
  "error": {
    "code": "USER_002",
    "message": "ユーザー名は既に使用されています"
  }
}
```

---

### PUT /users/:id - ユーザー更新

**リクエスト:**
```bash
curl -X PUT http://localhost:8080/api/v1/users/3 \
  -H "Authorization: Bearer {access_token}" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Updated Name",
    "department": "Marketing",
    "status": "active"
  }'
```

**リクエストボディ（全フィールドオプション）:**
```json
{
  "username": "updated_user",
  "email": "updated@example.com",
  "full_name": "Updated Name",
  "department": "Marketing",
  "role": "manager",
  "status": "active"
}
```

**レスポンス (200 OK):**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "user": {
      "id": 3,
      "username": "new_user",
      "email": "newuser@example.com",
      "full_name": "Updated Name",
      "department": "Marketing",
      "role": "user",
      "status": "active",
      "last_login": null,
      "created_at": "2024-01-16T11:00:00Z",
      "updated_at": "2024-01-16T11:30:00Z"
    }
  }
}
```

---

### DELETE /users/:id - ユーザー削除

**リクエスト:**
```bash
curl -X DELETE http://localhost:8080/api/v1/users/3 \
  -H "Authorization: Bearer {access_token}"
```

**レスポンス (200 OK):**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "message": "User deleted successfully"
  }
}
```

**権限不足エラー (403 Forbidden):**
```json
{
  "code": 403,
  "message": "error",
  "error": {
    "code": "AUTH_004",
    "message": "権限がありません"
  }
}
```

---

## ロール・権限API

### GET /roles - ロール一覧取得

**リクエスト:**
```bash
curl -X GET http://localhost:8080/api/v1/roles \
  -H "Authorization: Bearer {access_token}"
```

**レスポンス (200 OK):**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "admin",
      "display_name": "システム管理者",
      "description": "全ての操作が可能",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": 2,
      "name": "manager",
      "display_name": "マネージャー",
      "description": "ユーザー管理と閲覧が可能",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": 3,
      "name": "user",
      "display_name": "一般ユーザー",
      "description": "基本的な操作が可能",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": 4,
      "name": "viewer",
      "display_name": "閲覧者",
      "description": "読み取り専用",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

### GET /roles/:id/permissions - ロールの権限一覧

**リクエスト:**
```bash
curl -X GET http://localhost:8080/api/v1/roles/1/permissions \
  -H "Authorization: Bearer {access_token}"
```

**レスポンス (200 OK):**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "role": {
      "id": 1,
      "name": "admin",
      "display_name": "システム管理者"
    },
    "permissions": [
      {
        "id": 1,
        "name": "users:read",
        "display_name": "ユーザー閲覧",
        "resource": "users",
        "action": "read"
      },
      {
        "id": 2,
        "name": "users:create",
        "display_name": "ユーザー作成",
        "resource": "users",
        "action": "create"
      },
      {
        "id": 3,
        "name": "users:update",
        "display_name": "ユーザー更新",
        "resource": "users",
        "action": "update"
      },
      {
        "id": 4,
        "name": "users:delete",
        "display_name": "ユーザー削除",
        "resource": "users",
        "action": "delete"
      }
    ]
  }
}
```

---

## 組織API

### GET /organizations - 組織ツリー取得

**リクエスト:**
```bash
curl -X GET http://localhost:8080/api/v1/organizations \
  -H "Authorization: Bearer {access_token}"
```

**レスポンス (200 OK):**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "株式会社Effisio",
      "code": "EFFISIO",
      "description": "会社全体",
      "parent_id": null,
      "level": 0,
      "path": "/1/",
      "status": "active",
      "children": [
        {
          "id": 2,
          "name": "開発部",
          "code": "DEV",
          "parent_id": 1,
          "level": 1,
          "path": "/1/2/",
          "status": "active",
          "children": [
            {
              "id": 4,
              "name": "フロントエンドチーム",
              "code": "DEV_FE",
              "parent_id": 2,
              "level": 2,
              "path": "/1/2/4/",
              "status": "active",
              "children": []
            },
            {
              "id": 5,
              "name": "バックエンドチーム",
              "code": "DEV_BE",
              "parent_id": 2,
              "level": 2,
              "path": "/1/2/5/",
              "status": "active",
              "children": []
            }
          ]
        },
        {
          "id": 3,
          "name": "営業部",
          "code": "SALES",
          "parent_id": 1,
          "level": 1,
          "path": "/1/3/",
          "status": "active",
          "children": []
        }
      ]
    }
  ]
}
```

---

## ダッシュボードAPI

### GET /dashboard/overview - ダッシュボード概要

**リクエスト:**
```bash
curl -X GET http://localhost:8080/api/v1/dashboard/overview \
  -H "Authorization: Bearer {access_token}"
```

**レスポンス (200 OK):**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_users": 150,
    "active_users": 145,
    "inactive_users": 3,
    "suspended_users": 2,
    "users_by_role": {
      "admin": 5,
      "manager": 15,
      "user": 125,
      "viewer": 5
    },
    "recent_logins": [
      {
        "user_id": 1,
        "username": "john_doe",
        "full_name": "John Doe",
        "login_time": "2024-01-16T10:30:00Z"
      },
      {
        "user_id": 2,
        "username": "jane_smith",
        "full_name": "Jane Smith",
        "login_time": "2024-01-16T09:15:00Z"
      }
    ],
    "last_updated": "2024-01-16T11:00:00Z"
  }
}
```

---

## 監査ログAPI

### GET /audit-logs - 監査ログ一覧

**リクエスト:**
```bash
curl -X GET "http://localhost:8080/api/v1/audit-logs?page=1&per_page=20&user_id=1&action=login&from=2024-01-01&to=2024-01-31" \
  -H "Authorization: Bearer {access_token}"
```

**クエリパラメータ:**
- `page` (int): ページ番号
- `per_page` (int): 1ページあたりのアイテム数
- `user_id` (int): ユーザーIDでフィルタ
- `action` (string): アクションでフィルタ (login, create, update, delete等)
- `resource_type` (string): リソース種別でフィルタ (users, roles等)
- `from` (date): 開始日時（ISO 8601形式）
- `to` (date): 終了日時（ISO 8601形式）

**レスポンス (200 OK):**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "username": "admin",
      "action": "login",
      "resource_type": "auth",
      "resource_id": null,
      "old_values": null,
      "new_values": null,
      "ip_address": "192.168.1.100",
      "user_agent": "Mozilla/5.0...",
      "created_at": "2024-01-16T10:30:00Z"
    },
    {
      "id": 2,
      "user_id": 1,
      "username": "admin",
      "action": "create",
      "resource_type": "users",
      "resource_id": 5,
      "old_values": null,
      "new_values": {
        "username": "new_user",
        "email": "newuser@example.com",
        "role": "user",
        "status": "active"
      },
      "ip_address": "192.168.1.100",
      "user_agent": "Mozilla/5.0...",
      "created_at": "2024-01-16T10:35:00Z"
    },
    {
      "id": 3,
      "user_id": 1,
      "username": "admin",
      "action": "update",
      "resource_type": "users",
      "resource_id": 5,
      "old_values": {
        "role": "user",
        "status": "active"
      },
      "new_values": {
        "role": "manager",
        "status": "active"
      },
      "ip_address": "192.168.1.100",
      "user_agent": "Mozilla/5.0...",
      "created_at": "2024-01-16T10:40:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 150,
    "total_pages": 8,
    "has_next": true,
    "has_prev": false
  }
}
```

---

## エラーコード一覧

### 認証エラー (AUTH_xxx)

| コード | HTTPステータス | 説明 |
|-------|--------------|------|
| AUTH_001 | 401 | 認証失敗 |
| AUTH_002 | 401 | トークン無効 |
| AUTH_003 | 401 | トークン期限切れ |
| AUTH_004 | 403 | 権限不足 |

### ユーザーエラー (USER_xxx)

| コード | HTTPステータス | 説明 |
|-------|--------------|------|
| USER_001 | 404 | ユーザーが見つかりません |
| USER_002 | 409 | ユーザー名は既に使用されています |
| USER_003 | 409 | メールアドレスは既に使用されています |

### バリデーションエラー (VALIDATION_xxx)

| コード | HTTPステータス | 説明 |
|-------|--------------|------|
| VALIDATION_001 | 422 | バリデーションエラー |
| VALIDATION_002 | 400 | データ形式が正しくありません |

### サーバーエラー (SERVER_xxx)

| コード | HTTPステータス | 説明 |
|-------|--------------|------|
| SERVER_001 | 500 | サーバー内部エラー |
| SERVER_002 | 500 | データベースエラー |

---

## Postman Collection

開発を効率化するために、Postman Collectionをエクスポートして提供することを推奨します。

**ファイル名**: `effisio-api.postman_collection.json`

**使い方**:
1. Postmanを開く
2. File > Import からJSONファイルをインポート
3. 環境変数を設定（`base_url`, `access_token`）
4. リクエストを実行

---

## 変更履歴

| 日付 | バージョン | 変更内容 |
|------|-----------|---------|
| 2024-01-16 | 1.0.0 | 初版作成 |
