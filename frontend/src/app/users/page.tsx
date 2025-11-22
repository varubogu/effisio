'use client';

import { useUsers } from '@/hooks/useUsers';
import { UserList } from '@/components/users/UserList';

export default function UsersPage() {
  const { data, isLoading, error } = useUsers();

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <div className="mb-4 inline-block h-12 w-12 animate-spin rounded-full border-4 border-solid border-primary-600 border-r-transparent"></div>
          <p className="text-gray-600">読み込み中...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center text-danger-600">
          <p className="text-xl font-bold">エラーが発生しました</p>
          <p className="mt-2">{error.message}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold">ユーザー管理</h1>
        <p className="mt-2 text-gray-600">
          システムに登録されているユーザーの一覧
        </p>
        {data && data.pagination && (
          <p className="mt-1 text-sm text-gray-500">
            全{data.pagination.total}件中 {data.data.length}件を表示
          </p>
        )}
      </div>

      <UserList users={data?.data || []} />
    </div>
  );
}
