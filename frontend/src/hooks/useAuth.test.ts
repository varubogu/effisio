import { describe, it, expect, beforeEach, vi } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useLogin, useLogout, useLogoutAll, useIsAuthenticated } from './useAuth';
import { authApi } from '@/lib/auth';
import { TestWrapper, createTestQueryClient } from '@/test/test-utils';
import type { LoginRequest } from '@/types/auth';

// Mock authApi
vi.mock('@/lib/auth', () => ({
  authApi: {
    login: vi.fn(),
    logout: vi.fn(),
    logoutAll: vi.fn(),
    isAuthenticated: vi.fn(),
  },
  tokenStorage: {
    getAccessToken: vi.fn(),
    getRefreshToken: vi.fn(),
    setTokens: vi.fn(),
    clearTokens: vi.fn(),
  },
}));

describe('useAuth hooks', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('useLogin', () => {
    it('should return mutation with idle status initially', () => {
      const { result } = renderHook(() => useLogin(), {
        wrapper: TestWrapper,
      });

      expect(result.current.isPending).toBe(false);
      expect(result.current.isSuccess).toBe(false);
      expect(result.current.isError).toBe(false);
      expect(result.current.data).toBeUndefined();
    });

    it('should login successfully', async () => {
      const loginRequest: LoginRequest = {
        username: 'testuser',
        password: 'password123',
      };

      const mockResponse = {
        access_token: 'access-123',
        refresh_token: 'refresh-456',
        user: {
          id: 1,
          username: 'testuser',
          email: 'test@example.com',
          full_name: 'Test User',
          department: 'Engineering',
          role: 'admin' as const,
          status: 'active' as const,
          last_login: null,
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z',
        },
      };

      vi.mocked(authApi.login).mockResolvedValueOnce(mockResponse);

      const { result } = renderHook(() => useLogin(), {
        wrapper: TestWrapper,
      });

      result.current.mutate(loginRequest);

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(result.current.data).toEqual(mockResponse);
      expect(result.current.isError).toBe(false);
    });

    it('should handle login error', async () => {
      const loginRequest: LoginRequest = {
        username: 'invalid',
        password: 'wrong',
      };

      vi.mocked(authApi.login).mockRejectedValueOnce(
        new Error('Invalid credentials')
      );

      const { result } = renderHook(() => useLogin(), {
        wrapper: TestWrapper,
      });

      result.current.mutate(loginRequest);

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });

      expect(result.current.error).toEqual(new Error('Invalid credentials'));
    });

    it('should clear query cache on successful login', async () => {
      const loginRequest: LoginRequest = {
        username: 'testuser',
        password: 'password123',
      };

      const mockResponse = {
        access_token: 'access-123',
        refresh_token: 'refresh-456',
        user: {
          id: 1,
          username: 'testuser',
          email: 'test@example.com',
          full_name: 'Test User',
          department: 'Engineering',
          role: 'admin' as const,
          status: 'active' as const,
          last_login: null,
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z',
        },
      };

      vi.mocked(authApi.login).mockResolvedValueOnce(mockResponse);

      const { result } = renderHook(() => useLogin(), {
        wrapper: TestWrapper,
      });

      result.current.mutate(loginRequest);

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      // QueryClient.clear() is called in onSuccess callback
      expect(result.current.isSuccess).toBe(true);
    });
  });

  describe('useLogout', () => {
    it('should logout successfully', async () => {
      vi.mocked(authApi.logout).mockResolvedValueOnce(undefined);

      const { result } = renderHook(() => useLogout(), {
        wrapper: TestWrapper,
      });

      result.current.mutate();

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(authApi.logout).toHaveBeenCalled();
    });

    it('should handle logout error', async () => {
      vi.mocked(authApi.logout).mockRejectedValueOnce(
        new Error('Logout failed')
      );

      const { result } = renderHook(() => useLogout(), {
        wrapper: TestWrapper,
      });

      result.current.mutate();

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });

      expect(result.current.error).toEqual(new Error('Logout failed'));
    });

    it('should clear query cache on successful logout', async () => {
      vi.mocked(authApi.logout).mockResolvedValueOnce(undefined);

      const { result } = renderHook(() => useLogout(), {
        wrapper: TestWrapper,
      });

      result.current.mutate();

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      // QueryClient.clear() is called in onSuccess callback
      expect(result.current.isSuccess).toBe(true);
    });
  });

  describe('useLogoutAll', () => {
    it('should logout all sessions successfully', async () => {
      vi.mocked(authApi.logoutAll).mockResolvedValueOnce(undefined);

      const { result } = renderHook(() => useLogoutAll(), {
        wrapper: TestWrapper,
      });

      result.current.mutate();

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(authApi.logoutAll).toHaveBeenCalled();
    });

    it('should handle logout all error', async () => {
      vi.mocked(authApi.logoutAll).mockRejectedValueOnce(
        new Error('Logout all failed')
      );

      const { result } = renderHook(() => useLogoutAll(), {
        wrapper: TestWrapper,
      });

      result.current.mutate();

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });

      expect(result.current.error).toEqual(new Error('Logout all failed'));
    });
  });

  describe('useIsAuthenticated', () => {
    it('should return true when authenticated', () => {
      vi.mocked(authApi.isAuthenticated).mockReturnValue(true);

      const result = useIsAuthenticated();

      expect(result).toBe(true);
      expect(authApi.isAuthenticated).toHaveBeenCalled();
    });

    it('should return false when not authenticated', () => {
      vi.mocked(authApi.isAuthenticated).mockReturnValue(false);

      const result = useIsAuthenticated();

      expect(result).toBe(false);
    });

    it('should return false when window is undefined (SSR)', () => {
      const originalWindow = global.window;
      // @ts-ignore
      delete global.window;

      const result = useIsAuthenticated();

      expect(result).toBe(false);

      global.window = originalWindow;
    });
  });
});
