-- 修复房屋和户型关联关系的数据库迁移脚本
-- 执行时间: 2024年12月
-- 说明: 添加房屋与户型的关联关系，优化数据模型结构

-- ===================================
-- 第一部分：修改 sys_houses 表结构
-- ===================================

-- 1. 为 sys_houses 表添加 house_type_id 字段
ALTER TABLE sys_houses 
ADD COLUMN house_type_id INT UNSIGNED AFTER building_id COMMENT '所属户型ID';

-- 2. 添加索引
ALTER TABLE sys_houses 
ADD INDEX idx_house_type_id (house_type_id);

-- 3. 重命名字段以明确其含义
ALTER TABLE sys_houses 
CHANGE COLUMN area actual_area DECIMAL(8,2) COMMENT '实际建筑面积(平方米)',
CHANGE COLUMN usableArea actual_usable_area DECIMAL(8,2) COMMENT '实际使用面积(平方米)',
CHANGE COLUMN orientation actual_orientation VARCHAR(50) COMMENT '实际朝向',
CHANGE COLUMN view actual_view VARCHAR(100) COMMENT '实际景观',
CHANGE COLUMN salePrice actual_sale_price DECIMAL(12,2) DEFAULT 0 COMMENT '实际售价(元)',
CHANGE COLUMN rentPrice actual_rent_price DECIMAL(8,2) DEFAULT 0 COMMENT '实际月租金(元)';

-- 4. 添加价格调整相关字段
ALTER TABLE sys_houses 
ADD COLUMN price_adjustment DECIMAL(8,2) DEFAULT 0 COMMENT '价格调整金额(元)' AFTER actual_rent_price,
ADD COLUMN price_adjustment_reason VARCHAR(200) COMMENT '价格调整原因' AFTER price_adjustment;

-- 5. 删除冗余字段（这些信息现在由户型提供）
ALTER TABLE sys_houses 
DROP COLUMN IF EXISTS rooms,
DROP COLUMN IF EXISTS halls,
DROP COLUMN IF EXISTS bathrooms,
DROP COLUMN IF EXISTS balconies,
DROP COLUMN IF EXISTS salePricePer,
DROP COLUMN IF EXISTS rentPricePer;

-- ===================================
-- 第二部分：修改 sys_house_types 表结构
-- ===================================

-- 1. 重命名字段以明确其作为基准/标准的含义
ALTER TABLE sys_house_types 
CHANGE COLUMN area standard_area DECIMAL(8,2) NOT NULL COMMENT '标准建筑面积(平方米)',
CHANGE COLUMN orientation standard_orientation VARCHAR(50) COMMENT '标准朝向',
CHANGE COLUMN view standard_view VARCHAR(100) COMMENT '标准景观',
CHANGE COLUMN sale_price base_sale_price DECIMAL(12,2) DEFAULT 0 COMMENT '基准售价(元)',
CHANGE COLUMN rent_price base_rent_price DECIMAL(8,2) DEFAULT 0 COMMENT '基准月租金(元)',
CHANGE COLUMN sale_price_per base_sale_price_per DECIMAL(8,2) DEFAULT 0 COMMENT '基准单价(元/平方米)',
CHANGE COLUMN rent_price_per base_rent_price_per DECIMAL(6,2) DEFAULT 0 COMMENT '基准租金单价(元/平方米/月)';

-- 2. 优化库存字段
ALTER TABLE sys_house_types 
CHANGE COLUMN sale_stock sold_stock INT DEFAULT 0 COMMENT '已售库存',
CHANGE COLUMN rent_stock rented_stock INT DEFAULT 0 COMMENT '已租库存',
ADD COLUMN available_stock INT DEFAULT 0 COMMENT '可用库存' AFTER total_stock;

-- 3. 添加新字段
ALTER TABLE sys_house_types 
ADD COLUMN floor_plan_url VARCHAR(500) COMMENT '户型图URL' AFTER main_image,
ADD COLUMN floor_height DECIMAL(4,2) COMMENT '标准层高(米)' AFTER balconies;

-- 4. 删除不再需要的状态字段
ALTER TABLE sys_house_types 
DROP COLUMN IF EXISTS sale_status,
DROP COLUMN IF EXISTS rent_status;

-- ===================================
-- 第三部分：创建示例户型数据
-- ===================================

-- 插入示例户型数据（假设楼盘ID为1）
INSERT INTO sys_house_types (
    name, code, description, building_id, 
    standard_area, rooms, halls, bathrooms, balconies, floor_height,
    standard_orientation, standard_view,
    base_sale_price, base_rent_price, base_sale_price_per, base_rent_price_per,
    status, created_at, updated_at
) VALUES 
(
    '经典一居', 'A1', '精致一居室，适合单身人士', 1,
    45.5, 1, 1, 1, 1, 2.8,
    '南向', '小区景观',
    800000, 3500, 17582, 77,
    'active', NOW(), NOW()
),
(
    '舒适两居', 'B2', '温馨两居室，适合小家庭', 1,
    78.5, 2, 1, 1, 1, 2.8,
    '南北', '小区景观', 
    1200000, 5500, 15287, 70,
    'active', NOW(), NOW()
),
(
    '宽敞三居', 'C3', '宽敞三居室，适合大家庭', 1,
    108.0, 3, 2, 2, 2, 2.8,
    '南北', '城市景观',
    1800000, 8000, 16667, 74,
    'active', NOW(), NOW()
);

