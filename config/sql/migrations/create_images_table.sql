-- 创建图片管理表
-- 用于管理七牛云存储的图片文件

-- 图片主表
CREATE TABLE IF NOT EXISTS `sys_images` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `name` varchar(200) NOT NULL DEFAULT '' COMMENT '图片名称',
    `description` varchar(500) DEFAULT '' COMMENT '图片描述',
    `file_name` varchar(255) NOT NULL COMMENT '原始文件名',
    `file_size` bigint NOT NULL DEFAULT 0 COMMENT '文件大小(字节)',
    `mime_type` varchar(100) NOT NULL COMMENT 'MIME类型',
    `extension` varchar(10) NOT NULL DEFAULT '' COMMENT '文件扩展名',

    -- 存储路径信息
    `key` varchar(500) NOT NULL COMMENT '七牛云存储Key',
    `url` varchar(1000) NOT NULL COMMENT '原始图片URL',
    `thumbnail_url` varchar(1000) DEFAULT '' COMMENT '缩略图URL',
    `medium_url` varchar(1000) DEFAULT '' COMMENT '中等尺寸URL',
    `large_url` varchar(1000) DEFAULT '' COMMENT '大图URL',

    -- 分类信息
    `category` varchar(50) NOT NULL DEFAULT 'default' COMMENT '图片分类(building/house/avatar/banner等)',
    `module` varchar(50) NOT NULL DEFAULT 'common' COMMENT '所属模块',
    `module_id` bigint unsigned DEFAULT 0 COMMENT '模块关联ID',

    -- 图片属性
    `width` int DEFAULT 0 COMMENT '图片宽度',
    `height` int DEFAULT 0 COMMENT '图片高度',
    `hash` varchar(100) DEFAULT '' COMMENT '文件Hash',

    -- 状态控制
    `is_public` tinyint(1) DEFAULT 1 COMMENT '是否公开访问',
    `is_main` tinyint(1) DEFAULT 0 COMMENT '是否为主图',
    `sort_order` int DEFAULT 0 COMMENT '排序序号',
    `status` varchar(20) DEFAULT 'active' COMMENT '状态(active/inactive/deleted)',

    -- 审计字段
    `created_by` bigint unsigned DEFAULT 0 COMMENT '创建者ID',
    `updated_by` bigint unsigned DEFAULT 0 COMMENT '更新者ID',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',

    PRIMARY KEY (`id`),
    KEY `idx_category` (`category`),
    KEY `idx_module` (`module`, `module_id`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='图片管理表';

-- 图片分类配置表
CREATE TABLE IF NOT EXISTS `sys_image_categories` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `code` varchar(50) NOT NULL COMMENT '分类编码',
    `name` varchar(100) NOT NULL COMMENT '分类名称',
    `description` varchar(200) DEFAULT '' COMMENT '分类描述',
    `max_size` bigint DEFAULT 5242880 COMMENT '最大文件大小',
    `allowed_types` json COMMENT '允许的文件类型',
    `max_count` int DEFAULT 10 COMMENT '最大上传数量',
    `is_required` tinyint(1) DEFAULT 0 COMMENT '是否必填',
    `status` varchar(20) DEFAULT 'active' COMMENT '状态',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='图片分类配置表';

-- 插入默认的图片分类配置
INSERT INTO `sys_image_categories` (`code`, `name`, `description`, `max_size`, `allowed_types`, `max_count`, `is_required`, `status`) VALUES
('building', '楼盘图片', '楼盘相关的图片文件', 5242880, '["image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"]', 20, 0, 'active'),
('house', '房屋图片', '房屋相关的图片文件', 5242880, '["image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"]', 15, 0, 'active'),
('avatar', '头像图片', '用户头像图片', 2097152, '["image/jpeg", "image/jpg", "image/png"]', 1, 0, 'active'),
('banner', '横幅图片', '网站横幅和广告图片', 3145728, '["image/jpeg", "image/jpg", "image/png"]', 5, 0, 'active'),
('floor_plan', '户型图', '房屋户型图纸', 5242880, '["image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"]', 10, 0, 'active'),
('certificate', '证件图片', '身份证、营业执照等证件', 2097152, '["image/jpeg", "image/jpg", "image/png"]', 5, 0, 'active'),
('default', '默认分类', '未分类的图片文件', 5242880, '["image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"]', 10, 0, 'active');
