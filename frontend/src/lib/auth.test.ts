import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { authApi, tokenStorage } from './auth';
import { api } from './api';
import type { LoginRequest, LoginResponse, RefreshTokenResponse } from '@/types/auth';

// Mock the api module
vi.mock('./api', () => ({
  api: {
    post: vi.fn(),
  },
}));

// Setup localStorage mock
const localStorageMock = (() => {
  let store: Record<string, string> = {};

  return {
    getItem: (key: string) => store[key] || null,
    setItem: (key: string, value: string) => {
      store[key] = value.toString();
    },
    removeItem: (key: string) => {
      delete store[key];
    },
    clear: () => {
      store = {};
    },
  };
})();

Object.defineProperty(window, 'localStorage', {
  value: localStorageMock,
});

describe('tokenStorage', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  describe('getAccessToken', () => {
    it('should return null when no token is stored', () => {
      const token = tokenStorage.getAccessToken();
      expect(token).toBeNull();
    });

    it('should return the access token when it is stored', () => {
      localStorage.setItem('access_token', 'test-access-token');
      const token = tokenStorage.getAccessToken();
      expect(token).toBe('test-access-token');
    });

    it('should return null when window is undefined (SSR)', () => {
      const originalWindow = global.window;
      // @ts-ignore
      delete global.window;

      const token = tokenStorage.getAccessToken();
      expect(token).toBeNull();

      global.window = originalWindow;
    });
  });

  describe('getRefreshToken', () => {
    it('should return null when no token is stored', () => {
      const token = tokenStorage.getRefreshToken();
      expect(token).toBeNull();
    });

    it('should return the refresh token when it is stored', () => {
      localStorage.setItem('refresh_token', 'test-refresh-token');
      const token = tokenStorage.getRefreshToken();
      expect(token).toBe('test-refresh-token');
    });
  });

  describe('setTokens', () => {
    it('should set both access and refresh tokens', () => {
      tokenStorage.setTokens('access-123', 'refresh-456');

      expect(localStorage.getItem('access_token')).toBe('access-123');
      expect(localStorage.getItem('refresh_token')).toBe('refresh-456');
    });

    it('should overwrite existing tokens', () => {
      localStorage.setItem('access_token', 'old-access');
      localStorage.setItem('refresh_token', 'old-refresh');

      tokenStorage.setTokens('new-access', 'new-refresh');

      expect(localStorage.getItem('access_token')).toBe('new-access');
      expect(localStorage.getItem('refresh_token')).toBe('new-refresh');
    });

    it('should not throw when window is undefined (SSR)', () => {
      const originalWindow = global.window;
      // @ts-ignore
      delete global.window;

      expect(() => {
        tokenStorage.setTokens('access', 'refresh');
      }).not.toThrow();

      global.window = originalWindow;
    });
  });

  describe('clearTokens', () => {
    it('should remove both tokens', () => {
      localStorage.setItem('access_token', 'access');
      localStorage.setItem('refresh_token', 'refresh');

      tokenStorage.clearTokens();

      expect(localStorage.getItem('access_token')).toBeNull();
      expect(localStorage.getItem('refresh_token')).toBeNull();
    });

    it('should not throw if tokens are not set', () => {
      expect(() => {
        tokenStorage.clearTokens();
      }).not.toThrow();
    });
  });
});

