-- 角色数据初始化
-- 创建时间: 2025-08-28
-- 说明: 系统角色基础数据

INSERT INTO sys_role (id, name, `key`, status, sort, flag, remark, admin, data_scope, params, created_at, updated_at) VALUES 
(1, '超级管理员', 'admin', 1, 1, '', '超级管理员', 1, '1', '', NOW(), NOW()),
(2, '普通用户', 'common', 1, 2, '', '普通用户', 0, '2', '', NOW(), NOW()),
(3, '租客', 'tenant', 1, 3, '', '租客用户', 0, '5', '', NOW(), NOW()),
(4, '房东', 'landlord', 1, 4, '', '房东用户', 0, '4', '', NOW(), NOW());