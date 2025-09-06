package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"rentPro/rentpro-admin/common/database"
	"rentPro/rentpro-admin/common/utils"

	"github.com/gin-gonic/gin"
)

// HouseTypeResponse 户型响应结构
type HouseTypeResponse struct {
	ID                  uint64    `json:"id" db:"id"`
	BuildingID          uint64    `json:"building_id" db:"building_id"`
	BuildingName        string    `json:"building_name" db:"building_name"`
	Name                string    `json:"name" db:"name"`
	Code                string    `json:"code" db:"code"`
	Rooms               int64     `json:"rooms" db:"rooms"`
	Halls               int64     `json:"halls" db:"halls"`
	Bathrooms           int64     `json:"bathrooms" db:"bathrooms"`
	Balconies           int64     `json:"balconies" db:"balconies"`
	MaidRooms           int64     `json:"maid_rooms" db:"maid_rooms"`
	StandardArea        float64   `json:"standard_area" db:"standard_area"`
	StandardOrientation string    `json:"standard_orientation" db:"standard_orientation"`
	FloorPlanUrl        string    `json:"floor_plan_url" db:"floor_plan_url"`
	HasFloorPlan        bool      `json:"has_floor_plan" db:"has_floor_plan"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy           string    `json:"created_by" db:"created_by"`
	UpdatedBy           string    `json:"updated_by" db:"updated_by"`
	EditorName          string    `json:"editor_name" db:"editor_name"`
}

// SetupHouseTypeRoutes 设置户型相关路由
func SetupHouseTypeRoutes(api *gin.RouterGroup) {
	// 获取楼盘的户型列表
	api.GET("/house-types/building/:buildingId", func(c *gin.Context) {
		buildingIdStr := c.Param("buildingId")
		buildingId, err := strconv.ParseUint(buildingIdStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的楼盘ID",
			})
			return
		}

		// 分页参数
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 100 {
			pageSize = 10
		}
		offset := (page - 1) * pageSize

		// 构建查询
		query := `SELECT ht.id, ht.building_id, ht.name, ht.code, ht.rooms, ht.halls, ht.bathrooms, 
				 COALESCE(ht.balconies, 0) as balconies,
				 COALESCE(ht.maid_rooms, 0) as maid_rooms,
				 ht.standard_area, 
				 COALESCE(ht.standard_orientation, '') as standard_orientation, 
				 COALESCE(ht.floor_plan_url, '') as floor_plan_url,
				 CASE WHEN ht.floor_plan_url IS NOT NULL AND ht.floor_plan_url != '' THEN true ELSE false END as has_floor_plan,
				 ht.created_at, ht.updated_at, 
				 COALESCE(ht.created_by, '') as created_by, 
				 COALESCE(ht.updated_by, '') as updated_by,
				 COALESCE(u_updated.nick_name, u_created.nick_name, ht.updated_by, ht.created_by, '系统') as editor_name
				 FROM sys_house_types ht
				 LEFT JOIN sys_user u_created ON ht.created_by = u_created.username
				 LEFT JOIN sys_user u_updated ON ht.updated_by = u_updated.username
				 WHERE ht.building_id = ? AND ht.deleted_at IS NULL
				 ORDER BY ht.created_at DESC
				 LIMIT ? OFFSET ?`

		var houseTypes []HouseTypeResponse
		err = database.DB.Raw(query, buildingId, pageSize, offset).Scan(&houseTypes).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "查询户型数据失败",
				"error":   err.Error(),
			})
			return
		}

		// 查询总数
		var total int64
		countQuery := "SELECT COUNT(*) FROM sys_house_types WHERE building_id = ? AND deleted_at IS NULL"
		err = database.DB.Raw(countQuery, buildingId).Scan(&total).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "查询户型总数失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取户型列表成功",
			"data":    houseTypes,
			"total":   total,
			"page":    page,
			"size":    pageSize,
		})
	})

	// 获取单个户型信息
	api.GET("/house-types/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的户型ID",
			})
			return
		}

		query := `SELECT ht.id, ht.building_id, ht.name, ht.code, ht.rooms, ht.halls, ht.bathrooms, 
				 COALESCE(ht.balconies, 0) as balconies,
				 COALESCE(ht.maid_rooms, 0) as maid_rooms,
				 ht.standard_area, 
				 COALESCE(ht.standard_orientation, '') as standard_orientation, 
				 COALESCE(ht.floor_plan_url, '') as floor_plan_url,
				 CASE WHEN ht.floor_plan_url IS NOT NULL AND ht.floor_plan_url != '' THEN true ELSE false END as has_floor_plan,
				 ht.created_at, ht.updated_at, 
				 COALESCE(ht.created_by, '') as created_by, 
				 COALESCE(ht.updated_by, '') as updated_by,
				 COALESCE(u_updated.nick_name, u_created.nick_name, ht.updated_by, ht.created_by, '系统') as editor_name
				 FROM sys_house_types ht
				 LEFT JOIN sys_user u_created ON ht.created_by = u_created.username
				 LEFT JOIN sys_user u_updated ON ht.updated_by = u_updated.username
				 WHERE ht.id = ? AND ht.deleted_at IS NULL`

		var houseType HouseTypeResponse
		err = database.DB.Raw(query, id).First(&houseType).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "户型不存在",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取户型信息成功",
			"data":    houseType,
		})
	})

	// 创建户型
	api.POST("/house-types", func(c *gin.Context) {
		var houseType struct {
			BuildingID          uint64  `json:"building_id" binding:"required"`
			Name                string  `json:"name" binding:"required"`
			Code                string  `json:"code"`
			Rooms               int64   `json:"rooms"`
			Halls               int64   `json:"halls"`
			Bathrooms           int64   `json:"bathrooms"`
			Balconies           int64   `json:"balconies"`
			MaidRooms           int64   `json:"maid_rooms"`
			StandardArea        float64 `json:"standard_area"`
			StandardOrientation string  `json:"standard_orientation"`
		}

		if err := c.ShouldBindJSON(&houseType); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
				"error":   err.Error(),
			})
			return
		}

		// 自动生成code（如果未提供）
		if houseType.Code == "" {
			houseType.Code = strings.ToUpper(strings.ReplaceAll(houseType.Name, " ", "_"))
		}

		// TODO: 从JWT token获取真实用户
		currentUser := "admin"

		result := database.DB.Exec(
			"INSERT INTO sys_house_types (building_id, name, code, rooms, halls, bathrooms, balconies, maid_rooms, standard_area, standard_orientation, created_by, updated_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())",
			houseType.BuildingID,
			houseType.Name,
			houseType.Code,
			houseType.Rooms,
			houseType.Halls,
			houseType.Bathrooms,
			houseType.Balconies,
			houseType.MaidRooms,
			houseType.StandardArea,
			houseType.StandardOrientation,
			currentUser,
			currentUser,
		)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "创建户型失败",
				"error":   result.Error.Error(),
			})
			return
		}

		// 在七牛云创建户型文件夹
		imageManager := utils.GetImageManager()
		if err := imageManager.CreateHouseTypeFolder(houseType.BuildingID, houseType.Name, houseType.StandardArea); err != nil {
			// 记录错误但不影响户型创建成功
			fmt.Printf("⚠️  创建户型文件夹失败: %v\n", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "创建户型成功",
		})
	})

	// 更新户型
	api.PUT("/house-types/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的户型ID",
			})
			return
		}

		var houseType struct {
			Name                string  `json:"name"`
			Code                string  `json:"code"`
			Rooms               int64   `json:"rooms"`
			Halls               int64   `json:"halls"`
			Bathrooms           int64   `json:"bathrooms"`
			Balconies           int64   `json:"balconies"`
			MaidRooms           int64   `json:"maid_rooms"`
			StandardArea        float64 `json:"standard_area"`
			StandardOrientation string  `json:"standard_orientation"`
		}

		if err := c.ShouldBindJSON(&houseType); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "参数错误",
				"error":   err.Error(),
			})
			return
		}

		// TODO: 从JWT token获取真实用户
		currentUser := "admin"

		var setParts []string
		var values []interface{}

		if houseType.Name != "" {
			setParts = append(setParts, "name = ?")
			values = append(values, houseType.Name)
		}
		if houseType.Code != "" {
			setParts = append(setParts, "code = ?")
			values = append(values, houseType.Code)
		}
		if houseType.Rooms > 0 {
			setParts = append(setParts, "rooms = ?")
			values = append(values, houseType.Rooms)
		}
		if houseType.Halls > 0 {
			setParts = append(setParts, "halls = ?")
			values = append(values, houseType.Halls)
		}
		if houseType.Bathrooms > 0 {
			setParts = append(setParts, "bathrooms = ?")
			values = append(values, houseType.Bathrooms)
		}
		if houseType.Balconies >= 0 {
			setParts = append(setParts, "balconies = ?")
			values = append(values, houseType.Balconies)
		}
		if houseType.MaidRooms >= 0 {
			setParts = append(setParts, "maid_rooms = ?")
			values = append(values, houseType.MaidRooms)
		}
		if houseType.StandardArea > 0 {
			setParts = append(setParts, "standard_area = ?")
			values = append(values, houseType.StandardArea)
		}
		if houseType.StandardOrientation != "" {
			setParts = append(setParts, "standard_orientation = ?")
			values = append(values, houseType.StandardOrientation)
		}

		// 总是更新 updated_at 和 updated_by
		setParts = append(setParts, "updated_at = ?")
		values = append(values, time.Now())
		setParts = append(setParts, "updated_by = ?")
		values = append(values, currentUser)

		if len(setParts) == 2 { // 只有时间和用户更新字段
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "没有提供要更新的字段",
			})
			return
		}

		// 添加WHERE条件的参数
		values = append(values, id)

		query := "UPDATE sys_house_types SET " + strings.Join(setParts, ", ") + " WHERE id = ?"

		result := database.DB.Exec(query, values...)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "更新户型失败",
				"error":   result.Error.Error(),
			})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "户型不存在",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "更新户型成功",
		})
	})

	// 删除户型（软删除）
	api.DELETE("/house-types/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的户型ID",
			})
			return
		}

		// TODO: 从JWT token获取真实用户
		currentUser := "admin"

		result := database.DB.Exec("UPDATE sys_house_types SET deleted_at = NOW(), updated_by = ? WHERE id = ? AND deleted_at IS NULL", currentUser, id)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "删除户型失败",
				"error":   result.Error.Error(),
			})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "户型不存在或已被删除",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "删除户型成功",
		})
	})

	// 获取已删除的户型列表（回收站）
	api.GET("/house-types/building/:buildingId/deleted", func(c *gin.Context) {
		buildingIdStr := c.Param("buildingId")
		buildingId, err := strconv.ParseUint(buildingIdStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的楼盘ID",
			})
			return
		}

		pageStr := c.DefaultQuery("page", "1")
		pageSizeStr := c.DefaultQuery("pageSize", "10")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 || pageSize > 100 {
			pageSize = 10
		}

		offset := (page - 1) * pageSize

		// 查询已删除的户型列表
		query := `SELECT ht.id, ht.building_id, COALESCE(b.name, '') as building_name, ht.name, ht.code, ht.rooms, ht.halls, ht.bathrooms, 
				 COALESCE(ht.balconies, 0) as balconies,
				 COALESCE(ht.maid_rooms, 0) as maid_rooms,
				 ht.standard_area, 
				 COALESCE(ht.standard_orientation, '') as standard_orientation, 
				 COALESCE(ht.floor_plan_url, '') as floor_plan_url,
				 CASE WHEN ht.floor_plan_url IS NOT NULL AND ht.floor_plan_url != '' THEN true ELSE false END as has_floor_plan,
				 ht.created_at, ht.updated_at, ht.deleted_at,
				 COALESCE(ht.created_by, '') as created_by, 
				 COALESCE(ht.updated_by, '') as updated_by,
				 COALESCE(u_updated.nick_name, u_created.nick_name, ht.updated_by, ht.created_by, '系统') as editor_name
				 FROM sys_house_types ht
				 LEFT JOIN sys_buildings b ON ht.building_id = b.id AND b.deleted_at IS NULL
				 LEFT JOIN sys_user u_created ON ht.created_by = u_created.username
				 LEFT JOIN sys_user u_updated ON ht.updated_by = u_updated.username
				 WHERE ht.building_id = ? AND ht.deleted_at IS NOT NULL
				 ORDER BY ht.deleted_at DESC
				 LIMIT ? OFFSET ?`

		var houseTypes []HouseTypeResponse
		err = database.DB.Raw(query, buildingId, pageSize, offset).Scan(&houseTypes).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "查询已删除户型失败",
				"error":   err.Error(),
			})
			return
		}

		// 查询总数
		var total int64
		countQuery := "SELECT COUNT(*) FROM sys_house_types WHERE building_id = ? AND deleted_at IS NOT NULL"
		err = database.DB.Raw(countQuery, buildingId).Scan(&total).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "查询已删除户型总数失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取已删除户型列表成功",
			"data":    houseTypes,
			"total":   total,
		})
	})

	// 恢复户型（取消软删除）
	api.POST("/house-types/:id/restore", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的户型ID",
			})
			return
		}

		// TODO: 从JWT token获取真实用户
		currentUser := "admin"

		result := database.DB.Exec("UPDATE sys_house_types SET deleted_at = NULL, updated_by = ? WHERE id = ? AND deleted_at IS NOT NULL", currentUser, id)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "恢复户型失败",
				"error":   result.Error.Error(),
			})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "户型不存在或未被删除",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "恢复户型成功",
		})
	})

	// 永久删除户型
	api.DELETE("/house-types/:id/permanent", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的户型ID",
			})
			return
		}

		result := database.DB.Exec("DELETE FROM sys_house_types WHERE id = ? AND deleted_at IS NOT NULL", id)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "永久删除户型失败",
				"error":   result.Error.Error(),
			})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "户型不存在或未被删除",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "永久删除户型成功",
		})
	})

	// 获取户型的所有户型图
	api.GET("/house-types/:id/floor-plans", func(c *gin.Context) {
		idStr := c.Param("id")
		houseTypeId, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的户型ID",
			})
			return
		}

		// 查询户型图片列表
		query := `SELECT id, name, description, file_name, file_size, url, thumbnail_url, 
				 medium_url, large_url, created_at, sort_order 
				 FROM sys_images 
				 WHERE module = 'house_floor_plan' AND module_id = ? AND deleted_at IS NULL 
				 ORDER BY sort_order ASC, created_at ASC`

		var images []map[string]interface{}
		err = database.DB.Raw(query, houseTypeId).Scan(&images).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "查询户型图片失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "查询成功",
			"data":    images,
		})
	})

	// 删除单张户型图
	api.DELETE("/house-types/:id/floor-plans/:imageId", func(c *gin.Context) {
		idStr := c.Param("id")
		houseTypeId, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的户型ID",
			})
			return
		}

		imageIdStr := c.Param("imageId")
		imageId, err := strconv.ParseUint(imageIdStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的图片ID",
			})
			return
		}

		// TODO: 从JWT token获取真实用户
		currentUser := "admin"

		// 软删除图片记录
		result := database.DB.Exec(
			"UPDATE sys_images SET deleted_at = NOW(), updated_by = ? WHERE id = ? AND module = 'house_floor_plan' AND module_id = ? AND deleted_at IS NULL",
			currentUser, imageId, houseTypeId)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "删除户型图片失败",
				"error":   result.Error.Error(),
			})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "户型图片不存在",
			})
			return
		}

		// 更新户型表的 floor_plan_url（如果删除的是第一张图片，则设置为下一张图片的URL）
		var firstImageUrl string
		err = database.DB.Raw(
			"SELECT url FROM sys_images WHERE module = 'house_floor_plan' AND module_id = ? AND deleted_at IS NULL ORDER BY sort_order ASC, created_at ASC LIMIT 1",
			houseTypeId).Scan(&firstImageUrl).Error

		if err == nil {
			// 更新户型表的 floor_plan_url
			database.DB.Exec("UPDATE sys_house_types SET floor_plan_url = ? WHERE id = ?", firstImageUrl, houseTypeId)
		} else {
			// 没有图片了，清空 floor_plan_url
			database.DB.Exec("UPDATE sys_house_types SET floor_plan_url = NULL WHERE id = ?", houseTypeId)
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "删除户型图片成功",
		})
	})
}
