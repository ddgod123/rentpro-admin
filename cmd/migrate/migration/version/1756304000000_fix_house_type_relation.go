package version

import (
	"fmt"
	"log"
	"rentPro/rentpro-admin/common/database"
	"rentPro/rentpro-admin/common/models/rental"
)

func init() {
	// 注册迁移
	migrationList = append(migrationList, Migration{
		Version:     "1756304000000",
		Description: "修复房屋和户型关联关系",
		Up:          up1756304000000,
		Down:        down1756304000000,
	})
}

// up1756304000000 执行数据库结构变更
func up1756304000000() error {
	log.Println("开始执行数据库迁移: 修复房屋和户型关联关系")

	db := database.DB
	if db == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	// 1. 为 sys_houses 表添加 house_type_id 字段
	log.Println("1. 为 sys_houses 表添加 house_type_id 字段...")
	if err := db.Exec(`
		ALTER TABLE sys_houses 
		ADD COLUMN house_type_id INT UNSIGNED AFTER building_id,
		ADD INDEX idx_house_type_id (house_type_id)
	`).Error; err != nil {
		// 如果字段已存在，忽略错误
		log.Printf("添加 house_type_id 字段时出现错误（可能已存在）: %v", err)
	}

	// 2. 重命名 sys_house_types 表中的字段
	log.Println("2. 重命名 sys_house_types 表字段...")
	renameColumns := []string{
		"ALTER TABLE sys_house_types CHANGE COLUMN area standard_area DECIMAL(8,2) NOT NULL COMMENT '标准建筑面积(平方米)'",
		"ALTER TABLE sys_house_types CHANGE COLUMN orientation standard_orientation VARCHAR(50) COMMENT '标准朝向'",
		"ALTER TABLE sys_house_types CHANGE COLUMN view standard_view VARCHAR(100) COMMENT '标准景观'",
		"ALTER TABLE sys_house_types CHANGE COLUMN sale_price base_sale_price DECIMAL(12,2) DEFAULT 0 COMMENT '基准售价(元)'",
		"ALTER TABLE sys_house_types CHANGE COLUMN rent_price base_rent_price DECIMAL(8,2) DEFAULT 0 COMMENT '基准月租金(元)'",
		"ALTER TABLE sys_house_types CHANGE COLUMN sale_price_per base_sale_price_per DECIMAL(8,2) DEFAULT 0 COMMENT '基准单价(元/平方米)'",
		"ALTER TABLE sys_house_types CHANGE COLUMN rent_price_per base_rent_price_per DECIMAL(6,2) DEFAULT 0 COMMENT '基准租金单价(元/平方米/月)'",
		"ALTER TABLE sys_house_types CHANGE COLUMN total_stock total_stock INT DEFAULT 0 COMMENT '总库存'",
		"ALTER TABLE sys_house_types CHANGE COLUMN sale_stock sold_stock INT DEFAULT 0 COMMENT '已售库存'",
		"ALTER TABLE sys_house_types CHANGE COLUMN rent_stock rented_stock INT DEFAULT 0 COMMENT '已租库存'",
		"ALTER TABLE sys_house_types CHANGE COLUMN reserved_stock reserved_stock INT DEFAULT 0 COMMENT '已预订库存'",
	}

	for _, sql := range renameColumns {
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("重命名字段时出现错误: %v, SQL: %s", err, sql)
		}
	}

	// 3. 为 sys_house_types 添加新字段
	log.Println("3. 为 sys_house_types 添加新字段...")
	newColumns := []string{
		"ALTER TABLE sys_house_types ADD COLUMN available_stock INT DEFAULT 0 COMMENT '可用库存' AFTER total_stock",
		"ALTER TABLE sys_house_types ADD COLUMN floor_plan_url VARCHAR(500) COMMENT '户型图URL' AFTER main_image",
		"ALTER TABLE sys_house_types ADD COLUMN floor_height DECIMAL(4,2) COMMENT '标准层高(米)' AFTER balconies",
	}

	for _, sql := range newColumns {
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("添加新字段时出现错误: %v, SQL: %s", err, sql)
		}
	}

	// 4. 重命名 sys_houses 表中的字段
	log.Println("4. 重命名 sys_houses 表字段...")
	houseRenameColumns := []string{
		"ALTER TABLE sys_houses CHANGE COLUMN area actual_area DECIMAL(8,2) COMMENT '实际建筑面积(平方米)'",
		"ALTER TABLE sys_houses CHANGE COLUMN usableArea actual_usable_area DECIMAL(8,2) COMMENT '实际使用面积(平方米)'",
		"ALTER TABLE sys_houses CHANGE COLUMN orientation actual_orientation VARCHAR(50) COMMENT '实际朝向'",
		"ALTER TABLE sys_houses CHANGE COLUMN view actual_view VARCHAR(100) COMMENT '实际景观'",
		"ALTER TABLE sys_houses CHANGE COLUMN salePrice actual_sale_price DECIMAL(12,2) DEFAULT 0 COMMENT '实际售价(元)'",
		"ALTER TABLE sys_houses CHANGE COLUMN rentPrice actual_rent_price DECIMAL(8,2) DEFAULT 0 COMMENT '实际月租金(元)'",
	}

	for _, sql := range houseRenameColumns {
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("重命名房屋字段时出现错误: %v, SQL: %s", err, sql)
		}
	}

	// 5. 为 sys_houses 添加价格调整字段
	log.Println("5. 为 sys_houses 添加价格调整字段...")
	housePriceColumns := []string{
		"ALTER TABLE sys_houses ADD COLUMN price_adjustment DECIMAL(8,2) DEFAULT 0 COMMENT '价格调整金额(元)' AFTER actual_rent_price",
		"ALTER TABLE sys_houses ADD COLUMN price_adjustment_reason VARCHAR(200) COMMENT '价格调整原因' AFTER price_adjustment",
	}

	for _, sql := range housePriceColumns {
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("添加价格调整字段时出现错误: %v, SQL: %s", err, sql)
		}
	}

	// 6. 删除 sys_houses 表中不再需要的字段
	log.Println("6. 删除 sys_houses 表中冗余字段...")
	dropColumns := []string{
		"ALTER TABLE sys_houses DROP COLUMN IF EXISTS rooms",
		"ALTER TABLE sys_houses DROP COLUMN IF EXISTS halls",
		"ALTER TABLE sys_houses DROP COLUMN IF EXISTS bathrooms",
		"ALTER TABLE sys_houses DROP COLUMN IF EXISTS balconies",
		"ALTER TABLE sys_houses DROP COLUMN IF EXISTS salePricePer",
		"ALTER TABLE sys_houses DROP COLUMN IF EXISTS rentPricePer",
	}

	for _, sql := range dropColumns {
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("删除冗余字段时出现错误: %v, SQL: %s", err, sql)
		}
	}

	// 7. 创建示例户型数据（如果不存在）
	log.Println("7. 创建示例户型数据...")
	if err := createSampleHouseTypes(); err != nil {
		log.Printf("创建示例户型数据时出现错误: %v", err)
	}

	// 8. 尝试建立现有房屋与户型的关联
	log.Println("8. 建立现有房屋与户型的关联...")
	if err := linkExistingHousesToTypes(); err != nil {
		log.Printf("建立房屋户型关联时出现错误: %v", err)
	}

	// 9. 添加外键约束（可选，如果数据完整性允许）
	log.Println("9. 添加外键约束...")
	if err := db.Exec(`
		ALTER TABLE sys_houses 
		ADD CONSTRAINT fk_houses_house_type_id 
		FOREIGN KEY (house_type_id) REFERENCES sys_house_types(id)
		ON DELETE SET NULL ON UPDATE CASCADE
	`).Error; err != nil {
		log.Printf("添加外键约束时出现错误: %v", err)
	}

	log.Println("数据库迁移完成!")
	return nil
}

