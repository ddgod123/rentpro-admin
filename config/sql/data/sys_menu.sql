-- 菜单数据初始化
-- 创建时间: 2025-09-02
-- 说明: 系统菜单基础数据 - 与前端原始层级结构保持一致

-- 清除现有菜单数据并重新插入
DELETE FROM sys_role_menu;
DELETE FROM sys_menu;

-- 重置自增ID
ALTER TABLE sys_menu AUTO_INCREMENT = 1;

-- 插入菜单数据 - 与前端原始层级结构完全一致
-- 父级菜单
INSERT INTO sys_menu (id, name, title, icon, path, redirect, component, permission, parent_id, type, sort, visible, is_frame, is_cache, menu_type, status, perms, created_at, updated_at) VALUES 
(1, 'System', '系统管理', 'Setting', '/system', '', 'Layout', 'system:view', 0, 'M', 1, '0', '1', '0', '', '0', '', NOW(), NOW()),
(2, 'Rental', '租赁管理', 'OfficeBuilding', '/rental', '', 'Layout', 'rental:view', 0, 'M', 2, '0', '1', '0', '', '0', '', NOW(), NOW());

-- 系统管理子菜单
INSERT INTO sys_menu (id, name, title, icon, path, redirect, component, permission, parent_id, type, sort, visible, is_frame, is_cache, menu_type, status, perms, created_at, updated_at) VALUES 
(11, 'User', '用户管理', 'User', '/system/user', '', 'system/user/index', 'system:user:view', 1, 'C', 1, '0', '1', '0', '', '0', '', NOW(), NOW()),
(12, 'Role', '角色管理', 'UserFilled', '/system/role', '', 'system/role/index', 'system:role:view', 1, 'C', 2, '0', '1', '0', '', '0', '', NOW(), NOW()),
(13, 'Menu', '菜单管理', 'Menu', '/system/menu', '', 'system/menu/index', 'system:menu:view', 1, 'C', 3, '0', '1', '0', '', '0', '', NOW(), NOW());

-- 租赁管理子菜单
INSERT INTO sys_menu (id, name, title, icon, path, redirect, component, permission, parent_id, type, sort, visible, is_frame, is_cache, menu_type, status, perms, created_at, updated_at) VALUES 
(21, 'Building', '楼盘管理', 'House', '/rental/building', '', 'rental/building/building-management', 'rental:building:view', 2, 'C', 1, '0', '1', '0', '', '0', '', NOW(), NOW()),
(22, 'House', '房屋管理', 'House', '/rental/house', '', 'rental/house/index', 'rental:house:view', 2, 'C', 2, '0', '1', '0', '', '0', '', NOW(), NOW()),
(23, 'Tenant', '租户管理', 'User', '/rental/tenant', '', 'rental/tenant/index', 'rental:tenant:view', 2, 'C', 3, '0', '1', '0', '', '0', '', NOW(), NOW()),
(24, 'Agent', '经纪人管理', 'UserFilled', '/rental/agent', '', 'rental/agent/index', 'rental:agent:view', 2, 'C', 4, '0', '1', '0', '', '0', '', NOW(), NOW()),
(25, 'Landlord', '房东管理', 'UserFilled', '/rental/landlord', '', 'rental/landlord/index', 'rental:landlord:view', 2, 'C', 5, '0', '1', '0', '', '0', '', NOW(), NOW()),
(26, 'Contract', '合同管理', 'Document', '/rental/contract', '', 'rental/contract/index', 'rental:contract:view', 2, 'C', 6, '0', '1', '0', '', '0', '', NOW(), NOW());

-- 重新建立角色菜单关联
-- 超级管理员拥有所有菜单权限
INSERT INTO sys_role_menu (sys_role_id, sys_menu_id) VALUES 
(1, 1), (1, 2), (1, 11), (1, 12), (1, 13), (1, 21), (1, 22), (1, 23), (1, 24), (1, 25), (1, 26);

-- 普通用户只有租赁管理权限
INSERT INTO sys_role_menu (sys_role_id, sys_menu_id) VALUES 
(2, 2), (2, 21), (2, 22), (2, 23), (2, 24), (2, 25), (2, 26);