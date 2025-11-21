# Next.js フロントエンド開発ガイド

## 目次
1. [開発環境セットアップ](#開発環境セットアップ)
2. [プロジェクト構成](#プロジェクト構成)
3. [開発の流れ](#開発の流れ)
4. [Next.js App Router](#nextjs-app-router)
5. [TypeScript](#typescript)
6. [Tailwind CSS](#tailwind-css)
7. [状態管理](#状態管理)
8. [API 連携](#api-連携)
9. [テスト](#テスト)

---

## 開発環境セットアップ

### 前提条件
- Node.js 18.17 以上
- npm または yarn
- Go バックエンド実行中（http://localhost:8080）

### セットアップ手順

```bash
# プロジェクトディレクトリへ移動
cd internalsystem/frontend

# 依存関係インストール
npm install

# 環境変数設定
cp .env.example .env.local

# 開発サーバー起動
npm run dev

# ブラウザで http://localhost:3000 を開く
```

### 環境変数設定

`.env.local` ファイル：

```bash
# API エンドポイント
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1

# アプリケーション設定
NEXT_PUBLIC_APP_NAME="Internal Web System"
NEXT_PUBLIC_APP_VERSION="1.0.0"

# ロギング
NEXT_PUBLIC_LOG_LEVEL=debug
```

---

## プロジェクト構成

```
frontend/
├── app/                           # App Router
│   ├── (auth)/
│   │   ├── login/
│   │   │   └── page.tsx          # ログインページ
│   │   └── logout/
│   │       └── route.ts          # ログアウト API Route
│   ├── (dashboard)/
│   │   ├── layout.tsx            # ダッシュボード共通レイアウト
│   │   ├── page.tsx              # ダッシュボードホーム
│   │   ├── users/
│   │   │   ├── page.tsx          # ユーザー一覧
│   │   │   └── [id]/
│   │   │       └── page.tsx      # ユーザー詳細
│   │   └── settings/
│   │       └── page.tsx          # 設定ページ
│   ├── api/                       # API Routes
│   │   └── auth/
│   │       └── refresh/
│   │           └── route.ts      # トークンリフレッシュ
│   ├── layout.tsx                # ルート レイアウト
│   └── page.tsx                  # ホームページ
│
├── components/                    # 再利用可能なコンポーネント
│   ├── ui/                       # 基本 UI コンポーネント
│   │   ├── Button.tsx
│   │   ├── Card.tsx
│   │   ├── Form.tsx
│   │   ├── Table.tsx
│   │   ├── Modal.tsx
│   │   └── ...
│   ├── layout/
│   │   ├── Header.tsx
│   │   ├── Sidebar.tsx
│   │   └── Footer.tsx
│   ├── auth/
│   │   ├── LoginForm.tsx
│   │   └── ProtectedRoute.tsx
│   └── dashboard/
│       ├── UsersList.tsx
│       ├── DashboardStats.tsx
│       └── ...
│
├── lib/                           # ユーティリティ関数
│   ├── api.ts                    # API クライアント
│   ├── auth.ts                   # 認証ヘルパー
│   ├── storage.ts                # LocalStorage ハンドリング
│   └── utils.ts                  # その他ユーティリティ
│
├── hooks/                         # カスタム React Hooks
│   ├── useAuth.ts                # 認証フック
│   ├── useUsers.ts               # ユーザー取得フック
│   └── ...
│
├── store/                         # Zustand ストア
│   ├── authStore.ts              # 認証ストア
│   ├── uiStore.ts                # UI 状態ストア
│   └── ...
│
├── types/                         # TypeScript 型定義
│   ├── api.ts                    # API レスポンス型
│   ├── models.ts                 # ドメインモデル型
│   └── ...
│
├── styles/
│   └── globals.css               # グローバルスタイル
│
├── public/                        # 静的ファイル
├── package.json
├── tsconfig.json
├── next.config.js
├── tailwind.config.ts
├── postcss.config.js
├── .eslintrc.json
└── README.md
```

---

## 開発の流れ

### 1. ページコンポーネント作成（app/）

```typescript
// app/(dashboard)/users/page.tsx
"use client"; // クライアントコンポーネント

import { useUsers } from "@/hooks/useUsers";
import { UserList } from "@/components/dashboard/UsersList";

export default function UsersPage() {
    const { data: users, isLoading, error } = useUsers();

    if (isLoading) return <div>読み込み中...</div>;
    if (error) return <div>エラー: {error.message}</div>;

    return (
        <div className="p-6">
            <h1 className="text-3xl font-bold mb-6">ユーザー管理</h1>
            <UserList users={users} />
        </div>
    );
}
```

### 2. コンポーネント作成（components/）

```typescript
// components/dashboard/UsersList.tsx
"use client";

import { User } from "@/types/models";
import { Table } from "@/components/ui/Table";

interface UsersListProps {
    users: User[];
}

export function UserList({ users }: UsersListProps) {
    const columns = [
        { key: "username", label: "ユーザー名" },
        { key: "email", label: "メールアドレス" },
        { key: "role", label: "ロール" },
        { key: "status", label: "ステータス" },
    ];

    return (
        <Table columns={columns} data={users} />
    );
}
```

### 3. UI コンポーネント（components/ui/）

```typescript
// components/ui/Button.tsx
import { ReactNode } from "react";

interface ButtonProps {
    children: ReactNode;
    onClick?: () => void;
    variant?: "primary" | "secondary" | "danger";
    disabled?: boolean;
}

export function Button({
    children,
    onClick,
    variant = "primary",
    disabled = false,
}: ButtonProps) {
    const variants = {
        primary: "bg-blue-600 text-white hover:bg-blue-700",
        secondary: "bg-gray-300 text-gray-800 hover:bg-gray-400",
        danger: "bg-red-600 text-white hover:bg-red-700",
    };

    return (
        <button
            onClick={onClick}
            disabled={disabled}
            className={`
                px-4 py-2 rounded-lg font-medium transition-colors
                ${variants[variant]}
                ${disabled ? "opacity-50 cursor-not-allowed" : ""}
            `}
        >
            {children}
        </button>
    );
}
```

### 4. カスタムフック作成（hooks/）

```typescript
// hooks/useUsers.ts
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { User } from "@/types/models";

export function useUsers(page = 1, limit = 20) {
    return useQuery({
        queryKey: ["users", page, limit],
        queryFn: async () => {
            const response = await api.get<{data: User[]}>("/users", {
                params: { page, limit },
            });
            return response.data.data;
        },
    });
}
```

### 5. 状態管理（store/）

```typescript
// store/authStore.ts
import { create } from "zustand";
import { User } from "@/types/models";

interface AuthState {
    user: User | null;
    token: string | null;
    isLoggedIn: boolean;
    setUser: (user: User | null) => void;
    setToken: (token: string | null) => void;
    logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
    user: null,
    token: null,
    isLoggedIn: false,

    setUser: (user) => set({ user, isLoggedIn: user !== null }),
    setToken: (token) => set({ token }),

    logout: () => set({
        user: null,
        token: null,
        isLoggedIn: false,
    }),
}));
```

### 6. API クライアント（lib/）

```typescript
// lib/api.ts
import axios, { AxiosInstance } from "axios";
import { useAuthStore } from "@/store/authStore";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

export const api: AxiosInstance = axios.create({
    baseURL: API_URL,
    headers: {
        "Content-Type": "application/json",
    },
});

// リクエストインターセプター
api.interceptors.request.use((config) => {
    const token = useAuthStore.getState().token;
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// レスポンスインターセプター
api.interceptors.response.use(
    (response) => response,
    async (error) => {
        if (error.response?.status === 401) {
            // トークン期限切れ → リフレッシュ
            try {
                const response = await axios.post("/api/auth/refresh");
                const newToken = response.data.token;
                useAuthStore.getState().setToken(newToken);

                // 元のリクエストを再試行
                error.config.headers.Authorization = `Bearer ${newToken}`;
                return api(error.config);
            } catch {
                useAuthStore.getState().logout();
            }
        }
        return Promise.reject(error);
    }
);
```

---

## Next.js App Router

### ファイル構成

```typescript
// app/layout.tsx - ルート レイアウト
export default function RootLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return (
        <html lang="ja">
            <body>
                <Header />
                <main>{children}</main>
                <Footer />
            </body>
        </html>
    );
}

// app/(dashboard)/layout.tsx - グループ レイアウト
export default function DashboardLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return (
        <div className="flex">
            <Sidebar />
            <div className="flex-1">{children}</div>
        </div>
    );
}
```

### 動的ルート

```typescript
// app/(dashboard)/users/[id]/page.tsx
interface UserDetailPageProps {
    params: {
        id: string;
    };
}

export default function UserDetailPage({ params }: UserDetailPageProps) {
    return <div>ユーザーID: {params.id}</div>;
}
```

---

## TypeScript

### 型定義ファイル

```typescript
// types/models.ts
export interface User {
    id: number;
    username: string;
    email: string;
    fullName: string;
    role: Role;
    status: "active" | "inactive" | "suspended";
    createdAt: string;
}

export interface Role {
    id: number;
    name: "admin" | "manager" | "user" | "viewer";
}

// types/api.ts
export interface ApiResponse<T> {
    code: number;
    message: string;
    data: T;
}

export interface ApiError {
    code: number;
    message: string;
    errors?: Array<{
        field: string;
        message: string;
    }>;
}
```

### 厳密性設定

`tsconfig.json` で strict mode を有効化：

```json
{
    "compilerOptions": {
        "strict": true,
        "strictNullChecks": true,
        "strictFunctionTypes": true,
        "noImplicitAny": true,
        "noUnusedLocals": true,
        "noUnusedParameters": true,
        "noImplicitReturns": true
    }
}
```

---

## Tailwind CSS

### ユーティリティクラス

```typescript
// components/ui/Card.tsx
export function Card({ children }: { children: React.ReactNode }) {
    return (
        <div className="bg-white rounded-lg shadow-md p-6 mb-4">
            {children}
        </div>
    );
}
```

### レスポンシブデザイン

```typescript
<div className="
    grid
    grid-cols-1 md:grid-cols-2 lg:grid-cols-3
    gap-4
">
    {/* レスポンシブグリッド */}
</div>
```

---

## 状態管理

### Zustand（クライアント状態）

```typescript
// UI 状態管理
const useUIStore = create((set) => ({
    sidebarOpen: true,
    toggleSidebar: () => set((state) => ({
        sidebarOpen: !state.sidebarOpen,
    })),
}));
```

### TanStack Query（サーバー状態）

```typescript
// データフェッチングと自動キャッシング
const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["users"],
    queryFn: () => api.get("/users"),
    staleTime: 1000 * 60 * 5, // 5分
});
```

---

## API 連携

### ログイン処理

```typescript
// components/auth/LoginForm.tsx
"use client";

import { useState } from "react";
import { api } from "@/lib/api";
import { useAuthStore } from "@/store/authStore";
import { useRouter } from "next/navigation";

export function LoginForm() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const router = useRouter();
    const { setUser, setToken } = useAuthStore();

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();

        try {
            const response = await api.post("/auth/login", {
                email,
                password,
            });

            setUser(response.data.data.user);
            setToken(response.data.data.token);

            // ダッシュボードへリダイレクト
            router.push("/dashboard");
        } catch (error) {
            console.error("ログイン失敗:", error);
        }
    };

    return (
        <form onSubmit={handleLogin}>
            <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="メールアドレス"
            />
            <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="パスワード"
            />
            <button type="submit">ログイン</button>
        </form>
    );
}
```

---

## テスト

### ユニットテスト（Vitest）

```typescript
// components/ui/__tests__/Button.test.ts
import { render, screen } from "@testing-library/react";
import { Button } from "../Button";

describe("Button コンポーネント", () => {
    it("テキストを表示する", () => {
        render(<Button>クリック</Button>);
        expect(screen.getByText("クリック")).toBeInTheDocument();
    });

    it("click イベントを発火する", async () => {
        const handleClick = vi.fn();
        render(<Button onClick={handleClick}>クリック</Button>);

        const button = screen.getByText("クリック");
        await user.click(button);

        expect(handleClick).toHaveBeenCalled();
    });
});
```

### E2E テスト（Cypress）

```typescript
// cypress/e2e/auth.cy.ts
describe("認証フロー", () => {
    it("ユーザーがログインできる", () => {
        cy.visit("/login");
        cy.get('input[type="email"]').type("test@example.com");
        cy.get('input[type="password"]').type("password");
        cy.get("button:contains('ログイン')").click();
        cy.url().should("include", "/dashboard");
    });
});
```

---

## 実行コマンド

```bash
# 開発サーバー起動
npm run dev

# 本番ビルド
npm run build

# 本番サーバー起動
npm start

# テスト実行
npm run test

# リント
npm run lint

# フォーマット
npm run format
```

---

このガイドを参考に、Next.js フロントエンドを開発してください。
