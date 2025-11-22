package util

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJWTService(t *testing.T) {
	secret := "test-secret-key-for-testing-purpose"
	accessExp := 15 * time.Minute
	refreshExp := 7 * 24 * time.Hour

	svc := NewJWTService(secret, accessExp, refreshExp)

	assert.NotNil(t, svc)
	assert.Equal(t, []byte(secret), svc.secret)
	assert.Equal(t, accessExp, svc.accessTokenExpiration)
	assert.Equal(t, refreshExp, svc.refreshTokenExpiration)
}

func TestGenerateAccessToken(t *testing.T) {
	svc := NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)

	tests := []struct {
		name        string
		userID      uint
		username    string
		role        string
		permissions []string
		expectError bool
	}{
		{
			name:        "Valid access token generation",
			userID:      1,
			username:    "testuser",
			role:        "admin",
			permissions: []string{"users:read", "users:write"},
			expectError: false,
		},
		{
			name:        "Access token with empty permissions",
			userID:      2,
			username:    "user2",
			role:        "user",
			permissions: []string{},
			expectError: false,
		},
		{
			name:        "Access token with nil permissions",
			userID:      3,
			username:    "user3",
			role:        "viewer",
			permissions: nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := svc.GenerateAccessToken(tt.userID, tt.username, tt.role, tt.permissions)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				// Token should have 3 parts separated by dots (header.payload.signature)
				parts := strings.Split(token, ".")
				assert.Equal(t, 3, len(parts))
			}
		})
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	svc := NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)

	tests := []struct {
		name        string
		userID      uint
		tokenID     string
		expectError bool
	}{
		{
			name:        "Valid refresh token generation",
			userID:      1,
			tokenID:     "token-id-123",
			expectError: false,
		},
		{
			name:        "Refresh token with different user ID",
			userID:      999,
			tokenID:     "token-id-999",
			expectError: false,
		},
		{
			name:        "Refresh token with empty token ID",
			userID:      1,
			tokenID:     "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := svc.GenerateRefreshToken(tt.userID, tt.tokenID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				parts := strings.Split(token, ".")
				assert.Equal(t, 3, len(parts))
			}
		})
	}
}

func TestValidateAccessToken(t *testing.T) {
	secret := "test-secret-key"
	svc := NewJWTService(secret, 15*time.Minute, 7*24*time.Hour)

	// Generate a valid token
	validToken, err := svc.GenerateAccessToken(1, "testuser", "admin", []string{"users:read"})
	require.NoError(t, err)

	tests := []struct {
		name            string
		token           string
		expectError     bool
		expectClaims    bool
		validateClaims  func(t *testing.T, claims *AccessTokenClaims)
	}{
		{
			name:         "Valid access token",
			token:        validToken,
			expectError:  false,
			expectClaims: true,
			validateClaims: func(t *testing.T, claims *AccessTokenClaims) {
				assert.Equal(t, uint(1), claims.UserID)
				assert.Equal(t, "testuser", claims.Username)
				assert.Equal(t, "admin", claims.Role)
				assert.Equal(t, []string{"users:read"}, claims.Permissions)
				assert.NotNil(t, claims.ExpiresAt)
				assert.NotNil(t, claims.IssuedAt)
			},
		},
		{
			name:         "Empty token string",
			token:        "",
			expectError:  true,
			expectClaims: false,
		},
		{
			name:         "Malformed token",
			token:        "invalid.token.format",
			expectError:  true,
			expectClaims: false,
		},
		{
			name:         "Token with wrong signature",
			token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwicm9sZSI6ImFkbWluIiwicGVybWlzc2lvbnMiOlsidXNlcnM6cmVhZCJdLCJleHAiOjk5OTk5OTk5OTksImlhdCI6MTYwMDAwMDAwMCwibmJmIjoxNjAwMDAwMDAwLCJpc3MiOiJlZmZpc2lvIiwic3ViIjoidGVzdHVzZXIifQ.wrongsignature",
			expectError:  true,
			expectClaims: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := svc.ValidateAccessToken(tt.token)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				if tt.expectClaims && tt.validateClaims != nil {
					tt.validateClaims(t, claims)
				}
			}
		})
	}
}

func TestValidateRefreshToken(t *testing.T) {
	secret := "test-secret-key"
	svc := NewJWTService(secret, 15*time.Minute, 7*24*time.Hour)

	// Generate a valid refresh token
	validToken, err := svc.GenerateRefreshToken(1, "token-123")
	require.NoError(t, err)

	tests := []struct {
		name            string
		token           string
		expectError     bool
		expectClaims    bool
		validateClaims  func(t *testing.T, claims *RefreshTokenClaims)
	}{
		{
			name:         "Valid refresh token",
			token:        validToken,
			expectError:  false,
			expectClaims: true,
			validateClaims: func(t *testing.T, claims *RefreshTokenClaims) {
				assert.Equal(t, uint(1), claims.UserID)
				assert.Equal(t, "token-123", claims.TokenID)
				assert.NotNil(t, claims.ExpiresAt)
				assert.NotNil(t, claims.IssuedAt)
			},
		},
		{
			name:         "Empty token string",
			token:        "",
			expectError:  true,
			expectClaims: false,
		},
		{
			name:         "Malformed token",
			token:        "invalid.token",
			expectError:  true,
			expectClaims: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := svc.ValidateRefreshToken(tt.token)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				if tt.expectClaims && tt.validateClaims != nil {
					tt.validateClaims(t, claims)
				}
			}
		})
	}
}

