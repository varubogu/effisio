import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { usersApi } from '@/lib/users';
import type { CreateUserRequest, UpdateUserRequest } from '@/types/user';

const USERS_QUERY_KEY = ['users'];

// ユーザー一覧を取得（ページネーション付き）
export function useUsers(page = 1, perPage = 10) {
  return useQuery({
    queryKey: [...USERS_QUERY_KEY, page, perPage],
    queryFn: () => usersApi.getUsers(page, perPage),
  });
}

// ユーザーをIDで取得
export function useUser(id: number) {
  return useQuery({
    queryKey: [...USERS_QUERY_KEY, id],
    queryFn: () => usersApi.getUserById(id),
    enabled: !!id,
  });
}

// ユーザーを作成
export function useCreateUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateUserRequest) => usersApi.createUser(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY });
    },
  });
}

// ユーザーを更新
export function useUpdateUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateUserRequest }) =>
      usersApi.updateUser(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY });
      queryClient.invalidateQueries({ queryKey: [...USERS_QUERY_KEY, variables.id] });
    },
  });
}

// ユーザーを削除
export function useDeleteUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: number) => usersApi.deleteUser(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY });
    },
  });
}
