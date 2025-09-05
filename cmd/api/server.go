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
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"rentPro/rentpro-admin/common/config"
	"rentPro/rentpro-admin/common/database"
	"rentPro/rentpro-admin/common/global"
	"rentPro/rentpro-admin/common/initialize"
	"rentPro/rentpro-admin/common/models/image"
	"rentPro/rentpro-admin/common/models/rental"
	"rentPro/rentpro-admin/common/models/system"
	"rentPro/rentpro-admin/common/utils"
)

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
		// 用户相关API
		api.GET("/users", func(c *gin.Context) {
			// 从数据库查询用户列表
			var users []map[string]interface{}
			database.DB.Raw("SELECT id, username, nick_name, email, phone, status, created_at FROM sys_user WHERE deleted_at IS NULL").Scan(&users)

			c.JSON(http.StatusOK, gin.H{
				"message": "用户列表API",
				"data":    users,
				"total":   len(users),
			})
		})

		// 认证相关API
		api.POST("/auth/login", func(c *gin.Context) {
			// 解析请求体
			var loginData struct {
				Username string `json:"username" binding:"required"`
				Password string `json:"password" binding:"required"`
			}

			if err := c.ShouldBindJSON(&loginData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "请求参数错误",
					"error":   err.Error(),
				})
				return
			}

			fmt.Printf("登录请求: username=%s, password=%s\n", loginData.Username, loginData.Password)

			// 从数据库验证用户
			var user system.SysUser
			result := database.DB.Where("username = ? AND deleted_at IS NULL", loginData.Username).First(&user)

			fmt.Printf("数据库查询结果: 错误=%v\n", result.Error)

			if result.Error != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "用户名或密码错误",
				})
				return
			}

			// 验证密码
			if !user.ComparePassword(loginData.Password) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "用户名或密码错误",
				})
				return
			}

			// 获取JWT配置
			config, err := loadConfig(configYml)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "服务器配置错误",
				})
				return
			}

			// 创建JWT工具
			jwtUtil := utils.NewJWT(utils.JWTConfig{
				Secret:  config.Settings.JWT.Secret,
				Timeout: int64(config.Settings.JWT.Timeout),
			})

			// 生成token
			token, err := jwtUtil.GenerateToken(user.ID, user.Username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "生成token失败",
					"error":   err.Error(),
				})
				return
			}

			// 获取用户角色信息
			var roles []map[string]interface{}
			if user.RoleID > 0 {
				database.DB.Raw("SELECT id, name, `key` FROM sys_role WHERE id = ?", user.RoleID).Scan(&roles)
			}

			// 获取用户权限信息
			var permissions []string
			if len(roles) > 0 {
				roleID := roles[0]["id"]
				var menuIDs []struct {
					SysMenuID uint64
				}
				database.DB.Raw("SELECT sys_menu_id FROM sys_role_menu WHERE sys_role_id = ?", roleID).Scan(&menuIDs)

				if len(menuIDs) > 0 {
					var perms []struct {
						Permission string `json:"permission"`
					}
					// 构造IN查询参数
					menuIDList := make([]uint64, len(menuIDs))
					for i, item := range menuIDs {
						menuIDList[i] = item.SysMenuID
					}

					database.DB.Raw("SELECT permission FROM sys_menu WHERE id IN (?) AND permission IS NOT NULL AND permission != ''", menuIDList).Scan(&perms)

					for _, perm := range perms {
						permissions = append(permissions, perm.Permission)
					}
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "登录成功",
				"data": gin.H{
					"token": token,
					"user": gin.H{
						"id":          user.ID,
						"username":    user.Username,
						"nick_name":   user.NickName,
						"avatar":      user.Avatar,
						"email":       user.Email,
						"phone":       user.Phone,
						"roles":       roles,
						"permissions": permissions,
					},
				},
			})
		})

		// 退出登录API
		api.POST("/auth/logout", func(c *gin.Context) {
			// 退出登录逻辑：前端清除token即可，后端无需特殊处理
			// 但我们可以在这里添加一些额外的清理逻辑，如记录登出日志等

			// 获取JWT配置
			config, err := loadConfig(configYml)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "服务器配置错误",
				})
				return
			}

			// 从请求头获取token
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"message": "退出成功",
				})
				return
			}

			// 解析token
			tokenString := authHeader[len("Bearer "):]

			// 创建JWT工具
			jwtUtil := utils.NewJWT(utils.JWTConfig{
				Secret:  config.Settings.JWT.Secret,
				Timeout: int64(config.Settings.JWT.Timeout),
			})

			// 解析token（主要用于验证token有效性）
			claims, err := jwtUtil.ParseToken(tokenString)
			if err != nil {
				// 即使token无效，我们也认为退出成功，因为客户端会清除token
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"message": "退出成功",
				})
				return
			}

			// 可以在这里添加登出日志记录
			fmt.Printf("用户 %s (ID: %d) 已退出登录\n", claims.Username, claims.UserID)

			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "退出成功",
			})
		})

		// 获取用户信息API
		api.GET("/auth/userinfo", func(c *gin.Context) {
			// 从请求头获取token
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "未提供认证信息",
				})
				return
			}

			// 解析token
			tokenString := authHeader[len("Bearer "):]

			// 获取JWT配置
			config, err := loadConfig(configYml)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "服务器配置错误",
				})
				return
			}

			// 创建JWT工具
			jwtUtil := utils.NewJWT(utils.JWTConfig{
				Secret:  config.Settings.JWT.Secret,
				Timeout: int64(config.Settings.JWT.Timeout),
			})

			// 解析token
			claims, err := jwtUtil.ParseToken(tokenString)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "token无效或已过期",
					"error":   err.Error(),
				})
				return
			}

			// 根据用户ID获取用户信息
			var user system.SysUser
			result := database.DB.Where("id = ? AND deleted_at IS NULL", claims.UserID).First(&user)
			if result.Error != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "用户不存在",
				})
				return
			}

			// 获取用户角色信息
			var roles []map[string]interface{}
			if user.RoleID > 0 {
				database.DB.Raw("SELECT id, name, `key` FROM sys_role WHERE id = ?", user.RoleID).Scan(&roles)
			}

			// 获取用户权限信息
			var permissions []string
			if len(roles) > 0 {
				roleID := roles[0]["id"]
				var menuIDs []struct {
					SysMenuID uint64
				}
				database.DB.Raw("SELECT sys_menu_id FROM sys_role_menu WHERE sys_role_id = ?", roleID).Scan(&menuIDs)

				if len(menuIDs) > 0 {
					var perms []struct {
						Permission string `json:"permission"`
					}
					// 构造IN查询参数
					menuIDList := make([]uint64, len(menuIDs))
					for i, item := range menuIDs {
						menuIDList[i] = item.SysMenuID
					}

					database.DB.Raw("SELECT permission FROM sys_menu WHERE id IN (?) AND permission IS NOT NULL AND permission != ''", menuIDList).Scan(&perms)

					for _, perm := range perms {
						permissions = append(permissions, perm.Permission)
					}
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"data": gin.H{
					"id":          user.ID,
					"username":    user.Username,
					"nick_name":   user.NickName,
					"avatar":      user.Avatar,
					"email":       user.Email,
					"phone":       user.Phone,
					"roles":       roles,
					"permissions": permissions,
				},
				"message": "获取用户信息成功",
			})
		})

		// 检查token有效性API
		api.GET("/auth/check", func(c *gin.Context) {
			// 从请求头获取token
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "未提供认证信息",
				})
				return
			}

			// 解析token
			tokenString := authHeader[len("Bearer "):]

			// 获取JWT配置
			config, err := loadConfig(configYml)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "服务器配置错误",
				})
				return
			}

			// 创建JWT工具
			jwtUtil := utils.NewJWT(utils.JWTConfig{
				Secret:  config.Settings.JWT.Secret,
				Timeout: int64(config.Settings.JWT.Timeout),
			})

			// 解析token
			claims, err := jwtUtil.ParseToken(tokenString)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "token无效或已过期",
					"error":   err.Error(),
				})
				return
			}

			// 根据用户ID获取用户信息
			var user system.SysUser
			result := database.DB.Where("id = ? AND deleted_at IS NULL", claims.UserID).First(&user)
			if result.Error != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "用户不存在",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "token有效",
				"data": gin.H{
					"user_id":    claims.UserID,
					"username":   claims.Username,
					"expires_at": claims.ExpiresAt.Time.Unix(),
				},
			})
		})

		// 系统信息API
		api.GET("/system/info", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"app_name": "rentpro-admin",
				"version":  global.Version,
				"mode":     gin.Mode(),
			})
		})

		// 楼盘相关API
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
			size, _ := strconv.Atoi(pageSize)

			if pageNum < 1 {
				pageNum = 1
			}
			if size < 1 {
				size = 10
			}
			if size > 100 {
				size = 100 // 限制最大页面大小
			}

			// 构造查询条件
			offset := (pageNum - 1) * size

			// 构造SQL查询
			query := "SELECT id, name, district, business_area, property_type, status, created_at FROM sys_buildings WHERE 1=1"
			countQuery := "SELECT COUNT(*) FROM sys_buildings WHERE 1=1"

			// 添加搜索条件
			var args []interface{}
			if name != "" {
				query += " AND name LIKE ?"
				countQuery += " AND name LIKE ?"
				args = append(args, "%"+name+"%")
			}
			if district != "" {
				query += " AND district LIKE ?"
				countQuery += " AND district LIKE ?"
				args = append(args, "%"+district+"%")
			}
			if businessArea != "" {
				query += " AND business_area LIKE ?"
				countQuery += " AND business_area LIKE ?"
				args = append(args, "%"+businessArea+"%")
			}
			if status != "" {
				query += " AND status = ?"
				countQuery += " AND status = ?"
				args = append(args, status)
			}

			// 添加排序和分页
			query += " ORDER BY id DESC LIMIT ? OFFSET ?"
			args = append(args, size, offset)

			// 执行查询
			var buildings []map[string]interface{}
			database.DB.Raw(query, args...).Scan(&buildings)

			// 查询总数
			var total int64
			database.DB.Raw(countQuery, args[:len(args)-2]...).Scan(&total)

			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "楼盘列表获取成功",
				"data":    buildings,
				"total":   total,
				"page":    pageNum,
				"size":    size,
			})
		})

		// 获取单个楼盘详情
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
				"message": "获取成功",
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
					fmt.Printf("⚠️ 楼盘文件夹初始化失败: %v\n", err)
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

		// 更新楼盘
		api.PUT("/buildings/:id", func(c *gin.Context) {
			id := c.Param("id")

			// 解析请求体
			var buildingData struct {
				Name         string `json:"name"`
				District     string `json:"district"`
				BusinessArea string `json:"businessArea"`
				PropertyType string `json:"propertyType"`
				Status       string `json:"status"`
				Description  string `json:"description"`
			}

			if err := c.ShouldBindJSON(&buildingData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "请求参数错误",
					"error":   err.Error(),
				})
				return
			}

			// 更新数据库
			result := database.DB.Exec(
				"UPDATE sys_buildings SET name = ?, district = ?, business_area = ?, property_type = ?, status = ?, description = ?, updated_at = NOW() WHERE id = ?",
				buildingData.Name,
				buildingData.District,
				buildingData.BusinessArea,
				buildingData.PropertyType,
				buildingData.Status,
				buildingData.Description,
				id,
			)

			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "更新楼盘失败",
					"error":   result.Error.Error(),
				})
				return
			}

			if result.RowsAffected == 0 {
				c.JSON(http.StatusNotFound, gin.H{
					"code":    404,
					"message": "楼盘不存在",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "更新成功",
			})
		})

		// 删除楼盘
		api.DELETE("/buildings/:id", func(c *gin.Context) {
			id := c.Param("id")

			// 删除数据库记录
			result := database.DB.Exec("DELETE FROM sys_buildings WHERE id = ?", id)

			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "删除楼盘失败",
					"error":   result.Error.Error(),
				})
				return
			}

			if result.RowsAffected == 0 {
				c.JSON(http.StatusNotFound, gin.H{
					"code":    404,
					"message": "楼盘不存在",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "删除成功",
			})
		})

		// 获取楼盘的户型列表
		api.GET("/house-types/building/:buildingId", func(c *gin.Context) {
			buildingID := c.Param("buildingId")
			page := c.DefaultQuery("page", "1")
			pageSize := c.DefaultQuery("pageSize", "10")

			// 转换分页参数
			pageNum, err := strconv.Atoi(page)
			if err != nil || pageNum < 1 {
				pageNum = 1
			}
			size, err := strconv.Atoi(pageSize)
			if err != nil || size < 1 {
				size = 10
			}

			// 构造查询
			offset := (pageNum - 1) * size

			query := `
				SELECT 
					id, name, code, description, building_id,
					standard_area, rooms, halls, bathrooms, balconies, floor_height,
					standard_orientation, standard_view,
					base_sale_price, base_rent_price, base_sale_price_per, base_rent_price_per,
					total_stock, available_stock, sold_stock, rented_stock, reserved_stock,
					status, is_hot, main_image, floor_plan_url,
					created_at, updated_at
				FROM sys_house_types 
				WHERE building_id = ? AND deleted_at IS NULL
				ORDER BY id DESC 
				LIMIT ? OFFSET ?`

			countQuery := `
				SELECT COUNT(*) 
				FROM sys_house_types 
				WHERE building_id = ? AND deleted_at IS NULL`

			// 执行查询
			var houseTypes []map[string]interface{}
			database.DB.Raw(query, buildingID, size, offset).Scan(&houseTypes)

			// 查询总数
			var total int64
			database.DB.Raw(countQuery, buildingID).Scan(&total)

			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "户型列表获取成功",
				"data":    houseTypes,
				"total":   total,
				"page":    pageNum,
				"size":    size,
			})
		})

		// 获取楼盘基础信息
		api.GET("/buildings/:id/info", func(c *gin.Context) {
			id := c.Param("id")

			var building map[string]interface{}
			result := database.DB.Raw(`
				SELECT id, name, district, business_area, property_type, 
					   detailed_address, property_company, status, is_hot
				FROM sys_buildings 
				WHERE id = ? AND deleted_at IS NULL`, id).Scan(&building)

			if result.Error != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"code":    404,
					"message": "楼盘不存在",
				})
				return
			}

			if len(building) == 0 {
				c.JSON(http.StatusNotFound, gin.H{
					"code":    404,
					"message": "楼盘不存在",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "获取成功",
				"data":    building,
			})
		})

		// 删除户型
		api.DELETE("/house-types/:id", func(c *gin.Context) {
			id := c.Param("id")

			// 检查是否有关联的房屋
			var houseCount int64
			database.DB.Raw("SELECT COUNT(*) FROM sys_houses WHERE house_type_id = ? AND deleted_at IS NULL", id).Scan(&houseCount)

			if houseCount > 0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": fmt.Sprintf("该户型下还有 %d 套房屋，无法删除", houseCount),
				})
				return
			}

			// 软删除户型
			result := database.DB.Exec("UPDATE sys_house_types SET deleted_at = NOW() WHERE id = ?", id)

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
					"message": "户型不存在",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "删除成功",
			})
		})

		// 新增户型
		api.POST("/house-types", func(c *gin.Context) {
			// 解析请求体
			var houseTypeData struct {
				Name         string  `json:"name" binding:"required,min=1,max=100"`
				Code         string  `json:"code" binding:"required,max=50"`
				StandardArea float64 `json:"standard_area" binding:"required,gt=0"`
				BuildingID   uint    `json:"building_id" binding:"required,gt=0"`

				// 选填字段
				Rooms               *int     `json:"rooms,omitempty"`
				Halls               *int     `json:"halls,omitempty"`
				Bathrooms           *int     `json:"bathrooms,omitempty"`
				Balconies           *int     `json:"balconies,omitempty"`
				FloorHeight         *float64 `json:"floor_height,omitempty"`
				StandardOrientation *string  `json:"standard_orientation,omitempty"`
				StandardView        *string  `json:"standard_view,omitempty"`
				BaseSalePrice       *float64 `json:"base_sale_price,omitempty"`
				BaseRentPrice       *float64 `json:"base_rent_price,omitempty"`
				Description         *string  `json:"description,omitempty"`
				Status              *string  `json:"status,omitempty"`
				IsHot               *bool    `json:"is_hot,omitempty"`
				MainImage           *string  `json:"main_image,omitempty"`
				FloorPlanUrl        *string  `json:"floor_plan_url,omitempty"`
			}

			if err := c.ShouldBindJSON(&houseTypeData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "请求参数错误",
					"error":   err.Error(),
				})
				return
			}

			// 验证楼盘是否存在
			var buildingExists int64
			database.DB.Raw("SELECT COUNT(*) FROM sys_buildings WHERE id = ? AND deleted_at IS NULL", houseTypeData.BuildingID).Scan(&buildingExists)
			if buildingExists == 0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "所选楼盘不存在",
				})
				return
			}

			// 检查户型编码在同一楼盘内是否唯一
			var codeExists int64
			database.DB.Raw("SELECT COUNT(*) FROM sys_house_types WHERE building_id = ? AND code = ? AND deleted_at IS NULL", houseTypeData.BuildingID, houseTypeData.Code).Scan(&codeExists)
			if codeExists > 0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "该楼盘下户型编码已存在",
				})
				return
			}

			// 设置默认值
			rooms := 1
			halls := 1
			bathrooms := 1
			balconies := 0
			status := "active"
			isHot := false

			if houseTypeData.Rooms != nil {
				rooms = *houseTypeData.Rooms
			}
			if houseTypeData.Halls != nil {
				halls = *houseTypeData.Halls
			}
			if houseTypeData.Bathrooms != nil {
				bathrooms = *houseTypeData.Bathrooms
			}
			if houseTypeData.Balconies != nil {
				balconies = *houseTypeData.Balconies
			}
			if houseTypeData.Status != nil {
				status = *houseTypeData.Status
			}
			if houseTypeData.IsHot != nil {
				isHot = *houseTypeData.IsHot
			}

			// 计算单价
			var baseSalePricePer, baseRentPricePer float64
			if houseTypeData.BaseSalePrice != nil {
				baseSalePricePer = *houseTypeData.BaseSalePrice / houseTypeData.StandardArea
			}
			if houseTypeData.BaseRentPrice != nil {
				baseRentPricePer = *houseTypeData.BaseRentPrice / houseTypeData.StandardArea
			}

			// 构造插入SQL
			insertSQL := `
				INSERT INTO sys_house_types (
					name, code, building_id, standard_area, rooms, halls, bathrooms, balconies,
					floor_height, standard_orientation, standard_view,
					base_sale_price, base_rent_price, base_sale_price_per, base_rent_price_per,
					description, status, is_hot, main_image, floor_plan_url,
					total_stock, available_stock, sold_stock, rented_stock, reserved_stock,
					created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, 0, 0, 0, 0, NOW(), NOW())`

			result := database.DB.Exec(insertSQL,
				houseTypeData.Name,
				houseTypeData.Code,
				houseTypeData.BuildingID,
				houseTypeData.StandardArea,
				rooms,
				halls,
				bathrooms,
				balconies,
				houseTypeData.FloorHeight,
				houseTypeData.StandardOrientation,
				houseTypeData.StandardView,
				houseTypeData.BaseSalePrice,
				houseTypeData.BaseRentPrice,
				baseSalePricePer,
				baseRentPricePer,
				houseTypeData.Description,
				status,
				isHot,
				houseTypeData.MainImage,
				houseTypeData.FloorPlanUrl,
			)

			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "创建户型失败",
					"error":   result.Error.Error(),
				})
				return
			}

			// 获取新创建的户型ID
			var newID int64
			database.DB.Raw("SELECT LAST_INSERT_ID()").Scan(&newID)

			// 返回新创建的户型信息
			var newHouseType map[string]interface{}
			database.DB.Raw(`
				SELECT 
					id, name, code, description, building_id,
					standard_area, rooms, halls, bathrooms, balconies, floor_height,
					standard_orientation, standard_view,
					base_sale_price, base_rent_price, base_sale_price_per, base_rent_price_per,
					total_stock, available_stock, sold_stock, rented_stock, reserved_stock,
					status, is_hot, main_image, floor_plan_url,
					created_at, updated_at
				FROM sys_house_types 
				WHERE id = ?`, newID).Scan(&newHouseType)

			c.JSON(http.StatusCreated, gin.H{
				"code":    201,
				"message": "户型创建成功",
				"data":    newHouseType,
			})
		})

		// 上传户型图
		api.POST("/upload/floor-plan", func(c *gin.Context) {
			// 获取用户ID
			userID, exists := c.Get("user_id")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "未授权访问",
				})
				return
			}

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
			img, err := imageManager.UploadBuildingFloorPlan(file, uint64(houseType.BuildingID), houseTypeID, userID.(uint64))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "上传户型图失败",
					"error":   err.Error(),
				})
				return
			}

			// 更新户型的floor_plan_url
			updateResult := database.DB.Model(&houseType).Update("floor_plan_url", img.URL)
			if updateResult.Error != nil {
				// 如果数据库更新失败，删除已上传的文件
				imageManager.DeleteImage(img.ID, userID.(uint64))
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
			return
		})

		// 删除户型图
		api.DELETE("/house-types/:id/floor-plan", func(c *gin.Context) {
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
			if qiniuService != nil {
				// 使用七牛云服务删除文件
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
			} else {
				// 如果七牛云服务不可用，删除本地文件
				if houseType.FloorPlanUrl != "" {
					// 构建文件路径
					filePath := strings.TrimPrefix(houseType.FloorPlanUrl, "/")
					if _, err := os.Stat(filePath); err == nil {
						os.Remove(filePath)
					}
				}
			}

			// 清空数据库中的户型图URL
			updateResult := database.DB.Model(&houseType).Update("floor_plan_url", "")
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
		})

		// 获取区域列表
		api.GET("/districts", func(c *gin.Context) {
			var districts []rental.District
			result := database.DB.Where("status = ?", "active").Order("sort ASC, id ASC").Find(&districts)

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

			var businessAreas []rental.BusinessArea
			query := database.DB.Where("status = ?", "active")

			if districtId != "" {
				query = query.Where("district_id = ?", districtId)
			}

			result := query.Order("sort ASC, id ASC").Find(&businessAreas)

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

		// 七牛云配置测试API
		api.GET("/qiniu/config", func(c *gin.Context) {
			qiniuConfig := config.GetQiniuConfig()
			if qiniuConfig == nil {
				c.JSON(http.StatusOK, gin.H{
					"message": "七牛云配置未加载",
					"config":  nil,
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "七牛云配置加载成功",
				"config": map[string]interface{}{
					"bucket":        qiniuConfig.Bucket,
					"domain":        qiniuConfig.Domain,
					"zone":          qiniuConfig.Zone,
					"use_https":     qiniuConfig.UseHTTPS,
					"use_cdn":       qiniuConfig.UseCdnDomains,
					"max_file_size": qiniuConfig.Upload.MaxFileSize,
					"allowed_types": qiniuConfig.Upload.AllowedTypes,
					"upload_dir":    qiniuConfig.Upload.UploadDir,
				},
			})
		})
	}

	// ===========================
	// 图片管理 API
	// ===========================
	{
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

// fileExists 检查文件是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetAPIVersion 获取API服务器版本信息
func GetAPIVersion() string {
	return global.Version
}
