# TypeScript/React コーディング規約

Effisioプロジェクトにおける TypeScript および React のコーディング規約です。

## 基本原則

1. **型安全性**: strict mode を有効にし、any の使用を避ける
2. **関数型プログラミング**: 純粋関数と不変性を重視
3. **コンポーネント設計**: 単一責任の原則に従う

---

## TypeScript

### 型定義

```typescript
// ✅ 良い例: 明示的な型定義
interface User {
  id: number;
  username: string;
  email: string;
  role: 'admin' | 'user' | 'viewer';
}

type CreateUserRequest = Omit<User, 'id'>;

// ❌ 悪い例: any の使用
const user: any = getUserData();
```

### 型ガード

```typescript
// ✅ 良い例: 型ガードの使用
function isUser(obj: unknown): obj is User {
  return (
    typeof obj === 'object' &&
    obj !== null &&
    'id' in obj &&
    'username' in obj
  );
}

if (isUser(data)) {
  console.log(data.username); // 型安全
}
```

### Enum vs Union Types

```typescript
// ✅ 良い例: Union Types（推奨）
type UserRole = 'admin' | 'manager' | 'user' | 'viewer';

// ⚠️ Enum は必要な場合のみ
enum HttpStatus {
  OK = 200,
  BadRequest = 400,
  Unauthorized = 401,
}
```

---

## React コンポーネント

### 関数コンポーネント

```typescript
// ✅ 良い例: 関数コンポーネント + TypeScript
interface ButtonProps {
  children: React.ReactNode;
  onClick?: () => void;
  variant?: 'primary' | 'secondary';
  disabled?: boolean;
}

export function Button({
  children,
  onClick,
  variant = 'primary',
  disabled = false,
}: ButtonProps) {
  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className={`btn btn-${variant}`}
    >
      {children}
    </button>
  );
}

// ❌ 悪い例: クラスコンポーネント（新規作成時は避ける）
class Button extends React.Component<ButtonProps> {
  render() {
    // ...
  }
}
```

### フック

```typescript
// ✅ 良い例: カスタムフックの型定義
function useUsers(page = 1, limit = 20) {
  return useQuery<User[], Error>({
    queryKey: ['users', page, limit],
    queryFn: async () => {
      const response = await api.get<ApiResponse<User[]>>('/users', {
        params: { page, limit },
      });
      return response.data.data;
    },
  });
}

// 使用例
const { data, isLoading, error } = useUsers();
```

### Props の分割代入

```typescript
// ✅ 良い例: 分割代入
function UserCard({ user, onEdit }: { user: User; onEdit: (id: number) => void }) {
  return (
    <div>
      <h3>{user.username}</h3>
      <button onClick={() => onEdit(user.id)}>Edit</button>
    </div>
  );
}

// ❌ 悪い例: props をそのまま使用
function UserCard(props: { user: User; onEdit: (id: number) => void }) {
  return (
    <div>
      <h3>{props.user.username}</h3>
      <button onClick={() => props.onEdit(props.user.id)}>Edit</button>
    </div>
  );
}
```

---

## ファイル構成

### コンポーネントファイル

```
components/
├── ui/
│   ├── Button.tsx
│   ├── Card.tsx
│   └── index.ts        # エクスポートをまとめる
├── layout/
│   ├── Header.tsx
│   ├── Sidebar.tsx
│   └── index.ts
└── dashboard/
    ├── UsersList.tsx
    ├── DashboardStats.tsx
    └── index.ts
```

### 命名規則

```
// コンポーネント: PascalCase
Button.tsx
UsersList.tsx
DashboardStats.tsx

// フック: use + PascalCase
useAuth.ts
useUsers.ts
usePagination.ts

// ユーティリティ: camelCase
api.ts
utils.ts
formatDate.ts
```

---

## 状態管理

### Zustand

```typescript
// ✅ 良い例: Zustand ストアの型定義
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
  logout: () => set({ user: null, token: null, isLoggedIn: false }),
}));

// 使用例
const { user, isLoggedIn, logout } = useAuthStore();
```

### TanStack Query

