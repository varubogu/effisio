# デプロイメント計画

## 環境構成

### 開発環境 (Development)
- **ホスト**: ローカルマシン
- **データベース**: Docker PostgreSQL
- **キャッシュ**: Docker Redis
- **起動方法**: `docker-compose up`

### ステージング環境 (Staging)
- **ホスト**: Linux サーバー
- **データベース**: PostgreSQL (外部管理)
- **キャッシュ**: Redis (外部管理)
- **デプロイ**: Docker コンテナ

### 本番環境 (Production)
- **ホスト**: クラウドインスタンス (AWS EC2 または同等)
- **データベース**: マネージドデータベース (RDS または同等)
- **キャッシュ**: マネージドキャッシュ (ElastiCache または同等)
- **ロードバランサー**: Nginx
- **コンテナオーケストレーション**: Docker または Kubernetes
- **CDN**: CloudFront (オプション)

---

## ビルドプロセス

### 1. ローカルビルド

```bash
# ユニットテスト実行
go test ./... -v

# ビルド
go build -o bin/server cmd/server/main.go

# バイナリサイズ最適化
go build -ldflags="-s -w" -o bin/server cmd/server/main.go
```

### 2. Docker ビルド

```dockerfile
# Dockerfile (マルチステージビルド)
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

---

## CI/CD パイプライン

### GitHub Actions (例)

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21

      - name: Run tests
        run: go test ./... -v -coverprofile=coverage.out

      - name: Run linter
        uses: golangci/golangci-lint-action@v3

      - name: Upload coverage
        uses: codecov/codecov-action@v2

  build:
    runs-on: ubuntu-latest
    needs: test
    if: github.event_name == 'push'
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2

      - name: Build Docker image
        run: docker build -t internal-system:${{ github.sha }} .

      - name: Push to registry
        run: |
          echo ${{ secrets.DOCKER_PASSWORD }} | docker login -u ${{ secrets.DOCKER_USERNAME }} --password-stdin
          docker push internal-system:${{ github.sha }}

  deploy:
    runs-on: ubuntu-latest
    needs: build
    if: github.ref == 'refs/heads/main'
    steps:
      - name: Deploy to production
        run: |
          # デプロイスクリプト実行
          ssh deploy@prod-server "docker pull internal-system:${{ github.sha }} && docker-compose up -d"
```

---

## リリースプロセス

### バージョニング
- **形式**: Semantic Versioning (v1.0.0)
- **タグ管理**: Git タグで管理

```bash
# リリースタグ作成
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

### リリースチェックリスト

- [ ] 機能完成・テスト完了
- [ ] CHANGELOG 更新
- [ ] バージョン番号更新 (cmd/server/main.go)
- [ ] ドキュメント更新
- [ ] Git タグ作成・プッシュ
- [ ] GitHub Release 作成
- [ ] ステージング環境でテスト
- [ ] 本番環境へデプロイ
- [ ] ヘルスチェック・ログ確認

---

## ダウンタイム最小化戦略

### ローリングデプロイ

```bash
# 古いサーバーを徐々に停止しながら新しいバージョンを起動
1. 新しい Docker イメージをビルド
2. ロードバランサーのヘルスチェック無効化
3. 新しいインスタンス起動
4. トラフィック段階的に転送
5. 古いインスタンス停止
```

### ブルーグリーンデプロイ

```
Blue (現在の本番)  →  Green (新バージョン)
    ↓
  テスト・検証
    ↓
ロードバランサー切り替え
```

---

## バックアップ・リカバリ

### データベースバックアップ

```bash
# 日次フルバックアップ
0 2 * * * pg_dump -h localhost -U user dbname > /backup/db-$(date +\%Y\%m\%d).sql

# 時間ごと差分バックアップ (WAL アーカイブ)
# PostgreSQL の archive_command で自動化
```

### リストア手順

```bash
# フルバックアップからリストア
psql -U user dbname < /backup/db-20240116.sql

# WAL アーカイブからリストア (PITR)
# postgresql.conf で recovery_target_timeline = 'latest' を指定
```

---

## モニタリング・アラート

### メトリクス収集

- **Prometheus**: メトリクス収集
- **Grafana**: ダッシュボード可視化

### アラート設定

```
- CPU 使用率 > 80%: 警告
- メモリ使用率 > 90%: 警告
- ディスク使用率 > 95%: 警告
- エラーレート > 1%: 警告
- 応答時間 > 1000ms: 警告
- データベース接続数 > 50: 警告
```

### ログ集約

- **ELK Stack** または **Splunk**
- 本番環境のログを一元管理
- エラーログは Slack に通知

---

## 障害対応

### ヘルスチェック

```go
// GET /health
// 期待レスポンス:
{
  "status": "healthy",
  "timestamp": "2024-01-16T15:00:00Z",
  "database": "ok",
  "cache": "ok"
}
```

### 緊急時の対応

1. **エラーレート急上昇**
   - ロードバランサーから古いバージョンにロールバック
   - ログ確認
   - チーム通知

2. **データベースダウン**
   - キャッシュからの読み取りに切り替え
   - リードレプリカへのフェイルオーバー
   - 手動リストア

3. **メモリリーク**
   - インスタンスを再起動
   - コードレビュー・修正
   - 再デプロイ

---

## セキュリティアップデート

- **OS パッチ**: 月 1 回 (金曜日に実施)
- **Go ライブラリ**: 脆弱性検出時は即座に対応
- **Docker イメージ**: 月 1 回更新

---

## 本番環境チェックリスト

- [ ] HTTPS 有効化
- [ ] 環境変数設定 (.env)
- [ ] データベース接続確認
- [ ] キャッシュ接続確認
- [ ] ログローテーション設定
- [ ] バックアップスケジュール確認
- [ ] モニタリング・アラート設定
- [ ] ヘルスチェック動作確認
- [ ] ファイアウォール設定
- [ ] セキュリティグループ設定
