.PHONY: help build run test lint clean migrate-up migrate-down migrate-create seed dev

# 変数定義
BINARY_NAME=effisio-api
MAIN_PATH=./cmd/server/main.go
MIGRATION_DIR=./migrations
DB_URL ?= postgres://postgres:postgres@localhost:5432/effisio_dev?sslmode=disable

# デフォルトターゲット
help:
	@echo "Backend - 利用可能なコマンド:"
	@echo ""
	@echo "  make build         - バイナリをビルド"
	@echo "  make run           - アプリケーションを実行"
	@echo "  make dev           - 開発モード（hot reload）"
	@echo "  make test          - テストを実行"
	@echo "  make test-coverage - カバレッジ付きでテストを実行"
	@echo "  make lint          - リンターを実行"
	@echo "  make clean         - ビルド成果物を削除"
	@echo ""
	@echo "  make migrate-up    - マイグレーション実行"
	@echo "  make migrate-down  - マイグレーションをロールバック"
	@echo "  make migrate-create NAME=<name> - 新しいマイグレーションファイルを作成"
	@echo "  make seed          - シードデータを投入"
	@echo ""
	@echo "  make deps          - 依存関係をインストール"
	@echo "  make tidy          - go.mod を整理"
	@echo ""

# ビルド
build:
	@echo "🏗️  ビルドしています..."
	go build -ldflags="-s -w" -o bin/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ ビルド完了: bin/$(BINARY_NAME)"

# 実行
run: build
	@echo "🚀 アプリケーションを起動しています..."
	./bin/$(BINARY_NAME)

# 開発モード（hot reload with Air）
dev:
	@echo "🔥 開発モードで起動しています（hot reload有効）..."
	@if ! command -v air > /dev/null; then \
		echo "⚠️  Air がインストールされていません。インストールしています..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	air -c .air.toml

# テスト
test:
	@echo "🧪 テストを実行しています..."
	go test -v -race -timeout 30s ./...

# カバレッジ付きテスト
test-coverage:
	@echo "🧪 カバレッジ付きでテストを実行しています..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ カバレッジレポート: coverage.html"

# リンター
lint:
	@echo "🔍 リンターを実行しています..."
	@if ! command -v golangci-lint > /dev/null; then \
		echo "⚠️  golangci-lint がインストールされていません。インストールしています..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run ./...
	@echo "✅ リンター完了"

# フォーマット
fmt:
	@echo "✨ コードをフォーマットしています..."
	go fmt ./...
	@echo "✅ フォーマット完了"

# クリーンアップ
clean:
	@echo "🧹 クリーンアップしています..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -cache -testcache
	@echo "✅ クリーンアップ完了"

# 依存関係のインストール
deps:
	@echo "📦 依存関係をインストールしています..."
	go mod download
	@echo "✅ 依存関係のインストール完了"

# go.mod を整理
tidy:
	@echo "🧹 go.mod を整理しています..."
	go mod tidy
	@echo "✅ go.mod の整理完了"

# マイグレーション実行
migrate-up:
	@echo "📊 マイグレーションを実行しています..."
	@if ! command -v migrate > /dev/null; then \
		echo "⚠️  golang-migrate がインストールされていません"; \
		echo "インストール方法: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; \
		exit 1; \
	fi
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" up
	@echo "✅ マイグレーション完了"

# マイグレーションをロールバック
migrate-down:
	@echo "⏪ マイグレーションをロールバックしています..."
	@if ! command -v migrate > /dev/null; then \
		echo "⚠️  golang-migrate がインストールされていません"; \
		echo "インストール方法: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; \
		exit 1; \
	fi
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" down 1
	@echo "✅ ロールバック完了"

# 新しいマイグレーションファイルを作成
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "❌ エラー: NAME パラメータが必要です"; \
		echo "使用例: make migrate-create NAME=create_users_table"; \
		exit 1; \
	fi
	@if ! command -v migrate > /dev/null; then \
		echo "⚠️  golang-migrate がインストールされていません"; \
		echo "インストール方法: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; \
		exit 1; \
	fi
	@echo "📝 マイグレーションファイルを作成しています: $(NAME)"
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(NAME)
	@echo "✅ マイグレーションファイルを作成しました"

# シードデータ投入
seed:
	@echo "🌱 シードデータを投入しています..."
	@if [ -f "./scripts/seed.sh" ]; then \
		bash ./scripts/seed.sh; \
	else \
		echo "⚠️  scripts/seed.sh が見つかりません"; \
	fi
	@echo "✅ シードデータの投入完了"

# Docker環境でのテスト（CI用）
test-ci:
	@echo "🧪 CI環境でテストを実行しています..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# セキュリティチェック
security:
	@echo "🔒 セキュリティチェックを実行しています..."
	@if ! command -v gosec > /dev/null; then \
		echo "⚠️  gosec がインストールされていません。インストールしています..."; \
		go install github.com/securego/gosec/v2/cmd/gosec@latest; \
	fi
	gosec ./...
	@echo "✅ セキュリティチェック完了"
