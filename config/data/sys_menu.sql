-- 菜单数据初始化
-- 创建时间: 2025-08-29
-- 说明: 系统菜单基础数据 - 完整版本（所有菜单为同一级别）

-- 清除现有菜单数据并重新插入
DELETE FROM sys_role_menu WHERE sys_menu_id IN (1,2,3,4,5,6,7);
DELETE FROM sys_menu WHERE id IN (1,2,3,4,5,6,7);

-- 所有菜单为同一级别（parent_id = 0）
INSERT INTO sys_menu (id, name, title, icon, path, redirect, component, permission, parent_id, type, sort, visible, is_frame, is_cache, menu_type, status, perms, created_at, updated_at) VALUES 
(1, 'Dashboard', '仪表板', 'dashboard', '/dashboard', '', 'dashboard/index', 'dashboard:view', 0, 'C', 1, '0', '1', '0', '', '0', '', NOW(), NOW()),
(4, 'Building', '楼盘管理', 'building', '/rental/building', '', 'rental/building/index', 'rental:building:view', 0, 'C', 2, '0', '1', '0', '', '0', '', NOW(), NOW()),
(5, 'Room', '房源管理', 'room', '/rental/room', '', 'rental/room/index', 'rental:room:view', 0, 'C', 3, '0', '1', '0', '', '0', '', NOW(), NOW()),
(6, 'Tenant', '租客管理', 'tenant', '/rental/tenant', '', 'rental/tenant/index', 'rental:tenant:view', 0, 'C', 4, '0', '1', '0', '', '0', '', NOW(), NOW()),
(7, 'Contract', '合同管理', 'contract', '/rental/contract', '', 'rental/contract/index', 'rental:contract:view', 0, 'C', 5, '0', '1', '0', '', '0', '', NOW(), NOW()),
(2, 'User', '用户管理', 'user', '/system/user', '', 'system/user/index', 'system:user:view', 0, 'C', 6, '0', '1', '0', '', '0', '', NOW(), NOW()),
(3, 'Role', '角色管理', 'role', '/system/role', '', 'system/role/index', 'system:role:view', 0, 'C', 7, '0', '1', '0', '', '0', '', NOW(), NOW());

-- 重新建立角色菜单关联
-- 超级管理员拥有所有菜单权限
INSERT INTO sys_role_menu (sys_role_id, sys_menu_id) VALUES 
(1, 1), (1, 2), (1, 3),
(1, 4), (1, 5), (1, 6), (1, 7);

-- 普通用户只有基本菜单权限（仪表板和租赁管理）
INSERT INTO sys_role_menu (sys_role_id, sys_menu_id) VALUES 
(2, 1),
(2, 4), (2, 5), (2, 6), (2, 7);