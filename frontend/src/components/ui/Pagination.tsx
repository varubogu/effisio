interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  isLoading?: boolean;
}

export function Pagination({
  currentPage,
  totalPages,
  onPageChange,
  isLoading = false,
}: PaginationProps) {
  const pages: (number | string)[] = [];

  // ページ番号を生成（例: 1, 2, 3, ..., 10）
  const startPage = Math.max(1, currentPage - 2);
  const endPage = Math.min(totalPages, currentPage + 2);

  if (startPage > 1) {
    pages.push(1);
    if (startPage > 2) {
      pages.push('...');
    }
  }

  for (let i = startPage; i <= endPage; i++) {
    pages.push(i);
  }

  if (endPage < totalPages) {
    if (endPage < totalPages - 1) {
      pages.push('...');
    }
    pages.push(totalPages);
  }

  return (
    <nav className="flex items-center justify-center gap-2" aria-label="ページネーション">
      {/* 前ページボタン */}
      <button
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage <= 1 || isLoading}
        className="rounded-lg border border-gray-300 bg-white px-3 py-2 font-semibold text-gray-900 hover:bg-gray-50 disabled:bg-gray-100 disabled:text-gray-400"
      >
        前へ
      </button>

      {/* ページ番号 */}
      {pages.map((page, index) =>
        page === '...' ? (
          <span key={`dots-${index}`} className="px-2 text-gray-500">
            ...
          </span>
        ) : (
          <button
            key={page}
            onClick={() => onPageChange(page as number)}
            disabled={page === currentPage || isLoading}
            className={`min-w-10 rounded-lg px-3 py-2 font-semibold ${
              page === currentPage
                ? 'bg-blue-600 text-white'
                : 'border border-gray-300 bg-white text-gray-900 hover:bg-gray-50'
            } disabled:bg-gray-100 disabled:text-gray-400`}
          >
            {page}
          </button>
        )
      )}

      {/* 次ページボタン */}
      <button
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage >= totalPages || isLoading}
        className="rounded-lg border border-gray-300 bg-white px-3 py-2 font-semibold text-gray-900 hover:bg-gray-50 disabled:bg-gray-100 disabled:text-gray-400"
      >
        次へ
      </button>
    </nav>
  );
}