-- ===================================
-- 第四部分：建立现有数据关联
-- ===================================

-- 为现有房屋分配户型（基于面积匹配）
UPDATE sys_houses h 
SET house_type_id = (
    SELECT ht.id 
    FROM sys_house_types ht 
    WHERE ht.building_id = h.building_id 
    AND ABS(ht.standard_area - COALESCE(h.actual_area, 0)) = (
        SELECT MIN(ABS(ht2.standard_area - COALESCE(h.actual_area, 0)))
        FROM sys_house_types ht2 
        WHERE ht2.building_id = h.building_id
    )
    LIMIT 1
)
WHERE h.house_type_id IS NULL OR h.house_type_id = 0;

-- ===================================
-- 第五部分：更新库存统计
-- ===================================

-- 更新户型库存统计
UPDATE sys_house_types ht SET 
    total_stock = (
        SELECT COUNT(*) FROM sys_houses h 
        WHERE h.house_type_id = ht.id AND h.deleted_at IS NULL
    ),
    available_stock = (
        SELECT COUNT(*) FROM sys_houses h 
        WHERE h.house_type_id = ht.id 
        AND h.status = 'available' AND h.deleted_at IS NULL
    ),
    sold_stock = (
        SELECT COUNT(*) FROM sys_houses h 
        WHERE h.house_type_id = ht.id 
        AND h.sale_status = 'sold' AND h.deleted_at IS NULL
    ),
    rented_stock = (
        SELECT COUNT(*) FROM sys_houses h 
        WHERE h.house_type_id = ht.id 
        AND h.rent_status = 'rented' AND h.deleted_at IS NULL
    ),
    reserved_stock = (
        SELECT COUNT(*) FROM sys_houses h 
        WHERE h.house_type_id = ht.id 
        AND (h.sale_status = 'reserved' OR h.rent_status = 'reserved') 
        AND h.deleted_at IS NULL
    );

-- ===================================
-- 第六部分：添加约束和触发器
-- ===================================

-- 添加外键约束（可选，根据数据完整性决定是否启用）
-- ALTER TABLE sys_houses 
-- ADD CONSTRAINT fk_houses_house_type_id 
-- FOREIGN KEY (house_type_id) REFERENCES sys_house_types(id)
-- ON DELETE SET NULL ON UPDATE CASCADE;

-- 创建触发器自动更新库存统计
DELIMITER $$

CREATE TRIGGER trg_update_house_type_stock_after_house_insert
AFTER INSERT ON sys_houses
FOR EACH ROW
BEGIN
    IF NEW.house_type_id IS NOT NULL THEN
        CALL update_house_type_stock(NEW.house_type_id);
    END IF;
END$$

CREATE TRIGGER trg_update_house_type_stock_after_house_update
AFTER UPDATE ON sys_houses
FOR EACH ROW
BEGIN
    -- 更新新户型的库存
    IF NEW.house_type_id IS NOT NULL THEN
        CALL update_house_type_stock(NEW.house_type_id);
    END IF;
    
    -- 如果户型ID发生变化，也要更新旧户型的库存
    IF OLD.house_type_id IS NOT NULL AND OLD.house_type_id != NEW.house_type_id THEN
        CALL update_house_type_stock(OLD.house_type_id);
    END IF;
END$$

CREATE TRIGGER trg_update_house_type_stock_after_house_delete
AFTER DELETE ON sys_houses
FOR EACH ROW
BEGIN
    IF OLD.house_type_id IS NOT NULL THEN
        CALL update_house_type_stock(OLD.house_type_id);
    END IF;
END$$

-- 创建存储过程用于更新库存统计
CREATE PROCEDURE update_house_type_stock(IN house_type_id_param INT UNSIGNED)
BEGIN
    UPDATE sys_house_types SET 
        total_stock = (
            SELECT COUNT(*) FROM sys_houses 
            WHERE house_type_id = house_type_id_param AND deleted_at IS NULL
        ),
        available_stock = (
            SELECT COUNT(*) FROM sys_houses 
            WHERE house_type_id = house_type_id_param 
            AND status = 'available' AND deleted_at IS NULL
        ),
        sold_stock = (
            SELECT COUNT(*) FROM sys_houses 
            WHERE house_type_id = house_type_id_param 
            AND sale_status = 'sold' AND deleted_at IS NULL
        ),
        rented_stock = (
            SELECT COUNT(*) FROM sys_houses 
            WHERE house_type_id = house_type_id_param 
            AND rent_status = 'rented' AND deleted_at IS NULL
        ),
        reserved_stock = (
            SELECT COUNT(*) FROM sys_houses 
            WHERE house_type_id = house_type_id_param 
            AND (sale_status = 'reserved' OR rent_status = 'reserved') 
            AND deleted_at IS NULL
        )
    WHERE id = house_type_id_param;
END$$

DELIMITER ;

-- ===================================
-- 完成信息
-- ===================================

-- 显示迁移完成信息
SELECT 'Database migration completed successfully!' AS message;
SELECT 
    '房屋与户型关联关系修复完成' AS status,
    (SELECT COUNT(*) FROM sys_houses WHERE house_type_id IS NOT NULL) AS linked_houses,
    (SELECT COUNT(*) FROM sys_house_types) AS total_house_types;
