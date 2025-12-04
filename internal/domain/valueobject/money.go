package valueobject

import (
	"fmt"

	"github.com/pkg/errors"
)

// Currency 货币类型
type Currency string

const (
	CurrencyCNY Currency = "CNY" // 人民币
	CurrencyUSD Currency = "USD" // 美元
	CurrencyEUR Currency = "EUR" // 欧元
)

// Money 金额值对象
// 使用整数分表示金额，避免浮点数精度问题
type Money struct {
	amount   int64    // 金额（单位：分）
	currency Currency // 货币类型
}

// NewMoney 创建金额值对象（单位：分）
func NewMoney(amountInCents int64, currency Currency) (Money, error) {
	if amountInCents < 0 {
		return Money{}, errors.New("金额不能为负数")
	}

	if currency == "" {
		currency = CurrencyCNY
	}

	return Money{
		amount:   amountInCents,
		currency: currency,
	}, nil
}

// NewMoneyFromYuan 从元创建金额值对象（人民币）
func NewMoneyFromYuan(yuan float64) (Money, error) {
	if yuan < 0 {
		return Money{}, errors.New("金额不能为负数")
	}
	// 转换为分（四舍五入）
	amountInCents := int64(yuan*100 + 0.5)
	return NewMoney(amountInCents, CurrencyCNY)
}

// Amount 获取金额（单位：分）
func (m Money) Amount() int64 {
	return m.amount
}

// AmountInYuan 获取金额（单位：元）
func (m Money) AmountInYuan() float64 {
	return float64(m.amount) / 100
}

// Currency 获取货币类型
func (m Money) Currency() Currency {
	return m.currency
}

// Add 加法运算
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("不同货币类型不能相加")
	}
	return Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}, nil
}

// Subtract 减法运算
func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("不同货币类型不能相减")
	}
	if m.amount < other.amount {
		return Money{}, errors.New("金额不足")
	}
	return Money{
		amount:   m.amount - other.amount,
		currency: m.currency,
	}, nil
}

// Multiply 乘法运算
func (m Money) Multiply(factor int64) Money {
	return Money{
		amount:   m.amount * factor,
		currency: m.currency,
	}
}

// IsZero 判断金额是否为零
func (m Money) IsZero() bool {
	return m.amount == 0
}

// GreaterThan 大于比较
func (m Money) GreaterThan(other Money) bool {
	return m.currency == other.currency && m.amount > other.amount
}

// LessThan 小于比较
func (m Money) LessThan(other Money) bool {
	return m.currency == other.currency && m.amount < other.amount
}

// Equals 相等比较
func (m Money) Equals(other Money) bool {
	return m.currency == other.currency && m.amount == other.amount
}

// String 格式化输出
func (m Money) String() string {
	switch m.currency {
	case CurrencyCNY:
		return fmt.Sprintf("¥%.2f", m.AmountInYuan())
	case CurrencyUSD:
		return fmt.Sprintf("$%.2f", m.AmountInYuan())
	case CurrencyEUR:
		return fmt.Sprintf("€%.2f", m.AmountInYuan())
	default:
		return fmt.Sprintf("%.2f %s", m.AmountInYuan(), m.currency)
	}
}