```typescript
// ✅ 良い例: TanStack Query の型定義
function useUsers() {
  return useQuery({
    queryKey: ['users'],
    queryFn: async () => {
      const response = await api.get<ApiResponse<User[]>>('/users');
      return response.data.data;
    },
    staleTime: 5 * 60 * 1000, // 5分
  });
}

// ミューテーション
function useCreateUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (user: CreateUserRequest) => {
      const response = await api.post<ApiResponse<User>>('/users', user);
      return response.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });
}
```

---

## 非同期処理

### async/await

```typescript
// ✅ 良い例: async/await + エラーハンドリング
async function loginUser(email: string, password: string): Promise<User> {
  try {
    const response = await api.post<ApiResponse<{ user: User; token: string }>>(
      '/auth/login',
      { email, password }
    );
    return response.data.data.user;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      throw new Error(error.response?.data.message || 'Login failed');
    }
    throw error;
  }
}
```

---

## イベントハンドラ

### 命名規則

```typescript
// ✅ 良い例: handle + Action
function LoginForm() {
  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    // ...
  };

  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value);
  };

  return (
    <form onSubmit={handleSubmit}>
      <input type="email" onChange={handleEmailChange} />
    </form>
  );
}
```

---

## スタイリング（Tailwind CSS）

### クラス名の整理

```typescript
// ✅ 良い例: clsx または cn ヘルパー関数を使用
import { clsx } from 'clsx';

function Button({ variant, disabled, className }: ButtonProps) {
  return (
    <button
      className={clsx(
        'px-4 py-2 rounded-lg font-medium transition-colors',
        {
          'bg-blue-600 text-white hover:bg-blue-700': variant === 'primary',
          'bg-gray-300 text-gray-800 hover:bg-gray-400': variant === 'secondary',
          'opacity-50 cursor-not-allowed': disabled,
        },
        className
      )}
    >
      {children}
    </button>
  );
}
```

---

## パフォーマンス最適化

### React.memo

```typescript
// ✅ 良い例: メモ化が必要な場合のみ使用
export const UserCard = React.memo(({ user }: { user: User }) => {
  return (
    <div className="card">
      <h3>{user.username}</h3>
      <p>{user.email}</p>
    </div>
  );
});
```

### useMemo / useCallback

```typescript
// ✅ 良い例: 計算コストが高い場合のみ useMemo
function UsersList({ users }: { users: User[] }) {
  const activeUsers = useMemo(
    () => users.filter((user) => user.status === 'active'),
    [users]
  );

  const handleUserClick = useCallback((id: number) => {
    console.log('User clicked:', id);
  }, []);

  return <div>{/* ... */}</div>;
}
```

---

## テスト

### コンポーネントテスト

```typescript
// ✅ 良い例: React Testing Library
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Button } from './Button';

describe('Button', () => {
  it('renders with children', () => {
    render(<Button>Click me</Button>);
    expect(screen.getByText('Click me')).toBeInTheDocument();
  });

  it('calls onClick when clicked', async () => {
    const handleClick = vi.fn();
    render(<Button onClick={handleClick}>Click me</Button>);

    const button = screen.getByRole('button');
    await userEvent.click(button);

    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('is disabled when disabled prop is true', () => {
    render(<Button disabled>Click me</Button>);
    expect(screen.getByRole('button')).toBeDisabled();
  });
});
```

---

## 禁止事項

❌ **any の乱用**: 型安全性が失われる
❌ **クラスコンポーネント**: 関数コンポーネントを使用
❌ **インラインスタイル**: Tailwind CSS を使用
❌ **useEffect の依存配列の省略**: 必ず指定
❌ **defaultProps**: TypeScript のデフォルト引数を使用

---

## 推奨ツール

- **ESLint**: リント（必須）
- **Prettier**: コード整形（必須）
- **TypeScript**: 型チェック（必須）
- **Vitest**: ユニットテスト
- **React Testing Library**: コンポーネントテスト
- **Cypress / Playwright**: E2Eテスト

## チェックコマンド

```bash
# 型チェック
npm run type-check

# リント
npm run lint

# フォーマット
npm run format

# テスト
npm run test

# ビルド
npm run build
```

---

詳細は [React TypeScript Cheatsheet](https://react-typescript-cheatsheet.netlify.app/) を参照してください。
