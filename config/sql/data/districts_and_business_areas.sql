-- 区域和商圈数据初始化
-- 创建时间: 2025-09-02
-- 说明: 为楼盘管理功能提供区域和商圈基础数据

-- 清除现有数据
DELETE FROM sys_business_areas;
DELETE FROM sys_districts;

-- 重置自增ID
ALTER TABLE sys_business_areas AUTO_INCREMENT = 1;
ALTER TABLE sys_districts AUTO_INCREMENT = 1;

-- 插入区域数据
INSERT INTO sys_districts (id, code, name, city_code, sort, status) VALUES
(1, 'BJ001', '朝阳区', 'BJ', 1, 'active'),
(2, 'BJ002', '海淀区', 'BJ', 2, 'active'),
(3, 'BJ003', '西城区', 'BJ', 3, 'active'),
(4, 'BJ004', '东城区', 'BJ', 4, 'active'),
(5, 'BJ005', '丰台区', 'BJ', 5, 'active'),
(6, 'BJ006', '石景山区', 'BJ', 6, 'active'),
(7, 'BJ007', '通州区', 'BJ', 7, 'active'),
(8, 'BJ008', '昌平区', 'BJ', 8, 'active'),
(9, 'BJ009', '大兴区', 'BJ', 9, 'active'),
(10, 'BJ010', '房山区', 'BJ', 10, 'active');

-- 插入商圈数据
INSERT INTO sys_business_areas (id, code, name, district_id, city_code, sort, status) VALUES
-- 朝阳区商圈
(1, 'BJ001001', '国贸商圈', 1, 'BJ', 1, 'active'),
(2, 'BJ001002', '三里屯商圈', 1, 'BJ', 2, 'active'),
(3, 'BJ001003', '望京商圈', 1, 'BJ', 3, 'active'),
(4, 'BJ001004', '亚运村商圈', 1, 'BJ', 4, 'active'),
(5, 'BJ001005', 'CBD商圈', 1, 'BJ', 5, 'active'),

-- 海淀区商圈
(6, 'BJ002001', '中关村商圈', 2, 'BJ', 1, 'active'),
(7, 'BJ002002', '五道口商圈', 2, 'BJ', 2, 'active'),
(8, 'BJ002003', '西二旗商圈', 2, 'BJ', 3, 'active'),
(9, 'BJ002004', '上地商圈', 2, 'BJ', 4, 'active'),
(10, 'BJ002005', '万柳商圈', 2, 'BJ', 5, 'active'),

-- 西城区商圈
(11, 'BJ003001', '金融街商圈', 3, 'BJ', 1, 'active'),
(12, 'BJ003002', '西单商圈', 3, 'BJ', 2, 'active'),
(13, 'BJ003003', '什刹海商圈', 3, 'BJ', 3, 'active'),
(14, 'BJ003004', '德胜门商圈', 3, 'BJ', 4, 'active'),

-- 东城区商圈
(15, 'BJ004001', '王府井商圈', 4, 'BJ', 1, 'active'),
(16, 'BJ004002', '东单商圈', 4, 'BJ', 2, 'active'),
(17, 'BJ004003', '前门商圈', 4, 'BJ', 3, 'active'),
(18, 'BJ004004', '崇文门商圈', 4, 'BJ', 4, 'active'),

-- 丰台区商圈
(19, 'BJ005001', '丽泽商圈', 5, 'BJ', 1, 'active'),
(20, 'BJ005002', '丰台科技园商圈', 5, 'BJ', 2, 'active'),
(21, 'BJ005003', '方庄商圈', 5, 'BJ', 3, 'active'),

-- 石景山区商圈
(22, 'BJ006001', '石景山万达商圈', 6, 'BJ', 1, 'active'),
(23, 'BJ006002', '八角商圈', 6, 'BJ', 2, 'active'),

-- 通州区商圈
(24, 'BJ007001', '通州万达商圈', 7, 'BJ', 1, 'active'),
(25, 'BJ007002', '运河商务区商圈', 7, 'BJ', 2, 'active'),

-- 昌平区商圈
(26, 'BJ008001', '回龙观商圈', 8, 'BJ', 1, 'active'),
(27, 'BJ008002', '天通苑商圈', 8, 'BJ', 2, 'active'),

-- 大兴区商圈
(28, 'BJ009001', '亦庄商圈', 9, 'BJ', 1, 'active'),
(29, 'BJ009002', '大兴新城商圈', 9, 'BJ', 2, 'active'),

-- 房山区商圈
(30, 'BJ010001', '良乡商圈', 10, 'BJ', 1, 'active'),
(31, 'BJ010002', '燕山商圈', 10, 'BJ', 2, 'active');
