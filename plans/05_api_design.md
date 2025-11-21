# API 設計

## API 仕様概要

### 基本情報
- **プロトコル**: HTTPS
- **ベースURL**: `https://api.internal-system.com/v1` (本番環境)
- **レスポンス形式**: JSON
- **認証方式**: Bearer Token (JWT)
- **タイムゾーン**: UTC

---

## レスポンス形式

### 成功レスポンス

```json
{
  "code": 200,
  "message": "success",
  "data": {
    // レスポンスデータ
  }
}
```

### エラーレスポンス

```json
{
  "code": 400,
  "message": "Invalid request",
  "errors": [
    {
      "field": "email",
      "message": "Email format is invalid"
    }
  ]
}
```

### ページング付きレスポンス

```json
{
  "code": 200,
  "message": "success",
  "data": [
    // アイテム配列
  ],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5,
    "has_next": true,
    "has_prev": false
  }
}
```

---

## HTTP ステータスコード

| コード | 意味 | 説明 |
|-------|------|------|
| 200 | OK | リクエスト成功 |
| 201 | Created | リソース作成成功 |
| 204 | No Content | 成功、ボディなし |
| 400 | Bad Request | リクエスト形式エラー |
| 401 | Unauthorized | 認証エラー |
| 403 | Forbidden | 権限不足 |
| 404 | Not Found | リソース未検出 |
| 409 | Conflict | リソース重複 (例: ユーザー名重複) |
| 422 | Unprocessable Entity | バリデーションエラー |
| 429 | Too Many Requests | レート制限 |
| 500 | Internal Server Error | サーバーエラー |
| 503 | Service Unavailable | メンテナンス中 |

---

## 認証 API

### POST /auth/login
ユーザーがログインする

**リクエスト**:
```json
{
  "username": "john_doe",
  "password": "password123"
}
```

**レスポンス (成功)**:
```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "expires_in": 3600,
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "full_name": "John Doe",
      "role": "user"
    }
  }
}
```

---

### POST /auth/refresh
トークンをリフレッシュする

**リクエスト**:
```json
{
  "refresh_token": "eyJhbGc..."
}
```

---

### POST /auth/logout
ログアウト

**レスポンス**:
```json
{
  "code": 200,
  "message": "Logout successful"
}
```

---

## ユーザー API

### GET /users
ユーザー一覧を取得

**クエリパラメータ**:
- `page` (int): ページ番号（デフォルト: 1）
- `per_page` (int): 1ページあたりのアイテム数（デフォルト: 20）
- `role` (string): ロールでフィルタ
- `status` (string): ステータスでフィルタ
- `search` (string): 名前またはメールで検索

**レスポンス**:
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
      "role": "user",
      "status": "active",
      "last_login": "2024-01-15T10:30:00Z",
      "created_at": "2024-01-01T00:00:00Z"
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

### POST /users
新規ユーザーを作成 (管理者のみ)

**リクエスト**:
```json
{
  "username": "jane_doe",
  "email": "jane@example.com",
  "password": "securepassword",
  "full_name": "Jane Doe",
  "department": "HR",
  "role": "user"
}
```

**レスポンス**:
```json
{
  "code": 201,
  "message": "User created successfully",
  "data": {
    "id": 2,
    "username": "jane_doe",
    "email": "jane@example.com",
    "full_name": "Jane Doe",
    "role": "user",
    "created_at": "2024-01-16T14:00:00Z"
  }
}
```

---

### GET /users/{id}
特定のユーザー情報を取得

**レスポンス**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "department": "Engineering",
    "role": "user",
    "status": "active",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### PUT /users/{id}
ユーザー情報を更新 (管理者または本人)

**リクエスト**:
```json
{
  "full_name": "John Updated",
  "department": "Sales",
  "status": "active"
}
```

---

### DELETE /users/{id}
ユーザーを削除 (管理者のみ)

**レスポンス**:
```json
{
  "code": 200,
  "message": "User deleted successfully"
}
```

---

## ダッシュボード API

### GET /dashboard/overview
ダッシュボード概要情報を取得

**レスポンス**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_users": 150,
    "active_users": 145,
    "total_resources": 1200,
    "last_updated": "2024-01-16T15:00:00Z"
  }
}
```

---

## エラーハンドリング

### 標準エラーコード

```
AUTH_001 - 認証失敗
AUTH_002 - トークン無効
AUTH_003 - トークン期限切れ
USER_001 - ユーザー未検出
USER_002 - ユーザー名重複
USER_003 - メール重複
PERMISSION_001 - 権限不足
VALIDATION_001 - バリデーションエラー
SERVER_001 - 内部サーバーエラー
```

---

## レート制限

- **デフォルト**: 100 リクエスト / 分
- **認証済みユーザー**: 1000 リクエスト / 分

レスポンスヘッダ:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1642340400
```

---

## セキュリティヘッダー

すべてのレスポンスに以下を含める：

```
Content-Type: application/json
Content-Security-Policy: default-src 'self'
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000; includeSubDomains
```

---

## API ドキュメント

### Swagger/OpenAPI 生成

```bash
swag init -g cmd/server/main.go
```

生成されたドキュメント: `/docs/swagger.json`
Swagger UI: `http://localhost:8080/swagger/index.html`

---

## 実装チェックリスト

- [ ] 認証 API 実装
- [ ] ユーザー API 実装
- [ ] ダッシュボード API 実装
- [ ] エラーハンドリング実装
- [ ] レート制限実装
- [ ] ロギング実装
- [ ] API ドキュメント生成
