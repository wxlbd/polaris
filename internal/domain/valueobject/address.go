package valueobject

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Address 地址值对象
// 封装完整的中国地址信息
type Address struct {
	province string // 省/直辖市/自治区
	city     string // 市
	district string // 区/县
	street   string // 街道/详细地址
	zipCode  string // 邮政编码
}

// NewAddress 创建地址值对象
func NewAddress(province, city, district, street, zipCode string) (Address, error) {
	province = strings.TrimSpace(province)
	city = strings.TrimSpace(city)
	district = strings.TrimSpace(district)
	street = strings.TrimSpace(street)
	zipCode = strings.TrimSpace(zipCode)

	if province == "" {
		return Address{}, errors.New("省份不能为空")
	}
	if city == "" {
		return Address{}, errors.New("城市不能为空")
	}

	return Address{
		province: province,
		city:     city,
		district: district,
		street:   street,
		zipCode:  zipCode,
	}, nil
}

// Province 获取省份
func (a Address) Province() string {
	return a.province
}

// City 获取城市
func (a Address) City() string {
	return a.city
}

// District 获取区县
func (a Address) District() string {
	return a.district
}

// Street 获取街道详细地址
func (a Address) Street() string {
	return a.street
}

// ZipCode 获取邮政编码
func (a Address) ZipCode() string {
	return a.zipCode
}

// FullAddress 获取完整地址
func (a Address) FullAddress() string {
	parts := []string{a.province, a.city}
	if a.district != "" {
		parts = append(parts, a.district)
	}
	if a.street != "" {
		parts = append(parts, a.street)
	}
	return strings.Join(parts, "")
}

// ShortAddress 获取简短地址（省市区）
func (a Address) ShortAddress() string {
	parts := []string{a.province, a.city}
	if a.district != "" {
		parts = append(parts, a.district)
	}
	return strings.Join(parts, "")
}

// String 实现 Stringer 接口
func (a Address) String() string {
	addr := a.FullAddress()
	if a.zipCode != "" {
		addr = fmt.Sprintf("%s (%s)", addr, a.zipCode)
	}
	return addr
}

// Equals 比较两个地址是否相等
func (a Address) Equals(other Address) bool {
	return a.province == other.province &&
		a.city == other.city &&
		a.district == other.district &&
		a.street == other.street &&
		a.zipCode == other.zipCode
}

// IsEmpty 检查地址是否为空
func (a Address) IsEmpty() bool {
	return a.province == "" && a.city == ""
}

// IsSameCity 判断是否在同一城市
func (a Address) IsSameCity(other Address) bool {
	return a.province == other.province && a.city == other.city
}

// IsSameProvince 判断是否在同一省份
func (a Address) IsSameProvince(other Address) bool {
	return a.province == other.province
}
