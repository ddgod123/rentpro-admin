package v1

import (
	"rentPro/rentpro-admin/common/api"
	"rentPro/rentpro-admin/common/database"

	"github.com/gin-gonic/gin"
)

// RegisterBuildingRoutes 注册楼盘管理路由
func RegisterBuildingRoutes(r *gin.RouterGroup, authController *api.AuthController) {
	buildingController := NewBuildingController(database.DB)

	buildingGroup := r.Group("/building")
	buildingGroup.Use(authController.JWTAuthMiddleware()) // 添加JWT认证中间件
	{
		// 获取楼盘列表
		buildingGroup.GET("", buildingController.GetBuildings)
		// 获取楼盘详情
		buildingGroup.GET("/:id", buildingController.GetBuilding)
		// 创建楼盘
		buildingGroup.POST("", buildingController.CreateBuilding)
		// 更新楼盘
		buildingGroup.PUT("/:id", buildingController.UpdateBuilding)
		// 删除楼盘
		buildingGroup.DELETE("/:id", buildingController.DeleteBuilding)
	}
}