describe('authApi', () => {
  beforeEach(() => {
    localStorage.clear();
    vi.clearAllMocks();
  });

  describe('login', () => {
    it('should login successfully and store tokens', async () => {
      const loginRequest: LoginRequest = {
        username: 'testuser',
        password: 'password123',
      };

      const mockLoginResponse: LoginResponse = {
        access_token: 'access-token-123',
        refresh_token: 'refresh-token-456',
        user: {
          id: 1,
          username: 'testuser',
          email: 'test@example.com',
          full_name: 'Test User',
          department: 'Engineering',
          role: 'admin',
          status: 'active',
          last_login: null,
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z',
        },
      };

      vi.mocked(api.post).mockResolvedValueOnce({
        data: {
          code: 200,
          message: 'OK',
          data: mockLoginResponse,
        },
      });

      const result = await authApi.login(loginRequest);

      expect(result).toEqual(mockLoginResponse);
      expect(localStorage.getItem('access_token')).toBe('access-token-123');
      expect(localStorage.getItem('refresh_token')).toBe('refresh-token-456');
      expect(api.post).toHaveBeenCalledWith('/auth/login', loginRequest);
    });

    it('should throw error on login failure', async () => {
      const loginRequest: LoginRequest = {
        username: 'invalid',
        password: 'wrong',
      };

      vi.mocked(api.post).mockRejectedValueOnce(
        new Error('Invalid credentials')
      );

      await expect(authApi.login(loginRequest)).rejects.toThrow(
        'Invalid credentials'
      );
      expect(localStorage.getItem('access_token')).toBeNull();
    });
  });

  describe('refreshToken', () => {
    it('should refresh token successfully', async () => {
      localStorage.setItem('refresh_token', 'old-refresh-token');

      const mockRefreshResponse: RefreshTokenResponse = {
        access_token: 'new-access-token',
        refresh_token: 'new-refresh-token',
      };

      vi.mocked(api.post).mockResolvedValueOnce({
        data: {
          code: 200,
          message: 'OK',
          data: mockRefreshResponse,
        },
      });

      const result = await authApi.refreshToken();

      expect(result).toEqual(mockRefreshResponse);
      expect(localStorage.getItem('access_token')).toBe('new-access-token');
      expect(localStorage.getItem('refresh_token')).toBe('new-refresh-token');
      expect(api.post).toHaveBeenCalledWith('/auth/refresh', {
        refresh_token: 'old-refresh-token',
      });
    });

    it('should throw error when no refresh token available', async () => {
      await expect(authApi.refreshToken()).rejects.toThrow(
        'No refresh token available'
      );
    });

    it('should throw error on refresh failure', async () => {
      localStorage.setItem('refresh_token', 'valid-token');

      vi.mocked(api.post).mockRejectedValueOnce(
        new Error('Token expired')
      );

      await expect(authApi.refreshToken()).rejects.toThrow('Token expired');
    });
  });

  describe('logout', () => {
    it('should logout successfully and clear tokens', async () => {
      localStorage.setItem('refresh_token', 'token-to-logout');
      localStorage.setItem('access_token', 'access-token');

      vi.mocked(api.post).mockResolvedValueOnce({ data: null });

      await authApi.logout();

      expect(localStorage.getItem('access_token')).toBeNull();
      expect(localStorage.getItem('refresh_token')).toBeNull();
      expect(api.post).toHaveBeenCalledWith('/auth/logout', {
        refresh_token: 'token-to-logout',
      });
    });

    it('should clear tokens even if logout request fails', async () => {
      localStorage.setItem('refresh_token', 'token-to-logout');
      localStorage.setItem('access_token', 'access-token');

      vi.mocked(api.post).mockRejectedValueOnce(
        new Error('Server error')
      );

      // Should not throw
      await authApi.logout();

      expect(localStorage.getItem('access_token')).toBeNull();
      expect(localStorage.getItem('refresh_token')).toBeNull();
    });

    it('should clear tokens when no refresh token exists', async () => {
      localStorage.setItem('access_token', 'access-token');

      await authApi.logout();

      expect(localStorage.getItem('access_token')).toBeNull();
      expect(api.post).not.toHaveBeenCalled();
    });
  });

  describe('logoutAll', () => {
    it('should logout from all sessions and clear tokens', async () => {
      localStorage.setItem('access_token', 'access-token');
      localStorage.setItem('refresh_token', 'refresh-token');

      vi.mocked(api.post).mockResolvedValueOnce({ data: null });

      await authApi.logoutAll();

      expect(localStorage.getItem('access_token')).toBeNull();
      expect(localStorage.getItem('refresh_token')).toBeNull();
      expect(api.post).toHaveBeenCalledWith('/auth/logout-all');
    });

    it('should throw error on logout failure', async () => {
      vi.mocked(api.post).mockRejectedValueOnce(
        new Error('Server error')
      );

      await expect(authApi.logoutAll()).rejects.toThrow('Server error');
    });
  });

  describe('isAuthenticated', () => {
    it('should return true when access token exists', () => {
      localStorage.setItem('access_token', 'token-123');

      expect(authApi.isAuthenticated()).toBe(true);
    });

    it('should return false when no access token exists', () => {
      localStorage.clear();

      expect(authApi.isAuthenticated()).toBe(false);
    });

    it('should return false when access token is empty string', () => {
      localStorage.setItem('access_token', '');

      expect(authApi.isAuthenticated()).toBe(false);
    });
  });
});
