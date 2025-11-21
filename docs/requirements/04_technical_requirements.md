# 技術要件書

## 1. 開発環境要件

### 1.1 開発言語・ランタイム
**要件ID**: TR-LANG-001

- **Go**: 1.21 以上
- **データベース**: PostgreSQL 13.0 以上
- **キャッシュ**: Redis 6.0 以上
- **コンテナ**: Docker 20.10 以上
- ** орхестレーション**: Docker Compose 1.29 以上（開発環境）

### 1.2 開発ツール
**要件ID**: TR-TOOL-001

| ツール | バージョン | 用途 |
|--------|-----------|------|
| Go | 1.21+ | 開発言語 |
| golangci-lint | v1.55+ | リント・静的解析 |
| testify | v1.8+ | テスティング |
| sqlc | v1.24+ | SQL型安全性 |
| swag | v1.16+ | API ドキュメント |

---

## 2. ランタイム環境要件

### 2.1 サーバーOS
**要件ID**: TR-OS-001

- **本番**: Ubuntu 20.04 LTS 以上 / CentOS 7.9 以上 / Amazon Linux 2
- **開発**: Linux / macOS / Windows (WSL2)

### 2.2 ハードウェア要件
**要件ID**: TR-HW-001

#### 開発環境
- CPU: 2 コア以上
- メモリ: 4GB 以上
- ストレージ: 20GB 以上

#### ステージング環境
- CPU: 2 コア以上
- メモリ: 8GB 以上
- ストレージ: 50GB 以上

#### 本番環境
- CPU: 4 コア以上（スケーラブル）
- メモリ: 16GB 以上（スケーラブル）
- ストレージ: 100GB 以上（増設可能）

---

## 3. バックエンド技術要件

### 3.1 Webフレームワーク
**要件ID**: TR-BACK-001

**候補**: Gin, Echo, 標準ライブラリ

**推奨**: Gin または Echo
- 理由：高速、シンプル、ミドルウェア豊富、保守性が高い

**選定基準**:
- [ ] レスポンス時間の比較テスト
- [ ] ミドルウェア機能の充実度
- [ ] コミュニティサポート
- [ ] ドキュメント品質

### 3.2 データベース
**要件ID**: TR-DB-001

- **DBMS**: PostgreSQL 13.0 以上
- **ドライバ**: pq または pgx
- **接続プール**: pgbouncer（本番環境）
- **最大接続数**: 50（開発/ステージング）、100+（本番）

### 3.3 ORM/Query Builder
**要件ID**: TR-ORM-001

**候補**: GORM, sqlc, 標準 database/sql

| ツール | 特徴 | 推奨用途 |
|--------|------|---------|
| GORM | 高機能、開発効率重視 | 開発速度優先の場合 |
| sqlc | 型安全、パフォーマンス重視 | 型安全性と性能重視 |
| sql | シンプル、完全制御 | 複雑なクエリが必要な場合 |

**選定方針**: 型安全性とパフォーマンスを重視する場合は sqlc、開発速度を重視する場合は GORM

### 3.4 認証・認可
**要件ID**: TR-AUTH-001

- **JWT ライブラリ**: jwt-go または golang-jwt
- **パスワードハッシング**: golang.org/x/crypto/bcrypt
- **トークン署名**: RS256 推奨（非対称暗号化）
- **トークン有効期限**: Access Token 1時間、Refresh Token 30日

### 3.5 キャッシング
**要件ID**: TR-CACHE-001

- **ライブラリ**: go-redis/redis/v8
- **キャッシュ層**: Redis 6.0 以上
- **キャッシュ戦略**: TTL ベース + イベントベース無効化
- **キャッシュサイズ**: 1GB 以上（本番環境）

### 3.6 ロギング
**要件ID**: TR-LOG-001

| ツール | 特徴 | 推奨 |
|--------|------|------|
| zap | 高速、構造化ログ | ★★★ |
| logrus | 普及度高、柔軟 | ★★ |
| slog | Go 1.21 標準 | ★★ |

**選定**: zap （高性能、構造化ログ）

**ログレベル**: DEBUG, INFO, WARN, ERROR
**ログ形式**: JSON（本番）
**ログ出力先**: stdout（Docker環境）、ファイル（オンプレ）

### 3.7 バリデーション
**要件ID**: TR-VALID-001

- **ライブラリ**: github.com/go-playground/validator/v10
- **バリデーションタイプ**: required, email, min, max, alphanumeric など

### 3.8 HTTP クライアント
**要件ID**: TR-HTTP-001

- **ライブラリ**: net/http（標準） または go-resty/resty（高機能時）
- **タイムアウト**: 30秒

