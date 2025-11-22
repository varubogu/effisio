'use client';

import Link from 'next/link';
import { useState } from 'react';

import { useDeleteUser } from '@/hooks/useUsers';
import { ConfirmDialog } from '@/components/ui/ConfirmDialog';
import type { User } from '@/types/user';

interface UserListProps {
  users: User[];
}

export function UserList({ users }: UserListProps) {
  const [deleteUserId, setDeleteUserId] = useState<number | null>(null);
  const deleteUserMutation = useDeleteUser();

  const handleDeleteClick = (userId: number) => {
    setDeleteUserId(userId);
  };

  const handleDeleteConfirm = async () => {
    if (deleteUserId === null) return;

    try {
      await deleteUserMutation.mutateAsync(deleteUserId);
      setDeleteUserId(null);
    } catch (error) {
      console.error('ユーザー削除失敗:', error);
    }
  };

  const handleDeleteCancel = () => {
    setDeleteUserId(null);
  };
  if (users.length === 0) {
    return (
      <div className="rounded-lg border border-gray-200 bg-white p-8 text-center">
        <p className="text-gray-500">ユーザーが見つかりません</p>
      </div>
    );
  }

  return (
    <div className="overflow-hidden rounded-lg border border-gray-200 bg-white shadow">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
              ID
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
              ユーザー名
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
              氏名
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
              部署
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
              メールアドレス
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
              ロール
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
              ステータス
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
              最終ログイン
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
              アクション
            </th>
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-200 bg-white">
          {users.map((user) => (
            <tr key={user.id} className="hover:bg-gray-50">
              <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-900">
                {user.id}
              </td>
              <td className="whitespace-nowrap px-6 py-4 text-sm font-medium text-gray-900">
                {user.username}
              </td>
              <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-900">
                {user.full_name || '-'}
              </td>
              <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
                {user.department || '-'}
              </td>
              <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
                {user.email}
              </td>
              <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
                <span
                  className={`inline-flex rounded-full px-2 text-xs font-semibold leading-5 ${
                    user.role === 'admin'
                      ? 'bg-purple-100 text-purple-800'
                      : user.role === 'manager'
                        ? 'bg-blue-100 text-blue-800'
                        : user.role === 'user'
                          ? 'bg-green-100 text-green-800'
                          : 'bg-gray-100 text-gray-800'
                  }`}
                >
                  {user.role}
                </span>
              </td>
              <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
                <span
                  className={`inline-flex rounded-full px-2 text-xs font-semibold leading-5 ${
                    user.status === 'active'
                      ? 'bg-green-100 text-green-800'
                      : user.status === 'inactive'
                        ? 'bg-gray-100 text-gray-800'
                        : 'bg-red-100 text-red-800'
                  }`}
                >
                  {user.status === 'active'
                    ? 'アクティブ'
                    : user.status === 'inactive'
                      ? '非アクティブ'
                      : '停止中'}
                </span>
              </td>
              <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
                {user.last_login
                  ? new Date(user.last_login).toLocaleString('ja-JP')
                  : '未ログイン'}
              </td>
              <td className="whitespace-nowrap px-6 py-4 text-sm">
                <div className="flex space-x-2">
                  <Link
                    href={`/users/${user.id}`}
                    className="text-blue-600 hover:text-blue-800"
                  >
                    詳細
                  </Link>
                  <Link
                    href={`/users/${user.id}/edit`}
                    className="text-gray-600 hover:text-gray-800"
                  >
                    編集
                  </Link>
                  <button
                    onClick={() => handleDeleteClick(user.id)}
                    className="text-red-600 hover:text-red-800"
                  >
                    削除
                  </button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {/* 削除確認ダイアログ */}
      <ConfirmDialog
        isOpen={deleteUserId !== null}
        title="ユーザーを削除しますか？"
        message="このユーザーを削除すると、関連するすべてのデータも削除されます。この操作は元に戻せません。"
        confirmText="削除"
        cancelText="キャンセル"
        isDangerous
        isLoading={deleteUserMutation.isPending}
        onConfirm={handleDeleteConfirm}
        onCancel={handleDeleteCancel}
      />
    </div>
  );
}
