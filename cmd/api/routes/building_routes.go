package routes

import (
	"net/http"
	"strconv"

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
		district := c.Query("district")
		businessArea := c.Query("business_area")
		status := c.Query("status")

		// 转换分页参数
		pageNum, _ := strconv.Atoi(page)
		pageSizeNum, _ := strconv.Atoi(pageSize)
		offset := (pageNum - 1) * pageSizeNum

		// 构建查询条件
		query := "SELECT id, name, district, business_area, property_type, status, created_at FROM sys_buildings WHERE 1=1"
		args := []interface{}{}

		if name != "" {
			query += " AND name LIKE ?"
			args = append(args, "%"+name+"%")
		}
		if district != "" {
			query += " AND district = ?"
			args = append(args, district)
		}
		if businessArea != "" {
			query += " AND business_area = ?"
			args = append(args, businessArea)
		}
		if status != "" {
			query += " AND status = ?"
			args = append(args, status)
		}

		query += " ORDER BY id DESC LIMIT ? OFFSET ?"
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

		// 获取总数
		countQuery := "SELECT COUNT(*) FROM sys_buildings WHERE 1=1"
		countArgs := []interface{}{}

		if name != "" {
			countQuery += " AND name LIKE ?"
			countArgs = append(countArgs, "%"+name+"%")
		}
		if district != "" {
			countQuery += " AND district = ?"
			countArgs = append(countArgs, district)
		}
		if businessArea != "" {
			countQuery += " AND business_area = ?"
			countArgs = append(countArgs, businessArea)
		}
		if status != "" {
			countQuery += " AND status = ?"
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
		result := database.DB.Table("sys_buildings").Where("id = ?", id).First(&building)

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
			Name            string `json:"name" binding:"required"`
			Developer       string `json:"developer"`
			DetailedAddress string `json:"detailedAddress" binding:"required"`
			City            string `json:"city" binding:"required"`
			District        string `json:"district" binding:"required"`
			BusinessArea    string `json:"businessArea"`
			SubDistrict     string `json:"subDistrict"`
			PropertyType    string `json:"propertyType"`
			PropertyCompany string `json:"propertyCompany"`
			Description     string `json:"description"`
			Status          string `json:"status"`
			IsHot           bool   `json:"isHot"`
		}

		if err := c.ShouldBindJSON(&buildingData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"error":   err.Error(),
			})
			return
		}

		// 插入数据库
		result := database.DB.Exec(
			"INSERT INTO sys_buildings (name, developer, detailed_address, city, district, business_area, sub_district, property_type, property_company, description, status, is_hot, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())",
			buildingData.Name,
			buildingData.Developer,
			buildingData.DetailedAddress,
			buildingData.City,
			buildingData.District,
			buildingData.BusinessArea,
			buildingData.SubDistrict,
			buildingData.PropertyType,
			buildingData.PropertyCompany,
			buildingData.Description,
			buildingData.Status,
			buildingData.IsHot,
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
				"id":              newBuildingID,
				"name":            buildingData.Name,
				"developer":       buildingData.Developer,
				"detailedAddress": buildingData.DetailedAddress,
				"city":            buildingData.City,
				"district":        buildingData.District,
				"businessArea":    buildingData.BusinessArea,
				"subDistrict":     buildingData.SubDistrict,
				"propertyType":    buildingData.PropertyType,
				"propertyCompany": buildingData.PropertyCompany,
				"description":     buildingData.Description,
				"status":          buildingData.Status,
				"isHot":           buildingData.IsHot,
			},
		})
	})

	// 更新楼盘信息
	api.PUT("/buildings/:id", func(c *gin.Context) {
		id := c.Param("id")

		// 解析请求体
		var buildingData struct {
			Name            string `json:"name"`
			Developer       string `json:"developer"`
			DetailedAddress string `json:"detailedAddress"`
			City            string `json:"city"`
			District        string `json:"district"`
			BusinessArea    string `json:"businessArea"`
			SubDistrict     string `json:"subDistrict"`
			PropertyType    string `json:"propertyType"`
			PropertyCompany string `json:"propertyCompany"`
			Description     string `json:"description"`
			Status          string `json:"status"`
			IsHot           bool   `json:"isHot"`
		}

		if err := c.ShouldBindJSON(&buildingData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"error":   err.Error(),
			})
			return
		}

		// 构建更新数据
		updateData := make(map[string]interface{})
		if buildingData.Name != "" {
			updateData["name"] = buildingData.Name
		}
		if buildingData.Developer != "" {
			updateData["developer"] = buildingData.Developer
		}
		if buildingData.DetailedAddress != "" {
			updateData["detailed_address"] = buildingData.DetailedAddress
		}
		if buildingData.City != "" {
			updateData["city"] = buildingData.City
		}
		if buildingData.District != "" {
			updateData["district"] = buildingData.District
		}
		if buildingData.BusinessArea != "" {
			updateData["business_area"] = buildingData.BusinessArea
		}
		if buildingData.SubDistrict != "" {
			updateData["sub_district"] = buildingData.SubDistrict
		}
		if buildingData.PropertyType != "" {
			updateData["property_type"] = buildingData.PropertyType
		}
		if buildingData.PropertyCompany != "" {
			updateData["property_company"] = buildingData.PropertyCompany
		}
		if buildingData.Description != "" {
			updateData["description"] = buildingData.Description
		}
		if buildingData.Status != "" {
			updateData["status"] = buildingData.Status
		}
		updateData["is_hot"] = buildingData.IsHot
		updateData["updated_at"] = "NOW()"

		result := database.DB.Table("sys_buildings").Where("id = ?", id).Updates(updateData)

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

		// 删除数据库记录
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
}
