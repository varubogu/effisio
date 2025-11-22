import { api } from './api';
import type {
  LoginRequest,
  LoginResponse,
  RefreshTokenRequest,
  RefreshTokenResponse,
  ApiResponse,
} from '@/types/auth';

// ローカルストレージのキー
const ACCESS_TOKEN_KEY = 'access_token';
const REFRESH_TOKEN_KEY = 'refresh_token';

// トークンストレージ
export const tokenStorage = {
  getAccessToken: (): string | null => {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem(ACCESS_TOKEN_KEY);
  },

  getRefreshToken: (): string | null => {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem(REFRESH_TOKEN_KEY);
  },

  setTokens: (accessToken: string, refreshToken: string): void => {
    if (typeof window === 'undefined') return;
    localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
    localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
  },

  clearTokens: (): void => {
    if (typeof window === 'undefined') return;
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
  },
};

export const authApi = {
  // ログイン
  async login(data: LoginRequest): Promise<LoginResponse> {
    const response = await api.post<ApiResponse<LoginResponse>>('/auth/login', data);
    const loginData = response.data.data;

    // トークンを保存
    tokenStorage.setTokens(loginData.access_token, loginData.refresh_token);

    return loginData;
  },

  // トークンをリフレッシュ
  async refreshToken(): Promise<RefreshTokenResponse> {
    const refreshToken = tokenStorage.getRefreshToken();
    if (!refreshToken) {
      throw new Error('No refresh token available');
    }

    const data: RefreshTokenRequest = { refresh_token: refreshToken };
    const response = await api.post<ApiResponse<RefreshTokenResponse>>(
      '/auth/refresh',
      data
    );
    const refreshData = response.data.data;

    // 新しいトークンを保存
    tokenStorage.setTokens(refreshData.access_token, refreshData.refresh_token);

    return refreshData;
  },

  // ログアウト
  async logout(): Promise<void> {
    const refreshToken = tokenStorage.getRefreshToken();
    if (refreshToken) {
      try {
        const data: RefreshTokenRequest = { refresh_token: refreshToken };
        await api.post('/auth/logout', data);
      } catch (error) {
        // ログアウトエラーは無視（トークンが既に無効な場合など）
        console.warn('Logout request failed:', error);
      }
    }

    // トークンをクリア
    tokenStorage.clearTokens();
  },

  // 全セッションからログアウト
  async logoutAll(): Promise<void> {
    await api.post('/auth/logout-all');
    tokenStorage.clearTokens();
  },

  // 認証状態をチェック
  isAuthenticated: (): boolean => {
    return !!tokenStorage.getAccessToken();
  },
};
