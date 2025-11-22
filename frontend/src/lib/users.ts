import { api } from './api';
import type {
  User,
  PaginatedResponse,
  ApiResponse,
  CreateUserRequest,
  UpdateUserRequest,
} from '@/types/user';

export const usersApi = {
  // ユーザー一覧を取得（ページネーション付き）
  async getUsers(page = 1, perPage = 10): Promise<PaginatedResponse<User>> {
    const response = await api.get<PaginatedResponse<User>>('/users', {
      params: { page, per_page: perPage },
    });
    return response.data;
  },

  // ユーザーをIDで取得
  async getUserById(id: number): Promise<User> {
    const response = await api.get<ApiResponse<{ user: User }>>(`/users/${id}`);
    return response.data.data.user;
  },

  // ユーザーを作成
  async createUser(data: CreateUserRequest): Promise<User> {
    const response = await api.post<ApiResponse<{ user: User }>>('/users', data);
    return response.data.data.user;
  },

  // ユーザーを更新
  async updateUser(id: number, data: UpdateUserRequest): Promise<User> {
    const response = await api.put<ApiResponse<{ user: User }>>(`/users/${id}`, data);
    return response.data.data.user;
  },

  // ユーザーを削除
  async deleteUser(id: number): Promise<void> {
    await api.delete(`/users/${id}`);
  },
};
