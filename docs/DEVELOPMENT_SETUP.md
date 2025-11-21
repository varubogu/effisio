# 開発環境セットアップガイド

このガイドでは、Effisioプロジェクトの開発環境をセットアップする手順を説明します。

## 目次

1. [前提条件](#前提条件)
2. [セットアップ手順](#セットアップ手順)
3. [OS別の詳細手順](#os別の詳細手順)
4. [動作確認](#動作確認)
5. [トラブルシューティング](#トラブルシューティング)

---

## 前提条件

### 必須ソフトウェア

以下のソフトウェアがインストールされていることを確認してください：

| ソフトウェア | 最小バージョン | 推奨バージョン | 用途 |
|-------------|--------------|--------------|------|
| **Git** | 2.30+ | 最新 | バージョン管理 |
| **Docker** | 20.10+ | 最新 | コンテナ実行 |
| **Docker Compose** | 1.29+ | 2.x | マルチコンテナ管理 |
| **Go** | 1.21+ | 1.21+ | バックエンド開発 |
| **Node.js** | 18.17+ | 20.x LTS | フロントエンド開発 |
| **npm** | 9.0+ | 最新 | パッケージ管理 |

### 推奨開発ツール

| ツール | 用途 |
|--------|------|
| **Visual Studio Code** | エディタ（推奨拡張機能あり） |
| **GoLand** / **WebStorm** | JetBrains IDE |
| **Postman** / **Insomnia** | API テスト |
| **DBeaver** / **TablePlus** | データベース管理 |

### ハードウェア要件

- **CPU**: 4コア以上推奨
- **メモリ**: 8GB以上必須、16GB以上推奨
- **ストレージ**: 20GB以上の空き容量

---

## セットアップ手順

### 1. リポジトリのクローン

```bash
# HTTPSでクローン
git clone https://github.com/yourusername/effisio.git

# SSHでクローン（推奨）
git clone git@github.com:yourusername/effisio.git

# プロジェクトディレクトリに移動
cd effisio
```

### 2. 環境変数の設定

#### バックエンド

```bash
cd backend
cp .env.example .env.local
```

`.env.local` を編集して、必要に応じて設定を変更してください：

```bash
# 最低限の設定例
DB_PASSWORD=your_secure_password
JWT_SECRET=$(openssl rand -base64 64)
REDIS_PASSWORD=your_redis_password
```

#### フロントエンド

```bash
cd ../frontend
cp .env.example .env.local
```

`.env.local` を編集（デフォルトのままで問題なし）：

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

### 3. Docker環境の起動

プロジェクトルートディレクトリで：

```bash
cd ..

# 基本起動（バックエンド + フロントエンド + DB + Redis）
docker-compose up -d

# 管理ツール付きで起動（Adminer, Redis Commander, Mailhog）
docker-compose --profile tools up -d

# Nginx付きで起動
docker-compose --profile with-nginx up -d
```

初回起動時は、Dockerイメージのビルドとダウンロードに時間がかかります（5-10分程度）。

### 4. データベースのマイグレーション

```bash
# マイグレーションツールのインストール（初回のみ）
brew install golang-migrate  # macOS
# または
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# マイグレーション実行
cd backend
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/effisio_dev?sslmode=disable" up
```

### 5. 依存関係のインストール

#### バックエンド（ローカル開発する場合）

```bash
cd backend
go mod download
go mod verify
```

#### フロントエンド（ローカル開発する場合）

```bash
cd frontend
npm ci
```

### 6. サービスの起動確認

```bash
# 全サービスのステータス確認
docker-compose ps

# ログ確認
docker-compose logs -f

# 特定のサービスのログ確認
docker-compose logs -f backend
docker-compose logs -f frontend
```

すべてのサービスが `Up` 状態になっていることを確認してください。

### 7. アクセス確認

ブラウザで以下のURLにアクセスして動作を確認：

| サービス | URL | 説明 |
|---------|-----|------|
| フロントエンド | http://localhost:3000 | Next.js アプリケーション |
| バックエンドAPI | http://localhost:8080/api/v1 | Go API サーバー |
| Swagger UI | http://localhost:8080/swagger/index.html | API ドキュメント |
| ヘルスチェック | http://localhost:8080/health | API ヘルスチェック |
| Adminer | http://localhost:8081 | データベース管理ツール |
| Redis Commander | http://localhost:8082 | Redis 管理ツール |
| Mailhog | http://localhost:8025 | メール送信テスト |
| Prometheus | http://localhost:9090 | メトリクス収集 |

---

## OS別の詳細手順

### macOS

#### 1. Homebrewのインストール

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

#### 2. 必要なソフトウェアのインストール

```bash
# Git
brew install git

# Docker Desktop
brew install --cask docker

# Go
brew install go@1.21

# Node.js
brew install node@20

# golang-migrate
brew install golang-migrate

# その他の便利ツール
brew install make jq
```

#### 3. Docker Desktopの起動

アプリケーションフォルダから Docker Desktop を起動し、完全に起動するまで待ちます。

#### 4. セットアップ続行

上記の「セットアップ手順」に従ってください。

### Windows

#### 1. WSL2のセットアップ（推奨）

Windows環境では WSL2 + Ubuntu を使用することを強く推奨します。

```powershell
# PowerShellを管理者権限で実行
wsl --install
```

再起動後、Ubuntuのセットアップを完了させます。

#### 2. Docker Desktop for Windowsのインストール

1. [Docker Desktop](https://www.docker.com/products/docker-desktop/) をダウンロード
2. インストーラーを実行
3. WSL2 バックエンドを有効化
4. Docker Desktop を起動

#### 3. WSL2内での環境構築

WSL2のUbuntuターミナルで：

```bash
# システムアップデート
sudo apt update && sudo apt upgrade -y

# Git
sudo apt install git -y

# Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Node.js（nvm経由）
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc
nvm install 20
nvm use 20

# golang-migrate
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# その他のツール
sudo apt install make jq -y
```

#### 4. セットアップ続行

上記の「セットアップ手順」に従ってください。

### Linux（Ubuntu/Debian）

#### 1. 必要なソフトウェアのインストール

```bash
# システムアップデート
sudo apt update && sudo apt upgrade -y

# Git
sudo apt install git -y

# Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Docker Compose
sudo apt install docker-compose -y

# Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Node.js（nvm経由）
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc
nvm install 20
nvm use 20

# golang-migrate
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# その他のツール
sudo apt install make jq -y
```

#### 2. ログアウト・再ログイン

Dockerグループの変更を反映させるため、一度ログアウトして再ログインしてください。

#### 3. セットアップ続行

上記の「セットアップ手順」に従ってください。

---

## 動作確認

### 1. バックエンドAPIのテスト

```bash
# ヘルスチェック
curl http://localhost:8080/health

# 期待される応答:
# {"status":"healthy","timestamp":"2024-01-16T15:00:00Z"}
```

### 2. データベース接続確認

```bash
# PostgreSQLに接続
docker-compose exec postgres psql -U postgres -d effisio_dev

# SQLコマンドで確認
\dt  # テーブル一覧
\q   # 終了
```

### 3. Redisの確認

```bash
# Redisに接続
docker-compose exec redis redis-cli

# コマンド実行
PING  # 応答: PONG
KEYS *
exit
```

### 4. フロントエンドの確認

ブラウザで http://localhost:3000 を開き、ページが正常に表示されることを確認。

### 5. ログの確認

```bash
# すべてのサービスのログ
docker-compose logs

# リアルタイムでログを監視
docker-compose logs -f

# 特定のサービスのログ
docker-compose logs backend
docker-compose logs frontend
```

---

## トラブルシューティング

### 問題1: ポート競合エラー

**症状**: `Bind for 0.0.0.0:8080 failed: port is already allocated`

**解決方法**:

```bash
# 使用中のポートを確認
# macOS/Linux
lsof -i :8080

# Windows
netstat -ano | findstr :8080

# プロセスを終了するか、docker-compose.yml でポートを変更
```

### 問題2: Docker Composeが起動しない

**症状**: `docker-compose: command not found`

**解決方法**:

```bash
# Docker Composeのバージョン確認
docker compose version  # V2 の場合（スペースあり）
docker-compose version  # V1 の場合（ハイフンあり）

# V2 を使用する場合は、docker-compose を docker compose に読み替えてください
```

### 問題3: データベース接続エラー

**症状**: `connection refused` または `database does not exist`

**解決方法**:

```bash
# PostgreSQLコンテナが起動しているか確認
docker-compose ps postgres

# コンテナのログを確認
docker-compose logs postgres

# コンテナを再起動
docker-compose restart postgres

# データベースが存在するか確認
docker-compose exec postgres psql -U postgres -c "\l"
```

### 問題4: マイグレーションエラー

**症状**: `no change` または `dirty database version`

**解決方法**:

```bash
# マイグレーションバージョンを確認
migrate -path ./backend/migrations -database "postgres://postgres:postgres@localhost:5432/effisio_dev?sslmode=disable" version

# 強制的にバージョンをリセット（注意: データが失われる可能性あり）
migrate -path ./backend/migrations -database "postgres://postgres:postgres@localhost:5432/effisio_dev?sslmode=disable" force VERSION

# すべてダウンして再度アップ
migrate -path ./backend/migrations -database "postgres://postgres:postgres@localhost:5432/effisio_dev?sslmode=disable" down
migrate -path ./backend/migrations -database "postgres://postgres:postgres@localhost:5432/effisio_dev?sslmode=disable" up
```

### 問題5: Goモジュールのダウンロードエラー

**症状**: `go: module ... not found`

**解決方法**:

```bash
# Go プロキシを確認
go env GOPROXY

# モジュールキャッシュをクリア
go clean -modcache

# 再度ダウンロード
go mod download
```

### 問題6: npm install エラー

**症状**: `EACCES: permission denied`

**解決方法**:

```bash
# node_modules を削除
rm -rf node_modules package-lock.json

# npm キャッシュをクリア
npm cache clean --force

# 再インストール
npm install
```

### 問題7: Dockerのディスク容量不足

**症状**: `no space left on device`

**解決方法**:

```bash
# 使用していないイメージ・コンテナを削除
docker system prune -a

# ボリュームも削除（注意: データが失われます）
docker system prune -a --volumes

# ディスク使用量を確認
docker system df
```

### 問題8: ホットリロードが動作しない

**症状**: コードを変更してもブラウザに反映されない

**解決方法（フロントエンド）**:

```bash
# Next.jsの.nextディレクトリを削除
cd frontend
rm -rf .next

# コンテナを再起動
docker-compose restart frontend
```

**解決方法（バックエンド）**:

```bash
# Airの設定を確認
cat backend/.air.toml

# tmpディレクトリを削除
rm -rf backend/tmp

# コンテナを再起動
docker-compose restart backend
```

---

## 便利なコマンド集

### Docker関連

```bash
# すべてのサービスを起動
docker-compose up -d

# すべてのサービスを停止
docker-compose down

# すべてのサービスを停止してボリュームも削除
docker-compose down -v

# 特定のサービスのみ起動
docker-compose up -d backend

# 特定のサービスのみ再起動
docker-compose restart backend

# ログをリアルタイムで監視
docker-compose logs -f

# コンテナ内でコマンド実行
docker-compose exec backend sh
docker-compose exec frontend sh

# イメージを再ビルド
docker-compose build --no-cache
```

### データベース関連

```bash
# PostgreSQLに接続
docker-compose exec postgres psql -U postgres -d effisio_dev

# データベースバックアップ
docker-compose exec postgres pg_dump -U postgres effisio_dev > backup.sql

# データベースリストア
cat backup.sql | docker-compose exec -T postgres psql -U postgres -d effisio_dev

# テーブル一覧
docker-compose exec postgres psql -U postgres -d effisio_dev -c "\dt"
```

### 開発関連

```bash
# Go のテスト実行
cd backend
go test ./... -v

# Go のカバレッジ確認
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Go のリント実行
golangci-lint run ./...

# Next.jsのビルド
cd frontend
npm run build

# Next.jsのテスト
npm run test
```

---

## VS Code 推奨拡張機能

開発効率を上げるため、以下の拡張機能のインストールを推奨します：

### Go開発

- **Go** (golang.go) - Go言語サポート
- **Go Test Explorer** - テスト実行
- **Go Outline** - アウトライン表示

### TypeScript/React開発

- **ES7+ React/Redux/React-Native snippets** - Reactスニペット
- **Prettier - Code formatter** - コード整形
- **ESLint** - リント
- **Tailwind CSS IntelliSense** - Tailwind補完

### Docker

- **Docker** (ms-azuretools.vscode-docker) - Docker管理

### データベース

- **PostgreSQL** (ckolkman.vscode-postgres) - PostgreSQL管理

### その他

- **GitLens** - Git拡張
- **REST Client** - APIテスト
- **Error Lens** - エラー表示強化

---

## 次のステップ

開発環境のセットアップが完了したら：

1. [Git運用ルール](./GIT_WORKFLOW.md) を確認
2. [コーディング規約](./CODING_GUIDELINES.md) を確認
3. [開発ガイド](./BACKEND_DEVELOPMENT_GUIDE.md) を確認
4. 最初のタスクに取り組む

質問や問題がある場合は、チームのSlackチャンネルで質問してください。
