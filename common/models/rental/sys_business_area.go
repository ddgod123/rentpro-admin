package rental

// SysBusinessArea 商圈数据模型
type SysBusinessArea struct {
	ID         uint64 `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	Code       string `json:"code" gorm:"type:varchar(20);uniqueIndex;not null;comment:商圈代码"`
	Name       string `json:"name" gorm:"type:varchar(100);not null;comment:商圈名称"`
	DistrictID uint64 `json:"district_id" gorm:"index:idx_district_id;comment:区域ID"`
	CityCode   string `json:"city_code" gorm:"type:varchar(20);not null;index:idx_city_code;comment:城市代码"`
	Sort       int64  `json:"sort" gorm:"default:0;comment:排序"`
	Status     string `json:"status" gorm:"type:varchar(20);default:active;comment:状态"`

	// 关联关系
	District SysDistrict `json:"district,omitempty" gorm:"foreignKey:DistrictID;references:ID"`
}

// TableName 指定表名
func (SysBusinessArea) TableName() string {
	return "sys_business_areas"
}

// BusinessAreaOption 商圈选项结构（用于前端下拉选择）
type BusinessAreaOption struct {
	ID         uint64 `json:"id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	DistrictID uint64 `json:"district_id"`
}

// ToBusinessAreaOption 转换为商圈选项
func (b *SysBusinessArea) ToBusinessAreaOption() BusinessAreaOption {
	return BusinessAreaOption{
		ID:         b.ID,
		Code:       b.Code,
		Name:       b.Name,
		DistrictID: b.DistrictID,
	}
}
