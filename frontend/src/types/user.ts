export type UserRole = 'admin' | 'manager' | 'user' | 'viewer';

export interface User {
  id: number;
  username: string;
  email: string;
  role: UserRole;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateUserRequest {
  username: string;
  email: string;
  password: string;
  role?: UserRole;
}

export interface UpdateUserRequest {
  username?: string;
  email?: string;
  role?: UserRole;
  is_active?: boolean;
}

export interface UsersResponse {
  users: User[];
}

export interface UserResponse {
  user: User;
}
