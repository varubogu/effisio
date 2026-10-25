# トラブルシューティングガイド

このドキュメントでは、Effisioプロジェクトでよく遭遇する問題とその解決方法を説明します。

## 目次

- [環境セットアップの問題](#環境セットアップの問題)
- [Docker/コンテナの問題](#dockerコンテナの問題)
- [データベースの問題](#データベースの問題)
- [バックエンドの問題](#バックエンドの問題)
- [フロントエンドの問題](#フロントエンドの問題)
- [ネットワーク/APIの問題](#ネットワークapiの問題)
- [パフォーマンスの問題](#パフォーマンスの問題)
- [ビルド/デプロイの問題](#ビルドデプロイの問題)
- [完全リセット手順](#完全リセット手順)

---

## 環境セットアップの問題

### 問題1: `make setup` が失敗する

**症状:**
```
make: *** No rule to make target 'setup'. Stop.
```

**原因:**
- Makefileが存在しない、または壊れている
- 間違ったディレクトリで実行している

**解決方法:**
```bash
# プロジェクトルートにいることを確認
pwd
# 期待される出力: /path/to/effisio

# Makefileの存在確認
ls -la Makefile

# Makefileが存在しない場合、Gitから再取得
git checkout Makefile
```

---

### 問題2: Goのバージョンが古い

**症状:**
```
go: go.mod file indicates go 1.21, but using go version go1.20
```

**原因:**
Go 1.21以上が必要だが、古いバージョンがインストールされている

**解決方法:**

**macOS:**
```bash
# Homebrewでアップデート
brew update
brew upgrade go

# バージョン確認
go version
# 期待される出力: go version go1.21.x darwin/amd64
```

**Linux:**
```bash
# 公式サイトから最新版をダウンロード
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz

# 既存のGoを削除
sudo rm -rf /usr/local/go

# 新しいバージョンをインストール
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz

# バージョン確認
go version
```

**Windows:**
- https://go.dev/dl/ から最新のインストーラーをダウンロード
- インストーラーを実行

---

### 問題3: Node.jsのバージョンが古い

**症状:**
```
Error: The engine "node" is incompatible with this module. Expected version ">=22.23.4 <23.0.0".
```

**原因:**
Node.js 22.23.4 以上（22 系）が必要だが、古いバージョンがインストールされている。Node 20 以前には CVE-2026-111497 の修正がない。

**解決方法:**

**macOS/Linux (nvmを使用):**
```bash
# nvmをインストール（まだの場合）
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash

# 修正済みの 22 系をインストール
nvm install 22.23.4
nvm use 22.23.4
nvm alias default 22.23.4

# バージョン確認
node --version
# 期待される出力: v22.23.4
```

**Windows:**
- https://nodejs.org/dist/v22.23.4/ から 22.23.4 をインストールする
- 22 系の新しい修正版でもよい。24 系は Next.js 14 では使わない

---

### 問題4: Dockerが起動していない

**症状:**
```
Cannot connect to the Docker daemon at unix:///var/run/docker.sock
```

**原因:**
Dockerデーモンが起動していない

**解決方法:**

**macOS:**
```bash
# Docker Desktopを起動
open -a Docker

# 起動を待つ（約30秒）
sleep 30

# 確認
docker ps
```

**Linux:**
```bash
# Dockerサービスを起動
sudo systemctl start docker

# 自動起動を有効化
sudo systemctl enable docker

# 確認
docker ps
```

**Windows:**
- Docker Desktopを起動
- タスクバーのDockerアイコンが緑色になるまで待つ

---

## Docker/コンテナの問題

### 問題5: ポート番号が既に使用されている

**症状:**
```
Error starting userland proxy: listen tcp 0.0.0.0:3000: bind: address already in use
```

**原因:**
指定されたポート（3000, 8080, 5432など）が他のプロセスで使用されている

**解決方法:**

**使用中のポートを確認:**
```bash
# macOS/Linux
lsof -i :3000
lsof -i :8080
lsof -i :5432

# 期待される出力:
# COMMAND   PID   USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
# node    12345  user   23u  IPv4 0x...      0t0  TCP *:3000 (LISTEN)
```

**プロセスを終了:**
```bash
# PIDを確認して終了
kill -9 12345

# または、全てのNode.jsプロセスを終了
killall node
```

**または、ポート番号を変更:**
```bash
# docker-compose.yml を編集
vim docker-compose.yml

# ポート番号を変更
services:
  frontend:
    ports:
      - "3001:3000"  # 3000 → 3001 に変更
```

---

### 問題6: Dockerコンテナが起動しない

**症状:**
```
ERROR: for backend  Container "xxx" is unhealthy
```

**原因:**
コンテナ内のアプリケーションが正常に起動していない

**解決方法:**

**ログを確認:**
```bash
# 全コンテナのログを確認
docker-compose logs

# 特定のコンテナのログを確認
docker-compose logs backend

# リアルタイムでログを確認
docker-compose logs -f backend
```

**コンテナを再起動:**
```bash
# 特定のコンテナを再起動
docker-compose restart backend

# 全コンテナを再起動
docker-compose restart

# 完全に再構築
docker-compose down
docker-compose up -d --build
```

---

### 問題7: Dockerイメージのビルドが失敗する

**症状:**
```
ERROR [backend 5/8] RUN go build -o bin/server cmd/server/main.go
```

**原因:**
- Goのコンパイルエラー
- 依存関係の問題

**解決方法:**

**キャッシュをクリアして再ビルド:**
```bash
# キャッシュを使わずにビルド
docker-compose build --no-cache backend

# または全体を再ビルド
docker-compose down
docker system prune -a
docker-compose up -d --build
```

**ローカルでビルドして確認:**
```bash
cd backend
go mod tidy
go build -o bin/server cmd/server/main.go

# エラーメッセージを確認して修正
```

---

### 問題8: Dockerボリュームの権限エラー

**症状:**
```
ERROR: for postgres  Cannot start service postgres:
  OCI runtime create failed: container_linux.go:380:
  starting container process caused: process_linux.go:545:
  container init caused: rootfs_stat(/var/lib/postgresql/data): permission denied
```

**原因:**
Dockerボリュームのファイル権限が不正

**解決方法:**
```bash
# 全コンテナとボリュームを削除
docker-compose down -v

# Docker volumeを完全削除
docker volume prune -a

# 再起動
docker-compose up -d

# マイグレーション実行
make migrate-up
make seed
```

---

## データベースの問題

### 問題9: データベースに接続できない

**症状:**
```
Error: failed to connect to postgres: dial tcp 127.0.0.1:5432: connect: connection refused
```

**原因:**
- PostgreSQLコンテナが起動していない
- 接続情報が間違っている

**解決方法:**

**コンテナの状態を確認:**
```bash
# PostgreSQLコンテナが起動しているか確認
docker-compose ps postgres

# 期待される出力:
#       Name                     Command              State           Ports
# ---------------------------------------------------------------------------------
# effisio_postgres_1   docker-entrypoint.sh postgres   Up      0.0.0.0:5432->5432/tcp
```

**接続情報を確認:**
```bash
# backend/.env を確認
cat backend/.env | grep DB_

# 期待される出力:
# DB_HOST=postgres
# DB_PORT=5432
# DB_USER=postgres
# DB_PASSWORD=postgres
# DB_NAME=effisio_dev
```

**psqlで直接接続してみる:**
```bash
docker-compose exec postgres psql -U postgres -d effisio_dev

# 接続できたら成功
# effisio_dev=#
```

---

### 問題10: マイグレーションが失敗する

**症状:**
```
error: Dirty database version 1. Fix and force version.
```

**原因:**
マイグレーションが途中で失敗し、データベースが不整合な状態になっている

**解決方法:**

**方法1: マイグレーションをリセット**
```bash
# 全マイグレーションをロールバック
make migrate-down

# 再度マイグレーション
make migrate-up
```

**方法2: 強制的にバージョンを設定**
```bash
# 現在のマイグレーションバージョンを確認
docker-compose exec postgres psql -U postgres -d effisio_dev -c "SELECT * FROM schema_migrations;"

# バージョンを強制設定（例: バージョン1）
migrate -path backend/migrations \
  -database "postgresql://postgres:postgres@localhost:5432/effisio_dev?sslmode=disable" \
  force 1

# 再度マイグレーション
make migrate-up
```

**方法3: データベースを完全にリセット**
```bash
# ⚠️ 警告: 全データが削除されます
docker-compose down -v
docker-compose up -d postgres
sleep 5
make migrate-up
make seed
```

---

### 問題11: シードデータが投入されない

**症状:**
```
🌱 シードデータを投入しています...
ERROR:  relation "users" does not exist
```

**原因:**
マイグレーションが実行されていない

**解決方法:**
```bash
# マイグレーションを先に実行
make migrate-up

# その後シード投入
make seed
```

---

### 問題12: データベースのパフォーマンスが遅い

**症状:**
クエリの実行に時間がかかる

**原因:**
- インデックスが不足
- N+1クエリ問題
- 不適切なクエリ

**解決方法:**

**EXPLAIN ANALYZEでクエリを分析:**
```sql
EXPLAIN ANALYZE SELECT * FROM users WHERE username LIKE '%alice%';

-- Seq Scan が表示されたらインデックスが使われていない
-- Index Scan が表示されたらインデックスが使われている
```

**インデックスを追加:**
```sql
-- 部分一致検索にはGINインデックスが有効
CREATE INDEX idx_users_username_gin ON users USING gin(username gin_trgm_ops);
CREATE INDEX idx_users_email_gin ON users USING gin(email gin_trgm_ops);

-- pg_trgm拡張が必要
CREATE EXTENSION IF NOT EXISTS pg_trgm;
```

**N+1クエリを修正（Goコード）:**
```go
// ❌ Bad: N+1 query problem
users, _ := repo.FindAll()
for _, user := range users {
    tasks, _ := taskRepo.FindByUserID(user.ID) // N回のクエリ
}

// ✅ Good: Preload
var users []User
db.Preload("Tasks").Find(&users) // 1回のクエリ
```

---

## バックエンドの問題

### 問題13: `go mod download` が失敗する

**症状:**
```
go: github.com/gin-gonic/gin@v1.9.1: Get "https://proxy.golang.org/...": dial tcp: i/o timeout
```

**原因:**
- ネットワークの問題
- Goプロキシの問題

**解決方法:**

**Goプロキシを変更:**
```bash
# 中国のプロキシを使用
export GOPROXY=https://goproxy.cn,direct

# または日本のプロキシを使用
export GOPROXY=https://goproxy.io,direct

# 再度ダウンロード
go mod download
```

**キャッシュをクリア:**
```bash
go clean -modcache
go mod download
```

---

### 問題14: ホットリロード（Air）が動かない

**症状:**
コードを変更しても自動的にリロードされない

**原因:**
- Airがインストールされていない
- .air.toml の設定が間違っている

**解決方法:**

**Airを手動インストール:**
```bash
go install github.com/cosmtrek/air@latest

# PATHを確認
echo $GOPATH/bin
# このパスが$PATHに含まれているか確認

# 含まれていない場合は追加（~/.zshrc または ~/.bashrc）
export PATH=$PATH:$(go env GOPATH)/bin
source ~/.zshrc
```

**Airの設定を確認:**
```bash
# backend/.air.toml を確認
cat backend/.air.toml

# 正しい設定:
# [build]
#   cmd = "go build -o ./bin/server ./cmd/server"
#   bin = "bin/server"
```

**Dockerコンテナを再起動:**
```bash
docker-compose restart backend
docker-compose logs -f backend

# 期待される出力:
# backend_1  | Running...
```

---

### 問題15: bcrypt のハッシュ化が遅い

**症状:**
ユーザー登録に1秒以上かかる

**原因:**
bcryptのcostが高すぎる

**解決方法:**

**costを調整:**
```go
// backend/internal/service/user.go

// ❌ Bad: cost=14 は非常に遅い
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

// ✅ Good: cost=10 がバランスが良い（デフォルト）
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
```

**注意:** costを下げすぎるとセキュリティが低下します。10-12が推奨値です。

---

### 問題16: GORMのエラーが分かりにくい

**症状:**
```
Error 1062: Duplicate entry 'alice' for key 'users.username'
```

**原因:**
GORMのエラーがそのまま返されている

**解決方法:**

**カスタムエラーハンドリングを追加:**
```go
// backend/internal/repository/user.go

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
    err := r.db.WithContext(ctx).Create(user).Error
    if err != nil {
        // MySQLのエラーコードをチェック
        if strings.Contains(err.Error(), "Duplicate entry") {
            if strings.Contains(err.Error(), "username") {
                return errors.New("username already exists")
            }
            if strings.Contains(err.Error(), "email") {
                return errors.New("email already exists")
            }
        }
        return err
    }
    return nil
}
```

---

## フロントエンドの問題

### 問題17: `npm install` が失敗する

**症状:**
```
npm ERR! code ERESOLVE
npm ERR! ERESOLVE unable to resolve dependency tree
```

**原因:**
依存関係の競合

**解決方法:**

**方法1: legacy peer depsを使用:**
```bash
cd frontend
npm install --legacy-peer-deps
```

**方法2: キャッシュをクリアして再インストール:**
```bash
cd frontend

# キャッシュをクリア
npm cache clean --force

# node_modulesとpackage-lock.jsonを削除
rm -rf node_modules package-lock.json

# 再インストール
npm install
```

**方法3: Node.jsのバージョンを確認:**
```bash
node --version
# v18以上であることを確認
```

---

### 問題18: TypeScriptの型エラー

**症状:**
```
Type 'string | undefined' is not assignable to type 'string'.
```

**原因:**
strictモードで undefined が許可されていない

**解決方法:**

**方法1: optional chaining と nullish coalescing を使用:**
```typescript
// ❌ Bad
const name: string = user.full_name; // full_name は string | undefined

// ✅ Good
const name: string = user.full_name ?? ''; // デフォルト値を設定
const name: string | undefined = user.full_name; // 型を明示
```

**方法2: 型ガードを使用:**
```typescript
if (user.full_name) {
    const name: string = user.full_name; // この中では string
}
```

---

### 問題19: Next.jsのビルドが失敗する

**症状:**
```
Error: Build failed because of webpack errors
```

**原因:**
- TypeScriptエラー
- 未使用のimport
- 構文エラー

**解決方法:**

**型チェックを実行:**
```bash
cd frontend
npm run type-check

# エラーメッセージを確認して修正
```

**リンターを実行:**
```bash
npm run lint

# 自動修正
npm run lint -- --fix
```

**キャッシュをクリアして再ビルド:**
```bash
rm -rf .next
npm run build
```

---

### 問題20: React Queryのキャッシュが更新されない

**症状:**
データを更新してもUIに反映されない

**原因:**
ミューテーション後にキャッシュを無効化していない

**解決方法:**

**onSuccessでキャッシュを無効化:**
```typescript
// ❌ Bad
export function useCreateUser() {
  return useMutation({
    mutationFn: usersApi.createUser,
    // キャッシュが更新されない
  });
}

// ✅ Good
export function useCreateUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: usersApi.createUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] }); // キャッシュを無効化
    },
  });
}
```

---

### 問題21: CORSエラー

**症状:**
```
Access to XMLHttpRequest at 'http://localhost:8080/api/v1/users' from origin 'http://localhost:3000'
has been blocked by CORS policy
```

**原因:**
バックエンドでCORSが正しく設定されていない

**解決方法:**

**バックエンドのCORS設定を確認:**
```go
// backend/internal/middleware/cors.go

func CORS() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"}, // フロントエンドのURL
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    })
}
```

**開発環境では全てのオリジンを許可:**
```go
AllowOrigins: []string{"*"},  // 開発環境のみ
```

---

## ネットワーク/APIの問題

### 問題22: APIが404エラーを返す

**症状:**
```
GET http://localhost:8080/api/v1/users 404 Not Found
```

**原因:**
- ルーティングが正しく設定されていない
- バックエンドが起動していない

**解決方法:**

**バックエンドが起動しているか確認:**
```bash
docker-compose ps backend

# 期待される出力: State が Up
```

**ルーティングを確認:**
```bash
# backend/cmd/server/main.go を確認
grep -n "GET.*users" backend/cmd/server/main.go

# 期待される出力:
# users.GET("", handler.List)
```

**curlで直接アクセスしてみる:**
```bash
curl http://localhost:8080/api/v1/users

# エラーメッセージを確認
```

---

### 問題23: APIが500エラーを返す

**症状:**
```
GET http://localhost:8080/api/v1/users 500 Internal Server Error
```

**原因:**
バックエンドでエラーが発生している

**解決方法:**

**バックエンドのログを確認:**
```bash
docker-compose logs -f backend

# エラースタックトレースを確認
```

**よくある原因:**
1. データベース接続エラー
2. nil pointer dereference
3. 型変換エラー
4. バリデーションエラー

**デバッグログを追加:**
```go
s.logger.Debug("Processing request",
    zap.Any("request", req),
    zap.String("user_id", userID),
)
```

---

### 問題24: APIのレスポンスが遅い

**症状:**
APIリクエストに3秒以上かかる

**原因:**
- N+1クエリ問題
- インデックスの不足
- 不要なデータの取得

**解決方法:**

**ログで実行時間を計測:**
```go
start := time.Now()
defer func() {
    s.logger.Info("Request completed",
        zap.Duration("duration", time.Since(start)),
    )
}()
```

**データベースクエリを最適化:**
```go
// ❌ Bad: 全カラム取得
db.Find(&users)

// ✅ Good: 必要なカラムのみ
db.Select("id", "username", "email", "role", "status").Find(&users)
```

**ページネーションを実装:**
```go
// ❌ Bad: 全件取得
var users []User
db.Find(&users)

// ✅ Good: ページネーション
var users []User
db.Offset(offset).Limit(perPage).Find(&users)
```

---

## パフォーマンスの問題

### 問題25: フロントエンドの初回ロードが遅い

**症状:**
ページの初回読み込みに5秒以上かかる

**原因:**
バンドルサイズが大きすぎる

**解決方法:**

**バンドルサイズを分析:**
```bash
cd frontend

# バンドル分析ツールをインストール
npm install --save-dev @next/bundle-analyzer

# next.config.js に追加
const withBundleAnalyzer = require('@next/bundle-analyzer')({
  enabled: process.env.ANALYZE === 'true',
})

module.exports = withBundleAnalyzer(nextConfig)

# 分析実行
ANALYZE=true npm run build
```

**解決策:**
1. Dynamic importを使用
2. 不要なライブラリを削除
3. Tree-shakingを活用
4. コード分割を実装

```typescript
// ❌ Bad: 全部インポート
import { Button, Table, Modal, Form } from 'antd';

// ✅ Good: 個別にインポート
import Button from 'antd/lib/button';
import Table from 'antd/lib/table';
```

---

### 問題26: メモリリークが発生する

**症状:**
長時間使用するとブラウザが遅くなる

**原因:**
React hooksのクリーンアップ不足

**解決方法:**

**useEffect のクリーンアップを実装:**
```typescript
// ❌ Bad: クリーンアップなし
useEffect(() => {
  const interval = setInterval(() => {
    fetchData();
  }, 5000);
}, []);

// ✅ Good: クリーンアップあり
useEffect(() => {
  const interval = setInterval(() => {
    fetchData();
  }, 5000);

  return () => clearInterval(interval); // クリーンアップ
}, []);
```

**React Query のキャッシュ設定を調整:**
```typescript
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000,     // 5分
      cacheTime: 10 * 60 * 1000,    // 10分
      refetchOnWindowFocus: false,  // 不要な再フェッチを防ぐ
    },
  },
});
```

---

## ビルド/デプロイの問題

### 問題27: Dockerイメージのサイズが大きい

**症状:**
Dockerイメージが1GB以上

**原因:**
マルチステージビルドを使用していない

**解決方法:**

**Dockerfileをマルチステージビルドに変更:**
```dockerfile
# ❌ Bad: 全部入り（1.5GB）
FROM golang:1.21
WORKDIR /app
COPY . .
RUN go build -o server cmd/server/main.go
CMD ["./server"]

# ✅ Good: マルチステージ（50MB）
FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
CMD ["./server"]
```

---

### 問題28: プロダクションビルドが失敗する

**症状:**
```
npm run build
> next build
Error: Minified React error #...
```

**原因:**
環境変数が設定されていない

**解決方法:**

**.env.productionを作成:**
```bash
cd frontend

# .env.production を作成
cat > .env.production <<EOF
NEXT_PUBLIC_API_URL=https://api.example.com
NODE_ENV=production
EOF

# ビルド実行
npm run build
```

---

## 完全リセット手順

全てがうまくいかない場合の最終手段です。**全てのデータが削除されます。**

### ステップ1: Docker環境を完全削除

```bash
# 全コンテナを停止・削除
docker-compose down -v

# 全未使用リソースを削除
docker system prune -a --volumes

# 確認メッセージに 'y' を入力
# WARNING! This will remove:
#   - all stopped containers
#   - all networks not used by at least one container
#   - all volumes not used by at least one container
#   - all images without at least one container associated to them
# Are you sure you want to continue? [y/N] y
```

### ステップ2: ローカルファイルをクリーンアップ

```bash
# backendをクリーンアップ
cd backend
rm -rf bin/ coverage.out coverage.html
go clean -cache -testcache -modcache
cd ..

# frontendをクリーンアップ
cd frontend
rm -rf .next node_modules package-lock.json .turbo
npm cache clean --force
cd ..

# Gitの未追跡ファイルを削除（注意: カスタムファイルも削除されます）
git clean -fdx
```

### ステップ3: 再セットアップ

```bash
# 初回セットアップ実行
make setup

# 開発環境起動
make dev

# 別のターミナルで:
# マイグレーション実行
make migrate-up

# シードデータ投入
make seed
```

### ステップ4: 動作確認

```bash
# バックエンド疎通確認
curl http://localhost:8080/api/v1/ping

# 期待される出力:
# {"message":"pong"}

# ユーザー一覧取得
curl http://localhost:8080/api/v1/users | jq

# ブラウザで確認
# http://localhost:3000
```

---

## サポート

このガイドで解決しない場合:

1. **ログを確認:**
   ```bash
   docker-compose logs -f
   ```

2. **GitHub Issuesを検索:**
   https://github.com/varubogu/effisio/issues

3. **新しいIssueを作成:**
   - エラーメッセージ全文
   - 実行したコマンド
   - 環境情報（OS、Docker version、Go version、Node.js version）
   - ログファイル

4. **チームに相談:**
   - Slackの #effisio-support チャンネル

---

## よく使うデバッグコマンド

```bash
# Dockerコンテナの状態確認
docker-compose ps

# ログ確認
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f postgres

# コンテナに入る
docker-compose exec backend sh
docker-compose exec frontend sh
docker-compose exec postgres psql -U postgres -d effisio_dev

# リソース使用量確認
docker stats

# ネットワーク確認
docker network ls
docker network inspect effisio_default

# ボリューム確認
docker volume ls
docker volume inspect effisio_postgres_data

# システム全体の状況確認
docker system df
```

---

## まとめ

問題に遭遇したら:

1. **エラーメッセージを読む** - ほとんどの場合、原因が書いてある
2. **ログを確認する** - `docker-compose logs -f`
3. **このガイドで検索** - Ctrl+F でキーワード検索
4. **完全リセット** - 最終手段として環境をリセット

**重要:** トラブルシューティングで問題を解決したら、このドキュメントに追記してチームで共有しましょう。

---

**最終更新**: 2025-01-20
