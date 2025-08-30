-- 岗位数据初始化
-- 创建时间: 2025-08-28
-- 说明: 系统岗位基础数据

INSERT INTO sys_post (id, post_code, post_name, sort, status, remark, created_at, updated_at) VALUES 
(1, 'ceo', '董事长', 1, '0', '董事长', NOW(), NOW()),
(2, 'se', '项目经理', 2, '0', '项目经理', NOW(), NOW()),
(3, 'hr', '人力资源', 3, '0', '人力资源', NOW(), NOW()),
(4, 'user', '普通员工', 4, '0', '普通员工', NOW(), NOW());