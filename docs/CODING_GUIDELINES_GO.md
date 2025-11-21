# Go コーディング規約

Effisioプロジェクトにおける Go 言語のコーディング規約です。

## 基本原則

1. **標準に従う**: [Effective Go](https://go.dev/doc/effective_go) を遵守
2. **シンプルに**: 複雑さを避け、読みやすいコードを書く
3. **一貫性**: プロジェクト全体で統一された書き方を維持

---

## 命名規則

### パッケージ名

```go
// ✅ 良い例: 小文字、短く、明確
package user
package auth
package repository

// ❌ 悪い例: 大文字、アンダースコア、長すぎ
package User
package user_management
package userauthenticationmodule
```

### 変数名・関数名

```go
// ✅ 良い例: camelCase、明確な名前
var userName string
var userCount int
func getUserByID(id int) (*User, error)

// ❌ 悪い例: スネークケース、略語の乱用
var user_name string
var usrCnt int
func getUsrByID(id int) (*User, error)
```

### 定数名

```go
// ✅ 良い例: PascalCase または camelCase
const MaxRetryCount = 3
const defaultTimeout = 30 * time.Second

// ❌ 悪い例: ALL_CAPS（C言語スタイル）
const MAX_RETRY_COUNT = 3
```

### インターフェース名

```go
// ✅ 良い例: -er で終わる
type Reader interface{}
type Writer interface{}
type UserRepository interface{}

// ❌ 悪い例: I- で始まる（Java/C#スタイル）
type IUserRepository interface{}
```

---

## ファイル構成

### ディレクトリ構造

```
backend/
├── cmd/
│   └── server/
│       └── main.go           # エントリーポイント
├── internal/                 # 非公開パッケージ
│   ├── config/
│   │   └── config.go
│   ├── models/
│   │   ├── user.go
│   │   └── role.go
│   ├── repository/
│   │   ├── user_repository.go
│   │   └── role_repository.go
│   ├── service/
│   │   ├── auth_service.go
│   │   └── user_service.go
│   ├── handler/
│   │   ├── auth_handler.go
│   │   └── user_handler.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── cors.go
│   └── utils/
│       ├── errors.go
│       └── response.go
├── pkg/                      # 公開パッケージ（他プロジェクトから利用可）
└── tests/
```

### ファイル名

```
// ✅ 良い例: スネークケース
user_repository.go
auth_service.go
user_repository_test.go

// ❌ 悪い例: camelCase
userRepository.go
AuthService.go
```

---

## コメント・ドキュメント

### 関数コメント

```go
// ✅ 良い例: 関数名で始まり、完全な文
// CreateUser はユーザーを新規作成します。
// メールアドレスの重複チェックを行い、パスワードをハッシュ化して保存します。
func CreateUser(req *CreateUserRequest) (*User, error) {
    // ...
}

// ❌ 悪い例: 不完全、関数名がない
// ユーザー作成
func CreateUser(req *CreateUserRequest) (*User, error) {
    // ...
}
```

### パッケージコメント

```go
// ✅ 良い例: パッケージの最初のファイル（doc.go または main.go）に記載
// Package repository はデータアクセス層を提供します。
// このパッケージは GORM を使用してデータベース操作を行います。
package repository
```

### TODO コメント

```go
// TODO(username): リファクタリング必要 - キャッシュ機能を追加
// FIXME: エラーハンドリングを改善
// HACK: 暫定対応 - 将来的に修正が必要
```

---

## エラーハンドリング

### エラーの返却

```go
// ✅ 良い例: エラーを最後の戻り値として返す
func GetUserByID(id int64) (*User, error) {
    user, err := repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    return user, nil
}

// ❌ 悪い例: エラーを握りつぶす
func GetUserByID(id int64) *User {
    user, _ := repo.FindByID(id)  // エラーを無視
    return user
}
```

### エラーのラップ

```go
// ✅ 良い例: fmt.Errorf + %w でエラーをラップ
if err != nil {
    return nil, fmt.Errorf("failed to create user: %w", err)
}

// ✅ 良い例: カスタムエラー型
var ErrUserNotFound = errors.New("user not found")

if user == nil {
    return ErrUserNotFound
}
```

### パニックの使用

```go
// ✅ 良い例: 回復不可能なエラーのみパニック
func init() {
    if config.DBHost == "" {
        panic("DB_HOST is required")
    }
}

// ❌ 悪い例: 通常のエラーでパニック
func GetUser(id int) *User {
    user, err := repo.FindByID(id)
    if err != nil {
        panic(err)  // ダメ！
    }
    return user
}
```

---

## 関数設計

### 引数の数

```go
// ✅ 良い例: 引数が多い場合は構造体にまとめる
type CreateUserRequest struct {
    Username string
    Email    string
    Password string
    FullName string
}

func CreateUser(req *CreateUserRequest) (*User, error) {
    // ...
}

// ❌ 悪い例: 引数が多すぎる
func CreateUser(username, email, password, fullName, department string, roleID int64) (*User, error) {
    // ...
}
```

### 早期リターン

```go
// ✅ 良い例: 早期リターンでネストを減らす
func ProcessUser(user *User) error {
    if user == nil {
        return errors.New("user is nil")
    }

    if !user.IsActive {
        return errors.New("user is not active")
    }

    // メイン処理
    return nil
}

// ❌ 悪い例: 深いネスト
func ProcessUser(user *User) error {
    if user != nil {
        if user.IsActive {
            // メイン処理
            return nil
        } else {
            return errors.New("user is not active")
        }
    }
    return errors.New("user is nil")
}
```

---

## 構造体

### 構造体定義

```go
// ✅ 良い例: フィールドにタグを付与、明確な型
type User struct {
    ID        int64     `json:"id" gorm:"primaryKey"`
    Username  string    `json:"username" gorm:"uniqueIndex;not null"`
    Email     string    `json:"email" gorm:"uniqueIndex;not null"`
    Password  string    `json:"-" gorm:"not null"`  // JSONから除外
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// ✅ 良い例: コンストラクタ関数
func NewUser(username, email, password string) *User {
    return &User{
        Username: username,
        Email:    email,
        Password: password,
    }
}
```

### メソッド

```go
// ✅ 良い例: ポインタレシーバー（変更を伴う場合）
func (u *User) SetPassword(password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashedPassword)
    return nil
}

// ✅ 良い例: 値レシーバー（変更を伴わない場合）
func (u User) IsActive() bool {
    return u.Status == "active"
}
```

---

## インターフェース

### インターフェース定義

```go
// ✅ 良い例: 小さく、目的が明確
type UserRepository interface {
    GetByID(id int64) (*User, error)
    Create(user *User) error
    Update(user *User) error
    Delete(id int64) error
}

// ✅ 良い例: 単一メソッドのインターフェース
type Validator interface {
    Validate() error
}
```

### 依存性注入

```go
// ✅ 良い例: インターフェースに依存
type UserService struct {
    repo UserRepository  // インターフェース
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

// ❌ 悪い例: 具体的な型に依存
type UserService struct {
    repo *UserRepositoryImpl  // 具体的な型
}
```

---

## 並行処理

### goroutine の使用

```go
// ✅ 良い例: context でタイムアウト制御
func ProcessUsers(ctx context.Context, users []*User) error {
    var wg sync.WaitGroup
    errCh := make(chan error, len(users))

    for _, user := range users {
        wg.Add(1)
        go func(u *User) {
            defer wg.Done()
            if err := processUser(ctx, u); err != nil {
                errCh <- err
            }
        }(user)
    }

    wg.Wait()
    close(errCh)

    for err := range errCh {
        if err != nil {
            return err
        }
    }
    return nil
}
```

### channelの使用

```go
// ✅ 良い例: バッファ付きチャネル、適切なクローズ
func producer(ctx context.Context) <-chan int {
    ch := make(chan int, 10)
    go func() {
        defer close(ch)
        for i := 0; i < 100; i++ {
            select {
            case <-ctx.Done():
                return
            case ch <- i:
            }
        }
    }()
    return ch
}
```

---

## テスト

### テスト関数名

```go
// ✅ 良い例: Test<Function>_<Scenario>
func TestCreateUser_Success(t *testing.T) {}
func TestCreateUser_DuplicateEmail(t *testing.T) {}
func TestGetUserByID_NotFound(t *testing.T) {}
```

### テーブル駆動テスト

```go
// ✅ 良い例: テーブル駆動テスト
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "test@example.com", false},
        {"invalid email", "invalid", true},
        {"empty email", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

---

## パフォーマンス

### 文字列連結

```go
// ✅ 良い例: strings.Builder
var sb strings.Builder
for _, s := range strings {
    sb.WriteString(s)
}
result := sb.String()

// ❌ 悪い例: + による連結（遅い）
result := ""
for _, s := range strings {
    result += s
}
```

### スライスの事前割り当て

```go
// ✅ 良い例: 容量を事前に確保
users := make([]*User, 0, 100)

// ❌ 悪い例: 容量なし（再割り当てが頻発）
users := make([]*User, 0)
```

---

## 禁止事項

❌ **グローバル変数**: 可能な限り避ける（設定値を除く）
❌ **init() の乱用**: 副作用を伴う処理は避ける
❌ **panic の乱用**: 通常のエラーは error で返す
❌ **naked return**: 明示的に return を書く
❌ **不要な else**: 早期リターンを使う

---

## 推奨ツール

- **gofmt**: コード整形（必須）
- **golangci-lint**: 静的解析（必須）
- **go vet**: 静的解析（必須）
- **goimports**: import 整理
- **staticcheck**: 高度な静的解析

## チェックコマンド

```bash
# フォーマット
go fmt ./...

# 静的解析
go vet ./...
golangci-lint run ./...

# テスト
go test ./... -v

# カバレッジ
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

詳細は [Effective Go](https://go.dev/doc/effective_go) を参照してください。
