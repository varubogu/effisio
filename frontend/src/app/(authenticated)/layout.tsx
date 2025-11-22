'use client';

import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

import { DashboardLayout } from '@/components/layout/DashboardLayout';
import { tokenStorage } from '@/lib/auth';

export default function AuthenticatedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const [isAuthorized, setIsAuthorized] = useState<boolean | null>(null);

  useEffect(() => {
    const isAuthenticated = !!tokenStorage.getAccessToken();

    if (!isAuthenticated) {
      router.push('/auth/login');
      return;
    }

    setIsAuthorized(true);
  }, [router]);

  // 認証確認中の表示
  if (isAuthorized === null) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="mb-4 inline-block">
            <div className="h-8 w-8 animate-spin rounded-full border-4 border-gray-300 border-t-blue-600" />
          </div>
          <p className="text-gray-600">読み込み中...</p>
        </div>
      </div>
    );
  }

  if (!isAuthorized) {
    return null;
  }

  return <DashboardLayout>{children}</DashboardLayout>;
}
