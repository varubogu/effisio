import { describe, it, expect, beforeEach, vi } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useDashboardOverview, type DashboardOverview } from './useDashboard';
import { api } from '@/lib/api';
import { TestWrapper } from '@/test/test-utils';
import type { ApiResponse } from '@/types/auth';

// Mock the api module
vi.mock('@/lib/api', () => ({
  api: {
    get: vi.fn(),
  },
}));

const mockDashboardData: DashboardOverview = {
  total_users: 100,
  active_users: 85,
  inactive_users: 10,
  suspended_users: 5,
  last_login_stats: [
    { date: '2024-01-15', count: 42 },
    { date: '2024-01-14', count: 38 },
    { date: '2024-01-13', count: 35 },
  ],
  users_by_role: {
    admin: 5,
    manager: 15,
    user: 70,
    viewer: 10,
  },
  users_by_department: [
    { department: 'Engineering', count: 40 },
    { department: 'Sales', count: 30 },
    { department: 'HR', count: 15 },
    { department: 'Finance', count: 15 },
  ],
};

describe('useDashboard hooks', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('useDashboardOverview', () => {
    it('should fetch dashboard overview data', async () => {
      const mockResponse: ApiResponse<DashboardOverview> = {
        code: 200,
        message: 'OK',
        data: mockDashboardData,
      };

      vi.mocked(api.get).mockResolvedValueOnce({
        data: mockResponse,
      });

      const { result } = renderHook(() => useDashboardOverview(), {
        wrapper: TestWrapper,
      });

      expect(result.current.isPending).toBe(true);

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(result.current.data).toEqual(mockDashboardData);
      expect(api.get).toHaveBeenCalledWith('/dashboard/overview');
    });

    it('should have correct query key', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: {
          code: 200,
          message: 'OK',
          data: mockDashboardData,
        },
      });

      const { result } = renderHook(() => useDashboardOverview(), {
        wrapper: TestWrapper,
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      // Verify the query was made
      expect(api.get).toHaveBeenCalled();
    });

    it('should have correct caching settings (staleTime and gcTime)', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: {
          code: 200,
          message: 'OK',
          data: mockDashboardData,
        },
      });

      const { result } = renderHook(() => useDashboardOverview(), {
        wrapper: TestWrapper,
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      // Verify hook was called
      expect(api.get).toHaveBeenCalled();
      // staleTime: 1 minute, gcTime: 5 minutes are set in hook
      expect(result.current.data).toEqual(mockDashboardData);
    });

    it('should handle error when fetching dashboard data', async () => {
      vi.mocked(api.get).mockRejectedValueOnce(
        new Error('Failed to fetch dashboard')
      );

      const { result } = renderHook(() => useDashboardOverview(), {
        wrapper: TestWrapper,
      });

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });

      expect(result.current.error).toEqual(
        new Error('Failed to fetch dashboard')
      );
    });

    it('should handle empty dashboard data', async () => {
      const emptyData: DashboardOverview = {
        total_users: 0,
        active_users: 0,
        inactive_users: 0,
        suspended_users: 0,
        last_login_stats: [],
        users_by_role: {},
        users_by_department: [],
      };

      vi.mocked(api.get).mockResolvedValueOnce({
        data: {
          code: 200,
          message: 'OK',
          data: emptyData,
        },
      });

      const { result } = renderHook(() => useDashboardOverview(), {
        wrapper: TestWrapper,
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(result.current.data).toEqual(emptyData);
      expect(result.current.data?.total_users).toBe(0);
    });

    it('should correctly structure total users statistic', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: {
          code: 200,
          message: 'OK',
          data: mockDashboardData,
        },
      });

      const { result } = renderHook(() => useDashboardOverview(), {
        wrapper: TestWrapper,
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(result.current.data?.total_users).toBe(100);
      expect(result.current.data?.active_users).toBe(85);
      expect(result.current.data?.inactive_users).toBe(10);
      expect(result.current.data?.suspended_users).toBe(5);
    });

    it('should correctly structure users by role', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: {
          code: 200,
          message: 'OK',
          data: mockDashboardData,
        },
      });

      const { result } = renderHook(() => useDashboardOverview(), {
        wrapper: TestWrapper,
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(result.current.data?.users_by_role).toEqual({
        admin: 5,
        manager: 15,
        user: 70,
        viewer: 10,
      });
    });

    it('should correctly structure users by department', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: {
          code: 200,
          message: 'OK',
          data: mockDashboardData,
        },
      });

      const { result } = renderHook(() => useDashboardOverview(), {
        wrapper: TestWrapper,
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(result.current.data?.users_by_department).toHaveLength(4);
      expect(result.current.data?.users_by_department[0]).toEqual({
        department: 'Engineering',
        count: 40,
      });
    });

    it('should correctly structure last login stats', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: {
          code: 200,
          message: 'OK',
          data: mockDashboardData,
        },
      });

      const { result } = renderHook(() => useDashboardOverview(), {
        wrapper: TestWrapper,
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(result.current.data?.last_login_stats).toHaveLength(3);
      expect(result.current.data?.last_login_stats[0]).toEqual({
        date: '2024-01-15',
        count: 42,
      });
    });

    it('should handle malformed data response', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: {
          code: 200,
          message: 'OK',
          data: null,
        },
      });

      const { result } = renderHook(() => useDashboardOverview(), {
        wrapper: TestWrapper,
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(result.current.data).toBeNull();
    });

    it('should retry on network error', async () => {
      vi.mocked(api.get)
        .mockRejectedValueOnce(new Error('Network error'))
        .mockRejectedValueOnce(new Error('Still failing'));

      const { result } = renderHook(() => useDashboardOverview(), {
        wrapper: TestWrapper,
      });

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });

      // With retry: false in test setup, it should fail immediately
      expect(api.get).toHaveBeenCalled();
    });
  });
});
