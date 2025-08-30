-- 部门数据初始化
-- 创建时间: 2025-08-28
-- 说明: 系统部门基础数据

INSERT INTO sys_dept (id, parent_id, dept_path, dept_name, sort, leader, phone, email, status, created_at, updated_at) VALUES 
(1, 0, '0,1', 'RentPro科技', 1, '系统管理员', '15888888888', 'admin@rentpro.com', '0', NOW(), NOW()),
(2, 1, '0,1,2', '技术部', 1, '技术总监', '15666666666', 'tech@rentpro.com', '0', NOW(), NOW()),
(3, 1, '0,1,3', '运营部', 2, '运营总监', '15777777777', 'ops@rentpro.com', '0', NOW(), NOW());