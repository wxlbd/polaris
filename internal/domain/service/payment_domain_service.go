package service

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/wxlbd/polaris/internal/domain/valueobject"
)

// PaymentStatus 支付状态
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"   // 待支付
	PaymentStatusPaid      PaymentStatus = "paid"      // 已支付
	PaymentStatusFailed    PaymentStatus = "failed"    // 支付失败
	PaymentStatusRefunded  PaymentStatus = "refunded"  // 已退款
	PaymentStatusCancelled PaymentStatus = "cancelled" // 已取消
)

// PaymentMethod 支付方式
type PaymentMethod string

const (
	PaymentMethodWechat  PaymentMethod = "wechat"  // 微信支付
	PaymentMethodAlipay  PaymentMethod = "alipay"  // 支付宝
	PaymentMethodBalance PaymentMethod = "balance" // 余额支付
)

// PaymentGateway 支付网关接口
// 定义与外部支付服务的交互抽象
type PaymentGateway interface {
	// CreatePayment 创建支付订单
	CreatePayment(ctx context.Context, orderID string, amount valueobject.Money, method PaymentMethod) (string, error)
	// QueryPayment 查询支付状态
	QueryPayment(ctx context.Context, transactionID string) (PaymentStatus, error)
	// RefundPayment 退款
	RefundPayment(ctx context.Context, transactionID string, amount valueobject.Money, reason string) error
}

// PaymentDomainService 支付领域服务
// 处理与支付相关的复杂领域逻辑
type PaymentDomainService struct {
	gateway PaymentGateway
}

// NewPaymentDomainService 创建支付领域服务
func NewPaymentDomainService(gateway PaymentGateway) *PaymentDomainService {
	return &PaymentDomainService{
		gateway: gateway,
	}
}

// PaymentRequest 支付请求
type PaymentRequest struct {
	OrderID     string            // 订单ID
	UserID      int64             // 用户ID
	Amount      valueobject.Money // 支付金额
	Method      PaymentMethod     // 支付方式
	Description string            // 支付描述
	ExpireTime  time.Duration     // 过期时间
}

// PaymentResult 支付结果
type PaymentResult struct {
	TransactionID string        // 交易ID
	PayURL        string        // 支付链接或二维码URL
	Status        PaymentStatus // 支付状态
	ExpireAt      time.Time     // 过期时间
}

// CreatePayment 创建支付
func (s *PaymentDomainService) CreatePayment(ctx context.Context, req PaymentRequest) (*PaymentResult, error) {
	// 验证支付金额
	if req.Amount.IsZero() {
		return nil, errors.New("支付金额不能为零")
	}

	// 验证订单ID
	if req.OrderID == "" {
		return nil, errors.New("订单ID不能为空")
	}

	// 验证支付方式
	if !s.isValidPaymentMethod(req.Method) {
		return nil, errors.Errorf("不支持的支付方式: %s", req.Method)
	}

	// 调用支付网关创建支付
	payURL, err := s.gateway.CreatePayment(ctx, req.OrderID, req.Amount, req.Method)
	if err != nil {
		return nil, errors.Wrap(err, "创建支付失败")
	}

	// 计算过期时间
	expireAt := time.Now().Add(req.ExpireTime)
	if req.ExpireTime == 0 {
		expireAt = time.Now().Add(15 * time.Minute) // 默认15分钟过期
	}

	return &PaymentResult{
		TransactionID: req.OrderID,
		PayURL:        payURL,
		Status:        PaymentStatusPending,
		ExpireAt:      expireAt,
	}, nil
}

// VerifyPayment 验证支付状态
func (s *PaymentDomainService) VerifyPayment(ctx context.Context, transactionID string) (PaymentStatus, error) {
	if transactionID == "" {
		return "", errors.New("交易ID不能为空")
	}

	status, err := s.gateway.QueryPayment(ctx, transactionID)
	if err != nil {
		return "", errors.Wrap(err, "查询支付状态失败")
	}

	return status, nil
}

// ProcessRefund 处理退款
func (s *PaymentDomainService) ProcessRefund(ctx context.Context, transactionID string, amount valueobject.Money, reason string) error {
	if transactionID == "" {
		return errors.New("交易ID不能为空")
	}

	if amount.IsZero() {
		return errors.New("退款金额不能为零")
	}

	if reason == "" {
		reason = "用户申请退款"
	}

	// 先查询支付状态
	status, err := s.gateway.QueryPayment(ctx, transactionID)
	if err != nil {
		return errors.Wrap(err, "查询支付状态失败")
	}

	// 只有已支付的订单才能退款
	if status != PaymentStatusPaid {
		return errors.Errorf("当前支付状态不支持退款: %s", status)
	}

	// 调用支付网关退款
	if err := s.gateway.RefundPayment(ctx, transactionID, amount, reason); err != nil {
		return errors.Wrap(err, "退款失败")
	}

	return nil
}

// CanRefund 判断是否可以退款
func (s *PaymentDomainService) CanRefund(ctx context.Context, transactionID string, refundAmount valueobject.Money, paidAmount valueobject.Money) (bool, string, error) {
	// 查询支付状态
	status, err := s.gateway.QueryPayment(ctx, transactionID)
	if err != nil {
		return false, "", errors.Wrap(err, "查询支付状态失败")
	}

	// 检查状态
	if status != PaymentStatusPaid {
		return false, "订单未支付或已退款", nil
	}

	// 检查退款金额
	if refundAmount.GreaterThan(paidAmount) {
		return false, "退款金额不能大于支付金额", nil
	}

	return true, "", nil
}

// isValidPaymentMethod 验证支付方式
func (s *PaymentDomainService) isValidPaymentMethod(method PaymentMethod) bool {
	switch method {
	case PaymentMethodWechat, PaymentMethodAlipay, PaymentMethodBalance:
		return true
	default:
		return false
	}
}

// CalculateRefundAmount 计算实际退款金额
// 根据业务规则计算可退款金额（如扣除手续费等）
func (s *PaymentDomainService) CalculateRefundAmount(originalAmount valueobject.Money, usedDays int, totalDays int) (valueobject.Money, error) {
	if usedDays < 0 || totalDays <= 0 {
		return valueobject.Money{}, errors.New("无效的天数参数")
	}

	if usedDays >= totalDays {
		// 已使用完毕，不退款
		return valueobject.NewMoney(0, originalAmount.Currency())
	}

	// 按比例计算退款金额
	remainingDays := totalDays - usedDays
	refundAmount := originalAmount.Amount() * int64(remainingDays) / int64(totalDays)

	return valueobject.NewMoney(refundAmount, originalAmount.Currency())
}
