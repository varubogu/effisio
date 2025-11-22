'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useCreateUser } from '@/hooks/useUsers';
import { Alert } from '@/components/ui/Alert';
import type { CreateUserRequest, UserRole } from '@/types/user';

type FormData = CreateUserRequest;

const roles: { value: UserRole; label: string }[] = [
  { value: 'admin', label: '管理者' },
  { value: 'manager', label: 'マネージャー' },
  { value: 'user', label: 'ユーザー' },
  { value: 'viewer', label: 'ビューア' },
];

export default function NewUserPage() {
  const router = useRouter();
  const [successMessage, setSuccessMessage] = useState('');
  const createUserMutation = useCreateUser();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({
    defaultValues: {
      username: '',
      email: '',
      full_name: '',
      department: '',
      password: '',
      role: 'user',
    },
  });

  const onSubmit = async (data: FormData) => {
    try {
      await createUserMutation.mutateAsync(data);
      setSuccessMessage('ユーザーを作成しました。');
      setTimeout(() => {
        router.push('/users');
      }, 1500);
    } catch (error) {
      console.error('ユーザー作成失敗:', error);
    }
  };

  return (
    <main className="bg-gray-50 px-4 py-8">
      <div className="mx-auto max-w-2xl">
        {/* ヘッダー */}
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900">新規ユーザー作成</h1>
          <p className="mt-2 text-gray-600">新しいユーザーアカウントを作成します</p>
        </div>

        {/* 成功メッセージ */}
        {successMessage && (
          <div className="mb-6">
            <Alert type="success" message={successMessage} />
          </div>
        )}

        {/* エラーメッセージ */}
        {createUserMutation.error && (
          <div className="mb-6">
            <Alert
              type="error"
              title="エラー"
              message={
                createUserMutation.error instanceof Error
                  ? createUserMutation.error.message
                  : 'ユーザーの作成に失敗しました'
              }
            />
          </div>
        )}

        {/* フォーム */}
        <form onSubmit={handleSubmit(onSubmit)} className="rounded-lg bg-white shadow">
          <div className="space-y-6 px-6 py-4">
            {/* ユーザー名 */}
            <div>
              <label htmlFor="username" className="block text-sm font-medium text-gray-700">
                ユーザー名 <span className="text-red-600">*</span>
              </label>
              <input
                {...register('username', {
                  required: 'ユーザー名は必須です',
                  minLength: { value: 3, message: '3文字以上で入力してください' },
                })}
                id="username"
                type="text"
                placeholder="例: john_doe"
                className={`mt-1 w-full rounded-lg border px-3 py-2 text-gray-900 focus:border-blue-500 focus:outline-none ${
                  errors.username ? 'border-red-500' : 'border-gray-300'
                }`}
                disabled={isSubmitting || createUserMutation.isPending}
              />
              {errors.username && (
                <p className="mt-1 text-sm text-red-600">{errors.username.message}</p>
              )}
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
                placeholder="例: john@example.com"
                className={`mt-1 w-full rounded-lg border px-3 py-2 text-gray-900 focus:border-blue-500 focus:outline-none ${
                  errors.email ? 'border-red-500' : 'border-gray-300'
                }`}
                disabled={isSubmitting || createUserMutation.isPending}
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
                disabled={isSubmitting || createUserMutation.isPending}
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
                disabled={isSubmitting || createUserMutation.isPending}
              />
            </div>

            {/* パスワード */}
            <div>
              <label htmlFor="password" className="block text-sm font-medium text-gray-700">
                パスワード <span className="text-red-600">*</span>
              </label>
              <input
                {...register('password', {
                  required: 'パスワードは必須です',
                  minLength: {
                    value: 8,
                    message: '8文字以上で入力してください',
                  },
                  pattern: {
                    value: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]/,
                    message: '大文字、小文字、数字、記号を含める必要があります',
                  },
                })}
                id="password"
                type="password"
                placeholder="8文字以上、大小文字・数字・記号を含む"
                className={`mt-1 w-full rounded-lg border px-3 py-2 text-gray-900 focus:border-blue-500 focus:outline-none ${
                  errors.password ? 'border-red-500' : 'border-gray-300'
                }`}
                disabled={isSubmitting || createUserMutation.isPending}
              />
              {errors.password && (
                <p className="mt-1 text-sm text-red-600">{errors.password.message}</p>
              )}
              <p className="mt-1 text-xs text-gray-500">
                パスワードは以下の条件を満たす必要があります：
                <ul className="mt-1 list-inside list-disc">
                  <li>8文字以上</li>
                  <li>大文字を含む</li>
                  <li>小文字を含む</li>
                  <li>数字を含む</li>
                  <li>記号を含む</li>
                </ul>
              </p>
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
                disabled={isSubmitting || createUserMutation.isPending}
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
          </div>

          {/* ボタン */}
          <div className="border-t border-gray-200 bg-gray-50 px-6 py-4">
            <div className="flex gap-3">
              <button
                type="submit"
                disabled={isSubmitting || createUserMutation.isPending}
                className="flex-1 rounded-lg bg-blue-600 px-4 py-2 font-semibold text-white hover:bg-blue-700 disabled:bg-gray-400"
              >
                {isSubmitting || createUserMutation.isPending ? '作成中...' : 'ユーザーを作成'}
              </button>
              <Link
                href="/users"
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
