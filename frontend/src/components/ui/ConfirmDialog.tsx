import { useEffect, useRef } from 'react';

interface ConfirmDialogProps {
  isOpen: boolean;
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  isDangerous?: boolean;
  isLoading?: boolean;
  onConfirm: () => void;
  onCancel: () => void;
}

export function ConfirmDialog({
  isOpen,
  title,
  message,
  confirmText = '確認',
  cancelText = 'キャンセル',
  isDangerous = false,
  isLoading = false,
  onConfirm,
  onCancel,
}: ConfirmDialogProps) {
  const dialogRef = useRef<HTMLDialogElement>(null);

  useEffect(() => {
    if (isOpen) {
      dialogRef.current?.showModal();
    } else {
      dialogRef.current?.close();
    }
  }, [isOpen]);

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Escape') {
      onCancel();
    }
  };

  return (
    <dialog
      ref={dialogRef}
      className="rounded-lg shadow-lg backdrop:bg-black backdrop:bg-opacity-50"
      onKeyDown={handleKeyDown}
    >
      <div className="w-96 max-w-sm rounded-lg bg-white p-6">
        {/* タイトル */}
        <h2 className="text-lg font-bold text-gray-900">{title}</h2>

        {/* メッセージ */}
        <p className="mt-3 text-gray-600">{message}</p>

        {/* ボタン */}
        <div className="mt-6 flex gap-3">
          <button
            onClick={onCancel}
            disabled={isLoading}
            className="flex-1 rounded-lg border border-gray-300 bg-white px-4 py-2 font-semibold text-gray-900 hover:bg-gray-50 disabled:bg-gray-100"
          >
            {cancelText}
          </button>
          <button
            onClick={onConfirm}
            disabled={isLoading}
            className={`flex-1 rounded-lg px-4 py-2 font-semibold text-white disabled:bg-gray-400 ${
              isDangerous
                ? 'bg-red-600 hover:bg-red-700'
                : 'bg-blue-600 hover:bg-blue-700'
            }`}
          >
            {isLoading ? '処理中...' : confirmText}
          </button>
        </div>
      </div>
    </dialog>
  );
}
