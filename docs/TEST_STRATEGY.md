# テスト戦略ドキュメント

このドキュメントでは、Effisioプロジェクトにおけるテスト戦略を定義します。

## テストピラミッド

```
        /\
       /  \     E2E Tests (10%)
      /    \    
     /------\   Integration Tests (20%)
    /        \  
   /          \ Unit Tests (70%)
  /__________\
```

## テストの種類

### 1. ユニットテスト (70%)

**目的**: 個別の関数・メソッドの動作を検証

**対象**:
- ビジネスロジック（service層）
- ユーティリティ関数
- バリデーション
- モデルのメソッド

**ツール**:
- **バックエンド**: Go標準 `testing` + `testify`
- **フロントエンド**: Vitest + React Testing Library

**カバレッジ目標**: 70%以上

**例（Go）**:
```go
func TestCreateUser_Success(t *testing.T) {
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)
    
    user, err := service.CreateUser(&CreateUserRequest{
        Username: "test",
        Email: "test@example.com",
    })
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "test", user.Username)
}
```

**例（TypeScript）**:
```typescript
describe('Button', () => {
  it('renders with children', () => {
    render(<Button>Click me</Button>);
    expect(screen.getByText('Click me')).toBeInTheDocument();
  });
});
```

### 2. 統合テスト (20%)

**目的**: 複数のコンポーネント間の連携を検証

**対象**:
- API エンドポイント
- データベース操作
- 外部サービス連携

**ツール**:
- **バックエンド**: `testing` + `httptest` + テスト用DB
- **フロントエンド**: Vitest + MSW (Mock Service Worker)

**例（Go）**:
```go
func TestLoginAPI(t *testing.T) {
    router := setupRouter()
    req := httptest.NewRequest("POST", "/api/v1/auth/login", 
        bytes.NewBufferString(`{"email":"test@example.com","password":"password"}`))
    w := httptest.NewRecorder()
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
}
```

### 3. E2Eテスト (10%)

**目的**: ユーザーシナリオ全体を検証

**対象**:
- ユーザー登録〜ログイン〜操作のフロー
- 主要な業務フロー

**ツール**:
- **Cypress** または **Playwright**

**例**:
```typescript
describe('Login Flow', () => {
  it('allows user to login', () => {
    cy.visit('/login');
    cy.get('input[type="email"]').type('test@example.com');
    cy.get('input[type="password"]').type('password');
    cy.get('button').contains('Login').click();
    cy.url().should('include', '/dashboard');
  });
});
```

## テストデータ管理

### Fixturesの使用

```go
// tests/fixtures/users.json
[
  {
    "username": "test_user_1",
    "email": "test1@example.com",
    "password": "password"
  }
]
```

### テスト用DB

- **ローカル**: Docker Compose で起動
- **CI**: GitHub Actions の services で起動

## テスト実行

### ローカル

```bash
# バックエンド
cd backend
go test ./... -v

# カバレッジ
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# フロントエンド
cd frontend
npm run test
npm run test:coverage
```

### CI/CD

GitHub Actions で自動実行（PR作成時・マージ時）

## ベストプラクティス

1. **AAA パターン**: Arrange, Act, Assert
2. **1テスト1アサーション**: テストは明確な目的を持つ
3. **テストの独立性**: テスト間に依存関係を作らない
4. **テストデータの分離**: 本番データを使わない
5. **モックの適切な使用**: 外部依存を分離

## 除外対象

以下はテストの優先度が低い：
- 自動生成コード
- 設定ファイル
- 簡単なgetter/setter
- サードパーティライブラリのラッパー

---

詳細は各プロジェクトのテストコードを参照してください。
