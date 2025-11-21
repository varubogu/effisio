import { api } from './api';
import type {
  User,
  UsersResponse,
  UserResponse,
  CreateUserRequest,
  UpdateUserRequest,
} from '@/types/user';

export const usersApi = {
  // ユーザー一覧を取得
  async getUsers(): Promise<User[]> {
    const response = await api.get<UsersResponse>('/users');
    return response.data.users;
  },

  // ユーザーをIDで取得
  async getUserById(id: number): Promise<User> {
    const response = await api.get<UserResponse>(`/users/${id}`);
    return response.data.user;
  },

  // ユーザーを作成
  async createUser(data: CreateUserRequest): Promise<User> {
    const response = await api.post<UserResponse>('/users', data);
    return response.data.user;
  },

  // ユーザーを更新
  async updateUser(id: number, data: UpdateUserRequest): Promise<User> {
    const response = await api.put<UserResponse>(`/users/${id}`, data);
    return response.data.user;
  },

  // ユーザーを削除
  async deleteUser(id: number): Promise<void> {
    await api.delete(`/users/${id}`);
  },
};
