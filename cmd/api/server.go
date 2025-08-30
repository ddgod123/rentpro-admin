// Package api 提供RentPro房源管理系统的API服务器功能
// 包含用户认证、权限管理、租赁管理等核心业务接口
package api

import (
	"fmt"
	"log"

	"rentPro/rentpro-admin/common/api"
	"rentPro/rentpro-admin/common/database"
	"rentPro/rentpro-admin/common/models/base"
	"rentPro/rentpro-admin/common/router"
	"rentPro/rentpro-admin/common/service"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// 全局变量定义
var (
	configFile string // 配置文件路径
	port       string // 服务监听端口
)

// StartCmd 启动API服务器命令
// 使用Cobra命令行框架，提供标准的CLI接口
var StartCmd = &cobra.Command{
	Use:     "api",
	Short:   "启动API服务器",
	Long:    "启动RentPro Admin的API服务器，提供RESTful API接口",
	Example: "rentpro-admin api -c config/settings.yml -p 8002",
	RunE:    run,
}

// init 初始化命令行参数
// 设置配置文件路径和服务端口等命令行选项
func init() {
	// 配置文件路径参数，默认值为 config/settings.yml
	StartCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config/settings.yml", "配置文件路径")
	// 服务端口参数，默认值为 8002
	StartCmd.PersistentFlags().StringVarP(&port, "port", "p", "8002", "服务端口")
}

// run 执行API服务器启动逻辑
// 这是Cobra命令的核心执行函数
func run(cmd *cobra.Command, args []string) error {
	fmt.Printf("正在启动API服务器...\n")
	fmt.Printf("配置文件: %s\n", configFile)
	fmt.Printf("服务端口: %s\n", port)

	// 初始化数据库连接和基础数据
	if err := initDatabase(); err != nil {
		return fmt.Errorf("数据库初始化失败: %v", err)
	}

	// 设置Gin框架为生产模式，提高性能
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin引擎实例
	r := gin.New()

	// 添加基础中间件
	r.Use(gin.Logger())   // 请求日志记录
	r.Use(gin.Recovery()) // 异常恢复处理

	// 添加CORS跨域中间件，支持前端跨域访问
	r.Use(corsMiddleware())

	// 注册所有API路由
	registerRoutes(r)

	// 启动HTTP服务器
	fmt.Printf("🚀 API服务器启动成功！\n")
	fmt.Printf("📡 监听地址: http://localhost:%s\n", port)
	fmt.Printf("📖 API文档: http://localhost:%s/swagger/index.html\n", port)
	fmt.Printf("💡 按 Ctrl+C 停止服务器\n\n")

	// 启动服务器并监听指定端口
	return r.Run(":" + port)
}

// initDatabase 初始化数据库连接和基础数据
// 包括数据库连接、表结构创建、默认数据插入等
func initDatabase() error {
	log.Printf("初始化数据库连接...")

	// 建立数据库连接
	database.Setup()

	// 初始化权限相关的数据表结构
	// 包括用户、角色、菜单、部门等基础表
	if err := base.InitAuthTables(database.DB); err != nil {
		return fmt.Errorf("初始化数据表失败: %v", err)
	}

	// 初始化系统默认数据
	// 包括默认用户、角色、菜单、权限等基础数据
	if err := base.InitDefaultData(database.DB); err != nil {
		return fmt.Errorf("初始化默认数据失败: %v", err)
	}

	log.Printf("✅ 数据库初始化完成")
	return nil
}

// registerRoutes 注册所有API路由
// 按照功能模块组织路由结构，包括认证、系统管理、租赁管理等
func registerRoutes(r *gin.Engine) {
	// 创建认证控制器实例
	authController := api.NewAuthController()

	// API根路径分组，所有API都以 /api 开头
	apiGroup := r.Group("/api")

	// 认证相关路由组（无需认证即可访问）
	authGroup := apiGroup.Group("/auth")
	{
		// 用户登录接口
		authGroup.POST("/login", authController.Login)
	}

	// 需要JWT认证的路由组
	protectedGroup := apiGroup.Group("/auth")
	protectedGroup.Use(authController.JWTAuthMiddleware()) // 添加JWT认证中间件
	{
		// 获取当前用户信息
		protectedGroup.GET("/user-info", authController.GetUserInfo)
		// 用户登出
		protectedGroup.POST("/logout", authController.Logout)
		// 获取用户菜单权限
		protectedGroup.GET("/menus", authController.GetMenus)
		// 检查用户权限
		protectedGroup.GET("/check-permission", authController.CheckPermission)
	}

	// 管理员专用路由组
	adminGroup := apiGroup.Group("/admin")
	adminGroup.Use(authController.JWTAuthMiddleware()) // JWT认证
	adminGroup.Use(authController.AdminMiddleware())   // 管理员权限验证
	{
		// 管理员用户列表接口
		adminGroup.GET("/users", func(c *gin.Context) {
			c.JSON(200, gin.H{"msg": "用户列表 - 需要管理员权限"})
		})
	}

	// 注册动态路由（包含楼盘管理等所有业务路由）
	registerDynamicRoutes(r, authController)

	// 健康检查接口
	// 用于监控系统状态和负载均衡器健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "rentpro-admin",
			"version": "1.0.0",
		})
	})

	// 根路径接口
	// 提供API基本信息和服务状态
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "欢迎使用 RentPro Admin API",
			"version": "1.0.0",
			"docs":    "/swagger/index.html",
		})
	})
}

// registerDynamicRoutes 注册动态路由
func registerDynamicRoutes(r *gin.Engine, authController *api.AuthController) {
	fmt.Println("=== 开始注册动态路由 ===")

	// 创建菜单服务
	menuService := service.NewMenuService(database.DB)
	fmt.Println("菜单服务创建成功")

	// 创建动态路由生成器
	dynamicRouter := router.NewDynamicRouter(menuService, authController)
	fmt.Println("动态路由生成器创建成功")

	// 注册动态路由
	dynamicRouter.RegisterDynamicRoutes(r)
	fmt.Println("=== 动态路由注册完成 ===")
}

// corsMiddleware CORS跨域中间件
// 处理跨域请求，允许前端应用访问API
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许所有来源的跨域请求
		c.Header("Access-Control-Allow-Origin", "*")
		// 允许的HTTP方法
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		// 允许的请求头
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		// 暴露给客户端的响应头
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		// 允许发送Cookie等凭证信息
		c.Header("Access-Control-Allow-Credentials", "true")

		// 处理预检请求（OPTIONS方法）
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // 返回204状态码，表示请求成功但无内容
			return
		}

		// 继续处理下一个中间件或路由处理器
		c.Next()
	}
}
