'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useUser } from '@/hooks/useUsers';
import { Alert } from '@/components/ui/Alert';

interface UserDetailPageProps {
  params: {
    id: string;
  };
}

export default function UserDetailPage({ params }: UserDetailPageProps) {
  const router = useRouter();
  const userId = parseInt(params.id, 10);
  const { data: user, isLoading, error } = useUser(userId);

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
          <div className="mb-6 flex items-center justify-between">
            <h1 className="text-3xl font-bold text-gray-900">ユーザー詳細</h1>
            <Link
              href="/users"
              className="rounded-lg bg-gray-600 px-4 py-2 font-semibold text-white hover:bg-gray-700"
            >
              戻る
            </Link>
          </div>

          <Alert
            type="error"
            title="エラー"
            message={error?.message || 'ユーザーが見つかりません'}
          />
        </div>
      </main>
    );
  }

  const roleLabels: Record<string, string> = {
    admin: '管理者',
    manager: 'マネージャー',
    user: 'ユーザー',
    viewer: 'ビューア',
  };

  const statusLabels: Record<string, string> = {
    active: 'アクティブ',
    inactive: '非アクティブ',
    suspended: '停止中',
  };

  return (
    <main className="bg-gray-50 px-4 py-8">
      <div className="mx-auto max-w-2xl">
        {/* ヘッダー */}
        <div className="mb-6 flex items-center justify-between">
          <h1 className="text-3xl font-bold text-gray-900">ユーザー詳細</h1>
          <div className="space-x-2">
            <Link
              href={`/users/${user.id}/edit`}
              className="rounded-lg bg-blue-600 px-4 py-2 font-semibold text-white hover:bg-blue-700"
            >
              編集
            </Link>
            <Link
              href="/users"
              className="rounded-lg bg-gray-600 px-4 py-2 font-semibold text-white hover:bg-gray-700"
            >
              戻る
            </Link>
          </div>
        </div>

        {/* ユーザー情報カード */}
        <div className="rounded-lg bg-white shadow">
          <div className="border-b border-gray-200 px-6 py-4">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-2xl font-bold text-gray-900">{user.full_name || user.username}</h2>
                <p className="mt-1 text-gray-600">ユーザーID: {user.id}</p>
              </div>
              <span
                className={`inline-flex rounded-full px-3 py-1 text-sm font-semibold leading-5 ${
                  user.status === 'active'
                    ? 'bg-green-100 text-green-800'
                    : user.status === 'inactive'
                      ? 'bg-gray-100 text-gray-800'
                      : 'bg-red-100 text-red-800'
                }`}
              >
                {statusLabels[user.status] || user.status}
              </span>
            </div>
          </div>

          {/* 基本情報 */}
          <div className="px-6 py-4">
            <h3 className="text-lg font-semibold text-gray-900">基本情報</h3>
            <dl className="mt-4 space-y-4">
              <div className="flex items-center justify-between border-b border-gray-200 pb-4">
                <dt className="font-medium text-gray-700">ユーザー名</dt>
                <dd className="text-gray-900">{user.username}</dd>
              </div>
              <div className="flex items-center justify-between border-b border-gray-200 pb-4">
                <dt className="font-medium text-gray-700">メールアドレス</dt>
                <dd className="text-gray-900">{user.email}</dd>
              </div>
              <div className="flex items-center justify-between border-b border-gray-200 pb-4">
                <dt className="font-medium text-gray-700">氏名</dt>
                <dd className="text-gray-900">{user.full_name || '-'}</dd>
              </div>
              <div className="flex items-center justify-between pb-4">
                <dt className="font-medium text-gray-700">部署</dt>
                <dd className="text-gray-900">{user.department || '-'}</dd>
              </div>
            </dl>
          </div>

          {/* ロール・権限 */}
          <div className="border-t border-gray-200 px-6 py-4">
            <h3 className="text-lg font-semibold text-gray-900">ロール・権限</h3>
            <dl className="mt-4 space-y-4">
              <div className="flex items-center justify-between border-b border-gray-200 pb-4">
                <dt className="font-medium text-gray-700">ロール</dt>
                <dd>
                  <span
                    className={`inline-flex rounded-full px-2 py-1 text-xs font-semibold ${
                      user.role === 'admin'
                        ? 'bg-purple-100 text-purple-800'
                        : user.role === 'manager'
                          ? 'bg-blue-100 text-blue-800'
                          : user.role === 'user'
                            ? 'bg-green-100 text-green-800'
                            : 'bg-gray-100 text-gray-800'
                    }`}
                  >
                    {roleLabels[user.role] || user.role}
                  </span>
                </dd>
              </div>
              <div className="flex items-center justify-between pb-4">
                <dt className="font-medium text-gray-700">ステータス</dt>
                <dd>
                  <span
                    className={`inline-flex rounded-full px-2 py-1 text-xs font-semibold ${
                      user.status === 'active'
                        ? 'bg-green-100 text-green-800'
                        : user.status === 'inactive'
                          ? 'bg-gray-100 text-gray-800'
                          : 'bg-red-100 text-red-800'
                    }`}
                  >
                    {statusLabels[user.status] || user.status}
                  </span>
                </dd>
              </div>
            </dl>
          </div>

          {/* ログイン履歴 */}
          <div className="border-t border-gray-200 px-6 py-4">
            <h3 className="text-lg font-semibold text-gray-900">ログイン履歴</h3>
            <dl className="mt-4 space-y-4">
              <div className="flex items-center justify-between border-b border-gray-200 pb-4">
                <dt className="font-medium text-gray-700">最終ログイン</dt>
                <dd className="text-gray-900">
                  {user.last_login
                    ? new Date(user.last_login).toLocaleString('ja-JP')
                    : '未ログイン'}
                </dd>
              </div>
              <div className="flex items-center justify-between pb-4">
                <dt className="font-medium text-gray-700">作成日時</dt>
                <dd className="text-gray-900">
                  {new Date(user.created_at).toLocaleString('ja-JP')}
                </dd>
              </div>
            </dl>
          </div>
        </div>

        {/* アクション */}
        <div className="mt-6 space-y-4">
          <div className="flex gap-2">
            <Link
              href={`/users/${user.id}/edit`}
              className="flex-1 rounded-lg bg-blue-600 px-4 py-3 text-center font-semibold text-white hover:bg-blue-700"
            >
              情報を編集
            </Link>
            <button
              onClick={() => router.back()}
              className="flex-1 rounded-lg bg-gray-600 px-4 py-3 font-semibold text-white hover:bg-gray-700"
            >
              戻る
            </button>
          </div>
        </div>
      </div>
    </main>
  );
}
