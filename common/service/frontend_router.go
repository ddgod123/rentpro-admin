package service

import (
	"encoding/json"
	"fmt"
	"rentPro/rentpro-admin/common/models/system"
	"strings"
)

// FrontendRoute 前端路由配置
type FrontendRoute struct {
	Path      string            `json:"path"`
	Name      string            `json:"name"`
	Component string            `json:"component"`
	Meta      FrontendRouteMeta `json:"meta"`
	Children  []FrontendRoute   `json:"children,omitempty"`
	Redirect  string            `json:"redirect,omitempty"`
}

// FrontendRouteMeta 前端路由元数据
type FrontendRouteMeta struct {
	Title       string `json:"title"`
	Icon        string `json:"icon"`
	Permission  string `json:"permission"`
	KeepAlive   bool   `json:"keepAlive"`
	RequireAuth bool   `json:"requireAuth"`
	Hidden      bool   `json:"hidden"`
}

// FrontendRouterGenerator 前端路由生成器
type FrontendRouterGenerator struct {
	menuService *MenuService
}

// NewFrontendRouterGenerator 创建前端路由生成器实例
func NewFrontendRouterGenerator(menuService *MenuService) *FrontendRouterGenerator {
	return &FrontendRouterGenerator{
		menuService: menuService,
	}
}

// GenerateUserRoutes 生成用户前端路由
func (g *FrontendRouterGenerator) GenerateUserRoutes(userID uint) ([]FrontendRoute, error) {
	// 获取用户菜单
	menus, err := g.menuService.GetMenusByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户菜单失败: %v", err)
	}

	// 转换为前端路由
	routes := g.convertToFrontendRoutes(menus)

	// 构建路由树
	tree := g.buildRouteTree(routes)

	return tree, nil
}

// GenerateAllRoutes 生成所有前端路由（管理员用）
func (g *FrontendRouterGenerator) GenerateAllRoutes() ([]FrontendRoute, error) {
	// 获取所有菜单
	menus, err := g.menuService.GetAllMenus()
	if err != nil {
		return nil, fmt.Errorf("获取菜单失败: %v", err)
	}

	// 转换为前端路由
	routes := g.convertToFrontendRoutes(menus)

	// 构建路由树
	tree := g.buildRouteTree(routes)

	return tree, nil
}

// GenerateRouteConfig 生成路由配置文件（JSON格式）
func (g *FrontendRouterGenerator) GenerateRouteConfig(userID uint) (string, error) {
	routes, err := g.GenerateUserRoutes(userID)
	if err != nil {
		return "", err
	}

	// 转换为JSON
	config, err := json.MarshalIndent(routes, "", "  ")
	if err != nil {
		return "", fmt.Errorf("序列化路由配置失败: %v", err)
	}

	return string(config), nil
}

// GenerateVueRouterConfig 生成Vue Router配置
func (g *FrontendRouterGenerator) GenerateVueRouterConfig(userID uint) (string, error) {
	routes, err := g.GenerateUserRoutes(userID)
	if err != nil {
		return "", err
	}

	// 生成Vue Router配置代码
	config := g.generateVueRouterCode(routes)

	return config, nil
}

// convertToFrontendRoutes 将菜单转换为前端路由
func (g *FrontendRouterGenerator) convertToFrontendRoutes(menus []system.SysMenu) []FrontendRoute {
	var routes []FrontendRoute

	for _, menu := range menus {
		route := g.convertToFrontendRoute(menu)
		routes = append(routes, route)
	}

	return routes
}

// convertToFrontendRoute 将单个菜单转换为前端路由
func (g *FrontendRouterGenerator) convertToFrontendRoute(menu system.SysMenu) FrontendRoute {
	route := FrontendRoute{
		Path:      menu.Path,
		Name:      menu.Name,
		Component: g.getComponentPath(menu),
		Meta: FrontendRouteMeta{
			Title:       menu.Title,
			Icon:        menu.Icon,
			Permission:  menu.Permission,
			KeepAlive:   menu.IsCache == "0",
			RequireAuth: true,
			Hidden:      menu.Visible == "1",
		},
		Children: []FrontendRoute{},
	}

	// 处理重定向
	if menu.Redirect != "" {
		route.Redirect = menu.Redirect
	}

	// 处理子路由
	if len(menu.Children) > 0 {
		for _, child := range menu.Children {
			childRoute := g.convertToFrontendRoute(child)
			route.Children = append(route.Children, childRoute)
		}
	}

	return route
}

