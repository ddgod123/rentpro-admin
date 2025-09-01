package base

import (
	"fmt"
	"rentPro/rentpro-admin/common/models/system"

	"gorm.io/gorm"
)

// InitAuthTables 初始化权限相关表结构
func InitAuthTables(db *gorm.DB) error {
	// 按照依赖关系顺序创建表
	models := []interface{}{
		&system.SysDept{},
		&system.SysPost{},
		&system.SysRole{},
		&system.SysMenu{},
		&system.SysUser{},
	}

	// 自动创建表结构
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}

	return nil
}

// InitDefaultData 初始化默认数据（从SQL文件读取）
func InitDefaultData(db *gorm.DB) error {
	// 创建SQL文件加载器
	sqlLoader := NewSQLFileLoader("config/sql/data")

	// 按照依赖关系顺序加载数据文件
	dataFiles := []string{
		"sys_dept.sql",      // 部门数据
		"sys_post.sql",      // 岗位数据
		"sys_role.sql",      // 角色数据
		"sys_menu.sql",      // 菜单数据
		"sys_user.sql",      // 用户数据
		"sys_role_menu.sql", // 角色菜单关联数据
	}

	// 依次加载并执行SQL文件
	for _, filename := range dataFiles {
		if err := sqlLoader.LoadAndExecuteSQL(db, filename); err != nil {
			return fmt.Errorf("加载数据文件 %s 失败: %v", filename, err)
		}
	}

	return nil
}

// ================================================================
// 以下函数已弃用，改为从SQL文件读取数据，保留作为备份参考
// ================================================================

// initDefaultDepts 初始化默认部门（已弃用，改为从SQL文件读取）
// Deprecated: 请使用SQL文件 config/sql/data/sys_dept.sql 替代
func initDefaultDepts(db *gorm.DB) error {
	depts := []system.SysDept{
		{
			ID:       1,
			ParentID: 0,
			DeptPath: "0,1",
			DeptName: "RentPro科技",
			Sort:     1,
			Leader:   "系统管理员",
			Phone:    "15888888888",
			Email:    "admin@rentpro.com",
			Status:   "0",
		},
		{
			ID:       2,
			ParentID: 1,
			DeptPath: "0,1,2",
			DeptName: "技术部",
			Sort:     1,
			Leader:   "技术总监",
			Phone:    "15666666666",
			Email:    "tech@rentpro.com",
			Status:   "0",
		},
		{
			ID:       3,
			ParentID: 1,
			DeptPath: "0,1,3",
			DeptName: "运营部",
			Sort:     2,
			Leader:   "运营总监",
			Phone:    "15777777777",
			Email:    "ops@rentpro.com",
			Status:   "0",
		},
	}

	for _, dept := range depts {
		var count int64
		db.Model(&system.SysDept{}).Where("id = ?", dept.ID).Count(&count)
		if count == 0 {
			if err := db.Create(&dept).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// initDefaultPosts 初始化默认岗位（已弃用，改为从SQL文件读取）
// Deprecated: 请使用SQL文件 config/sql/data/sys_post.sql 替代
func initDefaultPosts(db *gorm.DB) error {
	posts := []system.SysPost{
		{
			ID:       1,
			PostCode: "ceo",
			PostName: "董事长",
			Sort:     1,
			Status:   "0",
			Remark:   "董事长",
		},
		{
			ID:       2,
			PostCode: "se",
			PostName: "项目经理",
			Sort:     2,
			Status:   "0",
			Remark:   "项目经理",
		},
		{
			ID:       3,
			PostCode: "hr",
			PostName: "人力资源",
			Sort:     3,
			Status:   "0",
			Remark:   "人力资源",
		},
		{
			ID:       4,
			PostCode: "user",
			PostName: "普通员工",
			Sort:     4,
			Status:   "0",
			Remark:   "普通员工",
		},
	}

	for _, post := range posts {
		var count int64
		db.Model(&system.SysPost{}).Where("id = ?", post.ID).Count(&count)
		if count == 0 {
			if err := db.Create(&post).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
