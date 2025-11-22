import type { UserRole } from '@/types/user';

interface RoleGuardProps {
  allowedRoles: UserRole[];
  userRole?: UserRole;
  children: React.ReactNode;
  fallback?: React.ReactNode;
}

export function RoleGuard({
  allowedRoles,
  userRole,
  children,
  fallback = null,
}: RoleGuardProps) {
  // ロールが指定されていない場合は表示しない
  if (!userRole) {
    return fallback;
  }

  // ユーザーのロールが許可されたロールに含まれているかチェック
  if (!allowedRoles.includes(userRole)) {
    return fallback;
  }

  return <>{children}</>;
}
