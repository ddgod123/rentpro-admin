package rental

// District 区域结构
type District struct {
	ID            uint           `json:"id" gorm:"primaryKey;autoIncrement" comment:"主键ID"`
	Code          string         `json:"code" gorm:"size:20;not null;uniqueIndex" comment:"区域编码"`
	Name          string         `json:"name" gorm:"size:50;not null" comment:"区域名称"`
	CityCode      string         `json:"city_code" gorm:"size:20;not null;index" comment:"城市编码"`
	Sort          int            `json:"sort" gorm:"default:0" comment:"排序"`
	Status        string         `json:"status" gorm:"size:20;default:'active'" comment:"状态"`
	BusinessAreas []BusinessArea `json:"business_areas,omitempty" gorm:"foreignKey:DistrictID" comment:"商圈列表"`
}

// BusinessArea 商圈结构
type BusinessArea struct {
	ID         uint     `json:"id" gorm:"primaryKey;autoIncrement" comment:"主键ID"`
	Code       string   `json:"code" gorm:"size:20;not null;uniqueIndex" comment:"商圈编码"`
	Name       string   `json:"name" gorm:"size:100;not null" comment:"商圈名称"`
	DistrictID uint     `json:"district_id" gorm:"not null;index" comment:"所属区域ID"`
	District   District `json:"district,omitempty" gorm:"foreignKey:DistrictID" comment:"所属区域"`
	CityCode   string   `json:"city_code" gorm:"size:20;not null;index" comment:"城市编码"`
	Sort       int      `json:"sort" gorm:"default:0" comment:"排序"`
	Status     string   `json:"status" gorm:"size:20;default:'active'" comment:"状态"`
}

// TableName 设置区域表名
func (District) TableName() string {
	return "sys_districts"
}

// TableName 设置商圈表名
func (BusinessArea) TableName() string {
	return "sys_business_areas"
}
