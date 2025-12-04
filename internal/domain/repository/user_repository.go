package repository

import (
	"context"

	"github.com/wxlbd/polaris/internal/domain/entity"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *entity.User) error
	// FindByOpenID 根据OpenID查找用户
	FindByOpenID(ctx context.Context, openID string) (*entity.User, error)
	// FindByID 根据ID查找用户
	FindByID(ctx context.Context, userID int64) (*entity.User, error)
	// Update 更新用户
	Update(ctx context.Context, user *entity.User) error
	// UpdateLastLoginTime 更新最后登录时间
	UpdateLastLoginTime(ctx context.Context, openID string) error
}
