package service

import (
	"errors"
	"log"
	"sort"
	"time"

	"rentPro/rentpro-admin/common/models/system"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// JWTClaims JWT声明结构
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	RoleID   uint   `json:"role_id"`
	RoleKey  string `json:"role_key"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// AuthService 认证服务
type AuthService struct {
	db        *gorm.DB
	jwtSecret string
	jwtExpiry time.Duration
}

// NewAuthService 创建认证服务实例
func NewAuthService(db *gorm.DB, jwtSecret string, jwtExpiry time.Duration) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Token    string            `json:"token"`
	UserInfo *UserInfoResponse `json:"user_info"`
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	ID          uint           `json:"id"`
	Username    string         `json:"username"`
	NickName    string         `json:"nick_name"`
	Avatar      string         `json:"avatar"`
	Email       string         `json:"email"`
	Phone       string         `json:"phone"`
	IsAdmin     bool           `json:"is_admin"`
	Role        *RoleResponse  `json:"role,omitempty"`
	Dept        *DeptResponse  `json:"dept,omitempty"`
	Permissions []string       `json:"permissions"`
	Menus       []MenuResponse `json:"menus"`
}

// RoleResponse 角色响应
type RoleResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Key   string `json:"key"`
	Admin bool   `json:"admin"`
}

// DeptResponse 部门响应
type DeptResponse struct {
	ID       uint   `json:"id"`
	DeptName string `json:"dept_name"`
	DeptPath string `json:"dept_path"`
}

// MenuResponse 菜单响应
type MenuResponse struct {
	ID         uint           `json:"id"`
	Name       string         `json:"name"`
	Title      string         `json:"title"`
	Icon       string         `json:"icon"`
	Path       string         `json:"path"`
	Component  string         `json:"component"`
	Permission string         `json:"permission"`
	ParentID   uint           `json:"parent_id"`
	Type       string         `json:"type"`
	Sort       int            `json:"sort"`
	Visible    string         `json:"visible"`
	Children   []MenuResponse `json:"children,omitempty"`
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 查找用户
	var user system.SysUser
	err := s.db.Preload("Role").Preload("Dept").
		Where("username = ? AND status = 1", req.Username).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if !user.ComparePassword(req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if !user.IsActive() {
		return nil, errors.New("用户已被禁用")
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	s.db.Save(&user)

	// 生成JWT Token
	token, err := s.generateToken(&user)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	userInfo, err := s.GetUserInfo(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:    token,
		UserInfo: userInfo,
	}, nil
}

// GetUserInfo 获取用户详细信息
func (s *AuthService) GetUserInfo(userID uint) (*UserInfoResponse, error) {
	var user system.SysUser
	err := s.db.Preload("Role").Preload("Dept").
		Where("id = ? AND status = 1", userID).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	userInfo := &UserInfoResponse{
		ID:       user.ID,
		Username: user.Username,
		NickName: user.NickName,
		Avatar:   user.Avatar,
		Email:    user.Email,
		Phone:    user.Phone,
		IsAdmin:  user.IsAdmin,
	}

	// 角色信息
	if user.Role != nil {
		userInfo.Role = &RoleResponse{
			ID:    user.Role.ID,
			Name:  user.Role.Name,
			Key:   user.Role.Key,
			Admin: user.Role.Admin,
		}
	}

	// 部门信息
	if user.Dept != nil {
		userInfo.Dept = &DeptResponse{
			ID:       user.Dept.ID,
			DeptName: user.Dept.DeptName,
			DeptPath: user.Dept.DeptPath,
		}
	}

	// 获取用户权限和菜单
	permissions, menus, err := s.getUserPermissionsAndMenus(user.RoleID)
	if err != nil {
		return nil, err
	}

	userInfo.Permissions = permissions
	userInfo.Menus = menus

	return userInfo, nil
}

// generateToken 生成JWT Token
func (s *AuthService) generateToken(user *system.SysUser) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RoleID:   user.RoleID,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// 添加角色信息
	if user.Role != nil {
		claims.RoleKey = user.Role.Key
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ValidateToken 验证JWT Token
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	// 检查token是否在黑名单中（简单实现，实际项目中应该使用Redis）
	if s.isTokenBlacklisted(tokenString) {
		return nil, errors.New("token已被注销")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// isTokenBlacklisted 检查token是否在黑名单中
func (s *AuthService) isTokenBlacklisted(tokenString string) bool {
	// 这里应该检查Redis或其他存储中的黑名单
	// 暂时返回false，实际项目中需要实现
	return false
}

// BlacklistToken 将token加入黑名单
func (s *AuthService) BlacklistToken(tokenString string) error {
	// 这里应该将token加入Redis或其他存储的黑名单
	// 暂时只记录日志，实际项目中需要实现
	log.Printf("Token已加入黑名单: %s", tokenString[:10]+"...")
	return nil
}

// getUserPermissionsAndMenus 获取用户权限和菜单
func (s *AuthService) getUserPermissionsAndMenus(roleID uint) ([]string, []MenuResponse, error) {
	var role system.SysRole
	err := s.db.Preload("Menus", "status = '0'").
		Where("id = ? AND status = 1", roleID).
		First(&role).Error

	if err != nil {
		return nil, nil, err
	}

	// 获取所有有效的菜单
	var allMenus []system.SysMenu
	err = s.db.Where("status = '0'").Order("sort ASC, id ASC").Find(&allMenus).Error
	if err != nil {
		return nil, nil, err
	}

	// 构建菜单ID映射，用于快速查找
	menuIDMap := make(map[uint]system.SysMenu)
	for _, menu := range allMenus {
		menuIDMap[menu.ID] = menu
	}

	// 如果是管理员角色，获取所有菜单
	var menus []system.SysMenu
	if role.IsAdmin() {
		// 管理员获取所有菜单
		menus = allMenus
	} else {
		// 非管理员角色，只获取有权限的菜单
		menus = role.Menus
	}

	// 构建权限列表
	permissions := make([]string, 0)
	for _, menu := range menus {
		if menu.Permission != "" {
			permissions = append(permissions, menu.Permission)
		}
	}

	// 将所有菜单转换为MenuResponse格式（所有菜单都是一级）
	menuResponses := make([]MenuResponse, 0)
	for _, menu := range menus {
		menuResp := MenuResponse{
			ID:         menu.ID,
			Name:       menu.Name,
			Title:      menu.Title,
			Icon:       menu.Icon,
			Path:       menu.Path,
			Component:  menu.Component,
			Permission: menu.Permission,
			ParentID:   0, // 强制设置为0，表示所有菜单都是一级
			Type:       menu.Type,
			Sort:       menu.Sort,
			Visible:    menu.Visible,
			Children:   make([]MenuResponse, 0), // 空数组，没有子菜单
		}
		menuResponses = append(menuResponses, menuResp)
	}

	// 对菜单排序
	sort.Slice(menuResponses, func(i, j int) bool {
		return menuResponses[i].Sort < menuResponses[j].Sort
	})

	return permissions, menuResponses, nil
}

// HasPermission 检查用户是否有指定权限
func (s *AuthService) HasPermission(userID uint, permission string) (bool, error) {
	var user system.SysUser
	err := s.db.Preload("Role.Menus").Where("id = ?", userID).First(&user).Error
	if err != nil {
		return false, err
	}

	// 管理员拥有所有权限
	if user.IsAdmin || (user.Role != nil && user.Role.IsAdmin()) {
		return true, nil
	}

	// 检查角色权限
	if user.Role != nil {
		return user.Role.HasPermission(permission), nil
	}

	return false, nil
}
