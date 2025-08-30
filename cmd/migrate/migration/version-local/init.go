package version_local

import (
	"rentPro/rentpro-admin/cmd/migrate/migration"
	commonModels "rentPro/rentpro-admin/common/models/base"
	"rentPro/rentpro-admin/common/models/rental"

	"gorm.io/gorm"
)

func init() {
	migration.Migrate.SetVersion("1700000000001", createRentalTables)
}

func createRentalTables(db *gorm.DB, version string) error {
	// 1. 创建楼盘表结构
	err := db.AutoMigrate(&rental.SysBuildings{})
	if err != nil {
		return err
	}

	// 2. 检查是否已有数据，如果没有则插入初始数据
	var count int64
	db.Model(&rental.SysBuildings{}).Count(&count)

	if count == 0 {
		// 插入初始楼盘数据
		initialBuildings := []rental.SysBuildings{
			{
				Name:            "滨江一号",
				Developer:       "滨江地产",
				DetailedAddress: "上海市浦东新区张江高科技园区博云路2号",
				City:            "上海市",
				District:        "浦东新区",
				PropertyType:    "住宅",
				PropertyCompany: "滨江物业",
				Description:     "滨江一号是滨江地产打造的高品质住宅项目，位于浦东新区核心位置，配套设施完善。",
				Status:          "active",
				IsHot:           true,
				CreatedBy:       "admin",
				UpdatedBy:       "admin",
			},
			{
				Name:            "城市之光",
				Developer:       "城市发展",
				DetailedAddress: "北京市朝阳区建国路88号",
				City:            "北京市",
				District:        "朝阳区",
				PropertyType:    "商业",
				PropertyCompany: "城市物业",
				Description:     "城市之光是集购物、餐饮、娱乐于一体的综合性商业中心，地处北京市中心繁华地段。",
				Status:          "active",
				IsHot:           false,
				CreatedBy:       "admin",
				UpdatedBy:       "admin",
			},
			{
				Name:            "科技园大厦",
				Developer:       "科技地产",
				DetailedAddress: "深圳市南山区科技园南区",
				City:            "深圳市",
				District:        "南山区",
				PropertyType:    "办公",
				PropertyCompany: "科技物业",
				Description:     "科技园大厦是专为高科技企业打造的现代化办公大楼，配备高速网络和智能办公系统。",
				Status:          "active",
				IsHot:           false,
				CreatedBy:       "admin",
				UpdatedBy:       "admin",
			},
			{
				Name:            "湖景花园",
				Developer:       "湖景地产",
				DetailedAddress: "杭州市西湖区文三路478号",
				City:            "杭州市",
				District:        "西湖区",
				PropertyType:    "住宅",
				PropertyCompany: "湖景物业",
				Description:     "湖景花园依湖而建，环境优美，空气清新，是理想的居住选择。",
				Status:          "active",
				IsHot:           true,
				CreatedBy:       "admin",
				UpdatedBy:       "admin",
			},
			{
				Name:            "金融中心",
				Developer:       "金融地产",
				DetailedAddress: "广州市天河区珠江新城冼村路",
				City:            "广州市",
				District:        "天河区",
				PropertyType:    "商业",
				PropertyCompany: "金融物业",
				Description:     "金融中心是广州市标志性建筑，吸引了众多金融机构和企业入驻。",
				Status:          "pending",
				IsHot:           false,
				CreatedBy:       "admin",
				UpdatedBy:       "admin",
			},
		}

		// 批量插入初始数据
		for _, building := range initialBuildings {
			if err := db.Create(&building).Error; err != nil {
				return err
			}
		}
	}

	// 记录迁移完成
	return db.Create(&commonModels.Migration{
		Version: version,
		Name:    "创建租赁管理系统基础表结构并初始化楼盘数据",
		Status:  "completed",
	}).Error
}