// getComponentPath 获取组件路径
func (g *FrontendRouterGenerator) getComponentPath(menu system.SysMenu) string {
	if menu.Component == "" {
		return ""
	}

	// 如果是外部链接，返回空
	if menu.IsFrame == "0" {
		return ""
	}

	// 处理组件路径
	component := menu.Component
	if !strings.HasPrefix(component, "/") {
		component = "/" + component
	}

	// 添加.vue后缀（如果需要）
	if !strings.HasSuffix(component, ".vue") && !strings.HasSuffix(component, ".js") {
		component += ".vue"
	}

	return component
}

// buildRouteTree 构建路由树
func (g *FrontendRouterGenerator) buildRouteTree(routes []FrontendRoute) []FrontendRoute {
	// 创建路径到路由的映射
	routeMap := make(map[string]*FrontendRoute)
	var rootRoutes []FrontendRoute

	// 初始化映射
	for i := range routes {
		routeMap[routes[i].Path] = &routes[i]
	}

	// 构建树结构
	for i := range routes {
		route := &routes[i]
		parentPath := g.getParentPath(route.Path)

		if parentPath == "" || parentPath == "/" {
			// 根路由
			rootRoutes = append(rootRoutes, *route)
		} else {
			// 子路由
			if parent, exists := routeMap[parentPath]; exists {
				parent.Children = append(parent.Children, *route)
			}
		}
	}

	return rootRoutes
}

// getParentPath 获取父路径
func (g *FrontendRouterGenerator) getParentPath(path string) string {
	if path == "/" || path == "" {
		return ""
	}

	// 移除末尾的斜杠
	path = strings.TrimSuffix(path, "/")

	// 查找最后一个斜杠
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == 0 {
		return "/"
	}

	if lastSlash > 0 {
		return path[:lastSlash]
	}

	return ""
}

// generateVueRouterCode 生成Vue Router配置代码
func (g *FrontendRouterGenerator) generateVueRouterCode(routes []FrontendRoute) string {
	var code strings.Builder

	code.WriteString("// 动态生成的路由配置\n")
	code.WriteString("export const dynamicRoutes = [\n")

	for _, route := range routes {
		g.generateRouteCode(&code, route, 1)
	}

	code.WriteString("];\n\n")

	// 添加路由守卫
	code.WriteString("// 路由守卫\n")
	code.WriteString("export function setupRouteGuard(router) {\n")
	code.WriteString("  router.beforeEach((to, from, next) => {\n")
	code.WriteString("    const token = localStorage.getItem('token');\n")
	code.WriteString("    if (to.meta.requireAuth && !token) {\n")
	code.WriteString("      next('/login');\n")
	code.WriteString("    } else {\n")
	code.WriteString("      next();\n")
	code.WriteString("    }\n")
	code.WriteString("  });\n")
	code.WriteString("}\n")

	return code.String()
}

// generateRouteCode 生成单个路由的代码
func (g *FrontendRouterGenerator) generateRouteCode(code *strings.Builder, route FrontendRoute, indent int) {
	indentStr := strings.Repeat("  ", indent)

	code.WriteString(indentStr + "{\n")
	code.WriteString(indentStr + "  path: '" + route.Path + "',\n")
	code.WriteString(indentStr + "  name: '" + route.Name + "',\n")

	if route.Component != "" {
		code.WriteString(indentStr + "  component: () => import('" + route.Component + "'),\n")
	}

	if route.Redirect != "" {
		code.WriteString(indentStr + "  redirect: '" + route.Redirect + "',\n")
	}

	// 生成meta信息
	code.WriteString(indentStr + "  meta: {\n")
	code.WriteString(indentStr + "    title: '" + route.Meta.Title + "',\n")
	if route.Meta.Icon != "" {
		code.WriteString(indentStr + "    icon: '" + route.Meta.Icon + "',\n")
	}
	if route.Meta.Permission != "" {
		code.WriteString(indentStr + "    permission: '" + route.Meta.Permission + "',\n")
	}
	code.WriteString(indentStr + "    keepAlive: " + fmt.Sprintf("%t", route.Meta.KeepAlive) + ",\n")
	code.WriteString(indentStr + "    requireAuth: " + fmt.Sprintf("%t", route.Meta.RequireAuth) + ",\n")
	code.WriteString(indentStr + "    hidden: " + fmt.Sprintf("%t", route.Meta.Hidden) + ",\n")
	code.WriteString(indentStr + "  },\n")

	// 生成子路由
	if len(route.Children) > 0 {
		code.WriteString(indentStr + "  children: [\n")
		for _, child := range route.Children {
			g.generateRouteCode(code, child, indent+2)
		}
		code.WriteString(indentStr + "  ],\n")
	}

	code.WriteString(indentStr + "},\n")
}
