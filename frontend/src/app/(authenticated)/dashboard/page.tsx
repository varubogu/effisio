'use client';

import Link from 'next/link';

import { useDashboardOverview } from '@/hooks/useDashboard';
import { Alert } from '@/components/ui/Alert';

export default function DashboardPage() {
  const { data: overview, isLoading, error } = useDashboardOverview();

  if (isLoading) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="mb-4 inline-block">
            <div className="h-8 w-8 animate-spin rounded-full border-4 border-gray-300 border-t-blue-600" />
          </div>
          <p className="text-gray-600">読み込み中...</p>
        </div>
      </main>
    );
  }

  return (
    <main className="min-h-screen bg-gray-50 px-4 py-8">
      <div className="mx-auto max-w-7xl">
        {/* ヘッダー */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">ダッシュボード</h1>
          <p className="mt-2 text-gray-600">システム概要と統計情報</p>
        </div>

        {/* エラー表示 */}
        {error && (
          <div className="mb-6">
            <Alert
              type="error"
              title="エラー"
              message={error instanceof Error ? error.message : 'ダッシュボードの読み込みに失敗しました'}
            />
          </div>
        )}

        {/* 統計情報カード */}
        <div className="mb-8 grid gap-6 md:grid-cols-2 lg:grid-cols-4">
          {/* 総ユーザー数 */}
          <div className="rounded-lg bg-white p-6 shadow">
            <h3 className="text-sm font-medium text-gray-600">総ユーザー数</h3>
            <p className="mt-2 text-3xl font-bold text-gray-900">
              {overview?.total_users || 0}
            </p>
            <p className="mt-1 text-xs text-gray-500">全ユーザー</p>
          </div>

          {/* アクティブユーザー数 */}
          <div className="rounded-lg bg-white p-6 shadow">
            <h3 className="text-sm font-medium text-gray-600">アクティブユーザー</h3>
            <p className="mt-2 text-3xl font-bold text-green-600">
              {overview?.active_users || 0}
            </p>
            <p className="mt-1 text-xs text-gray-500">
              {overview
                ? Math.round((overview.active_users / overview.total_users) * 100)
                : 0}
              % がアクティブ
            </p>
          </div>

          {/* 非アクティブユーザー数 */}
          <div className="rounded-lg bg-white p-6 shadow">
            <h3 className="text-sm font-medium text-gray-600">非アクティブユーザー</h3>
            <p className="mt-2 text-3xl font-bold text-gray-600">
              {overview?.inactive_users || 0}
            </p>
            <p className="mt-1 text-xs text-gray-500">未ログイン状態</p>
          </div>

          {/* 停止中ユーザー数 */}
          <div className="rounded-lg bg-white p-6 shadow">
            <h3 className="text-sm font-medium text-gray-600">停止中ユーザー</h3>
            <p className="mt-2 text-3xl font-bold text-red-600">
              {overview?.suspended_users || 0}
            </p>
            <p className="mt-1 text-xs text-gray-500">アカウント停止</p>
          </div>
        </div>

        {/* ロール別ユーザー数 */}
        {overview && overview.users_by_role && (
          <div className="mb-8 rounded-lg bg-white p-6 shadow">
            <h2 className="text-lg font-semibold text-gray-900">ロール別ユーザー数</h2>
            <div className="mt-4 grid gap-4 md:grid-cols-4">
              {Object.entries(overview.users_by_role).map(([role, count]) => (
                <div key={role} className="rounded-lg bg-gray-50 p-4">
                  <p className="text-sm font-medium text-gray-600">{role}</p>
                  <p className="mt-1 text-2xl font-bold text-gray-900">{count}</p>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* 部門別ユーザー数 */}
        {overview && overview.users_by_department && overview.users_by_department.length > 0 && (
          <div className="mb-8 rounded-lg bg-white p-6 shadow">
            <h2 className="text-lg font-semibold text-gray-900">部門別ユーザー数</h2>
            <div className="mt-4 space-y-2">
              {overview.users_by_department.map((dept) => (
                <div key={dept.department} className="flex items-center justify-between border-b border-gray-200 pb-2">
                  <span className="text-sm text-gray-700">{dept.department}</span>
                  <span className="font-semibold text-gray-900">{dept.count}人</span>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* クイックリンク */}
        <div className="mb-8">
          <h2 className="mb-4 text-lg font-semibold text-gray-900">クイックアクション</h2>
          <div className="grid gap-4 md:grid-cols-3">
            <Link
              href="/users"
              className="flex items-center rounded-lg border border-gray-200 bg-white p-4 hover:border-blue-500 hover:bg-blue-50"
            >
              <div className="flex-1">
                <h3 className="font-semibold text-gray-900">ユーザー管理</h3>
                <p className="text-sm text-gray-600">ユーザーの追加・編集・削除</p>
              </div>
              <span className="text-gray-400">&gt;</span>
            </Link>

            <div className="flex items-center rounded-lg border border-gray-200 bg-gray-50 p-4 opacity-50">
              <div className="flex-1">
                <h3 className="font-semibold text-gray-900">レポート</h3>
                <p className="text-sm text-gray-600">システムレポート（準備中）</p>
              </div>
              <span className="text-gray-400">&gt;</span>
            </div>

            <div className="flex items-center rounded-lg border border-gray-200 bg-gray-50 p-4 opacity-50">
              <div className="flex-1">
                <h3 className="font-semibold text-gray-900">設定</h3>
                <p className="text-sm text-gray-600">システム設定（準備中）</p>
              </div>
              <span className="text-gray-400">&gt;</span>
            </div>
          </div>
        </div>

        {/* 最近のアクティビティ（監査ログ） */}
        <div>
          <h2 className="mb-4 text-lg font-semibold text-gray-900">最近のアクティビティ</h2>
          <div className="rounded-lg bg-white shadow">
            <div className="border-b border-gray-200 px-6 py-4">
              <p className="text-sm text-gray-600">
                監査ログ機能は Phase 4 で実装予定です。
              </p>
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}
