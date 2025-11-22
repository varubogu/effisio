# コード規約の具体例

このドキュメントでは、Effisioプロジェクトで使用するコーディング規約の具体的な実装例を提供します。

## 目次

- [Goコード例](#goコード例)
- [TypeScript/Reactコード例](#typescriptreactコード例)
- [テストコード例](#テストコード例)
- [よくある間違いと正しい実装](#よくある間違いと正しい実装)

---

## Goコード例

### 1. ハンドラーの標準パターン

```go
package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    "github.com/varubogu/effisio/backend/internal/model"
    "github.com/varubogu/effisio/backend/internal/service"
    "github.com/varubogu/effisio/backend/internal/util"
)

// UserHandler はユーザー関連のHTTPハンドラー
type UserHandler struct {
    service *service.UserService
    logger  *zap.Logger
}

// NewUserHandler は新しいUserHandlerを作成
func NewUserHandler(service *service.UserService, logger *zap.Logger) *UserHandler {
    return &UserHandler{
        service: service,
        logger:  logger,
    }
}

// List はユーザー一覧を取得
func (h *UserHandler) List(c *gin.Context) {
    // 1. ページネーションパラメータ取得
    params := util.GetPaginationParams(c)

    // 2. サービス層呼び出し
    users, pagination, err := h.service.List(c.Request.Context(), params)
    if err != nil {
        h.handleError(c, err)
        return
    }

    // 3. 成功レスポンス返却
    util.Paginated(c, users, pagination)
}

// GetByID はIDでユーザーを取得
func (h *UserHandler) GetByID(c *gin.Context) {
    // 1. パスパラメータ検証
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        util.Error(c, http.StatusBadRequest, util.ErrCodeInvalidFormat, "Invalid user ID")
        return
    }

    // 2. サービス層呼び出し
    user, err := h.service.GetByID(c.Request.Context(), uint(id))
    if err != nil {
        h.handleError(c, err)
        return
    }

    // 3. 成功レスポンス
    util.Success(c, gin.H{"user": user})
}

// Create は新しいユーザーを作成
func (h *UserHandler) Create(c *gin.Context) {
    // 1. リクエストボディをバインド
    var req model.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        util.ValidationError(c, util.ParseValidationErrors(err))
        return
    }

    // 2. サービス層呼び出し
    user, err := h.service.Create(c.Request.Context(), &req)
    if err != nil {
        h.handleError(c, err)
        return
    }

    // 3. 成功レスポンス (201 Created)
    util.Created(c, gin.H{"user": user})
}

// Update はユーザー情報を更新
func (h *UserHandler) Update(c *gin.Context) {
    // 1. パスパラメータ検証
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        util.Error(c, http.StatusBadRequest, util.ErrCodeInvalidFormat, "Invalid user ID")
        return
    }

    // 2. リクエストボディをバインド
    var req model.UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        util.ValidationError(c, util.ParseValidationErrors(err))
        return
    }

    // 3. サービス層呼び出し
    user, err := h.service.Update(c.Request.Context(), uint(id), &req)
    if err != nil {
        h.handleError(c, err)
        return
    }

    // 4. 成功レスポンス
    util.Success(c, gin.H{"user": user})
}

// Delete はユーザーを削除
func (h *UserHandler) Delete(c *gin.Context) {
    // 1. パスパラメータ検証
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        util.Error(c, http.StatusBadRequest, util.ErrCodeInvalidFormat, "Invalid user ID")
        return
    }

    // 2. サービス層呼び出し
    if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
        h.handleError(c, err)
        return
    }

    // 3. 成功レスポンス
    util.Success(c, gin.H{"message": "User deleted successfully"})
}

// handleError はエラーを適切なHTTPレスポンスに変換
func (h *UserHandler) handleError(c *gin.Context, err error) {
    if appErr, ok := err.(*util.AppError); ok {
        util.Error(c, appErr.StatusCode, appErr.Code, appErr.Message)
        return
    }

    // 予期しないエラー
    h.logger.Error("Unexpected error", zap.Error(err))
    util.Error(c, http.StatusInternalServerError, util.ErrCodeInternal, "Internal server error")
}
```

### 2. サービス層の標準パターン

```go
package service

import (
    "context"
    "errors"
    "fmt"

    "go.uber.org/zap"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"

    "github.com/varubogu/effisio/backend/internal/model"
    "github.com/varubogu/effisio/backend/internal/repository"
    "github.com/varubogu/effisio/backend/internal/util"
)

// UserService はユーザー関連のビジネスロジック
type UserService struct {
    repo   *repository.UserRepository
    logger *zap.Logger
}

// NewUserService は新しいUserServiceを作成
func NewUserService(repo *repository.UserRepository, logger *zap.Logger) *UserService {
    return &UserService{
        repo:   repo,
        logger: logger,
    }
}

// List はユーザー一覧を取得
func (s *UserService) List(ctx context.Context, params *util.PaginationParams) ([]*model.UserResponse, *util.PaginationInfo, error) {
    // 1. リポジトリ層呼び出し
    users, total, err := s.repo.FindAll(ctx, params)
    if err != nil {
        s.logger.Error("Failed to fetch users", zap.Error(err))
        return nil, nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
    }

    // 2. レスポンス変換
    responses := make([]*model.UserResponse, len(users))
    for i, user := range users {
        responses[i] = user.ToResponse()
    }

    // 3. ページネーション情報生成
    pagination := util.NewPaginationInfo(params.Page, params.PerPage, total)

    return responses, pagination, nil
}

// GetByID はIDでユーザーを取得
func (s *UserService) GetByID(ctx context.Context, id uint) (*model.UserResponse, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, util.NewNotFoundError(util.ErrCodeUserNotFound, "User not found")
        }
        s.logger.Error("Failed to fetch user", zap.Error(err), zap.Uint("id", id))
        return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
    }

    return user.ToResponse(), nil
}

// Create は新しいユーザーを作成
func (s *UserService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
    // 1. バリデーション（追加のビジネスルール）
    if err := s.validateCreateRequest(ctx, req); err != nil {
        return nil, err
    }

    // 2. パスワードハッシュ化
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        s.logger.Error("Failed to hash password", zap.Error(err))
        return nil, util.NewInternalError(util.ErrCodeInternal, err)
    }

    // 3. デフォルト値設定
    role := req.Role
    if role == "" {
        role = string(model.RoleUser)
    }

    // 4. エンティティ作成
    user := &model.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: string(hashedPassword),
        FullName:     &req.FullName,
        Department:   &req.Department,
        Role:         role,
        Status:       "active",
    }

    // 5. リポジトリ層呼び出し
    if err := s.repo.Create(ctx, user); err != nil {
        s.logger.Error("Failed to create user", zap.Error(err))
        return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
    }

    // 6. ログ記録
    s.logger.Info("User created successfully",
        zap.Uint("id", user.ID),
        zap.String("username", user.Username),
    )

    return user.ToResponse(), nil
}

// validateCreateRequest は作成リクエストの追加バリデーション
func (s *UserService) validateCreateRequest(ctx context.Context, req *model.CreateUserRequest) error {
    // ユーザー名の重複チェック
    existingUser, err := s.repo.FindByUsername(ctx, req.Username)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return util.NewInternalError(util.ErrCodeDatabaseError, err)
    }
    if existingUser != nil {
        return util.NewBadRequestError(util.ErrCodeUsernameDuplicate, "Username already exists")
    }

    // メールの重複チェック
    existingUser, err = s.repo.FindByEmail(ctx, req.Email)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return util.NewInternalError(util.ErrCodeDatabaseError, err)
    }
    if existingUser != nil {
        return util.NewBadRequestError(util.ErrCodeEmailDuplicate, "Email already exists")
    }

    // ロールの検証
    if req.Role != "" && !model.IsValidRole(req.Role) {
        return util.NewBadRequestError(util.ErrCodeValidation, "Invalid role")
    }

    return nil
}

// Update はユーザー情報を更新
func (s *UserService) Update(ctx context.Context, id uint, req *model.UpdateUserRequest) (*model.UserResponse, error) {
    // 1. 既存ユーザー取得
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, util.NewNotFoundError(util.ErrCodeUserNotFound, "User not found")
        }
        return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
    }

    // 2. 更新内容を適用
    if req.Username != "" {
        user.Username = req.Username
    }
    if req.Email != "" {
        user.Email = req.Email
    }
    if req.FullName != "" {
        user.FullName = &req.FullName
    }
    if req.Department != "" {
        user.Department = &req.Department
    }
    if req.Role != "" {
        if !model.IsValidRole(req.Role) {
            return nil, util.NewBadRequestError(util.ErrCodeValidation, "Invalid role")
        }
        user.Role = req.Role
    }
    if req.Status != "" {
        user.Status = req.Status
    }

    // 3. 更新実行
    if err := s.repo.Update(ctx, user); err != nil {
        s.logger.Error("Failed to update user", zap.Error(err), zap.Uint("id", id))
        return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
    }

    // 4. ログ記録
    s.logger.Info("User updated successfully", zap.Uint("id", id))

    return user.ToResponse(), nil
}

// Delete はユーザーを削除（ソフトデリート）
func (s *UserService) Delete(ctx context.Context, id uint) error {
    if err := s.repo.Delete(ctx, id); err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return util.NewNotFoundError(util.ErrCodeUserNotFound, "User not found")
        }
        s.logger.Error("Failed to delete user", zap.Error(err), zap.Uint("id", id))
        return util.NewInternalError(util.ErrCodeDatabaseError, err)
    }

    s.logger.Info("User deleted successfully", zap.Uint("id", id))

    return nil
}
```

### 3. リポジトリ層の標準パターン

```go
package repository

import (
    "context"

    "gorm.io/gorm"

    "github.com/varubogu/effisio/backend/internal/model"
    "github.com/varubogu/effisio/backend/internal/util"
)

// UserRepository はユーザーデータアクセス
type UserRepository struct {
    db *gorm.DB
}

// NewUserRepository は新しいUserRepositoryを作成
func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

// FindAll は全てのユーザーを取得（ページネーション付き）
func (r *UserRepository) FindAll(ctx context.Context, params *util.PaginationParams) ([]*model.User, int64, error) {
    var users []*model.User
    var total int64

    // 総件数取得
    if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // データ取得
    query := r.db.WithContext(ctx).
        Offset(params.Offset).
        Limit(params.PerPage).
        Order("created_at DESC")

    if err := query.Find(&users).Error; err != nil {
        return nil, 0, err
    }

    return users, total, nil
}

// FindByID はIDでユーザーを取得
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
    var user model.User
    if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

// FindByEmail はメールアドレスでユーザーを取得
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    var user model.User
    if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

// FindByUsername はユーザー名でユーザーを取得
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
    var user model.User
    if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

// Create は新しいユーザーを作成
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

// Update はユーザー情報を更新
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
    return r.db.WithContext(ctx).Save(user).Error
}

// Delete はユーザーを削除（ソフトデリート）
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}
```

---

## TypeScript/Reactコード例

### 1. カスタムフックの標準パターン

```typescript
// src/hooks/useUsers.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { usersApi } from '@/lib/users';
import type { CreateUserRequest, UpdateUserRequest } from '@/types/user';

const USERS_QUERY_KEY = ['users'] as const;

// ユーザー一覧を取得
export function useUsers(page: number = 1, perPage: number = 20) {
  return useQuery({
    queryKey: [...USERS_QUERY_KEY, { page, perPage }],
    queryFn: () => usersApi.getUsers(page, perPage),
    staleTime: 5 * 60 * 1000, // 5分
  });
}

// ユーザーをIDで取得
export function useUser(id: number) {
  return useQuery({
    queryKey: [...USERS_QUERY_KEY, id],
    queryFn: () => usersApi.getUserById(id),
    enabled: !!id, // idが存在する場合のみクエリ実行
  });
}

// ユーザーを作成
export function useCreateUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateUserRequest) => usersApi.createUser(data),
    onSuccess: () => {
      // ユーザー一覧のキャッシュを無効化
      queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY });
    },
    onError: (error) => {
      console.error('Failed to create user:', error);
    },
  });
}

// ユーザーを更新
export function useUpdateUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateUserRequest }) =>
      usersApi.updateUser(id, data),
    onSuccess: (_, variables) => {
      // 特定ユーザーのキャッシュを無効化
      queryClient.invalidateQueries({ queryKey: [...USERS_QUERY_KEY, variables.id] });
      // ユーザー一覧のキャッシュも無効化
      queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY });
    },
  });
}

// ユーザーを削除
export function useDeleteUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: number) => usersApi.deleteUser(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: USERS_QUERY_KEY });
    },
  });
}
```

### 2. コンポーネントの標準パターン

```typescript
// src/components/users/UserList.tsx
import type { FC } from 'react';
import type { User } from '@/types/user';
import { UserListItem } from './UserListItem';

interface UserListProps {
  users: User[];
  onEdit?: (user: User) => void;
  onDelete?: (userId: number) => void;
}

export const UserList: FC<UserListProps> = ({ users, onEdit, onDelete }) => {
  if (users.length === 0) {
    return (
      <div className="rounded-lg border border-gray-200 bg-white p-8 text-center">
        <p className="text-gray-500">ユーザーが見つかりません</p>
      </div>
    );
  }

  return (
    <div className="overflow-hidden rounded-lg border border-gray-200 bg-white shadow">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
              ID
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
              ユーザー名
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
              メールアドレス
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
              ロール
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
              ステータス
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
              アクション
            </th>
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-200 bg-white">
          {users.map((user) => (
            <UserListItem
              key={user.id}
              user={user}
              onEdit={onEdit}
              onDelete={onDelete}
            />
          ))}
        </tbody>
      </table>
    </div>
  );
};
```

```typescript
// src/components/users/UserListItem.tsx
import type { FC } from 'react';
import type { User } from '@/types/user';
import { RoleBadge } from '@/components/ui/RoleBadge';
import { StatusBadge } from '@/components/ui/StatusBadge';

interface UserListItemProps {
  user: User;
  onEdit?: (user: User) => void;
  onDelete?: (userId: number) => void;
}

export const UserListItem: FC<UserListItemProps> = ({ user, onEdit, onDelete }) => {
  return (
    <tr className="hover:bg-gray-50">
      <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-900">
        {user.id}
      </td>
      <td className="whitespace-nowrap px-6 py-4 text-sm font-medium text-gray-900">
        {user.username}
      </td>
      <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
        {user.email}
      </td>
      <td className="whitespace-nowrap px-6 py-4 text-sm">
        <RoleBadge role={user.role} />
      </td>
      <td className="whitespace-nowrap px-6 py-4 text-sm">
        <StatusBadge status={user.status} />
      </td>
      <td className="whitespace-nowrap px-6 py-4 text-sm">
        <div className="flex gap-2">
          {onEdit && (
            <button
              onClick={() => onEdit(user)}
              className="text-blue-600 hover:text-blue-900"
            >
              編集
            </button>
          )}
          {onDelete && (
            <button
              onClick={() => onDelete(user.id)}
              className="text-red-600 hover:text-red-900"
            >
              削除
            </button>
          )}
        </div>
      </td>
    </tr>
  );
};
```

### 3. フォームコンポーネント

```typescript
// src/components/users/UserForm.tsx
import type { FC } from 'react';
import { useForm } from 'react-hook-form';
import type { CreateUserRequest, User } from '@/types/user';

interface UserFormProps {
  user?: User;
  onSubmit: (data: CreateUserRequest) => void;
  onCancel?: () => void;
  isLoading?: boolean;
}

export const UserForm: FC<UserFormProps> = ({
  user,
  onSubmit,
  onCancel,
  isLoading = false,
}) => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateUserRequest>({
    defaultValues: user
      ? {
          username: user.username,
          email: user.email,
          full_name: user.full_name || '',
          department: user.department || '',
          role: user.role,
        }
      : undefined,
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div>
        <label htmlFor="username" className="block text-sm font-medium text-gray-700">
          ユーザー名 *
        </label>
        <input
          id="username"
          type="text"
          {...register('username', {
            required: 'ユーザー名は必須です',
            minLength: { value: 3, message: '3文字以上で入力してください' },
            maxLength: { value: 50, message: '50文字以内で入力してください' },
            pattern: {
              value: /^[a-zA-Z0-9_]+$/,
              message: '英数字とアンダースコアのみ使用できます',
            },
          })}
          className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-blue-500"
        />
        {errors.username && (
          <p className="mt-1 text-sm text-red-600">{errors.username.message}</p>
        )}
      </div>

      <div>
        <label htmlFor="email" className="block text-sm font-medium text-gray-700">
          メールアドレス *
        </label>
        <input
          id="email"
          type="email"
          {...register('email', {
            required: 'メールアドレスは必須です',
            pattern: {
              value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
              message: '有効なメールアドレスを入力してください',
            },
          })}
          className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-blue-500"
        />
        {errors.email && (
          <p className="mt-1 text-sm text-red-600">{errors.email.message}</p>
        )}
      </div>

      {!user && (
        <div>
          <label htmlFor="password" className="block text-sm font-medium text-gray-700">
            パスワード *
          </label>
          <input
            id="password"
            type="password"
            {...register('password', {
              required: 'パスワードは必須です',
              minLength: { value: 8, message: '8文字以上で入力してください' },
            })}
            className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-blue-500"
          />
          {errors.password && (
            <p className="mt-1 text-sm text-red-600">{errors.password.message}</p>
          )}
        </div>
      )}

      <div>
        <label htmlFor="full_name" className="block text-sm font-medium text-gray-700">
          氏名
        </label>
        <input
          id="full_name"
          type="text"
          {...register('full_name')}
          className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-blue-500"
        />
      </div>

      <div>
        <label htmlFor="department" className="block text-sm font-medium text-gray-700">
          部署
        </label>
        <input
          id="department"
          type="text"
          {...register('department')}
          className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-blue-500"
        />
      </div>

      <div>
        <label htmlFor="role" className="block text-sm font-medium text-gray-700">
          ロール
        </label>
        <select
          id="role"
          {...register('role')}
          className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-blue-500"
        >
          <option value="user">一般ユーザー</option>
          <option value="manager">マネージャー</option>
          <option value="admin">管理者</option>
          <option value="viewer">閲覧者</option>
        </select>
      </div>

      <div className="flex justify-end gap-3">
        {onCancel && (
          <button
            type="button"
            onClick={onCancel}
            className="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
          >
            キャンセル
          </button>
        )}
        <button
          type="submit"
          disabled={isLoading}
          className="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
        >
          {isLoading ? '保存中...' : user ? '更新' : '作成'}
        </button>
      </div>
    </form>
  );
};
```

---

## テストコード例

### 1. Goのユニットテスト

```go
package service_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "go.uber.org/zap"
    "gorm.io/gorm"

    "github.com/varubogu/effisio/backend/internal/model"
    "github.com/varubogu/effisio/backend/internal/service"
    "github.com/varubogu/effisio/backend/internal/util"
)

// モックリポジトリ
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
    args := m.Called(ctx, username)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    args := m.Called(ctx, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func TestUserService_GetByID(t *testing.T) {
    logger, _ := zap.NewDevelopment()
    mockRepo := new(MockUserRepository)
    userService := service.NewUserService(mockRepo, logger)

    t.Run("成功ケース", func(t *testing.T) {
        // モックの設定
        expectedUser := &model.User{
            ID:       1,
            Username: "testuser",
            Email:    "test@example.com",
            Role:     "user",
        }
        mockRepo.On("FindByID", mock.Anything, uint(1)).Return(expectedUser, nil)

        // テスト実行
        result, err := userService.GetByID(context.Background(), 1)

        // 検証
        assert.NoError(t, err)
        assert.NotNil(t, result)
        assert.Equal(t, "testuser", result.Username)
        mockRepo.AssertExpectations(t)
    })

    t.Run("ユーザーが見つからない", func(t *testing.T) {
        mockRepo := new(MockUserRepository)
        userService := service.NewUserService(mockRepo, logger)

        mockRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

        result, err := userService.GetByID(context.Background(), 999)

        assert.Error(t, err)
        assert.Nil(t, result)

        // エラーの種類を確認
        appErr, ok := err.(*util.AppError)
        assert.True(t, ok)
        assert.Equal(t, util.ErrCodeUserNotFound, appErr.Code)

        mockRepo.AssertExpectations(t)
    })
}

func TestUserService_Create(t *testing.T) {
    logger, _ := zap.NewDevelopment()

    t.Run("成功ケース", func(t *testing.T) {
        mockRepo := new(MockUserRepository)
        userService := service.NewUserService(mockRepo, logger)

        req := &model.CreateUserRequest{
            Username: "newuser",
            Email:    "new@example.com",
            Password: "SecurePass123!",
            Role:     "user",
        }

        // ユーザー名とメールの重複チェック
        mockRepo.On("FindByUsername", mock.Anything, "newuser").Return(nil, gorm.ErrRecordNotFound)
        mockRepo.On("FindByEmail", mock.Anything, "new@example.com").Return(nil, gorm.ErrRecordNotFound)
        mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)

        result, err := userService.Create(context.Background(), req)

        assert.NoError(t, err)
        assert.NotNil(t, result)
        assert.Equal(t, "newuser", result.Username)
        mockRepo.AssertExpectations(t)
    })

    t.Run("ユーザー名重複エラー", func(t *testing.T) {
        mockRepo := new(MockUserRepository)
        userService := service.NewUserService(mockRepo, logger)

        req := &model.CreateUserRequest{
            Username: "existinguser",
            Email:    "new@example.com",
            Password: "SecurePass123!",
        }

        existingUser := &model.User{
            ID:       1,
            Username: "existinguser",
        }

        mockRepo.On("FindByUsername", mock.Anything, "existinguser").Return(existingUser, nil)

        result, err := userService.Create(context.Background(), req)

        assert.Error(t, err)
        assert.Nil(t, result)

        appErr, ok := err.(*util.AppError)
        assert.True(t, ok)
        assert.Equal(t, util.ErrCodeUsernameDuplicate, appErr.Code)

        mockRepo.AssertExpectations(t)
    })
}
```

### 2. Reactコンポーネントのテスト

```typescript
// src/components/users/UserList.test.tsx
import { render, screen } from '@testing-library/react';
import { UserList } from './UserList';
import type { User } from '@/types/user';

const mockUsers: User[] = [
  {
    id: 1,
    username: 'john_doe',
    email: 'john@example.com',
    full_name: 'John Doe',
    department: 'Engineering',
    role: 'admin',
    status: 'active',
    last_login: '2024-01-16T10:30:00Z',
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-16T10:30:00Z',
  },
  {
    id: 2,
    username: 'jane_smith',
    email: 'jane@example.com',
    full_name: 'Jane Smith',
    department: 'HR',
    role: 'user',
    status: 'active',
    last_login: '2024-01-15T09:15:00Z',
    created_at: '2024-01-02T00:00:00Z',
    updated_at: '2024-01-15T09:15:00Z',
  },
];

describe('UserList', () => {
  it('ユーザーリストを表示する', () => {
    render(<UserList users={mockUsers} />);

    // ユーザー名が表示されているか確認
    expect(screen.getByText('john_doe')).toBeInTheDocument();
    expect(screen.getByText('jane_smith')).toBeInTheDocument();

    // メールアドレスが表示されているか確認
    expect(screen.getByText('john@example.com')).toBeInTheDocument();
    expect(screen.getByText('jane@example.com')).toBeInTheDocument();
  });

  it('ユーザーが0件の場合、メッセージを表示する', () => {
    render(<UserList users={[]} />);

    expect(screen.getByText('ユーザーが見つかりません')).toBeInTheDocument();
  });
});
```

---

## よくある間違いと正しい実装

### 1. エラーハンドリング

❌ **悪い例:**
```go
func (h *UserHandler) GetByID(c *gin.Context) {
    user, err := h.service.GetByID(c.Request.Context(), 1)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})  // エラーメッセージをそのまま返す（セキュリティリスク）
        return
    }
    c.JSON(200, user)  // 統一されていないレスポンス形式
}
```

✅ **良い例:**
```go
func (h *UserHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        util.Error(c, http.StatusBadRequest, util.ErrCodeInvalidFormat, "Invalid user ID")
        return
    }

    user, err := h.service.GetByID(c.Request.Context(), uint(id))
    if err != nil {
        h.handleError(c, err)  // 統一されたエラーハンドリング
        return
    }

    util.Success(c, gin.H{"user": user})  // 統一されたレスポンス形式
}
```

### 2. データベースクエリ

❌ **悪い例:**
```go
// N+1問題
func (r *UserRepository) FindAllWithDepartments(ctx context.Context) ([]*model.User, error) {
    var users []*model.User
    r.db.WithContext(ctx).Find(&users)

    for _, user := range users {
        // 各ユーザーごとにクエリ実行（N+1問題）
        var dept model.Department
        r.db.First(&dept, user.DepartmentID)
        user.Department = &dept
    }

    return users, nil
}
```

✅ **良い例:**
```go
// Preloadを使用
func (r *UserRepository) FindAllWithDepartments(ctx context.Context) ([]*model.User, error) {
    var users []*model.User
    if err := r.db.WithContext(ctx).
        Preload("Department").
        Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}
```

### 3. フロントエンド状態管理

❌ **悪い例:**
```typescript
// useEffectでデータ取得（古いパターン）
function UserList() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    setLoading(true);
    fetch('/api/users')
      .then(res => res.json())
      .then(data => {
        setUsers(data);
        setLoading(false);
      })
      .catch(err => {
        setError(err);
        setLoading(false);
      });
  }, []);

  // ...
}
```

✅ **良い例:**
```typescript
// TanStack Queryを使用
function UserList() {
  const { data: users, isLoading, error } = useUsers();

  if (isLoading) return <LoadingSpinner />;
  if (error) return <ErrorMessage error={error} />;

  return <UserListTable users={users || []} />;
}
```

---

## 変更履歴

| 日付 | バージョン | 変更内容 |
|------|-----------|---------|
| 2024-01-16 | 1.0.0 | 初版作成 |