### 3.9 マイグレーション
**要件ID**: TR-MIGRATION-001

- **ツール**: golang-migrate または sql-migrate
- **管理方法**: ファイルベース（migrations/ ディレクトリ）
- **バージョニング**: タイムスタンプベース（001_init.up.sql）

---

## 4. フロントエンド技術要件

### 4.1 初期フェーズ
**要件ID**: TR-FRONT-001

- **テンプレートエンジン**: html/template（Go標準）
- **理由**: サーバーサイド描画、シンプル

### 4.2 将来フェーズ
**要件ID**: TR-FRONT-002

- **フレームワーク**: React または Vue.js
- **言語**: TypeScript
- **スタイリング**: Tailwind CSS
- **ビルドツール**: Vite または Next.js

---

## 5. テストフレームワーク

### 5.1 ユニットテスト
**要件ID**: TR-TEST-001

- **ライブラリ**: testing（標準）+ testify/assert
- **モック**: github.com/stretchr/testify/mock または counterfeiter
- **テストデータ**: fixtures （JSON/YAML）

### 5.2 統合テスト
**要件ID**: TR-TEST-002

- **ツール**: testify/suite
- **DB**: テスト用PostgreSQL（Docker）
- **クリーンアップ**: 各テスト前後にロールバック

### 5.3 E2E テスト
**要件ID**: TR-TEST-003

- **ツール**: Cypress または Playwright
- **対象**: ブラウザベースのテスト
- **実行タイミング**: リリース前

---

## 6. ビルド・デプロイメント技術

### 6.1 ビルド
**要件ID**: TR-BUILD-001

```bash
# バイナリビルド
go build -ldflags="-s -w -X main.Version=1.0.0" -o bin/server cmd/server/main.go

# Docker イメージビルド
docker build -t internal-system:latest .
```

**ビルド最適化**:
- マルチステージビルド（開発イメージサイズ削減）
- ストリップ実行ファイル（-ldflags="-s -w"）
- キャッシング活用

### 6.2 CI/CD パイプライン
**要件ID**: TR-CICD-001

**プラットフォーム**: GitHub Actions （GitHub使用時）

**パイプラインステップ**:
1. チェックアウト
2. Go環境セットアップ
3. 依存関係インストール
4. リント実行（golangci-lint）
5. テスト実行（テストカバレッジ報告）
6. ビルド
7. Docker イメージビルド＆プッシュ
8. ステージング環境へデプロイ（develop ブランチ）
9. 本番環境へデプロイ（main ブランチ）

### 6.3 コンテナ化
**要件ID**: TR-CONTAINER-001

**ベースイメージ**: alpine:latest （軽量化）

**Dockerfile（マルチステージ）**:
```dockerfile
# ステージ1: ビルド
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o server cmd/server/main.go

# ステージ2: 実行
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

### 6.4 オーケストレーション
**要件ID**: TR-ORCH-001

**開発環境**: Docker Compose
- アプリケーション、PostgreSQL、Redis を一括管理

**本番環境**: Docker Compose または Kubernetes（将来）

---

## 7. インフラストラクチャ要件

### 7.1 ロードバランサー
**要件ID**: TR-INFRA-001

- **ツール**: Nginx または HAProxy
- **機能**: リバースプロキシ、SSL/TLS 終端、ヘルスチェック
- **ラウンドロビン**: ウェイト付きラウンドロビン

### 7.2 リバースプロキシ設定
**要件ID**: TR-INFRA-002

```nginx
upstream backend {
  server app1:8080 weight=1;
  server app2:8080 weight=1;
}

