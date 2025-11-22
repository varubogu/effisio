# 日々の開発ワークフロー

このドキュメントでは、Effisioプロジェクトでの日常的な開発作業の流れを説明します。

## 目次

- [朝の立ち上げルーチン](#朝の立ち上げルーチン)
- [機能開発のワークフロー](#機能開発のワークフロー)
- [テストのワークフロー](#テストのワークフロー)
- [デバッグのワークフロー](#デバッグのワークフロー)
- [コードレビューのワークフロー](#コードレビューのワークフロー)
- [コミット前のチェックリスト](#コミット前のチェックリスト)
- [終業時のルーチン](#終業時のルーチン)
- [よく使うコマンド集](#よく使うコマンド集)

---

## 朝の立ち上げルーチン

### 1. 最新コードの取得 (2分)

```bash
# プロジェクトディレクトリに移動
cd ~/effisio

# 最新のコードをpull
git checkout main
git pull origin main

# 作業ブランチに移動（または新規作成）
git checkout feature/your-feature-name
# または
git checkout -b feature/new-feature-name
```

### 2. 依存関係の更新確認 (1分)

```bash
# Go modulesの更新確認
cd backend
go mod download

# npm パッケージの更新確認
cd ../frontend
npm install
cd ..
```

### 3. 開発環境の起動 (2分)

```bash
# Docker環境を起動
make dev

# 起動確認
docker-compose ps

# 期待される出力（全てStateがUp）:
#        Name                      Command               State           Ports
# -----------------------------------------------------------------------------------
# effisio_postgres_1    docker-entrypoint.sh postgres   Up      0.0.0.0:5432->5432/tcp
# effisio_redis_1       docker-entrypoint.sh redis ...  Up      0.0.0.0:6379->6379/tcp
# effisio_backend_1     air -c .air.toml                Up      0.0.0.0:8080->8080/tcp
# effisio_frontend_1    npm run dev                     Up      0.0.0.0:3000->3000/tcp
```

### 4. 動作確認 (1分)

```bash
# バックエンドAPIの疎通確認
curl http://localhost:8080/api/v1/ping

# 期待される出力:
# {"message":"pong"}

# ブラウザで確認
# フロントエンド: http://localhost:3000
# Adminer: http://localhost:8081
```

**これで開発準備完了！** 所要時間: 約5-6分

---

## 機能開発のワークフロー

### 基本的な開発サイクル

```
1. Issue確認 → 2. ブランチ作成 → 3. 実装 → 4. テスト → 5. コミット → 6. PR作成
```

### ステップ1: Issue/タスクの確認

```bash
# GitHub IssueまたはJiraのタスクを確認
# 必要な情報:
# - 何を実装するか（機能要件）
# - 受け入れ条件（完了の定義）
# - 設計資料へのリンク
```

### ステップ2: ブランチの作成

```bash
# mainから最新を取得
git checkout main
git pull origin main

# 機能ブランチを作成
git checkout -b feature/issue-123-add-user-search

# ブランチ命名規則:
# - feature/xxx: 新機能
# - fix/xxx: バグ修正
# - refactor/xxx: リファクタリング
# - docs/xxx: ドキュメント更新
# - test/xxx: テスト追加
```

### ステップ3: バックエンド実装

**例: ユーザー検索機能を追加する場合**

#### 3-1. モデルの更新（必要な場合）

```bash
# backend/internal/model/user.go を編集
vim backend/internal/model/user.go
```

```go
// SearchUsersRequest 検索リクエスト
type SearchUsersRequest struct {
    Query      string `form:"q" binding:"max=100"`
    Role       string `form:"role" binding:"omitempty,oneof=admin manager user viewer"`
    Status     string `form:"status" binding:"omitempty,oneof=active inactive suspended"`
    Department string `form:"department" binding:"max=100"`
}
```

#### 3-2. リポジトリ層の実装

```bash
# backend/internal/repository/user.go を編集
vim backend/internal/repository/user.go
```

```go
// Search ユーザー検索
func (r *UserRepository) Search(ctx context.Context, req *model.SearchUsersRequest, params *util.PaginationParams) ([]*model.User, int64, error) {
    var users []*model.User
    var total int64

    query := r.db.WithContext(ctx).Model(&model.User{})

    // 検索条件を追加
    if req.Query != "" {
        query = query.Where("username LIKE ? OR email LIKE ? OR full_name LIKE ?",
            "%"+req.Query+"%", "%"+req.Query+"%", "%"+req.Query+"%")
    }
    if req.Role != "" {
        query = query.Where("role = ?", req.Role)
    }
    if req.Status != "" {
        query = query.Where("status = ?", req.Status)
    }
    if req.Department != "" {
        query = query.Where("department = ?", req.Department)
    }

    query.Count(&total)
    err := query.Offset(params.Offset).Limit(params.PerPage).Find(&users).Error

    return users, total, err
}
```

#### 3-3. サービス層の実装

```bash
vim backend/internal/service/user.go
```

```go
// Search ユーザー検索
func (s *UserService) Search(ctx context.Context, req *model.SearchUsersRequest, params *util.PaginationParams) (*util.PaginatedResponse, error) {
    users, total, err := s.repo.Search(ctx, req, params)
    if err != nil {
        s.logger.Error("Failed to search users", zap.Error(err))
        return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
    }

    responses := make([]*model.UserResponse, len(users))
    for i, user := range users {
        responses[i] = user.ToResponse()
    }

    return util.NewPaginatedResponse(responses, total, params), nil
}
```

#### 3-4. ハンドラー層の実装

```bash
vim backend/internal/handler/user.go
```

```go
// Search ユーザー検索
func (h *UserHandler) Search(c *gin.Context) {
    var req model.SearchUsersRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        util.ValidationError(c, util.ParseValidationErrors(err))
        return
    }

    params := util.GetPaginationParams(c)
    result, err := h.service.Search(c.Request.Context(), &req, params)
    if err != nil {
        util.HandleError(c, err)
        return
    }

    util.Paginated(c, result)
}
```

#### 3-5. ルーティングの追加

```bash
vim backend/cmd/server/main.go
```

```go
func setupUserRoutes(r *gin.Engine, handler *handler.UserHandler) {
    v1 := r.Group("/api/v1")
    {
        users := v1.Group("/users")
        {
            users.GET("", handler.List)
            users.GET("/search", handler.Search) // ➕ 追加
            users.GET("/:id", handler.GetByID)
            // ...
        }
    }
}
```

#### 3-6. 保存して自動リロード確認

```bash
# ファイルを保存すると、Airが自動的に検知してリロード
# ログを確認:
docker-compose logs -f backend

# 期待される出力:
# backend_1  | main.go has changed
# backend_1  | Building...
# backend_1  | Running...
# backend_1  | {"level":"info","ts":...,"msg":"Server starting","port":"8080"}
```

### ステップ4: フロントエンド実装

#### 4-1. API関数の追加

```bash
vim frontend/src/lib/users.ts
```

```typescript
export const usersApi = {
  // 既存の関数...

  async searchUsers(params: {
    q?: string;
    role?: string;
    status?: string;
    department?: string;
    page?: number;
    per_page?: number;
  }): Promise<PaginatedResponse<User>> {
    const response = await api.get<PaginatedResponse<User>>('/users/search', { params });
    return response.data;
  },
};
```

#### 4-2. カスタムフックの追加

```bash
vim frontend/src/hooks/useUsers.ts
```

```typescript
export function useSearchUsers(params: SearchUsersParams) {
  return useQuery({
    queryKey: ['users', 'search', params],
    queryFn: () => usersApi.searchUsers(params),
    enabled: !!params.q || !!params.role || !!params.status || !!params.department,
  });
}
```

#### 4-3. コンポーネントの実装

```bash
vim frontend/src/components/users/UserSearch.tsx
```

```typescript
'use client';

import { useState } from 'react';
import { useSearchUsers } from '@/hooks/useUsers';
import { UserList } from './UserList';

export function UserSearch() {
  const [query, setQuery] = useState('');
  const [role, setRole] = useState('');
  const [status, setStatus] = useState('');

  const { data, isLoading, error } = useSearchUsers({ q: query, role, status });

  return (
    <div className="space-y-4">
      <div className="flex gap-4">
        <input
          type="text"
          placeholder="ユーザー名、メール、氏名で検索"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="flex-1 px-4 py-2 border rounded"
        />
        <select value={role} onChange={(e) => setRole(e.target.value)} className="px-4 py-2 border rounded">
          <option value="">全てのロール</option>
          <option value="admin">管理者</option>
          <option value="manager">マネージャー</option>
          <option value="user">ユーザー</option>
          <option value="viewer">閲覧者</option>
        </select>
        <select value={status} onChange={(e) => setStatus(e.target.value)} className="px-4 py-2 border rounded">
          <option value="">全てのステータス</option>
          <option value="active">アクティブ</option>
          <option value="inactive">非アクティブ</option>
          <option value="suspended">停止中</option>
        </select>
      </div>

      {isLoading && <div>検索中...</div>}
      {error && <div className="text-red-600">エラーが発生しました</div>}
      {data && <UserList users={data.data} />}
    </div>
  );
}
```

#### 4-4. ページに組み込み

```bash
vim frontend/src/app/users/page.tsx
```

```typescript
import { UserSearch } from '@/components/users/UserSearch';

export default function UsersPage() {
  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-6">ユーザー管理</h1>
      <UserSearch />
    </div>
  );
}
```

#### 4-5. 保存して自動リロード確認

```bash
# ブラウザが自動的にリロードされる
# http://localhost:3000/users にアクセスして確認
```

---

## テストのワークフロー

### バックエンドテスト

#### 単体テストの作成

```bash
# テストファイルを作成
vim backend/internal/repository/user_test.go
```

```go
package repository

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    assert.NoError(t, err)
    db.AutoMigrate(&model.User{})
    return db
}

func TestUserRepository_Search(t *testing.T) {
    db := setupTestDB(t)
    repo := NewUserRepository(db)

    // テストデータを挿入
    users := []*model.User{
        {Username: "alice", Email: "alice@example.com", FullName: "Alice Smith", Role: "admin", Status: "active"},
        {Username: "bob", Email: "bob@example.com", FullName: "Bob Jones", Role: "user", Status: "active"},
    }
    for _, u := range users {
        db.Create(u)
    }

    // テスト実行
    req := &model.SearchUsersRequest{Query: "alice"}
    params := &util.PaginationParams{Page: 1, PerPage: 10}
    result, total, err := repo.Search(context.Background(), req, params)

    // 検証
    assert.NoError(t, err)
    assert.Equal(t, int64(1), total)
    assert.Len(t, result, 1)
    assert.Equal(t, "alice", result[0].Username)
}
```

#### テストの実行

```bash
cd backend

# 全テスト実行
make test

# 特定のパッケージのみ
go test ./internal/repository -v

# カバレッジ付き
make test-coverage

# カバレッジレポートを開く
open coverage.html
```

### フロントエンドテスト

#### コンポーネントテストの作成

```bash
vim frontend/src/components/users/UserSearch.test.tsx
```

```typescript
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { UserSearch } from './UserSearch';

// モックデータ
const mockUsers = [
  { id: 1, username: 'alice', email: 'alice@example.com', role: 'admin', status: 'active' },
];

// APIモック
vi.mock('@/hooks/useUsers', () => ({
  useSearchUsers: () => ({
    data: { data: mockUsers },
    isLoading: false,
    error: null,
  }),
}));

describe('UserSearch', () => {
  it('should render search form', () => {
    const queryClient = new QueryClient();
    render(
      <QueryClientProvider client={queryClient}>
        <UserSearch />
      </QueryClientProvider>
    );

    expect(screen.getByPlaceholderText('ユーザー名、メール、氏名で検索')).toBeInTheDocument();
  });

  it('should update query on input change', () => {
    const queryClient = new QueryClient();
    render(
      <QueryClientProvider client={queryClient}>
        <UserSearch />
      </QueryClientProvider>
    );

    const input = screen.getByPlaceholderText('ユーザー名、メール、氏名で検索');
    fireEvent.change(input, { target: { value: 'alice' } });

    expect(input).toHaveValue('alice');
  });
});
```

#### テストの実行

```bash
cd frontend

# 全テスト実行
npm test

# watchモードで実行
npm test -- --watch

# カバレッジ付き
npm test -- --coverage
```

### 統合テスト（E2E）

```bash
# curlでAPIを直接テスト
curl -X GET "http://localhost:8080/api/v1/users/search?q=alice" | jq

# 期待される出力:
# {
#   "code": 200,
#   "message": "success",
#   "data": [
#     {
#       "id": 1,
#       "username": "alice",
#       "email": "alice@example.com",
#       ...
#     }
#   ],
#   "pagination": {
#     "page": 1,
#     "per_page": 10,
#     "total": 1,
#     "total_pages": 1
#   }
# }
```

---

## デバッグのワークフロー

### バックエンドのデバッグ

#### ログの確認

```bash
# リアルタイムでログを表示
docker-compose logs -f backend

# 最新100行を表示
docker-compose logs --tail=100 backend

# エラーだけをフィルタ
docker-compose logs backend | grep ERROR
```

#### デバッグログの追加

```go
import "go.uber.org/zap"

func (s *UserService) Search(ctx context.Context, req *model.SearchUsersRequest, params *util.PaginationParams) (*util.PaginatedResponse, error) {
    // デバッグログを追加
    s.logger.Debug("Search users",
        zap.String("query", req.Query),
        zap.String("role", req.Role),
        zap.String("status", req.Status),
        zap.Int("page", params.Page),
    )

    // ... 処理
}
```

#### 環境変数でログレベルを変更

```bash
# .env ファイルを編集
LOG_LEVEL=debug  # info → debug に変更

# サービスを再起動
docker-compose restart backend
```

#### データベースクエリのデバッグ

```bash
# PostgreSQLに直接接続
docker-compose exec postgres psql -U postgres -d effisio_dev

# クエリを実行
SELECT * FROM users WHERE username LIKE '%alice%';

# EXPLAIN でクエリプランを確認
EXPLAIN ANALYZE SELECT * FROM users WHERE username LIKE '%alice%';

# 終了
\q
```

### フロントエンドのデバッグ

#### React Query DevToolsの使用

ブラウザで http://localhost:3000 を開くと、画面下部に React Query DevTools が表示されます。

- クエリの状態を確認
- キャッシュの内容を確認
- 手動でクエリを再実行

#### ブラウザコンソールでのデバッグ

```javascript
// ブラウザのコンソールを開く（F12）

// ローカルストレージの確認
console.log(localStorage.getItem('token'));

// APIリクエストの確認（Networkタブ）
// Filter: XHR を選択してAJAXリクエストのみ表示
```

#### Next.jsのソースマップ

```bash
# next.config.js で sourceMap を有効化（開発環境ではデフォルトで有効）
const nextConfig = {
  productionBrowserSourceMaps: true, // プロダクションでも有効化
};
```

---

## コードレビューのワークフロー

### Pull Requestの作成

```bash
# 変更をコミット
git add .
git commit -m "feat: Add user search functionality

- Add SearchUsersRequest model
- Implement Search method in repository, service, handler
- Add UserSearch component
- Add search API hook"

# プッシュ
git push origin feature/issue-123-add-user-search

# GitHubでPRを作成
# タイトル: [Issue #123] Add user search functionality
# 説明:
# ## 変更内容
# - ユーザー検索機能を追加
# - 検索条件: ユーザー名、メール、氏名、ロール、ステータス、部署
#
# ## テスト
# - [ ] ユニットテスト追加
# - [ ] 手動テスト完了
#
# ## スクリーンショット
# （スクリーンショットを添付）
```

### レビュー時のチェックポイント

**レビュアー向け:**

- [ ] コーディング規約に準拠しているか
- [ ] テストが十分にカバーされているか
- [ ] エラーハンドリングが適切か
- [ ] セキュリティ上の問題がないか（SQLインジェクション、XSSなど）
- [ ] パフォーマンスへの影響はないか
- [ ] ドキュメントが更新されているか

**レビュイー向け:**

- レビューコメントには24時間以内に返信
- 修正を加えたら「Fixed in commit abc1234」とコメント
- 議論が必要な場合は直接会話する

---

## コミット前のチェックリスト

### 自己チェック

```bash
# 1. リンター実行
cd backend && make lint
cd ../frontend && npm run lint

# 2. フォーマッター実行
cd backend && gofmt -w .
cd ../frontend && npm run format

# 3. テスト実行
cd backend && make test
cd ../frontend && npm test

# 4. ビルド確認
cd backend && make build
cd ../frontend && npm run build

# 5. 型チェック（フロントエンド）
cd frontend && npm run type-check
```

### コミットメッセージの書き方

**Conventional Commits形式:**

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Type:**
- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメントのみの変更
- `style`: コードの意味に影響しない変更（フォーマットなど）
- `refactor`: リファクタリング
- `test`: テストの追加・修正
- `chore`: ビルドプロセスやツールの変更

**例:**

```bash
git commit -m "feat(user): Add search functionality

- Add SearchUsersRequest model
- Implement Search in repository/service/handler layers
- Add UserSearch React component
- Add unit tests for search

Closes #123"
```

---

## 終業時のルーチン

### 1. 作業中の変更をコミット/スタッシュ

```bash
# 完成していない場合はスタッシュ
git stash save "WIP: user search implementation"

# または一旦コミット（後でamendする）
git add .
git commit -m "WIP: user search (incomplete)"
```

### 2. リモートにプッシュ

```bash
# 作業内容をバックアップとしてプッシュ
git push origin feature/your-feature-name
```

### 3. 開発環境の停止

```bash
# Dockerコンテナを停止（データは保持）
docker-compose down

# またはコンテナを残したまま終了
# → 次回の起動が速い
```

### 4. タスク管理の更新

- JiraやGitHub Issueのステータスを更新
- 進捗をチームに共有（Slackなど）

---

## よく使うコマンド集

### Docker関連

```bash
# 全サービス起動
docker-compose up -d

# 特定サービスのみ起動
docker-compose up -d postgres redis

# 全サービス停止
docker-compose down

# ボリュームも削除（完全リセット）
docker-compose down -v

# ログ確認
docker-compose logs -f
docker-compose logs -f backend
docker-compose logs --tail=100 backend

# サービス再起動
docker-compose restart backend

# コンテナに入る
docker-compose exec backend sh
docker-compose exec postgres psql -U postgres -d effisio_dev
```

### Git関連

```bash
# ブランチ一覧
git branch -a

# ブランチ切り替え
git checkout feature/xxx

# 新ブランチ作成
git checkout -b feature/xxx

# 変更の確認
git status
git diff
git diff --staged

# コミット
git add .
git commit -m "message"

# プッシュ
git push origin feature/xxx

# 最新を取得
git pull origin main

# スタッシュ
git stash save "message"
git stash list
git stash pop
git stash apply stash@{0}
```

### Make関連

```bash
# プロジェクトルート
make setup        # 初回セットアップ
make dev          # 開発環境起動
make test         # 全テスト実行
make lint         # 全リンター実行
make clean        # クリーンアップ
make build        # プロダクションビルド
make migrate-up   # マイグレーション実行
make migrate-down # マイグレーションロールバック
make seed         # シードデータ投入

# backend/
cd backend
make build        # バイナリビルド
make run          # ビルドして実行
make dev          # ホットリロードで実行
make test         # テスト実行
make test-coverage # カバレッジ付きテスト
make lint         # リンター実行

# frontend/
cd frontend
npm run dev       # 開発サーバー起動
npm run build     # プロダクションビルド
npm test          # テスト実行
npm run lint      # リンター実行
npm run format    # フォーマッター実行
npm run type-check # 型チェック
```

### データベース関連

```bash
# PostgreSQLに接続
docker-compose exec postgres psql -U postgres -d effisio_dev

# よく使うSQLコマンド
\l                 # データベース一覧
\c effisio_dev     # データベース切り替え
\dt                # テーブル一覧
\d users           # テーブル定義表示
\q                 # 終了

# データ確認
SELECT * FROM users;
SELECT COUNT(*) FROM users;
SELECT * FROM users WHERE role = 'admin';

# データ挿入
INSERT INTO users (username, email, password_hash, role, status)
VALUES ('test', 'test@example.com', 'hash', 'user', 'active');

# データ更新
UPDATE users SET status = 'inactive' WHERE id = 1;

# データ削除
DELETE FROM users WHERE id = 1;
```

### API確認（curl）

```bash
# Ping
curl http://localhost:8080/api/v1/ping

# ユーザー一覧
curl http://localhost:8080/api/v1/users | jq

# ユーザー検索
curl "http://localhost:8080/api/v1/users/search?q=alice" | jq

# ユーザー詳細
curl http://localhost:8080/api/v1/users/1 | jq

# ユーザー作成
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "email": "newuser@example.com",
    "password": "password123",
    "role": "user"
  }' | jq

# ユーザー更新
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newemail@example.com"
  }' | jq

# ユーザー削除
curl -X DELETE http://localhost:8080/api/v1/users/1
```

---

## トラブル時の対処

問題が発生した場合は、**[TROUBLESHOOTING.md](TROUBLESHOOTING.md)** を参照してください。

---

## まとめ

このワークフローに従うことで、効率的かつ一貫性のある開発作業が可能になります。

**重要なポイント:**

1. **朝のルーチン**: 最新コード取得 → 環境起動 → 動作確認
2. **開発サイクル**: ブランチ作成 → 実装 → テスト → コミット → PR
3. **テスト**: 単体テスト → 統合テスト → E2Eテスト
4. **コミット前**: リンター → フォーマッター → テスト → ビルド
5. **終業時**: コミット/スタッシュ → プッシュ → 環境停止

**次のステップ:**

実装中に問題が発生したら → **[TROUBLESHOOTING.md](TROUBLESHOOTING.md)**
Phase 1の具体的な実装手順 → **[PHASE1_IMPLEMENTATION_STEPS.md](PHASE1_IMPLEMENTATION_STEPS.md)**

---

**最終更新**: 2025-01-20
