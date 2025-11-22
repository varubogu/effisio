'use client';

import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

import { tokenStorage } from '@/lib/auth';

export default function Home() {
  const router = useRouter();

  useEffect(() => {
    // ログイン状況をチェック
    const isAuthenticated = !!tokenStorage.getAccessToken();

    if (isAuthenticated) {
      // ログイン済みの場合はダッシュボードへ
      router.push('/dashboard');
    } else {
      // 未ログインの場合はログインページへ
      router.push('/auth/login');
    }
  }, [router]);

  return (
    <main className="flex min-h-screen items-center justify-center bg-gray-50">
      <div className="text-center">
        <div className="mb-4 inline-block">
          <div className="h-8 w-8 animate-spin rounded-full border-4 border-gray-300 border-t-blue-600" />
        </div>
        <p className="text-gray-600">リダイレクト中...</p>
      </div>
    </main>
  );
}
