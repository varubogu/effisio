'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useUser, useUpdateUser } from '@/hooks/useUsers';
import { Alert } from '@/components/ui/Alert';
import type { UpdateUserRequest, UserRole, UserStatus } from '@/types/user';

interface EditUserPageProps {
  params: {
    id: string;
  };
}

type FormData = UpdateUserRequest;

const roles: { value: UserRole; label: string }[] = [
  { value: 'admin', label: '管理者' },
  { value: 'manager', label: 'マネージャー' },
  { value: 'user', label: 'ユーザー' },
  { value: 'viewer', label: 'ビューア' },
];

const statuses: { value: UserStatus; label: string }[] = [
  { value: 'active', label: 'アクティブ' },
  { value: 'inactive', label: '非アクティブ' },
  { value: 'suspended', label: '停止中' },
];

export default function EditUserPage({ params }: EditUserPageProps) {
  const router = useRouter();
  const userId = parseInt(params.id, 10);
  const [successMessage, setSuccessMessage] = useState('');

  const { data: user, isLoading, error } = useUser(userId);
  const updateUserMutation = useUpdateUser();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({
    values:
      user && !isLoading
        ? {
            email: user.email,
            full_name: user.full_name || '',
            department: user.department || '',
            role: user.role,
            status: user.status,
          }
        : {
            email: '',
            full_name: '',
            department: '',
            role: 'user',
            status: 'active',
          },
  });

  const onSubmit = async (data: FormData) => {
    try {
      await updateUserMutation.mutateAsync({ id: userId, data });
      setSuccessMessage('ユーザーを更新しました。');
      setTimeout(() => {
        router.push(`/users/${userId}`);
      }, 1500);
    } catch (err) {
      console.error('ユーザー更新失敗:', err);
    }
  };

  if (isLoading) {
    return (
      <main className="bg-gray-50 px-4 py-8">
        <div className="mx-auto max-w-2xl">
          <div className="flex min-h-96 items-center justify-center">
            <div className="text-center">
              <div className="mb-4 inline-block">
                <div className="h-8 w-8 animate-spin rounded-full border-4 border-gray-300 border-t-blue-600" />
              </div>
              <p className="text-gray-600">読み込み中...</p>
            </div>
          </div>
        </div>
      </main>
    );
  }

  if (error || !user) {
    return (
      <main className="bg-gray-50 px-4 py-8">
        <div className="mx-auto max-w-2xl">
          <div className="mb-6">
            <h1 className="text-3xl font-bold text-gray-900">ユーザー編集</h1>
          </div>

          <Alert
            type="error"
            title="エラー"
            message={error?.message || 'ユーザーが見つかりません'}
          />

          <div className="mt-6">
            <Link
              href="/users"
              className="rounded-lg bg-gray-600 px-4 py-2 font-semibold text-white hover:bg-gray-700"
            >
              戻る
            </Link>
          </div>
        </div>
      </main>
    );
  }

  return (
    <main className="bg-gray-50 px-4 py-8">
      <div className="mx-auto max-w-2xl">
        {/* ヘッダー */}
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900">ユーザー編集</h1>
          <p className="mt-2 text-gray-600">ユーザー情報を編集します: {user.username}</p>
        </div>

        {/* 成功メッセージ */}
        {successMessage && (
          <div className="mb-6">
            <Alert type="success" message={successMessage} />
          </div>
        )}

        {/* エラーメッセージ */}
        {updateUserMutation.error && (
          <div className="mb-6">
            <Alert
              type="error"
              title="エラー"
              message={
                updateUserMutation.error instanceof Error
                  ? updateUserMutation.error.message
                  : 'ユーザーの更新に失敗しました'
              }
            />
          </div>
        )}

        {/* フォーム */}
        <form onSubmit={handleSubmit(onSubmit)} className="rounded-lg bg-white shadow">
          <div className="space-y-6 px-6 py-4">
            {/* ユーザー名（読み取り専用） */}
            <div>
              <label className="block text-sm font-medium text-gray-700">ユーザー名</label>
              <div className="mt-1 rounded-lg border border-gray-300 bg-gray-50 px-3 py-2 text-gray-900">
                {user.username}
              </div>
              <p className="mt-1 text-xs text-gray-500">ユーザー名は変更できません</p>
            </div>

            {/* メールアドレス */}
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700">
                メールアドレス <span className="text-red-600">*</span>
              </label>
              <input
                {...register('email', {
                  required: 'メールアドレスは必須です',
                  pattern: {
                    value: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
                    message: '有効なメールアドレスを入力してください',
                  },
                })}
                id="email"
                type="email"
                className={`mt-1 w-full rounded-lg border px-3 py-2 text-gray-900 focus:border-blue-500 focus:outline-none ${
                  errors.email ? 'border-red-500' : 'border-gray-300'
                }`}
                disabled={isSubmitting || updateUserMutation.isPending}
              />
              {errors.email && (
                <p className="mt-1 text-sm text-red-600">{errors.email.message}</p>
              )}
            </div>

            {/* 氏名 */}
            <div>
              <label htmlFor="full_name" className="block text-sm font-medium text-gray-700">
                氏名
              </label>
              <input
                {...register('full_name')}
                id="full_name"
                type="text"
                placeholder="例: 田中太郎"
                className="mt-1 w-full rounded-lg border border-gray-300 px-3 py-2 text-gray-900 focus:border-blue-500 focus:outline-none"
                disabled={isSubmitting || updateUserMutation.isPending}
              />
            </div>

            {/* 部署 */}
            <div>
              <label htmlFor="department" className="block text-sm font-medium text-gray-700">
                部署
              </label>
              <input
                {...register('department')}
                id="department"
                type="text"
                placeholder="例: 営業部"
                className="mt-1 w-full rounded-lg border border-gray-300 px-3 py-2 text-gray-900 focus:border-blue-500 focus:outline-none"
                disabled={isSubmitting || updateUserMutation.isPending}
              />
            </div>

            {/* ロール */}
            <div>
              <label htmlFor="role" className="block text-sm font-medium text-gray-700">
                ロール <span className="text-red-600">*</span>
              </label>
              <select
                {...register('role', { required: 'ロールは必須です' })}
                id="role"
                className={`mt-1 w-full rounded-lg border px-3 py-2 text-gray-900 focus:border-blue-500 focus:outline-none ${
                  errors.role ? 'border-red-500' : 'border-gray-300'
                }`}
                disabled={isSubmitting || updateUserMutation.isPending}
              >
                {roles.map((role) => (
                  <option key={role.value} value={role.value}>
                    {role.label}
                  </option>
                ))}
              </select>
              {errors.role && (
                <p className="mt-1 text-sm text-red-600">{errors.role.message}</p>
              )}
            </div>

            {/* ステータス */}
            <div>
              <label htmlFor="status" className="block text-sm font-medium text-gray-700">
                ステータス <span className="text-red-600">*</span>
              </label>
              <select
                {...register('status', { required: 'ステータスは必須です' })}
                id="status"
                className={`mt-1 w-full rounded-lg border px-3 py-2 text-gray-900 focus:border-blue-500 focus:outline-none ${
                  errors.status ? 'border-red-500' : 'border-gray-300'
                }`}
                disabled={isSubmitting || updateUserMutation.isPending}
              >
                {statuses.map((status) => (
                  <option key={status.value} value={status.value}>
                    {status.label}
                  </option>
                ))}
              </select>
              {errors.status && (
                <p className="mt-1 text-sm text-red-600">{errors.status.message}</p>
              )}
            </div>
          </div>

          {/* ボタン */}
          <div className="border-t border-gray-200 bg-gray-50 px-6 py-4">
            <div className="flex gap-3">
              <button
                type="submit"
                disabled={isSubmitting || updateUserMutation.isPending}
                className="flex-1 rounded-lg bg-blue-600 px-4 py-2 font-semibold text-white hover:bg-blue-700 disabled:bg-gray-400"
              >
                {isSubmitting || updateUserMutation.isPending ? '更新中...' : '変更を保存'}
              </button>
              <Link
                href={`/users/${userId}`}
                className="flex-1 rounded-lg bg-gray-600 px-4 py-2 text-center font-semibold text-white hover:bg-gray-700"
              >
                キャンセル
              </Link>
            </div>
          </div>
        </form>
      </div>
    </main>
  );
}
