package valueobject

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// Phone 手机号值对象
// 封装手机号相关的验证和业务逻辑
type Phone struct {
	countryCode string // 国家区号，如 "+86"
	number      string // 手机号码
}

// 中国大陆手机号正则
var chinaPhoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

// NewPhone 创建一个新的 Phone 值对象
// 默认使用中国大陆区号 +86
func NewPhone(number string) (Phone, error) {
	return NewPhoneWithCountryCode("+86", number)
}

// NewPhoneWithCountryCode 使用指定区号创建 Phone 值对象
func NewPhoneWithCountryCode(countryCode, number string) (Phone, error) {
	countryCode = strings.TrimSpace(countryCode)
	number = strings.TrimSpace(number)

	if number == "" {
		return Phone{}, errors.New("手机号码不能为空")
	}

	// 移除所有空格和横线
	number = strings.ReplaceAll(number, " ", "")
	number = strings.ReplaceAll(number, "-", "")

	// 默认区号
	if countryCode == "" {
		countryCode = "+86"
	}

	// 确保区号以 + 开头
	if !strings.HasPrefix(countryCode, "+") {
		countryCode = "+" + countryCode
	}

	// 验证中国大陆手机号格式
	if countryCode == "+86" && !chinaPhoneRegex.MatchString(number) {
		return Phone{}, errors.New("手机号格式不正确")
	}

	return Phone{
		countryCode: countryCode,
		number:      number,
	}, nil
}

// CountryCode 获取国家区号
func (p Phone) CountryCode() string {
	return p.countryCode
}

// Number 获取手机号码
func (p Phone) Number() string {
	return p.number
}

// FullNumber 获取完整的手机号（含区号）
func (p Phone) FullNumber() string {
	return p.countryCode + p.number
}

// String 实现 Stringer 接口
func (p Phone) String() string {
	return p.FullNumber()
}

// Equals 比较两个 Phone 是否相等
func (p Phone) Equals(other Phone) bool {
	return p.countryCode == other.countryCode && p.number == other.number
}

// IsEmpty 检查手机号是否为空
func (p Phone) IsEmpty() bool {
	return p.number == ""
}

// Masked 获取脱敏后的手机号（如：138****8888）
func (p Phone) Masked() string {
	if len(p.number) < 7 {
		return p.number
	}
	return p.number[:3] + "****" + p.number[len(p.number)-4:]
}
