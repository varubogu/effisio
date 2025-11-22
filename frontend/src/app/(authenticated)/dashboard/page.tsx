'use client';

import Link from 'next/link';
import { useEffect, useState } from 'react';

interface DashboardStats {
  total_users: number;
  active_users: number;
}

export default function DashboardPage() {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    try {
      // ダミーデータを使用（バックエンド実装待ち）
      const mockStats: DashboardStats = {
        total_users: 42,
        active_users: 38,
      };
      setStats(mockStats);
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'ダッシュボードデータの取得に失敗しました';
      setError(errorMsg);
    } finally {
      setIsLoading(false);
    }
  }, []);

  if (isLoading) {
    return (
      <main className="flex min-h-screen items-center justify-center">
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
          <div className="mb-6 rounded-lg bg-red-50 p-4 text-sm text-red-800">
            {error}
          </div>
        )}

        {/* 統計情報カード */}
        <div className="mb-8 grid gap-6 md:grid-cols-2">
          {/* 総ユーザー数 */}
          <div className="rounded-lg bg-white p-6 shadow">
            <h3 className="text-sm font-medium text-gray-600">総ユーザー数</h3>
            <p className="mt-2 text-3xl font-bold text-gray-900">{stats?.total_users || 0}</p>
            <p className="mt-1 text-xs text-gray-500">全ユーザー</p>
          </div>

          {/* アクティブユーザー数 */}
          <div className="rounded-lg bg-white p-6 shadow">
            <h3 className="text-sm font-medium text-gray-600">アクティブユーザー</h3>
            <p className="mt-2 text-3xl font-bold text-green-600">{stats?.active_users || 0}</p>
            <p className="mt-1 text-xs text-gray-500">
              {stats ? Math.round((stats.active_users / stats.total_users) * 100) : 0}% のユーザーがアクティブ
            </p>
          </div>
        </div>

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

        {/* 最近のアクティビティ（ダミー） */}
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
