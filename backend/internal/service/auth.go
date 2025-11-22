package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/internal/repository"
	"github.com/varubogu/effisio/backend/pkg/util"
)

// AuthService は認証関連のビジネスロジックを提供します
type AuthService struct {
	userRepo         *repository.UserRepository
	refreshTokenRepo *repository.RefreshTokenRepository
	jwtService       *util.JWTService
	logger           *zap.Logger
}

// NewAuthService は新しいAuthServiceを作成します
func NewAuthService(
	userRepo *repository.UserRepository,
	refreshTokenRepo *repository.RefreshTokenRepository,
	jwtService *util.JWTService,
	logger *zap.Logger,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtService:       jwtService,
		logger:           logger,
	}
}

// LoginRequest はログインリクエストです
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse はログインレスポンスです
type LoginResponse struct {
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token"`
	User         *model.UserResponse `json:"user"`
}

// RefreshTokenRequest はリフレッシュトークンリクエストです
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse はリフレッシュトークンレスポンスです
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Login はユーザー名とパスワードで認証します
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// ユーザー名でユーザーを取得
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Login attempt with invalid username", zap.String("username", req.Username))
			return nil, util.NewUnauthorizedError(util.ErrCodeInvalidCredentials, errors.New("invalid credentials"))
		}
		s.logger.Error("Failed to find user", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// ユーザーのステータスをチェック
	if user.Status != model.UserStatusActive {
		s.logger.Warn("Login attempt by inactive user", zap.String("username", req.Username), zap.String("status", user.Status))
		return nil, util.NewForbiddenError(util.ErrCodeInsufficientPermission, errors.New("user account is not active"))
	}

	// パスワードを検証
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.logger.Warn("Login attempt with invalid password", zap.String("username", req.Username))
		return nil, util.NewUnauthorizedError(util.ErrCodeInvalidCredentials, errors.New("invalid credentials"))
	}

	// 権限リストを取得
	permissions := util.GetPermissionsForRole(user.Role)

	// アクセストークンを生成
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Username, user.Role, permissions)
	if err != nil {
		s.logger.Error("Failed to generate access token", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeInternalError, err)
	}

	// リフレッシュトークンを生成
	tokenID := uuid.New().String()
	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, tokenID)
	if err != nil {
		s.logger.Error("Failed to generate refresh token", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeInternalError, err)
	}

	// リフレッシュトークンをデータベースに保存
	refreshTokenModel := &model.RefreshToken{
		UserID:    user.ID,
		TokenID:   tokenID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7日間有効
		Revoked:   false,
	}
	if err := s.refreshTokenRepo.Create(ctx, refreshTokenModel); err != nil {
		s.logger.Error("Failed to save refresh token", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// 最終ログイン時刻を更新
	now := time.Now()
	user.LastLogin = &now
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Warn("Failed to update last login time", zap.Error(err))
		// ログイン時刻の更新失敗は致命的ではないので続行
	}

	s.logger.Info("User logged in successfully", zap.String("username", user.Username))

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
	}, nil
}

// RefreshToken はリフレッシュトークンで新しいアクセストークンを発行します
func (s *AuthService) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	// リフレッシュトークンを検証
	claims, err := s.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		s.logger.Warn("Invalid refresh token", zap.Error(err))
		return nil, util.NewUnauthorizedError(util.ErrCodeInvalidToken, err)
	}

	// データベースからリフレッシュトークンを取得
	refreshToken, err := s.refreshTokenRepo.FindByTokenID(ctx, claims.TokenID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Refresh token not found in database", zap.String("token_id", claims.TokenID))
			return nil, util.NewUnauthorizedError(util.ErrCodeInvalidToken, errors.New("invalid refresh token"))
		}
		s.logger.Error("Failed to find refresh token", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// トークンが有効かチェック
	if !refreshToken.IsValid() {
		s.logger.Warn("Revoked or expired refresh token used", zap.String("token_id", claims.TokenID))
		return nil, util.NewUnauthorizedError(util.ErrCodeTokenExpired, errors.New("refresh token is revoked or expired"))
	}

	// ユーザー情報を取得
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewUnauthorizedError(util.ErrCodeUserNotFound, err)
		}
		s.logger.Error("Failed to find user", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	// ユーザーのステータスをチェック
	if user.Status != model.UserStatusActive {
		s.logger.Warn("Refresh token used by inactive user", zap.Uint("user_id", user.ID))
		return nil, util.NewForbiddenError(util.ErrCodeInsufficientPermission, errors.New("user account is not active"))
	}

	// 権限リストを取得
	permissions := util.GetPermissionsForRole(user.Role)

	// 新しいアクセストークンを生成
	newAccessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Username, user.Role, permissions)
	if err != nil {
		s.logger.Error("Failed to generate access token", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeInternalError, err)
	}

	// 新しいリフレッシュトークンを生成（トークンローテーション）
	newTokenID := uuid.New().String()
	newRefreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, newTokenID)
	if err != nil {
		s.logger.Error("Failed to generate refresh token", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeInternalError, err)
	}

	// 古いリフレッシュトークンを無効化
	if err := s.refreshTokenRepo.Revoke(ctx, claims.TokenID); err != nil {
		s.logger.Error("Failed to revoke old refresh token", zap.Error(err))
		// 無効化失敗は致命的ではないので続行
	}

	// 新しいリフレッシュトークンをデータベースに保存
	newRefreshTokenModel := &model.RefreshToken{
		UserID:    user.ID,
		TokenID:   newTokenID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}
	if err := s.refreshTokenRepo.Create(ctx, newRefreshTokenModel); err != nil {
		s.logger.Error("Failed to save refresh token", zap.Error(err))
		return nil, util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	s.logger.Info("Tokens refreshed successfully", zap.Uint("user_id", user.ID))

	return &RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// Logout はリフレッシュトークンを無効化してログアウトします
func (s *AuthService) Logout(ctx context.Context, refreshTokenString string) error {
	// リフレッシュトークンを検証
	claims, err := s.jwtService.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		// トークンが無効でもログアウト処理は成功とする
		s.logger.Warn("Invalid refresh token on logout", zap.Error(err))
		return nil
	}

	// リフレッシュトークンを無効化
	if err := s.refreshTokenRepo.Revoke(ctx, claims.TokenID); err != nil {
		s.logger.Error("Failed to revoke refresh token", zap.Error(err))
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	s.logger.Info("User logged out successfully", zap.Uint("user_id", claims.UserID))
	return nil
}

// LogoutAll はユーザーの全リフレッシュトークンを無効化します
func (s *AuthService) LogoutAll(ctx context.Context, userID uint) error {
	if err := s.refreshTokenRepo.RevokeAllByUserID(ctx, userID); err != nil {
		s.logger.Error("Failed to revoke all refresh tokens", zap.Uint("user_id", userID), zap.Error(err))
		return util.NewInternalError(util.ErrCodeDatabaseError, err)
	}

	s.logger.Info("All sessions logged out", zap.Uint("user_id", userID))
	return nil
}
