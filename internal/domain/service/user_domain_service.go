package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/wxlbd/polaris/internal/domain/entity"
	"github.com/wxlbd/polaris/internal/domain/repository"
)

// PasswordHasher 密码哈希器接口
// 定义密码加密的抽象，允许切换不同的加密算法
type PasswordHasher interface {
	// Hash 对密码进行哈希
	Hash(password string) (string, error)
	// Verify 验证密码是否匹配
	Verify(password, hash string) bool
}

// UserDomainService 用户领域服务
// 领域服务用于处理：
// 1. 跨多个实体的业务逻辑
// 2. 不适合放在单个实体中的复杂业务规则
// 3. 需要调用仓储或外部服务的领域逻辑
type UserDomainService struct {
	userRepo       repository.UserRepository
	passwordHasher PasswordHasher
}

// NewUserDomainService 创建用户领域服务
func NewUserDomainService(userRepo repository.UserRepository, passwordHasher PasswordHasher) *UserDomainService {
	return &UserDomainService{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

// CheckUserExists 检查用户是否存在
// 这是一个领域服务方法，因为它需要访问仓储来验证业务规则
func (s *UserDomainService) CheckUserExists(ctx context.Context, openID string) (bool, error) {
	user, err := s.userRepo.FindByOpenID(ctx, openID)
	if err != nil {
		return false, errors.Wrap(err, "查询用户失败")
	}
	return user != nil, nil
}

// CanUserPerformAction 检查用户是否有权限执行某个操作
// 这是一个领域规则，根据用户状态和业务规则判断
func (s *UserDomainService) CanUserPerformAction(ctx context.Context, userID int64, action string) (bool, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, errors.Wrap(err, "查询用户失败")
	}

	if user == nil {
		return false, errors.New("用户不存在")
	}

	// 这里可以添加更复杂的权限检查逻辑
	// 例如：检查用户角色、会员状态、账户状态等
	switch action {
	case "create_content":
		// 所有登录用户都可以创建内容
		return true, nil
	case "admin_action":
		// 只有管理员可以执行管理操作（这里需要扩展User实体添加角色字段）
		return false, nil
	default:
		return false, nil
	}
}

// IsUserActive 检查用户是否活跃（最近30天有登录）
func (s *UserDomainService) IsUserActive(ctx context.Context, userID int64) (bool, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, errors.Wrap(err, "查询用户失败")
	}

	if user == nil {
		return false, errors.New("用户不存在")
	}

	// 检查最后登录时间是否在30天内
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).UnixMilli()
	return user.LastLoginTime > thirtyDaysAgo, nil
}

// MergeUserAccounts 合并用户账户
// 这是一个复杂的领域操作，涉及多个实体的修改
// 在实际项目中，这个操作可能需要事务支持
func (s *UserDomainService) MergeUserAccounts(ctx context.Context, primaryUserID, secondaryUserID int64) error {
	primaryUser, err := s.userRepo.FindByID(ctx, primaryUserID)
	if err != nil {
		return errors.Wrap(err, "查询主账户失败")
	}
	if primaryUser == nil {
		return errors.New("主账户不存在")
	}

	secondaryUser, err := s.userRepo.FindByID(ctx, secondaryUserID)
	if err != nil {
		return errors.Wrap(err, "查询次账户失败")
	}
	if secondaryUser == nil {
		return errors.New("次账户不存在")
	}

	// 这里应该实现账户合并的具体逻辑
	// 例如：合并用户数据、转移关联资源等
	// 实际实现需要根据业务需求来定义

	return nil
}

// GenerateUserToken 生成用户令牌
// 这是一个领域服务方法，封装令牌生成的业务逻辑
func (s *UserDomainService) GenerateUserToken(user *entity.User) (string, error) {
	if user == nil {
		return "", errors.New("用户对象不能为空")
	}

	// 生成随机字节
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", errors.Wrap(err, "生成随机数失败")
	}

	// 组合用户信息和随机数生成令牌
	data := fmt.Sprintf("%d:%s:%d:%s",
		user.ID,
		user.OpenID,
		time.Now().UnixNano(),
		hex.EncodeToString(randomBytes),
	)

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:]), nil
}
