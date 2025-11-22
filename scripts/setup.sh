#!/bin/bash

# Effisio プロジェクト初回セットアップスクリプト
# 使用方法: bash scripts/setup.sh

set -e

echo "========================================"
echo "  Effisio プロジェクトセットアップ"
echo "========================================"
echo ""

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 前提条件のチェック
echo "📋 前提条件をチェックしています..."

check_command() {
    if ! command -v $1 &> /dev/null; then
        echo -e "${RED}❌ $1 がインストールされていません${NC}"
        echo "   インストール方法については docs/DEVELOPMENT_SETUP.md を参照してください"
        exit 1
    else
        echo -e "${GREEN}✅ $1 が見つかりました${NC}"
    fi
}

check_command docker
check_command docker-compose
check_command go
check_command node
check_command npm

echo ""

# Go バージョンチェック
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
GO_MIN_VERSION="1.21"
if [ "$(printf '%s\n' "$GO_MIN_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$GO_MIN_VERSION" ]; then
    echo -e "${RED}❌ Go のバージョンが $GO_MIN_VERSION 以上である必要があります（現在: $GO_VERSION）${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Go バージョン: $GO_VERSION${NC}"

# Node.js バージョンチェック
NODE_VERSION=$(node --version | sed 's/v//')
NODE_MIN_VERSION="18.0.0"
if [ "$(printf '%s\n' "$NODE_MIN_VERSION" "$NODE_VERSION" | sort -V | head -n1)" != "$NODE_MIN_VERSION" ]; then
    echo -e "${RED}❌ Node.js のバージョンが $NODE_MIN_VERSION 以上である必要があります（現在: $NODE_VERSION）${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Node.js バージョン: $NODE_VERSION${NC}"

echo ""

# .envファイルの作成
echo "📝 環境設定ファイルを作成しています..."

if [ ! -f backend/.env ]; then
    cp backend/.env.example backend/.env
    echo -e "${GREEN}✅ backend/.env を作成しました${NC}"
else
    echo -e "${YELLOW}⚠️  backend/.env は既に存在します${NC}"
fi

if [ ! -f frontend/.env ]; then
    cp frontend/.env.example frontend/.env
    echo -e "${GREEN}✅ frontend/.env を作成しました${NC}"
else
    echo -e "${YELLOW}⚠️  frontend/.env は既に存在します${NC}"
fi

echo ""

# Backend セットアップ
echo "🔧 Backend をセットアップしています..."
cd backend

echo "  📦 Go モジュールをダウンロードしています..."
go mod download
go mod tidy

echo "  🔨 開発ツールをインストールしています..."
if ! command -v air &> /dev/null; then
    echo "     - Air (hot reload) をインストール中..."
    go install github.com/cosmtrek/air@latest
fi

if ! command -v golangci-lint &> /dev/null; then
    echo "     - golangci-lint をインストール中..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

echo -e "${GREEN}✅ Backend のセットアップが完了しました${NC}"
cd ..

echo ""

# Frontend セットアップ
echo "🎨 Frontend をセットアップしています..."
cd frontend

echo "  📦 npm パッケージをインストールしています..."
npm install

echo -e "${GREEN}✅ Frontend のセットアップが完了しました${NC}"
cd ..

echo ""

# Docker環境の準備
echo "🐳 Docker環境を準備しています..."
docker-compose pull

echo ""

# セットアップ完了
echo "========================================"
echo -e "${GREEN}  ✅ セットアップが完了しました！${NC}"
echo "========================================"
echo ""
echo "次のステップ:"
echo ""
echo "  1. 開発環境を起動:"
echo "     $ make dev"
echo ""
echo "  2. ブラウザでアクセス:"
echo "     - Frontend: http://localhost:3000"
echo "     - Backend API: http://localhost:8080"
echo "     - Adminer (DB管理): http://localhost:8081"
echo ""
echo "  3. データベースマイグレーション実行:"
echo "     $ make migrate-up"
echo ""
echo "  4. シードデータを投入:"
echo "     $ make seed"
echo ""
echo "詳細なドキュメントは docs/ フォルダを参照してください"
echo ""
