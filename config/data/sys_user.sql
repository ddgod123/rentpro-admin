-- 用户数据初始化
-- 创建时间: 2025-08-28
-- 说明: 系统用户基础数据
-- 注意: 密码字段将由系统自动加密处理

INSERT INTO sys_user (id, username, password, nick_name, avatar, email, phone, status, is_admin, remark, dept_id, post_id, role_id, salt, last_login_ip, last_login_at, created_at, updated_at) VALUES 
(1, 'admin', '123456', '超级管理员', '', 'admin@rentpro.com', '15888888888', 1, 1, '超级管理员账号', 1, 1, 1, '', '', NULL, NOW(), NOW()),
(2, 'test', '123456', '测试用户', '', 'test@rentpro.com', '15666666666', 1, 0, '测试用户账号', 2, 4, 2, '', '', NULL, NOW(), NOW());