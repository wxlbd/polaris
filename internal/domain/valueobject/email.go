package valueobject

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// Email 邮箱值对象
// 值对象的特点：
// 1. 不可变性 - 一旦创建就不能修改
// 2. 通过值进行比较 - 两个具有相同值的值对象是相等的
// 3. 自验证 - 创建时验证值的有效性
type Email struct {
	value string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail 创建一个新的 Email 值对象
// 创建时会验证邮箱格式的有效性
func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return Email{}, errors.New("邮箱地址不能为空")
	}

	if !emailRegex.MatchString(email) {
		return Email{}, errors.New("邮箱格式不正确")
	}

	return Email{value: strings.ToLower(email)}, nil
}

// Value 获取邮箱值
func (e Email) Value() string {
	return e.value
}

// String 实现 Stringer 接口
func (e Email) String() string {
	return e.value
}

// Equals 比较两个 Email 是否相等
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// IsEmpty 检查邮箱是否为空
func (e Email) IsEmpty() bool {
	return e.value == ""
}

// Domain 获取邮箱域名部分
func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// LocalPart 获取邮箱本地部分（@符号前的部分）
func (e Email) LocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}