// down1756304000000 回滚数据库变更
func down1756304000000() error {
	log.Println("开始回滚数据库迁移: 修复房屋和户型关联关系")

	db := database.DB
	if db == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	// 回滚操作（与up操作相反）
	log.Println("1. 删除外键约束...")
	db.Exec("ALTER TABLE sys_houses DROP FOREIGN KEY IF EXISTS fk_houses_house_type_id")

	log.Println("2. 删除 house_type_id 字段...")
	db.Exec("ALTER TABLE sys_houses DROP COLUMN IF EXISTS house_type_id")

	log.Println("3. 恢复原始字段名...")
	// 恢复 sys_house_types 字段名
	restoreHouseTypeColumns := []string{
		"ALTER TABLE sys_house_types CHANGE COLUMN standard_area area DECIMAL(8,2) NOT NULL",
		"ALTER TABLE sys_house_types CHANGE COLUMN standard_orientation orientation VARCHAR(50)",
		"ALTER TABLE sys_house_types CHANGE COLUMN standard_view view VARCHAR(100)",
		"ALTER TABLE sys_house_types CHANGE COLUMN base_sale_price sale_price DECIMAL(12,2) DEFAULT 0",
		"ALTER TABLE sys_house_types CHANGE COLUMN base_rent_price rent_price DECIMAL(8,2) DEFAULT 0",
		"ALTER TABLE sys_house_types CHANGE COLUMN base_sale_price_per sale_price_per DECIMAL(8,2) DEFAULT 0",
		"ALTER TABLE sys_house_types CHANGE COLUMN base_rent_price_per rent_price_per DECIMAL(6,2) DEFAULT 0",
		"ALTER TABLE sys_house_types CHANGE COLUMN sold_stock sale_stock INT DEFAULT 0",
		"ALTER TABLE sys_house_types CHANGE COLUMN rented_stock rent_stock INT DEFAULT 0",
	}

	for _, sql := range restoreHouseTypeColumns {
		db.Exec(sql)
	}

	// 恢复 sys_houses 字段名
	restoreHouseColumns := []string{
		"ALTER TABLE sys_houses CHANGE COLUMN actual_area area DECIMAL(8,2)",
		"ALTER TABLE sys_houses CHANGE COLUMN actual_usable_area usableArea DECIMAL(8,2)",
		"ALTER TABLE sys_houses CHANGE COLUMN actual_orientation orientation VARCHAR(50)",
		"ALTER TABLE sys_houses CHANGE COLUMN actual_view view VARCHAR(100)",
		"ALTER TABLE sys_houses CHANGE COLUMN actual_sale_price salePrice DECIMAL(12,2) DEFAULT 0",
		"ALTER TABLE sys_houses CHANGE COLUMN actual_rent_price rentPrice DECIMAL(8,2) DEFAULT 0",
	}

	for _, sql := range restoreHouseColumns {
		db.Exec(sql)
	}

	log.Println("4. 删除新增字段...")
	db.Exec("ALTER TABLE sys_house_types DROP COLUMN IF EXISTS available_stock")
	db.Exec("ALTER TABLE sys_house_types DROP COLUMN IF EXISTS floor_plan_url")
	db.Exec("ALTER TABLE sys_house_types DROP COLUMN IF EXISTS floor_height")
	db.Exec("ALTER TABLE sys_houses DROP COLUMN IF EXISTS price_adjustment")
	db.Exec("ALTER TABLE sys_houses DROP COLUMN IF EXISTS price_adjustment_reason")

	log.Println("数据库迁移回滚完成!")
	return nil
}

