package v1

import (
	"net/http"
	"rentPro/rentpro-admin/common/models/base"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BuildingController 楼盘管理控制器
type BuildingController struct {
	db *gorm.DB
}

// NewBuildingController 创建楼盘控制器
func NewBuildingController(db *gorm.DB) *BuildingController {
	return &BuildingController{db: db}
}

// GetBuildings 获取楼盘列表
func (c *BuildingController) GetBuildings(ctx *gin.Context) {
	var buildings []rental.SysBuildings
	
	// 分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize
	
	// 查询条件
	query := c.db.Model(&rental.SysBuildings{})
	
	// 搜索条件
	if name := ctx.Query("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if city := ctx.Query("city"); city != "" {
		query = query.Where("city = ?", city)
	}
	
	// 获取总数
	var total int64
	query.Count(&total)
	
	// 获取数据
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&buildings).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取楼盘列表失败",
			"error": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取楼盘列表成功",
		"data": gin.H{
			"list":     buildings,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GetBuilding 获取单个楼盘详情
func (c *BuildingController) GetBuilding(ctx *gin.Context) {
	id := ctx.Param("id")
	
	var building rental.SysBuildings
	err := c.db.Where("id = ?", id).First(&building).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "楼盘不存在",
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取楼盘详情成功",
		"data": building,
	})
}

// CreateBuilding 创建楼盘
func (c *BuildingController) CreateBuilding(ctx *gin.Context) {
	var building rental.SysBuildings
	if err := ctx.ShouldBindJSON(&building); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
			"error": err.Error(),
		})
		return
	}
	
	// 设置创建人
	building.CreatedBy = "admin" // 这里应该从JWT中获取当前用户
	
	err := c.db.Create(&building).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建楼盘失败",
			"error": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建楼盘成功",
		"data": building,
	})
}

// UpdateBuilding 更新楼盘
func (c *BuildingController) UpdateBuilding(ctx *gin.Context) {
	id := ctx.Param("id")
	
	var building rental.SysBuildings
	if err := ctx.ShouldBindJSON(&building); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
			"error": err.Error(),
		})
		return
	}
	
	// 设置更新人
	building.UpdatedBy = "admin" // 这里应该从JWT中获取当前用户
	
	err := c.db.Where("id = ?", id).Updates(&building).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新楼盘失败",
			"error": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新楼盘成功",
	})
}

// DeleteBuilding 删除楼盘
func (c *BuildingController) DeleteBuilding(ctx *gin.Context) {
	id := ctx.Param("id")
	
	err := c.db.Where("id = ?", id).Delete(&rental.SysBuildings{}).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除楼盘失败",
			"error": err.Error(),
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除楼盘成功",
	})
}
