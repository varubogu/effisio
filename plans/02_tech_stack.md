# 技術スタック計画（Go + Next.js）

## バックエンド（Go）

### 言語・ランタイム
- **Go**: 1.21以上
  - 理由：高速、並行処理が強力、単一バイナリでデプロイ可能、API サーバーに最適

### Webフレームワーク
- **Gin Gonic** ✅ （確定）
  - 理由：高速（Echo 同等の性能）、ミドルウェア豊富、ドキュメント充実、業界採用例多数
  - ベンチマーク：Echo との差はほぼなし、開発体験で Gin を選定

### データベース
- **PostgreSQL 13.0+** ✅ （確定）
  - 理由：ACID保証、複雑なクエリ対応、スケーラビリティ、JSON型対応

### ORM/クエリビルダ
- **GORM v1.25+** ✅ （確定）
  - 理由：開発生産性重視、リレーション管理が容易、フックシステムが強力
  - Next.js フロントから API 経由で常にシリアライズされたデータ受け取り

### キャッシュ
- **Redis 6.0+**: セッション、ダッシュボード統計キャッシュ層

### 認証・認可
- **JWT（RS256）**: ステートレス認証
  - Access Token: 1時間
  - Refresh Token: 30日（HttpOnly Cookie で安全に保存）
- **RBAC**: ロールベースアクセス制御（4ロール）

### ロギング
- **Zap v1.26+**: 構造化ログ（JSON形式）
  - 理由：高性能、本番環境対応、分散トレーシング対応

### バリデーション
- **github.com/go-playground/validator/v10**

---

## フロントエンド（Next.js）

### フレームワーク・言語
- **Next.js 14+** with **App Router** ✅
  - 理由：最新パターン、Server Components で SEO 対応、API Routes で簡易バック実装可
  - TypeScript を標準で使用

### 言語
- **TypeScript**: 型安全性重視
  - strict mode 有効

### スタイリング
- **Tailwind CSS 3.0+**: ユーティリティファースト CSS
  - 理由：SaaS 公開時のモダンな UI/UX に最適、カスタマイズ容易

### データフェッチング・キャッシング
- **TanStack Query（React Query） v5+** ✅
  - 理由：サーバー状態管理専門、キャッシング・再試行・同期が自動
  - Go API との連携に最適

### 状態管理
- **Zustand v4+** ✅
  - 理由：軽量、シンプル、TypeScript 対応、Redux より学習コスト低い
  - 用途：認証状態、UI状態（サイドバー開閉等）

### コンポーネント開発
- **Storybook 7+** （オプション）
  - 理由：コンポーネント設計ドキュメント化、単体テスト容易

### フォーム管理
- **React Hook Form v7+**
  - 理由：パフォーマンス優秀、バリデーション簡単、TypeScript 対応

### テスト
- **Vitest**: ユニットテスト
  - 理由：Vite ベース、高速、Jest 互換
- **React Testing Library**: コンポーネントテスト
- **Cypress** または **Playwright**: E2E テスト

### UI コンポーネントライブラリ
- **Shadcn/ui** または **Headless UI** + Tailwind
  - 理由：カスタマイズ可能、アクセシビリティ対応、Tailwind との相性良好

### HTTP クライアント
- **Axios** または **fetch API** (組み込み)
  - TanStack Query が上位でハンドル

---

## インフラストラクチャ

### コンテナ化
- **Docker**: 開発環境・本番環境の統一
  - バックエンド：Dockerfile（マルチステージビルド）
  - フロントエンド：Dockerfile（Node → Static Build）
- **Docker Compose**: ローカル開発環境
  - Go API、PostgreSQL、Redis、Next.js (dev server)

### CI/CD
- **GitHub Actions** ✅
  - ワークフロー：
    1. テスト（バック・フロント）
    2. リント・型チェック
    3. ビルド
    4. Docker イメージプッシュ
    5. ステージング デプロイ（develop ブランチ）
    6. 本番デプロイ（main ブランチ）

### デプロイメント
- **初期段階**: Docker Compose on Linux Server
- **将来**: Kubernetes / AWS ECS

### リバースプロキシ
- **Nginx**: フロント・バック の統合、SSL/TLS 終端

### モニタリング・ロギング
- **Prometheus**: メトリクス収集（Go サーバー）
- **Grafana**: ダッシュボード
- **ELK Stack** or **Datadog**: ログ集約

---

## 開発ツール（Go）

### テスト
- **testing**: ユニットテスト（標準）
- **testify v1.8+**: アサーション、モック
- **httptest**: HTTP テスト（標準）

### リント・フォーマット
- **golangci-lint v1.55+**: 複合リンター
  - 実行：`golangci-lint run ./...`
  - チェック項目：vet, errcheck, gosec, staticcheck 等

### デバッグ
- **Delve**: Go デバッガ

### API ドキュメンテーション
- **Swagger/OpenAPI 3.0**: **swaggo** で自動生成
  - エンドポイント：`/swagger/index.html`

---

## 開発ツール（Next.js）

### テスト
- **Vitest**: ユニットテスト
  - 理由：Vite ベース、Jest 互換、高速
- **React Testing Library**: コンポーネントテスト
- **Cypress** または **Playwright**: E2E テスト

### リント・フォーマット
- **ESLint**: JavaScript/TypeScript リント
- **Prettier**: コード整形
- **TypeScript**: 型チェック（`tsc --noEmit`）

### 開発サーバー
- **Next.js Dev Server**: `next dev`
  - Fast Refresh で高速開発体験

---

## 開発環境マシン要件

### CPU/メモリ
- **CPU**: 4 コア以上（推奨）
- **メモリ**: 8GB 以上（推奨）
  - Go ビルド、PostgreSQL、Redis、Next.js dev server を同時実行

### ストレージ
- **SSD**: 50GB 以上（node_modules 等が容量を使う）

### ネットワーク
- npm / cargo パッケージダウンロードに安定したインターネット接続が必要

---

## 学習ロードマップ

### バックエンド（Go）学習経路
1. **基礎**: Go 文法、ゴルーチン、チャネル
2. **Gin**: ルーティング、ミドルウェア、リクエスト/レスポンス処理
3. **GORM**: モデル定義、リレーション、フック
4. **JWT**: 認証フロー、トークン管理
5. **テスト**: ユニットテスト、モック、統合テスト

### フロントエンド（Next.js）学習経路
1. **React 基礎**: JSX、コンポーネント、フック
2. **Next.js App Router**: ファイルベースルーティング、Server Components
3. **TypeScript**: 型定義、厳密性
4. **Tailwind CSS**: ユーティリティクラス、レスポンシブデザイン
5. **TanStack Query**: サーバー状態管理、キャッシング
6. **Zustand**: クライアント状態管理

---

## 決定済み項目

- ✅ **Webフレームワーク**: Gin
- ✅ **ORM**: GORM
- ✅ **フロントエンド**: Next.js 14+ (App Router)
- ✅ **スタイリング**: Tailwind CSS
- ✅ **状態管理**: Zustand（クライアント）+ TanStack Query（サーバー）
- ✅ **デプロイ**: Docker Compose（初期） → Kubernetes（将来）
