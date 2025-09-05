package rental

import (
	"time"
)

// SysCity 城市数据模型
type SysCity struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	Code      string    `json:"code" gorm:"type:varchar(10);uniqueIndex;not null;comment:城市代码"`
	Name      string    `json:"name" gorm:"type:varchar(50);not null;comment:城市名称"`
	Sort      int64     `json:"sort" gorm:"default:0;comment:排序"`
	Status    string    `json:"status" gorm:"type:varchar(20);default:active;comment:状态"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:更新时间"`

	// 关联关系
	Districts []SysDistrict `json:"districts,omitempty" gorm:"foreignKey:CityID;references:ID"`
}

// TableName 指定表名
func (SysCity) TableName() string {
	return "sys_cities"
}

// CityOption 城市选项结构（用于前端下拉选择）
type CityOption struct {
	ID   uint64 `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// ToCityOption 转换为城市选项
func (c *SysCity) ToCityOption() CityOption {
	return CityOption{
		ID:   c.ID,
		Code: c.Code,
		Name: c.Name,
	}
}
