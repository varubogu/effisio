package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AccessTokenClaims はアクセストークンのクレームです
type AccessTokenClaims struct {
	UserID      uint     `json:"user_id"`
	Username    string   `json:"username"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// RefreshTokenClaims はリフレッシュトークンのクレームです
type RefreshTokenClaims struct {
	UserID  uint   `json:"user_id"`
	TokenID string `json:"token_id"`
	jwt.RegisteredClaims
}

// JWTService はJWT関連の処理を提供します
type JWTService struct {
	secret                  []byte
	accessTokenExpiration   time.Duration
	refreshTokenExpiration  time.Duration
}

// NewJWTService は新しいJWTServiceを作成します
func NewJWTService(secret string, accessTokenExpiration, refreshTokenExpiration time.Duration) *JWTService {
	return &JWTService{
		secret:                 []byte(secret),
		accessTokenExpiration:  accessTokenExpiration,
		refreshTokenExpiration: refreshTokenExpiration,
	}
}

// GenerateAccessToken はアクセストークンを生成します
func (s *JWTService) GenerateAccessToken(userID uint, username, role string, permissions []string) (string, error) {
	now := time.Now()
	claims := &AccessTokenClaims{
		UserID:      userID,
		Username:    username,
		Role:        role,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "effisio",
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// GenerateRefreshToken はリフレッシュトークンを生成します
func (s *JWTService) GenerateRefreshToken(userID uint, tokenID string) (string, error) {
	now := time.Now()
	claims := &RefreshTokenClaims{
		UserID:  userID,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "effisio",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// ValidateAccessToken はアクセストークンを検証します
func (s *JWTService) ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 署名方式の確認
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateRefreshToken はリフレッシュトークンを検証します
func (s *JWTService) ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 署名方式の確認
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*RefreshTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ExtractTokenFromAuthHeader はAuthorizationヘッダーからトークンを抽出します
func ExtractTokenFromAuthHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	// "Bearer <token>" の形式を想定
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) {
		return "", errors.New("invalid authorization header format")
	}

	if authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("authorization header must start with 'Bearer '")
	}

	return authHeader[len(bearerPrefix):], nil
}

// GetPermissionsForRole はロールに基づいて権限リストを返します
func GetPermissionsForRole(role string) []string {
	permissionMap := map[string][]string{
		"admin": {
			"users:read",
			"users:write",
			"users:delete",
			"tasks:read",
			"tasks:write",
			"tasks:delete",
			"settings:read",
			"settings:write",
		},
		"manager": {
			"users:read",
			"tasks:read",
			"tasks:write",
			"tasks:delete",
		},
		"user": {
			"tasks:read",
			"tasks:write",
		},
		"viewer": {
			"tasks:read",
		},
	}

	if permissions, ok := permissionMap[role]; ok {
		return permissions
	}

	return []string{} // デフォルトは権限なし
}
