package rental

// SysDistrict 区域数据模型
type SysDistrict struct {
	ID       uint64 `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	Code     string `json:"code" gorm:"type:varchar(20);uniqueIndex;not null;comment:区域代码"`
	Name     string `json:"name" gorm:"type:varchar(50);not null;comment:区域名称"`
	CityCode string `json:"city_code" gorm:"type:varchar(20);not null;index:idx_city_code;comment:城市代码"`
	CityID   uint64 `json:"city_id" gorm:"index:idx_city_id;comment:城市ID"`
	Sort     int64  `json:"sort" gorm:"default:0;comment:排序"`
	Status   string `json:"status" gorm:"type:varchar(20);default:active;comment:状态"`

	// 关联关系
	City          SysCity           `json:"city,omitempty" gorm:"foreignKey:CityID;references:ID"`
	BusinessAreas []SysBusinessArea `json:"business_areas,omitempty" gorm:"foreignKey:DistrictID;references:ID"`
}

// TableName 指定表名
func (SysDistrict) TableName() string {
	return "sys_districts"
}

// DistrictOption 区域选项结构（用于前端下拉选择）
type DistrictOption struct {
	ID     uint64 `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	CityID uint64 `json:"city_id"`
}

// ToDistrictOption 转换为区域选项
func (d *SysDistrict) ToDistrictOption() DistrictOption {
	return DistrictOption{
		ID:     d.ID,
		Code:   d.Code,
		Name:   d.Name,
		CityID: d.CityID,
	}
}
