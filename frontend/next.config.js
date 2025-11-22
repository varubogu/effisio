/** @type {import('next').NextConfig} */
const nextConfig = {
  // 厳格モード
  reactStrictMode: true,

  // 本番環境でソースマップを生成
  productionBrowserSourceMaps: false,

  // SWCミニファイを使用
  swcMinify: true,

  // 出力形式（standalone: Dockerに最適）
  output: 'standalone',

  // 環境変数の公開
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
  },

  // 画像最適化
  images: {
    domains: [],
    formats: ['image/avif', 'image/webp'],
    minimumCacheTTL: 60,
  },

  // ヘッダー設定
  async headers() {
    return [
      {
        source: '/:path*',
        headers: [
          {
            key: 'X-DNS-Prefetch-Control',
            value: 'on',
          },
          {
            key: 'X-Frame-Options',
            value: 'SAMEORIGIN',
          },
          {
            key: 'X-Content-Type-Options',
            value: 'nosniff',
          },
          {
            key: 'Referrer-Policy',
            value: 'origin-when-cross-origin',
          },
        ],
      },
    ];
  },

  // リダイレクト設定（必要に応じて追加）
  async redirects() {
    return [];
  },

  // リライト設定（APIプロキシ等）
  async rewrites() {
    return [
      // 開発環境でのAPIプロキシ
      {
        source: '/api/v1/:path*',
        destination: `${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'}/:path*`,
      },
    ];
  },

  // Webpack設定のカスタマイズ
  webpack: (config, { isServer }) => {
    // カスタムWebpack設定をここに追加
    return config;
  },

  // 実験的機能
  experimental: {
    // 必要に応じて有効化
    // serverActions: true,
  },

  // TypeScript設定
  typescript: {
    // 本番ビルド時に型エラーを無視しない
    ignoreBuildErrors: false,
  },

  // ESLint設定
  eslint: {
    // 本番ビルド時にESLintエラーを無視しない
    ignoreDuringBuilds: false,
  },

  // パワードバイヘッダーを削除（セキュリティ向上）
  poweredByHeader: false,

  // 圧縮を有効化
  compress: true,

  // ページ拡張子
  pageExtensions: ['tsx', 'ts', 'jsx', 'js'],
};

// バンドルアナライザー（オプション）
if (process.env.ANALYZE === 'true') {
  const withBundleAnalyzer = require('@next/bundle-analyzer')({
    enabled: true,
  });
  module.exports = withBundleAnalyzer(nextConfig);
} else {
  module.exports = nextConfig;
}
