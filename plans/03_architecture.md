# システムアーキテクチャ設計

## 全体構成図

```
┌────────────────────────────────────────────────────────────────┐
│                    クライアント層                               │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │          Next.js フロントエンド (TypeScript + React)      │ │
│  │  ┌───────────────────────────────────────────────────┐  │ │
│  │  │  Pages (App Router)                            │  │ │
│  │  │  - Dashboard    - User Management               │  │ │
│  │  │  - Login        - Settings                      │  │ │
│  │  ├───────────────────────────────────────────────────┤  │ │
│  │  │  Components (Tailwind CSS)                       │  │ │
│  │  │  - Buttons, Forms, Tables, Charts               │  │ │
│  │  ├───────────────────────────────────────────────────┤  │ │
│  │  │  State Management                                │  │ │
│  │  │  - Zustand (Auth, UI State)                      │  │ │
│  │  │  - TanStack Query (Server State)                 │  │ │
│  │  └───────────────────────────────────────────────────┘  │ │
│  └──────────────────────────────────────────────────────────┘ │
│         WebBrowser                                             │
└────────────────────────────────────────────────────────────────┘
                            ↓ JSON/REST API (HTTP/HTTPS)
┌────────────────────────────────────────────────────────────────┐
│           Nginx リバースプロキシ / ロードバランサー             │
│  - SSL/TLS 終端                                               │
│  - API/フロント ルーティング                                     │
│  - CORS ハンドリング                                           │
└────────────────────────────────────────────────────────────────┘
          ↙ API              ↘ Static Assets
         ↙                     ↘
┌──────────────────────────┐  ┌──────────────────────────┐
│   Go API サーバー層       │  │  フロント静的ファイル      │
│  (Gin Framework)         │  │ (Next.js Build Output)  │
│  ┌──────────────────────┐│  │  - HTML, JS, CSS        │
│  │  RESTful Endpoints   ││  │  - Images, Fonts        │
│  │  - /api/v1/auth      ││  │  - CDN Cache 対応       │
│  │  - /api/v1/users     ││  └──────────────────────────┘
│  │  - /api/v1/dashboard ││
│  │  - /api/v1/audit-logs││
│  └──────────────────────┘│
│  ┌──────────────────────┐│
│  │ ビジネスロジック層    ││
│  │  - Service          ││
│  │  - Repository (ORM) ││
│  │  - Middleware       ││
│  │  - Auth/RBAC        ││
│  └──────────────────────┘│
│  ┌──────────────────────┐│
│  │ データアクセス層      ││
│  │  - GORM Models      ││
│  │  - Queries          ││
│  └──────────────────────┘│
└──────────────────────────┘
         ↓
    ┌─────────────────────────────────────────┐
    │      インフラストラクチャ層              │
    │  ┌──────────────────────────────────────┤
    │  │ PostgreSQL    Redis     ログ収集       │
    │  │ (GORM)       (Cache)   (Zap/ELK)    │
    │  │ データベース   セッション・  監視ログ    │
    │  │             キャッシュ                │
    │  └──────────────────────────────────────┤
    └─────────────────────────────────────────┘
```

---

## レイヤー構成

### 1. API層 (Presentation Layer)
- HTTP エンドポイント定義
- リクエスト/レスポンス の変換
- バリデーション（軽い）
- **責務**: HTTP通信の処理

### 2. ビジネスロジック層 (Business Logic Layer)
- ユースケースの実装
- ドメイン知識の実装
- トランザクション制御
- **責務**: ビジネス規則の実装

### 3. データアクセス層 (Data Access Layer)
- データベースアクセス
- キャッシュ操作
- ファイル操作
- **責務**: データの永続化・取得

### 4. モデル層 (Models)
- ドメインモデル
- DTOs (Data Transfer Objects)
- リクエスト/レスポンス構造体

---

## ディレクトリ構成

