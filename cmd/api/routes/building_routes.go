package routes

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"rentPro/rentpro-admin/common/database"
	"rentPro/rentpro-admin/common/utils"

	"github.com/gin-gonic/gin"
)

// SetupBuildingRoutes 设置楼盘管理相关路由
func SetupBuildingRoutes(api *gin.RouterGroup) {
	// 获取楼盘列表
	api.GET("/buildings", func(c *gin.Context) {
		// 获取查询参数
		page := c.DefaultQuery("page", "1")
		pageSize := c.DefaultQuery("pageSize", "10")
		name := c.Query("name")
		city := c.Query("city")
		district := c.Query("district")
		businessArea := c.Query("business_area")
		status := c.Query("status")
		deleted := c.Query("deleted") // 是否查询软删除的数据

		// 转换分页参数
		pageNum, _ := strconv.Atoi(page)
		pageSizeNum, _ := strconv.Atoi(pageSize)
		offset := (pageNum - 1) * pageSizeNum

		// 构建查询条件 - 根据deleted参数决定查询条件，并关联用户表获取用户姓名
		var query string
		if deleted == "true" {
			// 查询软删除的数据，包含deleted_at字段
			query = `SELECT b.id, b.name, b.city, b.district, b.business_area, b.property_type, b.status, b.rent_count, 
					b.created_at, b.updated_at, b.created_by, b.updated_by, b.deleted_at,
					COALESCE(u_updated.nick_name, u_created.nick_name, b.updated_by, b.created_by) as editor_name
					FROM sys_buildings b
					LEFT JOIN sys_user u_created ON b.created_by = u_created.username
					LEFT JOIN sys_user u_updated ON b.updated_by = u_updated.username
					WHERE b.deleted_at IS NOT NULL`
		} else {
			// 查询正常数据
			query = `SELECT b.id, b.name, b.city, b.district, b.business_area, b.property_type, b.status, b.rent_count, 
					b.created_at, b.updated_at, b.created_by, b.updated_by,
					COALESCE(u_updated.nick_name, u_created.nick_name, b.updated_by, b.created_by) as editor_name
					FROM sys_buildings b
					LEFT JOIN sys_user u_created ON b.created_by = u_created.username
					LEFT JOIN sys_user u_updated ON b.updated_by = u_updated.username
					WHERE b.deleted_at IS NULL`
		}
		args := []interface{}{}

		if name != "" {
			query += " AND b.name LIKE ?"
			args = append(args, "%"+name+"%")
		}
		if city != "" {
			query += " AND b.city = ?"
			args = append(args, city)
		}
		if district != "" {
			query += " AND b.district = ?"
			args = append(args, district)
		}
		if businessArea != "" {
			query += " AND b.business_area = ?"
			args = append(args, businessArea)
		}
		if status != "" {
			query += " AND b.status = ?"
			args = append(args, status)
		}

		query += " ORDER BY b.rent_count DESC, b.created_at ASC LIMIT ? OFFSET ?"
		args = append(args, pageSizeNum, offset)

		var buildings []map[string]interface{}
		result := database.DB.Raw(query, args...).Scan(&buildings)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "查询楼盘列表失败",
				"error":   result.Error.Error(),
			})
			return
		}

		// 获取总数 - 根据deleted参数决定查询条件
		var countQuery string
		if deleted == "true" {
			countQuery = "SELECT COUNT(*) FROM sys_buildings b WHERE b.deleted_at IS NOT NULL"
		} else {
			countQuery = "SELECT COUNT(*) FROM sys_buildings b WHERE b.deleted_at IS NULL"
		}
		countArgs := []interface{}{}

		if name != "" {
			countQuery += " AND b.name LIKE ?"
			countArgs = append(countArgs, "%"+name+"%")
		}
		if city != "" {
			countQuery += " AND b.city = ?"
			countArgs = append(countArgs, city)
		}
		if district != "" {
			countQuery += " AND b.district = ?"
			countArgs = append(countArgs, district)
		}
		if businessArea != "" {
			countQuery += " AND b.business_area = ?"
			countArgs = append(countArgs, businessArea)
		}
		if status != "" {
			countQuery += " AND b.status = ?"
			countArgs = append(countArgs, status)
		}

		var total int64
		database.DB.Raw(countQuery, countArgs...).Scan(&total)

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取楼盘列表成功",
			"data":    buildings,
			"page":    pageNum,
			"size":    pageSizeNum,
			"total":   total,
		})
	})

	// 获取单个楼盘信息
	api.GET("/buildings/:id", func(c *gin.Context) {
		id := c.Param("id")

		var building map[string]interface{}
		result := database.DB.Table("sys_buildings").Where("id = ? AND deleted_at IS NULL", id).First(&building)

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "楼盘不存在",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取楼盘信息成功",
			"data":    building,
		})
	})

	// 创建楼盘
	api.POST("/buildings", func(c *gin.Context) {

		// 解析请求体
		var buildingData struct {
			Name         string `json:"name" binding:"required"`
			City         string `json:"city" binding:"required"`
			District     string `json:"district" binding:"required"`
			BusinessArea string `json:"businessArea"`
			PropertyType string `json:"propertyType"`
			Description  string `json:"description"`
			Status       string `json:"status"`
		}

		if err := c.ShouldBindJSON(&buildingData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"error":   err.Error(),
			})
			return
		}

		// 设置默认值
		if buildingData.Status == "" {
			buildingData.Status = "active"
		}

		// 获取当前用户 (暂时使用admin)
		currentUser := "admin" // TODO: 从JWT token或上下文中获取真实用户

		// 插入数据库
		result := database.DB.Exec(
			"INSERT INTO sys_buildings (name, city, district, business_area, property_type, description, status, created_by, updated_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())",
			buildingData.Name,
			buildingData.City,
			buildingData.District,
			buildingData.BusinessArea,
			buildingData.PropertyType,
			buildingData.Description,
			buildingData.Status,
			currentUser,
			currentUser,
		)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "创建楼盘失败",
				"error":   result.Error.Error(),
			})
			return
		}

		// 获取新创建的楼盘ID
		var newBuildingID int64
		database.DB.Raw("SELECT LAST_INSERT_ID()").Scan(&newBuildingID)

		// 初始化楼盘文件夹结构
		imageManager := utils.GetImageManager()
		if imageManager != nil {
			if err := imageManager.CreateBuildingFolder(uint64(newBuildingID), buildingData.Name); err != nil {
				// 文件夹创建失败不影响楼盘创建成功，只记录日志
				// 这里可以记录到日志系统
			}
		}

		c.JSON(http.StatusCreated, gin.H{
			"code":    201,
			"message": "楼盘创建成功",
			"data": gin.H{
				"id":           newBuildingID,
				"name":         buildingData.Name,
				"city":         buildingData.City,
				"district":     buildingData.District,
				"businessArea": buildingData.BusinessArea,
				"propertyType": buildingData.PropertyType,
				"description":  buildingData.Description,
				"status":       buildingData.Status,
			},
		})
	})

	// 更新楼盘信息
	api.PUT("/buildings/:id", func(c *gin.Context) {
		id := c.Param("id")

		// 解析请求体
		var buildingData struct {
			Name         string `json:"name"`
			City         string `json:"city"`
			District     string `json:"district"`
			BusinessArea string `json:"businessArea"`
			PropertyType string `json:"propertyType"`
			Description  string `json:"description"`
			Status       string `json:"status"`
		}

		if err := c.ShouldBindJSON(&buildingData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"error":   err.Error(),
			})
			return
		}

		// 构建SQL更新语句
		setParts := []string{}
		values := []interface{}{}

		if buildingData.Name != "" {
			setParts = append(setParts, "name = ?")
			values = append(values, buildingData.Name)
		}
		if buildingData.City != "" {
			setParts = append(setParts, "city = ?")
			values = append(values, buildingData.City)
		}
		if buildingData.District != "" {
			setParts = append(setParts, "district = ?")
			values = append(values, buildingData.District)
		}
		if buildingData.BusinessArea != "" {
			setParts = append(setParts, "business_area = ?")
			values = append(values, buildingData.BusinessArea)
		}
		if buildingData.PropertyType != "" {
			setParts = append(setParts, "property_type = ?")
			values = append(values, buildingData.PropertyType)
		}
		if buildingData.Description != "" {
			setParts = append(setParts, "description = ?")
			values = append(values, buildingData.Description)
		}
		if buildingData.Status != "" {
			setParts = append(setParts, "status = ?")
			values = append(values, buildingData.Status)
		}

		// 获取当前用户 (暂时使用admin)
		currentUser := "admin" // TODO: 从JWT token或上下文中获取真实用户

		// 总是更新 updated_at 和 updated_by
		setParts = append(setParts, "updated_at = ?")
		values = append(values, time.Now())
		setParts = append(setParts, "updated_by = ?")
		values = append(values, currentUser)
		// 注意：id 参数放在最后
		values = append(values, id)

		// 执行原生SQL更新
		sql := "UPDATE sys_buildings SET " + strings.Join(setParts, ", ") + " WHERE id = ?"
		result := database.DB.Exec(sql, values...)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "更新楼盘失败",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "更新楼盘成功",
		})
	})

	// 删除楼盘
	api.DELETE("/buildings/:id", func(c *gin.Context) {
		id := c.Param("id")

		// 检查楼盘是否存在
		var buildingExists int64
		database.DB.Raw("SELECT COUNT(*) FROM sys_buildings WHERE id = ? AND deleted_at IS NULL", id).Scan(&buildingExists)
		if buildingExists == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "楼盘不存在",
			})
			return
		}

		// 检查是否有关联的户型数据
		var houseTypeCount int64
		database.DB.Raw("SELECT COUNT(*) FROM sys_house_types WHERE building_id = ? AND deleted_at IS NULL", id).Scan(&houseTypeCount)

		if houseTypeCount > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "该楼盘下还有户型数据，无法删除",
				"data": gin.H{
					"house_type_count": houseTypeCount,
				},
			})
			return
		}

		// 删除数据库记录（软删除）
		result := database.DB.Exec("UPDATE sys_buildings SET deleted_at = NOW() WHERE id = ?", id)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "删除楼盘失败",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "删除楼盘成功",
		})
	})

	// 获取楼盘详细信息（包含户型统计等）
	api.GET("/buildings/:id/info", func(c *gin.Context) {
		id := c.Param("id")

		var building map[string]interface{}
		result := database.DB.Raw(`
			SELECT
				b.*,
				COUNT(ht.id) as house_type_count,
				COALESCE(SUM(ht.total_stock), 0) as total_stock,
				COALESCE(SUM(ht.available_stock), 0) as available_stock
			FROM sys_buildings b
			LEFT JOIN sys_house_types ht ON b.id = ht.building_id AND ht.deleted_at IS NULL
			WHERE b.id = ? AND b.deleted_at IS NULL
			GROUP BY b.id
		`, id).Scan(&building)

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "楼盘不存在",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取楼盘详细信息成功",
			"data":    building,
		})
	})

	// 获取楼盘图片列表
	api.GET("/buildings/images/:buildingId", func(c *gin.Context) {
		buildingIDStr := c.Param("buildingId")
		category := c.Query("category")

		buildingID, err := strconv.ParseUint(buildingIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "楼盘ID格式错误",
			})
			return
		}

		imageManager := utils.GetImageManager()
		if imageManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "图片管理器未初始化",
			})
			return
		}

		images, err := imageManager.GetBuildingImages(buildingID, category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取楼盘图片失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取楼盘图片成功",
			"data":    images,
		})
	})

	// 获取楼盘户型图列表
	api.GET("/buildings/floor-plans/:buildingId", func(c *gin.Context) {
		buildingIDStr := c.Param("buildingId")

		buildingID, err := strconv.ParseUint(buildingIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "楼盘ID格式错误",
			})
			return
		}

		imageManager := utils.GetImageManager()
		if imageManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "图片管理器未初始化",
			})
			return
		}

		images, err := imageManager.GetBuildingFloorPlans(buildingID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取楼盘户型图失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取楼盘户型图成功",
			"data":    images,
		})
	})

	// 注意：城市列表API已移动到 city_routes.go 中的 /api/v1/cities

	// 获取区域列表（支持按城市筛选）
	api.GET("/districts", func(c *gin.Context) {
		cityId := c.Query("cityId")

		query := "SELECT id, code, name, city_code, city_id, sort, status FROM sys_districts WHERE status = 'active'"
		args := []interface{}{}

		if cityId != "" {
			query += " AND city_id = ?"
			args = append(args, cityId)
		}

		query += " ORDER BY sort ASC"

		var districts []map[string]interface{}
		result := database.DB.Raw(query, args...).Scan(&districts)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取区域列表失败",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取区域列表成功",
			"data":    districts,
		})
	})

	// 获取商圈列表
	api.GET("/business-areas", func(c *gin.Context) {
		districtId := c.Query("districtId")

		query := "SELECT id, code, name, district_id, city_code, sort, status FROM sys_business_areas WHERE status = 'active'"
		args := []interface{}{}

		if districtId != "" {
			query += " AND district_id = ?"
			args = append(args, districtId)
		}

		query += " ORDER BY sort ASC"

		var businessAreas []map[string]interface{}
		result := database.DB.Raw(query, args...).Scan(&businessAreas)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取商圈列表失败",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取商圈列表成功",
			"data":    businessAreas,
		})
	})

	// 恢复楼盘（软删除恢复）
	api.POST("/buildings/:id/restore", func(c *gin.Context) {
		id := c.Param("id")

		// 验证ID
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "楼盘ID不能为空",
			})
			return
		}

		// 恢复楼盘（将deleted_at设置为NULL）
		result := database.DB.Exec("UPDATE sys_buildings SET deleted_at = NULL, updated_at = NOW() WHERE id = ? AND deleted_at IS NOT NULL", id)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "恢复楼盘失败",
				"error":   result.Error.Error(),
			})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "未找到可恢复的楼盘",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "恢复楼盘成功",
		})
	})

	// 永久删除楼盘
	api.DELETE("/buildings/:id/permanent", func(c *gin.Context) {
		id := c.Param("id")

		// 验证ID
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "楼盘ID不能为空",
			})
			return
		}

		// 检查楼盘是否已被软删除
		var count int64
		database.DB.Raw("SELECT COUNT(*) FROM sys_buildings WHERE id = ? AND deleted_at IS NOT NULL", id).Scan(&count)

		if count == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "只能永久删除已在回收站的楼盘",
			})
			return
		}

		// 永久删除楼盘（物理删除）
		result := database.DB.Exec("DELETE FROM sys_buildings WHERE id = ? AND deleted_at IS NOT NULL", id)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "永久删除楼盘失败",
				"error":   result.Error.Error(),
			})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "未找到可删除的楼盘",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "永久删除楼盘成功",
		})
	})
}
