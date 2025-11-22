'use client';

import Link from 'next/link';
import { useState } from 'react';

import { Pagination } from '@/components/ui/Pagination';
import { UserList } from '@/components/users/UserList';
import { useUsers } from '@/hooks/useUsers';

export default function UsersPage() {
  const [currentPage, setCurrentPage] = useState(1);
  const { data, isLoading, error } = useUsers(currentPage, 20);

  if (isLoading) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="mb-4 inline-block h-12 w-12 animate-spin rounded-full border-4 border-solid border-blue-600 border-r-transparent" />
          <p className="text-gray-600">読み込み中...</p>
        </div>
      </main>
    );
  }

  if (error) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-gray-50">
        <div className="text-center">
          <p className="text-xl font-bold text-red-600">エラーが発生しました</p>
          <p className="mt-2 text-gray-600">{error.message}</p>
        </div>
      </main>
    );
  }

  return (
    <main className="bg-gray-50 px-4 py-8">
      <div className="mx-auto max-w-7xl">
        <div className="mb-8 flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">ユーザー管理</h1>
            <p className="mt-2 text-gray-600">
              システムに登録されているユーザーの一覧
            </p>
            {data && data.pagination && (
              <p className="mt-1 text-sm text-gray-500">
                全{data.pagination.total}件中 {data.data.length}件を表示 (
                {data.pagination.page} / {data.pagination.total_pages} ページ)
              </p>
            )}
          </div>
          <Link
            href="/users/new"
            className="rounded-lg bg-blue-600 px-4 py-2 font-semibold text-white hover:bg-blue-700"
          >
            新規ユーザー作成
          </Link>
        </div>

        <UserList users={data?.data || []} />

        {/* ページネーション */}
        {data && data.pagination && data.pagination.total_pages > 1 && (
          <div className="mt-8">
            <Pagination
              currentPage={data.pagination.page}
              totalPages={data.pagination.total_pages}
              onPageChange={setCurrentPage}
              isLoading={isLoading}
            />
          </div>
        )}
      </div>
    </main>
  );
}
