import { useMutation, useQueryClient } from '@tanstack/react-query';
import { authApi } from '@/lib/auth';
import type { LoginRequest } from '@/types/auth';

// ログイン
export function useLogin() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: LoginRequest) => authApi.login(data),
    onSuccess: () => {
      // ログイン成功時にキャッシュをクリア
      queryClient.clear();
    },
  });
}

// ログアウト
export function useLogout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => authApi.logout(),
    onSuccess: () => {
      // ログアウト成功時にキャッシュをクリア
      queryClient.clear();
    },
  });
}

// 全セッションからログアウト
export function useLogoutAll() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => authApi.logoutAll(),
    onSuccess: () => {
      // ログアウト成功時にキャッシュをクリア
      queryClient.clear();
    },
  });
}

// 認証状態を確認
export function useIsAuthenticated(): boolean {
  if (typeof window === 'undefined') {
    return false;
  }
  return authApi.isAuthenticated();
}
