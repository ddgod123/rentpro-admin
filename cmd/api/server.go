// Package api 提供HTTP API服务器相关的命令行功能
// 用于启动 rentpro-admin 系统的HTTP API服务
package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"rentPro/rentpro-admin/common/database"
	"rentPro/rentpro-admin/common/global"
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
				Name         string `json:"name" binding:"required"`
				District     string `json:"district" binding:"required"`
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

			// 插入数据库
			result := database.DB.Exec(
				"INSERT INTO sys_buildings (name, district, business_area, property_type, status, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())",
				buildingData.Name,
				buildingData.District,
				buildingData.BusinessArea,
				buildingData.PropertyType,
				buildingData.Status,
				buildingData.Description,
			)

			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "创建楼盘失败",
					"error":   result.Error.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "创建成功",
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
