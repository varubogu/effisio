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
DELETE FROM users WHERE username IN ('admin', 'manager', 'testuser', 'viewer', 'suspended_user');

-- 管理者ユーザーを作成
-- パスワード: admin123 (bcryptハッシュ)
INSERT INTO users (username, email, full_name, department, password_hash, role, status, created_at, updated_at) VALUES
('admin', 'admin@example.com', '管理者 太郎', 'IT部', '\$2a\$10\$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin', 'active', NOW(), NOW());

-- マネージャーユーザーを作成
-- パスワード: manager123
INSERT INTO users (username, email, full_name, department, password_hash, role, status, created_at, updated_at) VALUES
('manager', 'manager@example.com', '管理 次郎', '営業部', '\$2a\$10\$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'manager', 'active', NOW(), NOW());

-- 一般ユーザーを作成
-- パスワード: user123
INSERT INTO users (username, email, full_name, department, password_hash, role, status, created_at, updated_at) VALUES
('testuser', 'user@example.com', 'テスト 三郎', '開発部', '\$2a\$10\$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user', 'active', NOW(), NOW());

-- 閲覧者ユーザーを作成
-- パスワード: viewer123
INSERT INTO users (username, email, full_name, department, password_hash, role, status, created_at, updated_at) VALUES
('viewer', 'viewer@example.com', '閲覧 四郎', '総務部', '\$2a\$10\$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'viewer', 'active', NOW(), NOW());

-- 停止中のユーザー（テスト用）
-- パスワード: suspended123
INSERT INTO users (username, email, full_name, department, password_hash, role, status, created_at, updated_at) VALUES
('suspended_user', 'suspended@example.com', '停止 五郎', 'なし', '\$2a\$10\$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user', 'suspended', NOW(), NOW());

EOF

echo "✅ シードデータの投入が完了しました"
echo ""
echo "作成されたユーザー:"
echo "  - admin (admin@example.com) - パスワード: admin123 - 管理者 太郎 (IT部)"
echo "  - manager (manager@example.com) - パスワード: manager123 - 管理 次郎 (営業部)"
echo "  - testuser (user@example.com) - パスワード: user123 - テスト 三郎 (開発部)"
echo "  - viewer (viewer@example.com) - パスワード: viewer123 - 閲覧 四郎 (総務部)"
echo "  - suspended_user (suspended@example.com) - パスワード: suspended123 - 停止 五郎 (停止中)"
