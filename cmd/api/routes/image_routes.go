package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"rentPro/rentpro-admin/common/database"
	"rentPro/rentpro-admin/common/models/image"
	"rentPro/rentpro-admin/common/utils"

	"github.com/gin-gonic/gin"
)

// SetupImageRoutes 设置图片管理相关路由
func SetupImageRoutes(api *gin.RouterGroup) {
	// 上传图片
	api.POST("/images/upload", func(c *gin.Context) {
		// 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权访问",
			})
			return
		}

		// 获取上传的文件
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "获取上传文件失败",
				"error":   err.Error(),
			})
			return
		}

		// 解析请求参数
		var req image.ImageUploadRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"error":   err.Error(),
			})
			return
		}

		// 获取图片管理器
		imageManager := utils.GetImageManager()
		if imageManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "图片管理器未初始化",
			})
			return
		}

		// 上传图片
		img, err := imageManager.UploadImage(file, &req, userID.(uint64))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "上传图片失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "图片上传成功",
			"data":    img,
		})
	})

	// 获取图片列表
	api.GET("/images", func(c *gin.Context) {
		var req image.ImageListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"error":   err.Error(),
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

		result, err := imageManager.ListImages(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取图片列表失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取图片列表成功",
			"data":    result,
		})
	})

	// 获取图片详情
	api.GET("/images/:id", func(c *gin.Context) {
		id := c.Param("id")
		imageID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "图片ID格式错误",
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

		img, err := imageManager.GetImage(imageID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取图片详情成功",
			"data":    img,
		})
	})

	// 更新图片信息
	api.PUT("/images/:id", func(c *gin.Context) {
		id := c.Param("id")
		imageID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "图片ID格式错误",
			})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权访问",
			})
			return
		}

		var req image.ImageUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"error":   err.Error(),
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

		if err := imageManager.UpdateImage(imageID, &req, userID.(uint64)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "更新图片信息失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "更新图片信息成功",
		})
	})

	// 删除图片
	api.DELETE("/images/:id", func(c *gin.Context) {
		id := c.Param("id")
		imageID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "图片ID格式错误",
			})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权访问",
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

		if err := imageManager.DeleteImage(imageID, userID.(uint64)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "删除图片失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "删除图片成功",
		})
	})

	// 批量删除图片
	api.DELETE("/images/batch", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权访问",
			})
			return
		}

		var req image.ImageBatchDeleteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"error":   err.Error(),
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

		if err := imageManager.BatchDeleteImages(req.IDs, userID.(uint64)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "批量删除图片失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "批量删除图片成功",
		})
	})

	// 获取模块图片
	api.GET("/images/module/:module/:moduleId", func(c *gin.Context) {
		module := c.Param("module")
		moduleIDStr := c.Param("moduleId")
		category := c.Query("category")

		moduleID, err := strconv.ParseUint(moduleIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "模块ID格式错误",
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

		images, err := imageManager.GetImagesByModule(module, moduleID, category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取模块图片失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取模块图片成功",
			"data":    images,
		})
	})

	// 设置主图
	api.PUT("/images/:id/set-main", func(c *gin.Context) {
		id := c.Param("id")
		imageID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "图片ID格式错误",
			})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权访问",
			})
			return
		}

		var req struct {
			Module   string `json:"module" binding:"required"`
			ModuleID uint64 `json:"moduleId" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"error":   err.Error(),
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

		if err := imageManager.SetMainImage(req.Module, req.ModuleID, imageID, userID.(uint64)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "设置主图失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "设置主图成功",
		})
	})

	// 获取图片统计信息
	api.GET("/images/stats", func(c *gin.Context) {
		imageManager := utils.GetImageManager()
		if imageManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "图片管理器未初始化",
			})
			return
		}

		stats, err := imageManager.GetImageStats()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取图片统计信息失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取图片统计信息成功",
			"data":    stats,
		})
	})

	// 上传楼盘户型图（兼容旧接口）
	api.POST("/upload/floor-plan", func(c *gin.Context) {
		// 从请求头获取token并验证
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未提供认证信息",
			})
			return
		}

		// 提取token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证格式错误",
			})
			return
		}

		// 解析token
		jwtInstance := &utils.JWT{
			Config: utils.JWTConfig{
				Secret:  "rentpro-admin-secret-key",
				Timeout: 86400, // 24 hours in seconds
			},
		}
		claims, err := jwtInstance.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "token无效",
				"error":   err.Error(),
			})
			return
		}

		// 检查token是否过期
		if claims.ExpiresAt != nil && time.Now().Unix() > claims.ExpiresAt.Unix() {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "token已过期",
			})
			return
		}

		// 获取用户ID
		userID := claims.UserID

		// 获取户型ID
		houseTypeIDStr := c.PostForm("house_type_id")
		if houseTypeIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "缺少户型ID参数",
			})
			return
		}

		houseTypeID, err := strconv.ParseUint(houseTypeIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "户型ID格式错误",
			})
			return
		}

		// 检查户型是否存在，并获取楼盘ID
		var houseType struct {
			BuildingID uint `json:"building_id"`
		}
		result := database.DB.Table("sys_house_types").Where("id = ? AND deleted_at IS NULL", houseTypeID).First(&houseType)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "户型不存在",
			})
			return
		}

		// 获取上传的文件
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "获取上传文件失败",
				"error":   err.Error(),
			})
			return
		}

		// 使用图片管理器上传楼盘户型图
		imageManager := utils.GetImageManager()
		if imageManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "图片管理器未初始化",
			})
			return
		}

		// 上传楼盘户型图
		img, err := imageManager.UploadBuildingFloorPlan(file, uint64(houseType.BuildingID), houseTypeID, uint64(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "上传户型图失败",
				"error":   err.Error(),
			})
			return
		}

		// 更新户型的floor_plan_url
		updateResult := database.DB.Model(&struct{}{}).Table("sys_house_types").Where("id = ?", houseTypeID).Update("floor_plan_url", img.URL)
		if updateResult.Error != nil {
			// 如果数据库更新失败，删除已上传的文件
			imageManager.DeleteImage(img.ID, uint64(userID))
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "更新数据库失败",
				"error":   updateResult.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "户型图上传成功",
			"data": gin.H{
				"image_id":      img.ID,
				"original_url":  img.URL,
				"thumbnail_url": img.ThumbnailURL,
				"medium_url":    img.MediumURL,
				"large_url":     img.LargeURL,
				"file_size":     img.FileSize,
				"building_id":   houseType.BuildingID,
				"house_type_id": houseTypeID,
			},
		})
	})

	// 批量上传户型图（新的多图上传API）
	api.POST("/upload/house-type-floor-plans", func(c *gin.Context) {
		// 从表单获取用户ID
		userIDStr := c.PostForm("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			userID = 1 // 默认用户ID，实际应该从JWT token获取
		}

		// 从表单获取户型ID
		houseTypeIDStr := c.PostForm("house_type_id")
		houseTypeID, err := strconv.ParseUint(houseTypeIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "户型ID格式错误",
			})
			return
		}

		// 检查户型是否存在
		var houseType struct {
			ID         uint64 `json:"id"`
			BuildingID uint   `json:"building_id"`
			Name       string `json:"name"`
		}
		result := database.DB.Table("sys_house_types").Where("id = ? AND deleted_at IS NULL", houseTypeID).First(&houseType)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "户型不存在",
			})
			return
		}

		// 获取上传的文件列表
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "获取上传文件失败",
				"error":   err.Error(),
			})
			return
		}

		files := form.File["files"]
		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请选择要上传的图片",
			})
			return
		}

		if len(files) > 5 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "一次最多只能上传5张图片",
			})
			return
		}

		// 使用图片管理器上传多张户型图
		imageManager := utils.GetImageManager()
		if imageManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "图片管理器未初始化",
			})
			return
		}

		// 批量上传户型图
		images, err := imageManager.UploadHouseTypeFloorPlans(files, houseTypeID, uint64(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "上传户型图失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": fmt.Sprintf("成功上传%d张户型图", len(images)),
			"data":    images,
		})
	})
}