func TestTokenExpiration(t *testing.T) {
	secret := "test-secret"
	// Create service with very short expiration
	svc := NewJWTService(secret, -1*time.Second, 7*24*time.Hour)

	// Generate an already-expired access token
	expiredToken, err := svc.GenerateAccessToken(1, "testuser", "user", nil)
	require.NoError(t, err)

	// The token should be generated, but validation should fail
	time.Sleep(100 * time.Millisecond) // Ensure token is expired
	claims, err := svc.ValidateAccessToken(expiredToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "expired")
}

func TestExtractTokenFromAuthHeader(t *testing.T) {
	tests := []struct {
		name        string
		authHeader  string
		expectError bool
		expectToken string
	}{
		{
			name:        "Valid Bearer token",
			authHeader:  "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.signature",
			expectError: false,
			expectToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.signature",
		},
		{
			name:        "Bearer with single character token",
			authHeader:  "Bearer x",
			expectError: false,
			expectToken: "x",
		},
		{
			name:        "Empty auth header",
			authHeader:  "",
			expectError: true,
			expectToken: "",
		},
		{
			name:        "Missing Bearer prefix",
			authHeader:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.signature",
			expectError: true,
			expectToken: "",
		},
		{
			name:        "Wrong prefix (Basic)",
			authHeader:  "Basic dXNlcm5hbWU6cGFzc3dvcmQ=",
			expectError: true,
			expectToken: "",
		},
		{
			name:        "Bearer with no space",
			authHeader:  "Bearertoken",
			expectError: true,
			expectToken: "",
		},
		{
			name:        "Lowercase bearer",
			authHeader:  "bearer token",
			expectError: true,
			expectToken: "",
		},
		{
			name:        "Just Bearer",
			authHeader:  "Bearer ",
			expectError: false,
			expectToken: "",
		},
		{
			name:        "Bearer with spaces in token",
			authHeader:  "Bearer token with spaces",
			expectError: false,
			expectToken: "token with spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := ExtractTokenFromAuthHeader(tt.authHeader)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectToken, token)
			}
		})
	}
}

func TestGetPermissionsForRole(t *testing.T) {
	tests := []struct {
		name            string
		role            string
		expectedMinPerms int
		expectedPerms   []string
	}{
		{
			name:            "Admin role permissions",
			role:            "admin",
			expectedMinPerms: 8,
			expectedPerms:   []string{"users:read", "users:write", "users:delete"},
		},
		{
			name:            "Manager role permissions",
			role:            "manager",
			expectedMinPerms: 4,
			expectedPerms:   []string{"users:read", "tasks:read", "tasks:write"},
		},
		{
			name:            "User role permissions",
			role:            "user",
			expectedMinPerms: 2,
			expectedPerms:   []string{"tasks:read", "tasks:write"},
		},
		{
			name:            "Viewer role permissions",
			role:            "viewer",
			expectedMinPerms: 1,
			expectedPerms:   []string{"tasks:read"},
		},
		{
			name:            "Unknown role",
			role:            "unknown",
			expectedMinPerms: 0,
			expectedPerms:   []string{},
		},
		{
			name:            "Empty role",
			role:            "",
			expectedMinPerms: 0,
			expectedPerms:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			perms := GetPermissionsForRole(tt.role)

			assert.Len(t, perms, tt.expectedMinPerms)
			for _, expectedPerm := range tt.expectedPerms {
				assert.Contains(t, perms, expectedPerm)
			}
		})
	}
}

func TestTokenRoundTrip(t *testing.T) {
	svc := NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)

	// Test access token round trip
	originalPerms := []string{"users:read", "users:write", "tasks:read"}
	token, err := svc.GenerateAccessToken(42, "johndoe", "manager", originalPerms)
	require.NoError(t, err)

	claims, err := svc.ValidateAccessToken(token)
	require.NoError(t, err)

	assert.Equal(t, uint(42), claims.UserID)
	assert.Equal(t, "johndoe", claims.Username)
	assert.Equal(t, "manager", claims.Role)
	assert.Equal(t, originalPerms, claims.Permissions)
	assert.Equal(t, "effisio", claims.Issuer)
	assert.Equal(t, "johndoe", claims.Subject)
}

func TestDifferentSecretsValidation(t *testing.T) {
	secret1 := "secret-key-1"
	secret2 := "secret-key-2"

	svc1 := NewJWTService(secret1, 15*time.Minute, 7*24*time.Hour)
	svc2 := NewJWTService(secret2, 15*time.Minute, 7*24*time.Hour)

	// Generate token with svc1
	token, err := svc1.GenerateAccessToken(1, "testuser", "user", nil)
	require.NoError(t, err)

	// Should validate with svc1
	claims, err := svc1.ValidateAccessToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	// Should NOT validate with svc2 (different secret)
	claims, err = svc2.ValidateAccessToken(token)
	assert.Error(t, err)
	assert.Nil(t, claims)
}
