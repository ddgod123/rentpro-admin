// Package api 提供HTTP API服务器相关的命令行功能
// 用于启动 rentpro-admin 系统的HTTP API服务
package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"rentPro/rentpro-admin/cmd/api/routes"
	"rentPro/rentpro-admin/common/database"
	"rentPro/rentpro-admin/common/global"
	"rentPro/rentpro-admin/common/initialize"
	"rentPro/rentpro-admin/common/utils"
)

// fileExists 检查文件是否存在
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// 配置数据结构，用于解析 settings.yml
type Config struct {
	Settings struct {
		Application struct {
			Mode         string `yaml:"mode"`
			Host         string `yaml:"host"`
			Name         string `yaml:"name"`
			Port         int    `yaml:"port"`
			ReadTimeout  int    `yaml:"readtimeout"`
			WriteTimeout int    `yaml:"writetimeout"`
			EnabledDP    bool   `yaml:"enabledp"`
		} `yaml:"application"`
		Logger struct {
			Path      string `yaml:"path"`
			Stdout    string `yaml:"stdout"`
			Level     string `yaml:"level"`
			EnabledDB bool   `yaml:"enableddb"`
		} `yaml:"logger"`
		JWT struct {
			Secret  string `yaml:"secret"`
			Timeout int    `yaml:"timeout"`
		} `yaml:"jwt"`
		Database struct {
			Driver string `yaml:"driver"`
			Source string `yaml:"source"`
		} `yaml:"database"`
	} `yaml:"settings"`
}

var (
	configYml   string
	port        int
	showVersion bool

	// StartCmd 定义了 api 子命令
	// 用于启动HTTP API服务器，支持以下功能：
	// 1. HTTP API服务启动和管理
	// 2. 数据库连接初始化
	// 3. 路由配置和中间件
	// 4. 优雅关闭
	// 命令注册：通过 rootCmd.AddCommand(api.StartCmd) 注册到根命令
	// 使用方式：
	//   - rentpro-admin api -c config/settings.yml : 使用指定配置文件启动API服务器
	//   - rentpro-admin api -p 8002                : 指定端口启动API服务器
	//   - rentpro-admin api -v                     : 显示版本信息
	// 版本信息来源：common/global/adm.go 中的 Version 常量
	StartCmd = &cobra.Command{
		Use:     "api",
		Short:   "启动HTTP API服务器",
		Long:    `rentpro-admin HTTP API服务器，提供完整的权限管理API、用户认证、JWT令牌等功能`,
		Example: "rentpro-admin api -c config/settings.yml -p 8002",
		PreRun: func(cmd *cobra.Command, args []string) {
			if showVersion {
				fmt.Printf("rentpro-admin api version: %s\n", global.Version)
				return
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				return nil
			}
			return run()
		},
	}
)

// init 初始化命令标志
func init() {
	// 添加版本标志支持
	StartCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "显示版本信息")

	// 配置文件路径标志
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.yml", "指定配置文件路径")

	// 端口标志
	StartCmd.PersistentFlags().IntVarP(&port, "port", "p", 0, "指定服务端口号")
}

// run 执行API服务器启动的核心逻辑
func run() error {
	fmt.Printf("=== rentpro-admin API服务器 v%s ===\n", global.Version)

	// 加载配置文件
	config, err := loadConfig(configYml)
	if err != nil {
		return fmt.Errorf("加载配置文件失败: %v", err)
	}

	// 初始化数据库连接
	fmt.Println("初始化数据库连接...")
	database.Setup()

	// 初始化七牛云服务
	fmt.Println("初始化七牛云服务...")
	// TODO: 取消注释以启用七牛云服务
	err = initialize.InitQiniu(config.Settings.Application.Mode)
	if err != nil {
		log.Printf("⚠️  七牛云服务初始化失败: %v", err)
		log.Println("将使用本地文件存储")
	} else {
		// 初始化图片管理器
		fmt.Println("初始化图片管理器...")
		err = utils.InitImageManager()
		if err != nil {
			log.Printf("⚠️  图片管理器初始化失败: %v", err)
		} else {
			log.Println("✅ 图片管理器初始化成功")
		}
	}

	// 设置Gin模式
	if config.Settings.Application.Mode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎
	router := gin.Default()

	// 设置中间件
	setupMiddleware(router)

	// 设置路由
	setupRoutes(router)

	// 确定端口
	serverPort := config.Settings.Application.Port
	if port > 0 {
		serverPort = port
	}

	// 创建HTTP服务器
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Settings.Application.Host, serverPort),
		Handler:      router,
		ReadTimeout:  time.Duration(config.Settings.Application.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Settings.Application.WriteTimeout) * time.Second,
	}

	// 启动服务器
	fmt.Printf("启动API服务器: %s:%d\n", config.Settings.Application.Host, serverPort)
	fmt.Printf("应用名称: %s\n", config.Settings.Application.Name)
	fmt.Printf("运行模式: %s\n", config.Settings.Application.Mode)

	// 在goroutine中启动服务器
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("服务器启动失败: %v\n", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("正在关闭服务器...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("服务器关闭失败: %v\n", err)
		return err
	}

	fmt.Println("✅ 服务器已优雅关闭")
	return nil
}

// loadConfig 加载和解析配置文件
func loadConfig(configPath string) (*Config, error) {
	// 检查文件是否存在
	if !fileExists(configPath) {
		return nil, fmt.Errorf("配置文件不存在: %s", configPath)
	}

	// 读取配置文件
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析 YAML 配置
	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}

// setupMiddleware 设置中间件
func setupMiddleware(router *gin.Engine) {
	// 添加CORS中间件
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 添加日志中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
}

func setupRoutes(router *gin.Engine) {
	// 静态文件服务
	router.Static("/uploads", "./uploads")

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": global.Version,
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// API版本路由组
	api := router.Group("/api/v1")
	{
		// 设置各个模块的路由
		routes.SetupAuthRoutes(api)      // 认证相关路由
		routes.SetupUserRoutes(api)      // 用户管理路由
		routes.SetupCityRoutes(api)      // 城市管理路由
		routes.SetupBuildingRoutes(api)  // 楼盘管理路由
		routes.SetupHouseTypeRoutes(api) // 户型管理路由
		routes.SetupImageRoutes(api)     // 图片管理路由
	}

	// 根路径
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "rentpro-admin API服务器",
			"version": global.Version,
			"docs":    "/api/v1",
		})
	})
}
