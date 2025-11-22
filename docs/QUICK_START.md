# クイックスタートガイド

このドキュメントでは、Effisioプロジェクトの開発環境を5分で起動し、すぐに開発を始められるようにします。

## 目次

- [前提条件](#前提条件)
- [5分でセットアップ](#5分でセットアップ)
- [初回起動](#初回起動)
- [動作確認](#動作確認)
- [よくある質問](#よくある質問)
- [次のステップ](#次のステップ)

---

## 前提条件

以下がインストールされていることを確認してください：

| ツール | 最小バージョン | 確認コマンド |
|--------|--------------|-------------|
| Docker | 20.10+ | `docker --version` |
| Docker Compose | 2.0+ | `docker-compose --version` |
| Go | 1.21+ | `go version` |
| Node.js | 18.0+ | `node --version` |
| npm | 9.0+ | `npm --version` |

**まだインストールしていない場合:**
- macOS: `brew install docker docker-compose go node`
- Windows: [Docker Desktop](https://www.docker.com/products/docker-desktop), [Go](https://go.dev/dl/), [Node.js](https://nodejs.org/)
- Linux: 各ディストリビューションのパッケージマネージャーを使用

---

## 5分でセットアップ

### 1. リポジトリのクローン (30秒)

```bash
# HTTPSでクローン
git clone https://github.com/varubogu/effisio.git
cd effisio

# またはSSHでクローン
git clone git@github.com:varubogu/effisio.git
cd effisio
```

### 2. 自動セットアップ実行 (3分)

```bash
# 初回セットアップスクリプトを実行
make setup
```

**このコマンドが自動的に実行すること:**
- 環境設定ファイル (.env) の作成
- Go依存関係のダウンロード
- npm パッケージのインストール
- 開発ツール（Air, golangci-lint）のインストール
- Dockerイメージのプル

### 3. 開発環境起動 (1分)

```bash
# Docker環境を起動
make dev
```

**起動するサービス:**
- PostgreSQL (localhost:5432)
- Redis (localhost:6379)
- Backend API (localhost:8080)
- Frontend (localhost:3000)
- Adminer (localhost:8081) - データベース管理UI
- Redis Commander (localhost:8082) - Redis管理UI

### 4. データベースセットアップ (30秒)

別のターミナルを開いて：

```bash
# マイグレーション実行
make migrate-up

# シードデータ投入
make seed
```

**これで完了です！**

---

## 初回起動

### アクセスURL

| サービス | URL | 説明 |
|---------|-----|------|
| フロントエンド | http://localhost:3000 | Next.jsアプリケーション |
| バックエンドAPI | http://localhost:8080 | Gin APIサーバー |
| API Ping | http://localhost:8080/api/v1/ping | 動作確認用エンドポイント |
| Adminer | http://localhost:8081 | PostgreSQL管理画面 |
| Redis Commander | http://localhost:8082 | Redis管理画面 |

### デフォルトログイン情報

**Adminer (PostgreSQL):**
- サーバー: `postgres`
- ユーザー名: `postgres`
- パスワード: `postgres`
- データベース: `effisio_dev`

**テストユーザー（シードデータ）:**
- 管理者: `admin` / `admin123`
- マネージャー: `manager` / `manager123`
- 一般ユーザー: `testuser` / `user123`
- 閲覧者: `viewer` / `viewer123`

---

## 動作確認

### 1. バックエンドAPI確認

```bash
# Ping エンドポイント
curl http://localhost:8080/api/v1/ping

# 期待されるレスポンス:
# {"message":"pong"}

# ユーザー一覧取得
curl http://localhost:8080/api/v1/users

# 期待されるレスポンス:
# {"users":[...]}
```

### 2. フロントエンド確認

ブラウザで http://localhost:3000 にアクセス

- トップページが表示される
- "ユーザー管理" リンクをクリック
- ユーザー一覧ページが表示される（シードデータの4ユーザー）

### 3. ホットリロード確認

**バックエンド:**
```bash
# backend/internal/handler/user.go を編集してみる
# 保存すると自動的にサーバーが再起動される
```

**フロントエンド:**
```bash
# frontend/src/app/page.tsx を編集してみる
# 保存すると自動的にブラウザがリロードされる
```

---

## よくある質問

### Q1. `make setup` でエラーが出る

**A1. Go/Node.jsのバージョンを確認**
```bash
go version  # Go 1.21以上必要
node --version  # Node.js 18以上必要
```

**A2. Dockerが起動していない**
```bash
docker ps  # エラーが出る場合はDockerを起動
```

### Q2. ポート番号が競合している

**エラー例:**
```
Error: bind: address already in use
```

**解決方法:**
```bash
# 使用中のポートを確認
lsof -i :3000  # フロントエンド
lsof -i :8080  # バックエンド
lsof -i :5432  # PostgreSQL

# プロセスを終了
kill -9 <PID>

# または docker-compose.yml でポート番号を変更
```

### Q3. マイグレーションが失敗する

**エラー例:**
```
error: Dirty database version 1. Fix and force version.
```

**解決方法:**
```bash
# データベースをリセット
make migrate-down
make migrate-up

# または完全にクリーンアップ
docker-compose down -v  # ボリュームも削除
make dev
make migrate-up
```

### Q4. npm install でエラーが出る

**解決方法:**
```bash
cd frontend

# キャッシュをクリア
npm cache clean --force

# node_modules削除して再インストール
rm -rf node_modules package-lock.json
npm install
```

### Q5. Air（ホットリロード）が動かない

**解決方法:**
```bash
# Airを手動インストール
go install github.com/cosmtrek/air@latest

# PATHを確認
echo $GOPATH/bin  # このパスが$PATHに含まれているか確認

# 含まれていない場合は追加（~/.zshrc または ~/.bashrc）
export PATH=$PATH:$(go env GOPATH)/bin
```

### Q6. データベースに接続できない

**確認事項:**
```bash
# 1. PostgreSQLコンテナが起動しているか
docker ps | grep postgres

# 2. 接続情報が正しいか
cat backend/.env | grep DB_

# 3. psqlで直接接続してみる
psql -h localhost -p 5432 -U postgres -d effisio_dev
# パスワード: postgres
```

### Q7. フロントエンドが真っ白

**解決方法:**
```bash
cd frontend

# Next.jsのキャッシュをクリア
rm -rf .next

# 再ビルド
npm run dev
```

---

## 開発コマンド一覧

### プロジェクト全体

```bash
make setup        # 初回セットアップ
make dev          # 開発環境起動
make test         # 全テスト実行
make lint         # 全リンター実行
make clean        # クリーンアップ
make build        # プロダクションビルド
```

### バックエンド

```bash
cd backend

make build        # バイナリビルド
make run          # ビルドして実行
make dev          # ホットリロードで実行（Air）
make test         # テスト実行
make lint         # リンター実行
make migrate-up   # マイグレーション実行
make migrate-down # マイグレーションロールバック
make seed         # シードデータ投入
```

### フロントエンド

```bash
cd frontend

npm run dev       # 開発サーバー起動
npm run build     # プロダクションビルド
npm run start     # プロダクションサーバー起動
npm test          # テスト実行
npm run lint      # リンター実行
npm run format    # コードフォーマット
```

### Docker

```bash
docker-compose up -d          # バックグラウンドで起動
docker-compose down           # 停止・削除
docker-compose down -v        # ボリュームも削除
docker-compose logs -f        # ログ表示
docker-compose logs -f backend # 特定サービスのログ
docker-compose restart backend # サービス再起動
docker-compose ps             # 起動中のサービス一覧
```

---

## トラブル時の完全リセット

全てがうまくいかない場合の最終手段：

```bash
# 1. Docker環境を完全削除
docker-compose down -v
docker system prune -a --volumes

# 2. ローカルファイルをクリーンアップ
cd backend
rm -rf bin/ coverage.out coverage.html
go clean -cache -testcache
cd ..

cd frontend
rm -rf .next node_modules package-lock.json
cd ..

# 3. 再セットアップ
make setup
make dev

# 4. マイグレーションとシード
make migrate-up
make seed
```

---

## 次のステップ

環境が正常に起動したら、以下のドキュメントを読んで開発を始めましょう：

### 1. 日々の開発作業
→ **[DAILY_WORKFLOW.md](DAILY_WORKFLOW.md)** を読む

### 2. Phase 1の実装を始める
→ **[PHASE1_IMPLEMENTATION_STEPS.md](PHASE1_IMPLEMENTATION_STEPS.md)** を読む

### 3. コーディング規約を理解する
→ **[CODING_GUIDELINES_GO.md](CODING_GUIDELINES_GO.md)** と **[CODING_GUIDELINES_TYPESCRIPT.md](CODING_GUIDELINES_TYPESCRIPT.md)** を読む

### 4. API実装方法を学ぶ
→ **[API_IMPLEMENTATION_GUIDE.md](API_IMPLEMENTATION_GUIDE.md)** を読む

---

## サポート

問題が解決しない場合：

1. **[TROUBLESHOOTING.md](TROUBLESHOOTING.md)** を確認
2. GitHub Issues で検索
3. 新しいIssueを作成（テンプレートに従って）

---

## 開発環境の停止

```bash
# Dockerコンテナを停止（データは保持）
docker-compose down

# 次回起動時
make dev
```

**データベースのデータは保持されます。** 完全にリセットしたい場合のみ `-v` フラグを使用してください。

---

## チェックリスト

初回セットアップが完了したら、以下を確認：

- [ ] `make dev` でエラーなく起動する
- [ ] http://localhost:3000 でフロントエンドが表示される
- [ ] http://localhost:8080/api/v1/ping で `{"message":"pong"}` が返る
- [ ] http://localhost:8080/api/v1/users でユーザー一覧が返る
- [ ] Adminer (localhost:8081) でデータベースにアクセスできる
- [ ] コード変更時にホットリロードが動作する
- [ ] `make test` でテストが実行される
- [ ] `make lint` でリンターが実行される

全てチェックできたら、開発環境のセットアップは完了です！
