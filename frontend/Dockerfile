# ========================================
# Effisio フロントエンド本番環境 Dockerfile
# ========================================

# ========================================
# ステージ1: 依存関係インストール
# ========================================
FROM node:20-alpine AS deps

WORKDIR /app

# 依存関係ファイルをコピー
COPY package.json package-lock.json* ./

# 依存関係インストール
RUN npm ci --only=production

# ========================================
# ステージ2: ビルド
# ========================================
FROM node:20-alpine AS builder

WORKDIR /app

# 依存関係をコピー
COPY --from=deps /app/node_modules ./node_modules

# ソースコードをコピー
COPY . .

# Next.js ビルド
# 環境変数は実行時に渡す
ENV NEXT_TELEMETRY_DISABLED 1
RUN npm run build

# ========================================
# ステージ3: 実行環境
# ========================================
FROM node:20-alpine AS runner

WORKDIR /app

# 本番環境設定
ENV NODE_ENV production
ENV NEXT_TELEMETRY_DISABLED 1

# 非rootユーザー作成
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# 必要なファイルのみコピー
COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

# 非rootユーザーに切り替え
USER nextjs

# ポート公開
EXPOSE 3000

ENV PORT 3000
ENV HOSTNAME "0.0.0.0"

# ヘルスチェック
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD node -e "require('http').get('http://localhost:3000/', (r) => {process.exit(r.statusCode === 200 ? 0 : 1)})"

# サーバー起動
CMD ["node", "server.js"]