server {
  listen 443 ssl http2;
  server_name api.internal-system.com;

  ssl_certificate /etc/ssl/cert.pem;
  ssl_certificate_key /etc/ssl/key.pem;

  location / {
    proxy_pass http://backend;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
  }
}
```

### 7.3 SSL/TLS
**要件ID**: TR-INFRA-003

- **プロトコル**: TLS 1.2 以上（推奨 1.3）
- **証明書**: Let's Encrypt （自動更新）
- **有効期限**: 90日（自動更新）

### 7.4 ファイアウォール設定
**要件ID**: TR-INFRA-004

| プロトコル | ポート | 許可元 | 説明 |
|-----------|--------|--------|------|
| HTTPS | 443 | 任意 | Webアクセス |
| HTTP | 80 | 任意 | HTTPS リダイレクト |
| SSH | 22 | 管理者IP | 管理用 |
| PostgreSQL | 5432 | アプリサーバーのみ | DB接続 |
| Redis | 6379 | アプリサーバーのみ | キャッシュ接続 |

---

## 8. モニタリング・ロギング技術

### 8.1 メトリクス収集
**要件ID**: TR-MONITOR-001

- **ツール**: Prometheus
- **クライアント**: prometheus/client_golang
- **スクレイプ間隔**: 15秒
- **保持期間**: 15日

**監視メトリクス**:
- HTTP リクエスト数、応答時間
- DB クエリ実行時間
- キャッシュヒット率
- エラーレート
- メモリ使用量、CPU使用率

### 8.2 ログ集約
**要件ID**: TR-MONITOR-002

- **ツール**: ELK Stack（Elasticsearch, Logstash, Kibana）
- **または**: Splunk
- **ログ形式**: JSON
- **保持期間**: 90日

### 8.3 可視化
**要件ID**: TR-MONITOR-003

- **ツール**: Grafana
- **ダッシュボード**: アプリケーション、インフラ、ビジネスメトリクス

### 8.4 アラート
**要件ID**: TR-MONITOR-004

- **ツール**: Prometheus Alertmanager
- **通知先**: Slack、メール
- **応答時間**: 5分以内

---

## 9. セキュリティツール

### 9.1 静的解析
**要件ID**: TR-SEC-001

- **ツール**: golangci-lint
- **チェック項目**:
  - vet：型安全
  - errcheck：エラーハンドリング
  - gosec：セキュリティ問題
  - staticcheck：品質問題

### 9.2 依存関係スキャン
**要件ID**: TR-SEC-002

- **ツール**: nancy または go-audit
- **実行タイミング**: CI/CD パイプライン
- **脆弱性DB更新**: 日次

### 9.3 コード品質
**要件ID**: TR-SEC-003

- **ツール**: SonarQube （オプション）
- **メトリクス**: カバレッジ、複雑度、脆弱性

---

## 10. バージョン管理

### 10.1 Git ワークフロー
**要件ID**: TR-VCS-001

**ブランチ戦略**: Git Flow
```
main (本番リリース)
  ├── develop (開発メインブランチ)
  │   ├── feature/... (機能開発)
  │   ├── hotfix/... (緊急修正)
  │   └── release/... (リリース準備)
```

**コミットメッセージ**: Conventional Commits
```
feat: ユーザー認証機能を追加
fix: パスワード検証ロジックを修正
docs: README を更新
test: ユーザーログインテストを追加
```

### 10.2 リリースプロセス
**要件ID**: TR-VCS-002

1. develop ブランチから release/* ブランチを作成
2. バージョン番号を更新（go.mod, main.go など）
3. CHANGELOG を更新
4. release/* ブランチから main ブランチへプルリクエスト
5. レビュー・マージ後、タグ作成（v1.0.0）
6. main から develop へもマージ

---

## 11. ドキュメンテーション技術

### 11.1 API ドキュメント
**要件ID**: TR-DOC-001

- **ツール**: Swagger/OpenAPI 3.0
- **生成**: swaggo (go-swagger)
- **URL**: /swagger/index.html
- **更新**: コード変更に自動反映

### 11.2 MarkdownドキュメントドM
**要件ID**: TR-DOC-002

- **リポジトリ**: docs/ ディレクトリ
- **形式**: Markdown
- **ホスティング**: GitHub Pages または Gitbook（将来）

### 11.3 コード内ドキュメント
**要件ID**: TR-DOC-003

- **言語**: 日本語 / 英語
- **スタイル**: Godoc形式
- **カバレッジ**: 公開パッケージ・関数は100%

---

## 12. ライブラリ依存関係管理

### 12.1 Go モジュール
**要件ID**: TR-DEP-001

```
go.mod, go.sum で管理
- go mod download: 依存関係ダウンロード
- go mod verify: チェックサム確認
- go mod tidy: 不要な依存関係削除
```

### 12.2 推奨ライブラリスタック
**要件ID**: TR-DEP-002

| 用途 | ライブラリ | バージョン |
|------|-----------|-----------|
| Web フレームワーク | gin-gonic/gin | v1.9+ |
| DB ドライバ | lib/pq | v1.10+ |
| ORM | GORM | v1.25+ |
| キャッシュ | go-redis/redis/v8 | v8.11+ |
| JWT | golang-jwt/jwt/v5 | v5.0+ |
| パスワード | golang.org/x/crypto | latest |
| ロギング | uber-go/zap | v1.26+ |
| テスト | stretchr/testify | v1.8+ |

### 12.3 依存関係ポリシー
**要件ID**: TR-DEP-003

- [ ] 外部パッケージは最小限に
- [ ] ライセンス確認（MIT, Apache 2.0 など）
- [ ] メンテナンス状況確認（更新頻度、イシュー対応）
- [ ] セキュリティアップデート：優先的に対応
