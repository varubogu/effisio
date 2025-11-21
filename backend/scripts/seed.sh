#!/bin/bash

# シードデータを投入するスクリプト
# 使用方法: bash scripts/seed.sh

set -e

echo "🌱 シードデータを投入しています..."

# 環境変数の読み込み
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

# デフォルト値の設定
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-effisio_dev}

# PostgreSQLに接続してシードデータを投入
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME << EOF

-- 既存のテストユーザーを削除（存在する場合）
DELETE FROM users WHERE username IN ('admin', 'manager', 'testuser', 'viewer');

-- 管理者ユーザーを作成
-- パスワード: admin123 (bcryptハッシュ)
INSERT INTO users (username, email, password, role, is_active) VALUES
('admin', 'admin@example.com', '\$2a\$10\$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin', true);

-- マネージャーユーザーを作成
-- パスワード: manager123
INSERT INTO users (username, email, password, role, is_active) VALUES
('manager', 'manager@example.com', '\$2a\$10\$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'manager', true);

-- 一般ユーザーを作成
-- パスワード: user123
INSERT INTO users (username, email, password, role, is_active) VALUES
('testuser', 'user@example.com', '\$2a\$10\$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user', true);

-- 閲覧者ユーザーを作成
-- パスワード: viewer123
INSERT INTO users (username, email, password, role, is_active) VALUES
('viewer', 'viewer@example.com', '\$2a\$10\$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'viewer', true);

EOF

echo "✅ シードデータの投入が完了しました"
echo ""
echo "作成されたユーザー:"
echo "  - admin (admin@example.com) - パスワード: admin123"
echo "  - manager (manager@example.com) - パスワード: manager123"
echo "  - testuser (user@example.com) - パスワード: user123"
echo "  - viewer (viewer@example.com) - パスワード: viewer123"
