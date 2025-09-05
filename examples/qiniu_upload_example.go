package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"rentPro/rentpro-admin/common/config"
	"rentPro/rentpro-admin/common/database"
	"rentPro/rentpro-admin/common/initialize"
	"rentPro/rentpro-admin/common/models/rental"
	"rentPro/rentpro-admin/common/utils"
)

// 这是一个使用七牛云上传户型图的示例
// 可以参考这个实现来替换现有的本地存储方案

func main() {
	// 初始化数据库
	database.Setup()

	// 初始化七牛云服务
	err := initialize.InitQiniu("development")
	if err != nil {
		log.Fatalf("初始化七牛云失败: %v", err)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 设置路由
	setupQiniuRoutes(r)

	// 启动服务器
	log.Println("启动示例服务器: http://localhost:8080")
	r.Run(":8080")
}

func setupQiniuRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		// 七牛云上传户型图
		api.POST("/qiniu/upload/floor-plan", qiniuUploadFloorPlan)

		// 删除七牛云户型图
		api.DELETE("/qiniu/house-types/:id/floor-plan", qiniuDeleteFloorPlan)

		// 获取七牛云配置信息
		api.GET("/qiniu/config", getQiniuConfig)

		// 健康检查
		api.GET("/qiniu/health", qiniuHealthCheck)
	}
}

// qiniuUploadFloorPlan 使用七牛云上传户型图
func qiniuUploadFloorPlan(c *gin.Context) {
	// 获取户型ID
	houseTypeID := c.PostForm("house_type_id")
	if houseTypeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少户型ID参数",
		})
		return
	}

	// 检查户型是否存在
	var houseType rental.SysHouseType
	result := database.DB.Where("id = ? AND deleted_at IS NULL", houseTypeID).First(&houseType)
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

	// 获取七牛云服务
	qiniuService := utils.GetQiniuService()
	if qiniuService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "七牛云服务未初始化",
		})
		return
	}

	// 验证文件
	err = qiniuService.ValidateFile(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 生成自定义文件名
	customKey := fmt.Sprintf("floor_plan_%s_%d.jpg", houseTypeID, time.Now().Unix())

	// 上传到七牛云
	uploadResult, err := qiniuService.UploadFile(file, customKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "上传到七牛云失败",
			"error":   err.Error(),
		})
		return
	}

	// 如果有旧的户型图，删除它
	if houseType.FloorPlanUrl != "" {
		oldKey := qiniuService.ExtractKeyFromURL(houseType.FloorPlanUrl)
		if oldKey != "" {
			err := qiniuService.DeleteFile(oldKey)
			if err != nil {
				log.Printf("删除旧户型图失败: %v", err)
			}
		}
	}

	// 更新数据库
	updateData := map[string]interface{}{
		"floor_plan_url": uploadResult.OriginalURL,
	}

	// 如果需要存储多个尺寸的URL，可以扩展数据库字段
	// updateData["floor_plan_thumbnail_url"] = uploadResult.ThumbnailURL
	// updateData["floor_plan_medium_url"] = uploadResult.MediumURL

	updateResult := database.DB.Model(&houseType).Updates(updateData)
	if updateResult.Error != nil {
		// 如果数据库更新失败，尝试删除已上传的文件
		qiniuService.DeleteFile(uploadResult.Key)
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
			"key":           uploadResult.Key,
			"hash":          uploadResult.Hash,
			"size":          uploadResult.Size,
			"original_url":  uploadResult.OriginalURL,
			"thumbnail_url": uploadResult.ThumbnailURL,
			"medium_url":    uploadResult.MediumURL,
			"large_url":     uploadResult.LargeURL,
			"styles":        uploadResult.Styles,
		},
	})
}

// qiniuDeleteFloorPlan 删除七牛云户型图
func qiniuDeleteFloorPlan(c *gin.Context) {
	houseTypeID := c.Param("id")

	// 检查户型是否存在
	var houseType rental.SysHouseType
	result := database.DB.Where("id = ? AND deleted_at IS NULL", houseTypeID).First(&houseType)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "户型不存在",
		})
		return
	}

	// 获取七牛云服务
	qiniuService := utils.GetQiniuService()
	if qiniuService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "七牛云服务未初始化",
		})
		return
	}

	// 删除七牛云文件
	if houseType.FloorPlanUrl != "" {
		key := qiniuService.ExtractKeyFromURL(houseType.FloorPlanUrl)
		if key != "" {
			err := qiniuService.DeleteFile(key)
			if err != nil {
				log.Printf("删除七牛云文件失败: %v", err)
				// 继续执行数据库更新，不因为云端删除失败而中断
			}
		}
	}

	// 清空数据库中的户型图URL
	updateResult := database.DB.Model(&houseType).Updates(map[string]interface{}{
		"floor_plan_url": "",
		// 如果有其他尺寸的URL字段，也要清空
		// "floor_plan_thumbnail_url": "",
		// "floor_plan_medium_url": "",
	})
	if updateResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新数据库失败",
			"error":   updateResult.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "户型图删除成功",
	})
}

// getQiniuConfig 获取七牛云配置信息（脱敏）
func getQiniuConfig(c *gin.Context) {
	qiniuConfig := config.GetQiniuConfig()
	if qiniuConfig == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "七牛云配置未初始化",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取配置成功",
		"data": gin.H{
			"bucket":        qiniuConfig.Bucket,
			"domain":        qiniuConfig.Domain,
			"zone":          qiniuConfig.Zone,
			"use_https":     qiniuConfig.UseHTTPS,
			"use_cdn":       qiniuConfig.UseCdnDomains,
			"max_file_size": qiniuConfig.Upload.MaxFileSize,
			"allowed_types": qiniuConfig.Upload.AllowedTypes,
			"image_styles":  qiniuConfig.ImageStyles,
		},
	})
}

// qiniuHealthCheck 七牛云健康检查
func qiniuHealthCheck(c *gin.Context) {
	err := initialize.CheckQiniuHealth()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "七牛云服务不可用",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "七牛云服务正常",
		"data": gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
		},
	})
}
