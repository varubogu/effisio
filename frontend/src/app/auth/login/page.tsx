'use client';

import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { useForm } from 'react-hook-form';

import { useLogin } from '@/hooks/useAuth';
import type { LoginRequest } from '@/types/auth';

export default function LoginPage() {
  const router = useRouter();
  const loginMutation = useLogin();
  const [errorMessage, setErrorMessage] = useState<string>('');
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginRequest>({
    defaultValues: {
      username: '',
      password: '',
    },
  });

  const onSubmit = async (data: LoginRequest) => {
    setErrorMessage('');
    try {
      await loginMutation.mutateAsync(data);
      router.push('/dashboard');
    } catch (error) {
      const errorMsg = error instanceof Error ? error.message : 'ログインに失敗しました';
      setErrorMessage(errorMsg);
    }
  };

  return (
    <main className="flex min-h-screen items-center justify-center bg-gray-50 px-4 py-12">
      <div className="w-full max-w-md rounded-lg bg-white p-8 shadow-md">
        <h1 className="mb-8 text-center text-3xl font-bold text-gray-900">Effisio</h1>
        <p className="mb-6 text-center text-sm text-gray-600">社内管理システムへのログイン</p>

        {errorMessage && (
          <div className="mb-4 rounded-lg bg-red-50 p-4 text-sm text-red-800">
            {errorMessage}
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          {/* ユーザー名入力 */}
          <div>
            <label htmlFor="username" className="block text-sm font-medium text-gray-700">
              ユーザー名
            </label>
            <input
              {...register('username', {
                required: 'ユーザー名は必須です',
              })}
              id="username"
              type="text"
              placeholder="ユーザー名を入力"
              className={`mt-1 w-full rounded-lg border px-3 py-2 text-gray-900 placeholder-gray-400 focus:border-blue-500 focus:outline-none ${
                errors.username ? 'border-red-500' : 'border-gray-300'
              }`}
              disabled={isSubmitting}
            />
            {errors.username && (
              <p className="mt-1 text-sm text-red-600">{errors.username.message}</p>
            )}
          </div>

          {/* パスワード入力 */}
          <div>
            <label htmlFor="password" className="block text-sm font-medium text-gray-700">
              パスワード
            </label>
            <input
              {...register('password', {
                required: 'パスワードは必須です',
              })}
              id="password"
              type="password"
              placeholder="パスワードを入力"
              className={`mt-1 w-full rounded-lg border px-3 py-2 text-gray-900 placeholder-gray-400 focus:border-blue-500 focus:outline-none ${
                errors.password ? 'border-red-500' : 'border-gray-300'
              }`}
              disabled={isSubmitting}
            />
            {errors.password && (
              <p className="mt-1 text-sm text-red-600">{errors.password.message}</p>
            )}
          </div>

          {/* ログインボタン */}
          <button
            type="submit"
            disabled={isSubmitting || loginMutation.isPending}
            className="mt-6 w-full rounded-lg bg-blue-600 px-4 py-2 font-semibold text-white hover:bg-blue-700 focus:outline-none disabled:bg-gray-400"
          >
            {isSubmitting || loginMutation.isPending ? 'ログイン中...' : 'ログイン'}
          </button>
        </form>

        <div className="mt-6 text-center text-xs text-gray-500">
          <p>開発環境テスト用ユーザー:</p>
          <p>ユーザー名: admin, パスワード: password123</p>
        </div>
      </div>
    </main>
  );
}