// createSampleHouseTypes 创建示例户型数据
func createSampleHouseTypes() error {
	db := database.DB

	// 检查是否已有户型数据
	var count int64
	if err := db.Model(&rental.SysHouseType{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("已存在户型数据，跳过创建示例数据")
		return nil
	}

	// 获取第一个楼盘ID作为示例
	var building rental.SysBuildings
	if err := db.First(&building).Error; err != nil {
		log.Println("未找到楼盘数据，跳过创建示例户型")
		return nil
	}

	// 创建示例户型
	sampleTypes := []rental.SysHouseType{
		{
			Name:                "经典一居",
			Code:                "A1",
			Description:         "精致一居室，适合单身人士",
			BuildingID:          building.ID,
			StandardArea:        45.5,
			Rooms:               1,
			Halls:               1,
			Bathrooms:           1,
			Balconies:           1,
			FloorHeight:         2.8,
			StandardOrientation: "南向",
			StandardView:        "小区景观",
			BaseSalePrice:       800000,
			BaseRentPrice:       3500,
			Status:              "active",
		},
		{
			Name:                "舒适两居",
			Code:                "B2",
			Description:         "温馨两居室，适合小家庭",
			BuildingID:          building.ID,
			StandardArea:        78.5,
			Rooms:               2,
			Halls:               1,
			Bathrooms:           1,
			Balconies:           1,
			FloorHeight:         2.8,
			StandardOrientation: "南北",
			StandardView:        "小区景观",
			BaseSalePrice:       1200000,
			BaseRentPrice:       5500,
			Status:              "active",
		},
		{
			Name:                "宽敞三居",
			Code:                "C3",
			Description:         "宽敞三居室，适合大家庭",
			BuildingID:          building.ID,
			StandardArea:        108.0,
			Rooms:               3,
			Halls:               2,
			Bathrooms:           2,
			Balconies:           2,
			FloorHeight:         2.8,
			StandardOrientation: "南北",
			StandardView:        "城市景观",
			BaseSalePrice:       1800000,
			BaseRentPrice:       8000,
			Status:              "active",
		},
	}

	for _, houseType := range sampleTypes {
		houseType.CalculateBasePricePer()
		if err := db.Create(&houseType).Error; err != nil {
			log.Printf("创建示例户型失败: %v", err)
		}
	}

	log.Println("示例户型数据创建完成")
	return nil
}

// linkExistingHousesToTypes 建立现有房屋与户型的关联
func linkExistingHousesToTypes() error {
	db := database.DB

	// 获取所有没有关联户型的房屋
	var houses []rental.SysHouse
	if err := db.Where("house_type_id IS NULL OR house_type_id = 0").Find(&houses).Error; err != nil {
		return err
	}

	log.Printf("找到 %d 个未关联户型的房屋", len(houses))

	for _, house := range houses {
		// 根据房屋属性查找最匹配的户型
		var houseType rental.SysHouseType
		err := db.Where("building_id = ?", house.BuildingID).
			Where("ABS(standard_area - ?) < 10", house.ActualArea). // 面积差异小于10平米
			First(&houseType).Error

		if err != nil {
			log.Printf("房屋 %s 未找到匹配的户型", house.Code)
			continue
		}

		// 更新房屋的户型关联
		if err := db.Model(&house).Update("house_type_id", houseType.ID).Error; err != nil {
			log.Printf("更新房屋 %s 的户型关联失败: %v", house.Code, err)
		} else {
			log.Printf("房屋 %s 已关联到户型 %s", house.Code, houseType.Code)
		}
	}

	return nil
}
