.PHONY: help setup dev test lint clean build deploy docker-up docker-down migrate-up migrate-down seed

# デフォルトターゲット
help:
	@echo "Effisio プロジェクト - 利用可能なコマンド:"
	@echo ""
	@echo "  make setup        - 初回セットアップ（依存関係インストール、.envファイル作成等）"
	@echo "  make dev          - 開発環境を起動（Docker Compose）"
	@echo "  make test         - 全テストを実行"
	@echo "  make lint         - リンター・フォーマッターを実行"
	@echo "  make clean        - ビルド成果物やキャッシュを削除"
	@echo "  make build        - プロダクション用ビルド"
	@echo ""
	@echo "  make docker-up    - Dockerコンテナを起動"
	@echo "  make docker-down  - Dockerコンテナを停止・削除"
	@echo ""
	@echo "  make migrate-up   - データベースマイグレーション実行"
	@echo "  make migrate-down - データベースマイグレーションをロールバック"
	@echo "  make seed         - シードデータを投入"
	@echo ""
	@echo "  make deploy       - デプロイ実行"
	@echo ""

# 初回セットアップ
setup:
	@echo "🚀 プロジェクトをセットアップしています..."
	@bash scripts/setup.sh

# 開発環境起動
dev: docker-up
	@echo "✅ 開発環境が起動しました"
	@echo ""
	@echo "  Backend:  http://localhost:8080"
	@echo "  Frontend: http://localhost:3000"
	@echo "  Adminer:  http://localhost:8081"
	@echo ""

# Dockerコンテナ起動
docker-up:
	@echo "🐳 Dockerコンテナを起動しています..."
	docker-compose up -d

# Dockerコンテナ停止
docker-down:
	@echo "🛑 Dockerコンテナを停止しています..."
	docker-compose down

# 全テスト実行
test:
	@echo "🧪 テストを実行しています..."
	@$(MAKE) -C backend test
	@$(MAKE) -C frontend test

# リンター・フォーマッター実行
lint:
	@echo "🔍 リンター・フォーマッターを実行しています..."
	@$(MAKE) -C backend lint
	@$(MAKE) -C frontend lint

# クリーンアップ
clean:
	@echo "🧹 クリーンアップしています..."
	@$(MAKE) -C backend clean
	@$(MAKE) -C frontend clean
	docker-compose down -v
	@echo "✅ クリーンアップ完了"

# プロダクション用ビルド
build:
	@echo "🏗️  プロダクション用ビルドを実行しています..."
	@$(MAKE) -C backend build
	@$(MAKE) -C frontend build
	@echo "✅ ビルド完了"

# データベースマイグレーション実行
migrate-up:
	@echo "📊 データベースマイグレーションを実行しています..."
	@$(MAKE) -C backend migrate-up

# データベースマイグレーションをロールバック
migrate-down:
	@echo "⏪ データベースマイグレーションをロールバックしています..."
	@$(MAKE) -C backend migrate-down

# シードデータ投入
seed:
	@echo "🌱 シードデータを投入しています..."
	@$(MAKE) -C backend seed

# デプロイ（環境に応じて調整が必要）
deploy:
	@echo "🚀 デプロイを実行しています..."
	@echo "⚠️  このコマンドは環境に応じてカスタマイズしてください"
	# docker-compose -f docker-compose.prod.yml up -d --build
