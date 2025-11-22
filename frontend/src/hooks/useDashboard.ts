import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api';
import type { ApiResponse } from '@/types/auth';

export interface DashboardOverview {
  total_users: number;
  active_users: number;
  inactive_users: number;
  suspended_users: number;
  last_login_stats: Array<{
    date: string;
    count: number;
  }>;
  users_by_role: Record<string, number>;
  users_by_department: Array<{
    department: string;
    count: number;
  }>;
}

const dashboardApi = {
  // ダッシュボード概要を取得
  async getOverview(): Promise<DashboardOverview> {
    const response = await api.get<ApiResponse<DashboardOverview>>('/dashboard/overview');
    return response.data.data;
  },
};

// ダッシュボード概要を取得
export function useDashboardOverview() {
  return useQuery({
    queryKey: ['dashboard', 'overview'],
    queryFn: () => dashboardApi.getOverview(),
    staleTime: 1 * 60 * 1000, // 1分間キャッシュ
    gcTime: 5 * 60 * 1000, // 5分間メモリに保持
  });
}
