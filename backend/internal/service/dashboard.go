package service

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/varubogu/effisio/backend/internal/model"
	"github.com/varubogu/effisio/backend/internal/repository"
)

// DashboardOverview はダッシュボード概要情報です
type DashboardOverview struct {
	TotalUsers      int64                    `json:"total_users"`
	ActiveUsers     int64                    `json:"active_users"`
	InactiveUsers   int64                    `json:"inactive_users"`
	SuspendedUsers  int64                    `json:"suspended_users"`
	LastLoginStats  []UserLoginStat          `json:"last_login_stats"`
	UsersByRole     map[string]int64         `json:"users_by_role"`
	UsersByDept     []DepartmentStat         `json:"users_by_department"`
}

// UserLoginStat はユーザーのログイン統計です
type UserLoginStat struct {
	Date   string `json:"date"`
	Count  int64  `json:"count"`
}

// DepartmentStat は部門別ユーザー統計です
type DepartmentStat struct {
	Department string `json:"department"`
	Count      int64  `json:"count"`
}

// DashboardService はダッシュボード関連のサービスです
type DashboardService struct {
	userRepo *repository.UserRepository
	logger   *zap.Logger
}

// NewDashboardService は新しいDashboardServiceを作成します
func NewDashboardService(userRepo *repository.UserRepository, logger *zap.Logger) *DashboardService {
	return &DashboardService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// GetOverview はダッシュボード概要を取得します
func (s *DashboardService) GetOverview(ctx context.Context) (*DashboardOverview, error) {
	// 全ユーザー数を取得
	totalUsers, err := s.userRepo.CountAll(ctx)
	if err != nil {
		s.logger.Error("failed to count all users", zap.Error(err))
		return nil, err
	}

	// ステータス別ユーザー数を取得
	activeUsers, err := s.userRepo.CountByStatus(ctx, model.UserStatusActive)
	if err != nil {
		s.logger.Error("failed to count active users", zap.Error(err))
		return nil, err
	}

	inactiveUsers, err := s.userRepo.CountByStatus(ctx, model.UserStatusInactive)
	if err != nil {
		s.logger.Error("failed to count inactive users", zap.Error(err))
		return nil, err
	}

	suspendedUsers, err := s.userRepo.CountByStatus(ctx, model.UserStatusSuspended)
	if err != nil {
		s.logger.Error("failed to count suspended users", zap.Error(err))
		return nil, err
	}

	// ロール別ユーザー数を取得
	usersByRole, err := s.userRepo.CountByRole(ctx)
	if err != nil {
		s.logger.Error("failed to count users by role", zap.Error(err))
		return nil, err
	}

	// 部門別ユーザー数を取得
	deptResults, err := s.userRepo.CountByDepartment(ctx)
	if err != nil {
		s.logger.Error("failed to count users by department", zap.Error(err))
		return nil, err
	}

	usersByDept := make([]DepartmentStat, len(deptResults))
	for i, result := range deptResults {
		usersByDept[i] = DepartmentStat{
			Department: result.Department,
			Count:      result.Count,
		}
	}

	// 過去7日のログイン統計を取得
	loginResults, err := s.userRepo.GetLastLoginStats(ctx, 7)
	if err != nil {
		s.logger.Error("failed to get login stats", zap.Error(err))
		return nil, err
	}

	loginStats := make([]UserLoginStat, len(loginResults))
	for i, result := range loginResults {
		loginStats[i] = UserLoginStat{
			Date:  result.Date,
			Count: result.Count,
		}
	}

	overview := &DashboardOverview{
		TotalUsers:     totalUsers,
		ActiveUsers:    activeUsers,
		InactiveUsers:  inactiveUsers,
		SuspendedUsers: suspendedUsers,
		UsersByRole:    usersByRole,
		UsersByDept:    usersByDept,
		LastLoginStats: loginStats,
	}

	return overview, nil
}
