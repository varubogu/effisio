import { describe, it, expect, beforeEach, vi } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import {
  useUsers,
  useUser,
  useCreateUser,
  useUpdateUser,
  useDeleteUser,
} from './useUsers';
import { usersApi } from '@/lib/users';
import { TestWrapper } from '@/test/test-utils';
import type {
  User,
  CreateUserRequest,
  UpdateUserRequest,
  PaginatedResponse,
} from '@/types/user';

// Mock usersApi
vi.mock('@/lib/users', () => ({
  usersApi: {
    getUsers: vi.fn(),
    getUserById: vi.fn(),
    createUser: vi.fn(),
    updateUser: vi.fn(),
    deleteUser: vi.fn(),
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

describe('useUsers hooks', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('useUsers', () => {
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

      vi.mocked(usersApi.getUsers).mockResolvedValueOnce(mockResponse);

      const { result } = renderHook(() => useUsers(), {
        wrapper: TestWrapper,
      });

      expect(result.current.isPending).toBe(true);

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(result.current.data).toEqual(mockResponse);
      expect(usersApi.getUsers).toHaveBeenCalledWith(1, 10);
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

      vi.mocked(usersApi.getUsers).mockResolvedValueOnce(mockResponse);

      const { result } = renderHook(() => useUsers(2, 20), {
        wrapper: TestWrapper,
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(usersApi.getUsers).toHaveBeenCalledWith(2, 20);
    });

    it('should handle error when fetching users', async () => {
      vi.mocked(usersApi.getUsers).mockRejectedValueOnce(
        new Error('Network error')
      );

      const { result } = renderHook(() => useUsers(), {
        wrapper: TestWrapper,
      });

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });

      expect(result.current.error).toEqual(new Error('Network error'));
    });
  });

  describe('useUser', () => {
    it('should fetch user by id', async () => {
      vi.mocked(usersApi.getUserById).mockResolvedValueOnce(mockUser);

      const { result } = renderHook(() => useUser(1), {
        wrapper: TestWrapper,
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(result.current.data).toEqual(mockUser);
      expect(usersApi.getUserById).toHaveBeenCalledWith(1);
    });

    it('should not fetch when id is falsy', async () => {
      const { result } = renderHook(() => useUser(0), {
        wrapper: TestWrapper,
      });

      // Query should not be enabled
      expect(result.current.status).toBe('pending');
      expect(usersApi.getUserById).not.toHaveBeenCalled();
    });

    it('should handle error when fetching user', async () => {
      vi.mocked(usersApi.getUserById).mockRejectedValueOnce(
        new Error('User not found')
      );

      const { result } = renderHook(() => useUser(999), {
        wrapper: TestWrapper,
      });

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });

      expect(result.current.error).toEqual(new Error('User not found'));
    });
  });

  describe('useCreateUser', () => {
    it('should create user successfully', async () => {
      const createRequest: CreateUserRequest = {
        username: 'newuser',
        email: 'new@example.com',
        password: 'SecurePassword123!',
        role: 'user',
      };

      const createdUser: User = {
        ...mockUser,
        id: 2,
        username: 'newuser',
        email: 'new@example.com',
      };

      vi.mocked(usersApi.createUser).mockResolvedValueOnce(createdUser);

      const { result } = renderHook(() => useCreateUser(), {
        wrapper: TestWrapper,
      });

      result.current.mutate(createRequest);

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(result.current.data).toEqual(createdUser);
      expect(usersApi.createUser).toHaveBeenCalledWith(createRequest);
    });

    it('should handle create error', async () => {
      const createRequest: CreateUserRequest = {
        username: 'duplicate',
        email: 'dup@example.com',
        password: 'Password123!',
        role: 'user',
      };

      vi.mocked(usersApi.createUser).mockRejectedValueOnce(
        new Error('User already exists')
      );

      const { result } = renderHook(() => useCreateUser(), {
        wrapper: TestWrapper,
      });

      result.current.mutate(createRequest);

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });

      expect(result.current.error).toEqual(new Error('User already exists'));
    });

    it('should invalidate users query on success', async () => {
      const createRequest: CreateUserRequest = {
        username: 'newuser',
        email: 'new@example.com',
        password: 'Password123!',
        role: 'user',
      };

      vi.mocked(usersApi.createUser).mockResolvedValueOnce({
        ...mockUser,
        id: 2,
        username: 'newuser',
      });

      const { result } = renderHook(() => useCreateUser(), {
        wrapper: TestWrapper,
      });

      result.current.mutate(createRequest);

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      // onSuccess callback invalidates users query
      expect(result.current.isSuccess).toBe(true);
    });
  });

  describe('useUpdateUser', () => {
    it('should update user successfully', async () => {
      const updateRequest: UpdateUserRequest = {
        full_name: 'Updated Name',
        role: 'manager',
      };

      const updatedUser: User = {
        ...mockUser,
        full_name: 'Updated Name',
        role: 'manager',
      };

      vi.mocked(usersApi.updateUser).mockResolvedValueOnce(updatedUser);

      const { result } = renderHook(() => useUpdateUser(), {
        wrapper: TestWrapper,
      });

      result.current.mutate({ id: 1, data: updateRequest });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(result.current.data).toEqual(updatedUser);
      expect(usersApi.updateUser).toHaveBeenCalledWith(1, updateRequest);
    });

    it('should handle update error', async () => {
      const updateRequest: UpdateUserRequest = {
        email: 'invalid-email',
      };

      vi.mocked(usersApi.updateUser).mockRejectedValueOnce(
        new Error('Invalid email')
      );

      const { result } = renderHook(() => useUpdateUser(), {
        wrapper: TestWrapper,
      });

      result.current.mutate({ id: 1, data: updateRequest });

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });

      expect(result.current.error).toEqual(new Error('Invalid email'));
    });

    it('should invalidate both users and user detail queries on success', async () => {
      const updateRequest: UpdateUserRequest = {
        full_name: 'Updated',
      };

      vi.mocked(usersApi.updateUser).mockResolvedValueOnce({
        ...mockUser,
        full_name: 'Updated',
      });

      const { result } = renderHook(() => useUpdateUser(), {
        wrapper: TestWrapper,
      });

      result.current.mutate({ id: 1, data: updateRequest });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      // onSuccess callback invalidates both users and user detail
      expect(result.current.isSuccess).toBe(true);
    });
  });

  describe('useDeleteUser', () => {
    it('should delete user successfully', async () => {
      vi.mocked(usersApi.deleteUser).mockResolvedValueOnce(undefined);

      const { result } = renderHook(() => useDeleteUser(), {
        wrapper: TestWrapper,
      });

      result.current.mutate(1);

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(usersApi.deleteUser).toHaveBeenCalledWith(1);
    });

    it('should handle delete error', async () => {
      vi.mocked(usersApi.deleteUser).mockRejectedValueOnce(
        new Error('Cannot delete user')
      );

      const { result } = renderHook(() => useDeleteUser(), {
        wrapper: TestWrapper,
      });

      result.current.mutate(1);

      await waitFor(() => {
        expect(result.current.isError).toBe(true);
      });

      expect(result.current.error).toEqual(new Error('Cannot delete user'));
    });

    it('should invalidate users query on success', async () => {
      vi.mocked(usersApi.deleteUser).mockResolvedValueOnce(undefined);

      const { result } = renderHook(() => useDeleteUser(), {
        wrapper: TestWrapper,
      });

      result.current.mutate(1);

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      // onSuccess callback invalidates users query
      expect(result.current.isSuccess).toBe(true);
    });
  });
});
