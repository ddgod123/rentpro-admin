package routes

import (
	"net/http"
	"strconv"

	"rentPro/rentpro-admin/common/database"
	"rentPro/rentpro-admin/common/models/rental"

	"github.com/gin-gonic/gin"
)

// SetupCityRoutes 设置城市相关路由
func SetupCityRoutes(api *gin.RouterGroup) {
	cityGroup := api.Group("/cities")
	{
		cityGroup.GET("", getCities)         // 获取城市列表
		cityGroup.GET("/:id", getCityByID)   // 根据ID获取城市
		cityGroup.POST("", createCity)       // 创建城市
		cityGroup.PUT("/:id", updateCity)    // 更新城市
		cityGroup.DELETE("/:id", deleteCity) // 删除城市
	}
}

// CityResponse 城市响应结构
type CityResponse struct {
	ID        uint64 `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Sort      int64  `json:"sort"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// getCities 获取城市列表
func getCities(c *gin.Context) {
	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "100"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 1000 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	// 构建查询
	db := database.DB
	query := db.Model(&rental.SysCity{})

	// 状态过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取城市总数失败",
			"error":   err.Error(),
		})
		return
	}

	// 获取城市列表
	var cities []rental.SysCity
	err := query.Order("sort ASC, id ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&cities).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取城市列表失败",
			"error":   err.Error(),
		})
		return
	}

	// 转换为响应格式
	var cityList []CityResponse
	for _, city := range cities {
		cityList = append(cityList, CityResponse{
			ID:        city.ID,
			Code:      city.Code,
			Name:      city.Name,
			Sort:      city.Sort,
			Status:    city.Status,
			CreatedAt: city.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: city.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      200,
		"message":   "获取城市列表成功",
		"data":      cityList, // 直接返回城市列表，前端期望的格式
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// getCityByID 根据ID获取城市
func getCityByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的城市ID",
		})
		return
	}

	var city rental.SysCity
	db := database.DB
	if err := db.First(&city, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "城市不存在",
		})
		return
	}

	cityResponse := CityResponse{
		ID:        city.ID,
		Code:      city.Code,
		Name:      city.Name,
		Sort:      city.Sort,
		Status:    city.Status,
		CreatedAt: city.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: city.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取城市成功",
		"data":    cityResponse,
	})
}

// CreateCityRequest 创建城市请求结构
type CreateCityRequest struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Sort   int64  `json:"sort"`
	Status string `json:"status"`
}

// createCity 创建城市
func createCity(c *gin.Context) {
	var req CreateCityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 设置默认状态
	if req.Status == "" {
		req.Status = "active"
	}

	city := rental.SysCity{
		Code:   req.Code,
		Name:   req.Name,
		Sort:   req.Sort,
		Status: req.Status,
	}

	db := database.DB
	if err := db.Create(&city).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建城市失败",
			"error":   err.Error(),
		})
		return
	}

	cityResponse := CityResponse{
		ID:        city.ID,
		Code:      city.Code,
		Name:      city.Name,
		Sort:      city.Sort,
		Status:    city.Status,
		CreatedAt: city.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: city.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "创建城市成功",
		"data":    cityResponse,
	})
}

// UpdateCityRequest 更新城市请求结构
type UpdateCityRequest struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Sort   int64  `json:"sort"`
	Status string `json:"status"`
}

// updateCity 更新城市
func updateCity(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的城市ID",
		})
		return
	}

	var req UpdateCityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	db := database.DB
	var city rental.SysCity
	if err := db.First(&city, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "城市不存在",
		})
		return
	}

	// 更新字段
	if req.Code != "" {
		city.Code = req.Code
	}
	if req.Name != "" {
		city.Name = req.Name
	}
	if req.Sort != 0 {
		city.Sort = req.Sort
	}
	if req.Status != "" {
		city.Status = req.Status
	}

	if err := db.Save(&city).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新城市失败",
			"error":   err.Error(),
		})
		return
	}

	cityResponse := CityResponse{
		ID:        city.ID,
		Code:      city.Code,
		Name:      city.Name,
		Sort:      city.Sort,
		Status:    city.Status,
		CreatedAt: city.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: city.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新城市成功",
		"data":    cityResponse,
	})
}

// deleteCity 删除城市
func deleteCity(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的城市ID",
		})
		return
	}

	db := database.DB
	var city rental.SysCity
	if err := db.First(&city, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "城市不存在",
		})
		return
	}

	if err := db.Delete(&city).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除城市失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除城市成功",
	})
}
