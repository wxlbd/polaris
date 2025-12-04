package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wxlbd/polaris/internal/domain/valueobject"
)

// NotificationChannel 通知渠道
type NotificationChannel string

const (
	ChannelSMS    NotificationChannel = "sms"    // 短信
	ChannelEmail  NotificationChannel = "email"  // 邮件
	ChannelPush   NotificationChannel = "push"   // APP推送
	ChannelWechat NotificationChannel = "wechat" // 微信消息
)

// NotificationSender 通知发送器接口
// 使用接口实现依赖倒置，具体实现在 infrastructure 层
type NotificationSender interface {
	// SendSMS 发送短信
	SendSMS(ctx context.Context, phone valueobject.Phone, content string) error
	// SendEmail 发送邮件
	SendEmail(ctx context.Context, email valueobject.Email, subject, content string) error
	// SendPush 发送APP推送
	SendPush(ctx context.Context, userID int64, title, content string) error
	// SendWechat 发送微信消息
	SendWechat(ctx context.Context, openID, templateID string, data map[string]interface{}) error
}

// NotificationDomainService 通知领域服务
// 处理与通知相关的领域逻辑
type NotificationDomainService struct {
	sender NotificationSender
}

// NewNotificationDomainService 创建通知领域服务
func NewNotificationDomainService(sender NotificationSender) *NotificationDomainService {
	return &NotificationDomainService{
		sender: sender,
	}
}

// NotificationRequest 通知请求
type NotificationRequest struct {
	UserID    int64                  // 用户ID
	Channel   NotificationChannel    // 通知渠道
	Phone     valueobject.Phone      // 手机号（短信渠道）
	Email     valueobject.Email      // 邮箱（邮件渠道）
	OpenID    string                 // 微信OpenID（微信渠道）
	Title     string                 // 标题
	Content   string                 // 内容
	ExtraData map[string]interface{} // 额外数据
}

// SendNotification 发送通知
// 根据渠道类型选择合适的发送方式
func (s *NotificationDomainService) SendNotification(ctx context.Context, req NotificationRequest) error {
	switch req.Channel {
	case ChannelSMS:
		if req.Phone.IsEmpty() {
			return errors.New("发送短信需要提供手机号")
		}
		return s.sender.SendSMS(ctx, req.Phone, req.Content)

	case ChannelEmail:
		if req.Email.IsEmpty() {
			return errors.New("发送邮件需要提供邮箱地址")
		}
		return s.sender.SendEmail(ctx, req.Email, req.Title, req.Content)

	case ChannelPush:
		if req.UserID == 0 {
			return errors.New("发送推送需要提供用户ID")
		}
		return s.sender.SendPush(ctx, req.UserID, req.Title, req.Content)

	case ChannelWechat:
		if req.OpenID == "" {
			return errors.New("发送微信消息需要提供OpenID")
		}
		templateID, ok := req.ExtraData["template_id"].(string)
		if !ok || templateID == "" {
			return errors.New("发送微信消息需要提供模板ID")
		}
		return s.sender.SendWechat(ctx, req.OpenID, templateID, req.ExtraData)

	default:
		return errors.Errorf("不支持的通知渠道: %s", req.Channel)
	}
}

// SendMultiChannelNotification 多渠道发送通知
// 尝试使用多个渠道发送通知，任一成功即返回
func (s *NotificationDomainService) SendMultiChannelNotification(ctx context.Context, channels []NotificationChannel, req NotificationRequest) error {
	var lastErr error

	for _, channel := range channels {
		req.Channel = channel
		if err := s.SendNotification(ctx, req); err != nil {
			lastErr = err
			continue
		}
		// 发送成功，返回
		return nil
	}

	if lastErr != nil {
		return errors.Wrap(lastErr, "所有渠道发送均失败")
	}

	return errors.New("没有指定通知渠道")
}

// ShouldNotify 判断是否应该发送通知
// 实现通知频率限制、用户偏好等领域规则
func (s *NotificationDomainService) ShouldNotify(ctx context.Context, userID int64, notificationType string) (bool, error) {
	// 这里可以实现更复杂的逻辑：
	// 1. 检查用户的通知偏好设置
	// 2. 检查通知频率限制
	// 3. 检查用户是否在免打扰时段
	// 4. 检查通知类型的重要程度

	// 示例：简单返回 true
	return true, nil
}
