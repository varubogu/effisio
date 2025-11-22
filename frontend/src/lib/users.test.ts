import { describe, it, expect, beforeEach, vi } from 'vitest';
import { usersApi } from './users';
import { api } from './api';
import type {
  User,
  CreateUserRequest,
  UpdateUserRequest,
  PaginatedResponse,
} from '@/types/user';

// Mock the api module
vi.mock('./api', () => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

const mockUser: User = {
  id: 1,
  username: 'testuser',
  email: 'test@example.com',
  full_name: 'Test User',
  department: 'Engineering',
  role: 'admin',
  status: 'active',
  last_login: '2024-01-15T10:30:00Z',
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-15T00:00:00Z',
};

describe('usersApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getUsers', () => {
    it('should fetch users with default pagination', async () => {
      const mockResponse: PaginatedResponse<User> = {
        code: 200,
        message: 'OK',
        data: [mockUser],
        pagination: {
          page: 1,
          per_page: 10,
          total: 1,
          total_pages: 1,
        },
      };

      vi.mocked(api.get).mockResolvedValueOnce({
        data: mockResponse,
      });

      const result = await usersApi.getUsers();

      expect(result).toEqual(mockResponse);
      expect(api.get).toHaveBeenCalledWith('/users', {
        params: { page: 1, per_page: 10 },
      });
    });

    it('should fetch users with custom pagination', async () => {
      const mockResponse: PaginatedResponse<User> = {
        code: 200,
        message: 'OK',
        data: [mockUser],
        pagination: {
          page: 2,
          per_page: 20,
          total: 1,
          total_pages: 1,
        },
      };

      vi.mocked(api.get).mockResolvedValueOnce({
        data: mockResponse,
      });

      const result = await usersApi.getUsers(2, 20);

      expect(result).toEqual(mockResponse);
      expect(api.get).toHaveBeenCalledWith('/users', {
        params: { page: 2, per_page: 20 },
      });
    });

    it('should handle empty users list', async () => {
      const mockResponse: PaginatedResponse<User> = {
        code: 200,
        message: 'OK',
        data: [],
        pagination: {
          page: 1,
          per_page: 10,
          total: 0,
          total_pages: 0,
        },
      };

      vi.mocked(api.get).mockResolvedValueOnce({
        data: mockResponse,
      });

      const result = await usersApi.getUsers();

      expect(result.data).toHaveLength(0);
      expect(result.pagination.total).toBe(0);
    });

    it('should throw error on fetch failure', async () => {
      vi.mocked(api.get).mockRejectedValueOnce(
        new Error('Network error')
      );

      await expect(usersApi.getUsers()).rejects.toThrow('Network error');
    });
  });

  describe('getUserById', () => {
    it('should fetch user by id', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: {
          code: 200,
          message: 'OK',
          data: { user: mockUser },
        },
      });

      const result = await usersApi.getUserById(1);

      expect(result).toEqual(mockUser);
      expect(api.get).toHaveBeenCalledWith('/users/1');
    });

    it('should throw error when user not found', async () => {
      vi.mocked(api.get).mockRejectedValueOnce(
        new Error('User not found')
      );

      await expect(usersApi.getUserById(999)).rejects.toThrow('User not found');
    });

    it('should handle different user data', async () => {
      const customUser: User = {
        ...mockUser,
        id: 5,
        username: 'customuser',
        role: 'user',
        status: 'inactive',
      };

      vi.mocked(api.get).mockResolvedValueOnce({
        data: {
          code: 200,
          message: 'OK',
          data: { user: customUser },
        },
      });

      const result = await usersApi.getUserById(5);

      expect(result).toEqual(customUser);
      expect(result.role).toBe('user');
      expect(result.status).toBe('inactive');
    });
  });

  describe('createUser', () => {
    it('should create user successfully', async () => {
      const createRequest: CreateUserRequest = {
        username: 'newuser',
        email: 'new@example.com',
        full_name: 'New User',
        department: 'Sales',
        password: 'SecurePassword123!',
        role: 'user',
      };

      const createdUser: User = {
        ...mockUser,
        id: 2,
        username: 'newuser',
        email: 'new@example.com',
      };

      vi.mocked(api.post).mockResolvedValueOnce({
        data: {
          code: 201,
          message: 'Created',
          data: { user: createdUser },
        },
      });

      const result = await usersApi.createUser(createRequest);

      expect(result).toEqual(createdUser);
      expect(api.post).toHaveBeenCalledWith('/users', createRequest);
    });

    it('should handle duplicate user error', async () => {
      const createRequest: CreateUserRequest = {
        username: 'existing',
        email: 'existing@example.com',
        password: 'Password123!',
        role: 'user',
      };

      vi.mocked(api.post).mockRejectedValueOnce(
        new Error('User already exists')
      );

      await expect(usersApi.createUser(createRequest)).rejects.toThrow(
        'User already exists'
      );
    });

    it('should create user without optional fields', async () => {
      const createRequest: CreateUserRequest = {
        username: 'simpleuser',
        email: 'simple@example.com',
        password: 'Password123!',
        role: 'viewer',
      };

      const createdUser: User = {
        ...mockUser,
        id: 3,
        username: 'simpleuser',
        email: 'simple@example.com',
        role: 'viewer',
      };

      vi.mocked(api.post).mockResolvedValueOnce({
        data: {
          code: 201,
          message: 'Created',
          data: { user: createdUser },
        },
      });

      const result = await usersApi.createUser(createRequest);

      expect(result).toEqual(createdUser);
    });
  });

  describe('updateUser', () => {
    it('should update user successfully', async () => {
      const updateRequest: UpdateUserRequest = {
        full_name: 'Updated Name',
        email: 'updated@example.com',
        role: 'manager',
      };

      const updatedUser: User = {
        ...mockUser,
        full_name: 'Updated Name',
        email: 'updated@example.com',
        role: 'manager',
      };

      vi.mocked(api.put).mockResolvedValueOnce({
        data: {
          code: 200,
          message: 'OK',
          data: { user: updatedUser },
        },
      });

      const result = await usersApi.updateUser(1, updateRequest);

      expect(result).toEqual(updatedUser);
      expect(api.put).toHaveBeenCalledWith('/users/1', updateRequest);
    });

    it('should update partial user fields', async () => {
      const updateRequest: UpdateUserRequest = {
        status: 'suspended',
      };

      const updatedUser: User = {
        ...mockUser,
        status: 'suspended',
      };

      vi.mocked(api.put).mockResolvedValueOnce({
        data: {
          code: 200,
          message: 'OK',
          data: { user: updatedUser },
        },
      });

      const result = await usersApi.updateUser(1, updateRequest);

      expect(result.status).toBe('suspended');
    });

    it('should throw error when user not found', async () => {
      const updateRequest: UpdateUserRequest = {
        full_name: 'New Name',
      };

      vi.mocked(api.put).mockRejectedValueOnce(
        new Error('User not found')
      );

      await expect(usersApi.updateUser(999, updateRequest)).rejects.toThrow(
        'User not found'
      );
    });
  });

  describe('deleteUser', () => {
    it('should delete user successfully', async () => {
      vi.mocked(api.delete).mockResolvedValueOnce({ data: null });

      await expect(usersApi.deleteUser(1)).resolves.toBeUndefined();

      expect(api.delete).toHaveBeenCalledWith('/users/1');
    });

    it('should throw error when user not found', async () => {
      vi.mocked(api.delete).mockRejectedValueOnce(
        new Error('User not found')
      );

      await expect(usersApi.deleteUser(999)).rejects.toThrow('User not found');
    });

    it('should handle delete error gracefully', async () => {
      vi.mocked(api.delete).mockRejectedValueOnce(
        new Error('Server error')
      );

      await expect(usersApi.deleteUser(1)).rejects.toThrow('Server error');
    });
  });
});
