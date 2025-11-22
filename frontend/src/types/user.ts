export type UserRole = 'admin' | 'manager' | 'user' | 'viewer';
export type UserStatus = 'active' | 'inactive' | 'suspended';

export interface User {
  id: number;
  username: string;
  email: string;
  full_name: string;
  department: string;
  role: UserRole;
  status: UserStatus;
  last_login: string | null;
  created_at: string;
  updated_at: string;
}

export interface CreateUserRequest {
  username: string;
  email: string;
  full_name?: string;
  department?: string;
  password: string;
  role: UserRole;
}

export interface UpdateUserRequest {
  email?: string;
  full_name?: string;
  department?: string;
  role?: UserRole;
  status?: UserStatus;
}

export interface PaginationInfo {
  page: number;
  per_page: number;
  total: number;
  total_pages: number;
}

export interface PaginatedResponse<T> {
  code: number;
  message: string;
  data: T[];
  pagination: PaginationInfo;
}

export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

export interface UsersResponse {
  users: User[];
}

export interface UserResponse {
  user: User;
}
