'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

import { useLogout } from '@/hooks/useAuth';

interface HeaderProps {
  onMenuToggle?: () => void;
}

export function Header({ onMenuToggle }: HeaderProps) {
  const router = useRouter();
  const logoutMutation = useLogout();
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);

  const handleLogout = async () => {
    try {
      await logoutMutation.mutateAsync();
      router.push('/auth/login');
    } catch (error) {
      console.error('ログアウト失敗:', error);
    }
  };

  return (
    <header className="border-b border-gray-200 bg-white shadow-sm">
      <div className="flex items-center justify-between px-4 py-4 sm:px-6 lg:px-8">
        {/* ロゴ・タイトル */}
        <div className="flex items-center">
          <button
            onClick={onMenuToggle}
            className="mr-4 inline-flex items-center justify-center rounded-lg p-2 text-gray-500 hover:bg-gray-100 md:hidden"
          >
            <svg
              className="h-6 w-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4 6h16M4 12h16M4 18h16"
              />
            </svg>
          </button>
          <Link href="/dashboard" className="flex items-center">
            <span className="text-xl font-bold text-gray-900">Effisio</span>
          </Link>
        </div>

        {/* ユーザーメニュー */}
        <div className="relative">
          <button
            onClick={() => setIsDropdownOpen(!isDropdownOpen)}
            className="flex items-center rounded-lg bg-gray-100 px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-200"
          >
            <svg
              className="mr-2 h-5 w-5"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fillRule="evenodd"
                d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
                clipRule="evenodd"
              />
            </svg>
            ユーザー
          </button>

          {/* ドロップダウン */}
          {isDropdownOpen && (
            <div className="absolute right-0 mt-2 w-48 rounded-lg bg-white shadow-lg">
              <button
                onClick={handleLogout}
                disabled={logoutMutation.isPending}
                className="block w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 disabled:opacity-50"
              >
                {logoutMutation.isPending ? 'ログアウト中...' : 'ログアウト'}
              </button>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}