```
internalsystem/                    # リポジトリルート
│
├── backend/                       # Go バックエンド
│   ├── cmd/
│   │   └── server/
│   │       └── main.go           # API サーバー エントリーポイント
│   ├── internal/
│   │   ├── config/               # 設定管理
│   │   ├── models/               # ドメインモデル（GORM）
│   │   ├── repository/           # データアクセス層（GORM クエリ）
│   │   ├── service/              # ビジネスロジック層
│   │   ├── handler/              # HTTP ハンドラ（Gin エンドポイント）
│   │   ├── middleware/           # JWT 認証、CORS、ロギング
│   │   └── utils/                # ユーティリティ関数
│   ├── migrations/               # DB マイグレーション（golang-migrate）
│   ├── tests/
│   │   ├── unit/                # ユニットテスト
│   │   ├── integration/         # 統合テスト
│   │   └── fixtures/            # テストデータ
│   ├── docs/
│   │   └── swagger.json         # API ドキュメント（swaggo 生成）
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile               # 本番環境用ビルド
│   ├── Dockerfile.dev           # 開発環境用
│   ├── .golangci.yml            # golangci-lint 設定
│   └── README.md
│
├── frontend/                      # Next.js フロントエンド
│   ├── app/                      # App Router (Pages + Layout)
│   │   ├── (auth)/
│   │   │   ├── login/
│   │   │   │   └── page.tsx
│   │   │   └── logout/
│   │   │       └── route.ts      # API Route（ログアウト処理）
│   │   ├── (dashboard)/
│   │   │   ├── layout.tsx        # ダッシュボードレイアウト
│   │   │   ├── page.tsx          # ダッシュボードページ
│   │   │   ├── users/
│   │   │   │   ├── page.tsx      # ユーザー一覧
│   │   │   │   └── [id]/
│   │   │   │       └── page.tsx  # ユーザー詳細
│   │   │   └── settings/
│   │   │       └── page.tsx      # 設定ページ
│   │   ├── api/                  # API Routes（トークン更新等）
│   │   │   └── auth/
│   │   │       └── refresh/
│   │   │           └── route.ts
│   │   ├── layout.tsx            # ルート レイアウト
│   │   └── page.tsx              # ホームページ
│   │
│   ├── components/               # 再利用可能なコンポーネント
│   │   ├── ui/                  # 基本 UI コンポーネント
│   │   │   ├── Button.tsx
│   │   │   ├── Card.tsx
│   │   │   ├── Form.tsx
│   │   │   ├── Table.tsx
│   │   │   ├── Modal.tsx
│   │   │   └── ...
│   │   ├── layout/              # レイアウトコンポーネント
│   │   │   ├── Header.tsx
│   │   │   ├── Sidebar.tsx
│   │   │   └── Footer.tsx
│   │   ├── auth/                # 認証関連コンポーネント
│   │   │   ├── LoginForm.tsx
│   │   │   └── ProtectedRoute.tsx
│   │   └── dashboard/           # ダッシュボードコンポーネント
│   │       ├── UsersList.tsx
│   │       ├── DashboardStats.tsx
│   │       └── ...
│   │
│   ├── lib/                      # ユーティリティ関数
│   │   ├── api.ts               # API クライアント（Axios + TanStack Query）
│   │   ├── auth.ts              # 認証ヘルパー
│   │   ├── storage.ts           # LocalStorage ハンドリング
│   │   └── utils.ts             # その他ユーティリティ
│   │
│   ├── hooks/                    # カスタム React Hooks
│   │   ├── useAuth.ts           # 認証フック
│   │   ├── useUsers.ts          # ユーザー取得フック
│   │   ├── usePagination.ts     # ページネーションフック
│   │   └── ...
│   │
│   ├── store/                    # Zustand ストア（クライアント状態）
│   │   ├── authStore.ts         # 認証ストア
│   │   ├── uiStore.ts           # UI 状態ストア
│   │   └── ...
│   │
│   ├── types/                    # TypeScript 型定義
│   │   ├── api.ts               # API レスポンス型
│   │   ├── models.ts            # ドメインモデル型
│   │   └── ...
│   │
│   ├── styles/                   # グローバルスタイル
│   │   └── globals.css          # Tailwind + カスタムCSS
│   │
│   ├── public/                   # 静的ファイル
│   │   ├── images/
│   │   ├── fonts/
│   │   └── icons/
│   │
│   ├── .env.local               # ローカル環境変数
│   ├── .env.example             # 環境変数テンプレート
│   ├── tsconfig.json            # TypeScript 設定
│   ├── next.config.js           # Next.js 設定
│   ├── tailwind.config.ts       # Tailwind CSS 設定
│   ├── postcss.config.js        # PostCSS 設定
│   ├── package.json
│   ├── package-lock.json
│   ├── Dockerfile               # 本番環境用ビルド
│   ├── Dockerfile.dev           # 開発環境用
│   ├── .eslintrc.json           # ESLint 設定
│   └── README.md
│
├── docs/
│   ├── requirements/            # 要件定義
│   ├── architecture/            # アーキテクチャドキュメント
│   ├── api/                     # API 仕様書
│   ├── deployment/              # デプロイメント手順
│   └── ...
│
├── plans/                        # 計画書
│   ├── 01_project_overview.md
│   ├── 02_tech_stack.md
│   ├── ...
│
├── docker-compose.yml           # 開発環境構築（フロント・バック・DB統合）
├── docker-compose.prod.yml      # 本番環境構成（参考）
├── nginx.conf                   # Nginx 設定（フロント・バック統合）
├── .github/
│   └── workflows/
│       ├── backend-test.yml     # Go テスト・リント CI
│       ├── frontend-test.yml    # Next.js テスト・ビルド CI
│       └── deploy.yml           # デプロイメント CI/CD
├── .gitignore
└── README.md                    # プロジェクト全体 README
```

---

## ディレクトリ詳細

### cmd/server/
- `main.go`: アプリケーションのエントリーポイント
  - 初期化処理
  - サーバー起動

### internal/
Go Best Practice に基づいて `internal/` ディレクトリに非公開パッケージを配置

- **config/**: 設定ファイル読み込み、環境変数処理
- **models/**: 構造体定義（ドメインモデル、API レスポンス）
- **repository/**: DB アクセス、データベース操作
- **service/**: ビジネスロジック、ユースケース実装
- **handler/**: HTTPハンドラ、API エンドポイント
- **middleware/**: CORS、認証、ロギング、エラーハンドリング
- **utils/**: ヘルパー関数、共通処理

---

## マイクロサービス化への布石

現在はモノリシック構成だが、将来的には以下のサービスに分割可能：

- **Auth Service**: 認証・認可
- **User Service**: ユーザー管理
- **Data Service**: データ管理
- **Report Service**: レポート生成
- **Notification Service**: 通知

各サービスは gRPC または REST API で通信
