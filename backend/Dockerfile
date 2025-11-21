# ========================================
# Effisio バックエンド本番環境 Dockerfile
# ========================================

# ========================================
# ステージ1: ビルド
# ========================================
FROM golang:1.21-alpine AS builder

# 作業ディレクトリ
WORKDIR /app

# ビルドに必要なパッケージ
RUN apk add --no-cache git make gcc musl-dev

# 依存関係のダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# バイナリビルド
# -ldflags: バイナリサイズを削減
# -trimpath: ファイルパス情報を削除（セキュリティ向上）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w -X main.Version=${VERSION:-dev}" \
    -trimpath \
    -o /app/bin/server \
    ./cmd/server/main.go

# ========================================
# ステージ2: 実行環境
# ========================================
FROM alpine:latest

# セキュリティアップデート & 証明書インストール
RUN apk --no-cache add ca-certificates tzdata

# 非rootユーザーを作成
RUN addgroup -g 1000 effisio && \
    adduser -D -u 1000 -G effisio effisio

# 作業ディレクトリ
WORKDIR /home/effisio

# ビルドしたバイナリをコピー
COPY --from=builder --chown=effisio:effisio /app/bin/server ./server

# マイグレーションファイルをコピー（必要に応じて）
COPY --from=builder --chown=effisio:effisio /app/migrations ./migrations

# 非rootユーザーに切り替え
USER effisio

# ポート公開
EXPOSE 8080

# ヘルスチェック
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# サーバー起動
CMD ["./server"]
