-- ================================================
-- 城市、区域、商圈数据重构脚本
-- 创建时间: 2025-09-05
-- 说明: 重新设计城市-区域-商圈三级级联数据结构
-- ================================================

-- 1. 创建城市表
DROP TABLE IF EXISTS sys_cities;
CREATE TABLE sys_cities (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE COMMENT '城市代码',
    name VARCHAR(50) NOT NULL COMMENT '城市名称',
    sort BIGINT DEFAULT 0 COMMENT '排序',
    status VARCHAR(20) DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='城市表';

-- 2. 插入城市数据
INSERT INTO sys_cities (id, code, name, sort, status) VALUES
(1, 'BJ', '北京', 1, 'active'),
(2, 'SH', '上海', 2, 'active'),
(3, 'GZ', '广州', 3, 'active'),
(4, 'SZ', '深圳', 4, 'active');

-- 3. 修改区域表结构，添加城市ID关联
ALTER TABLE sys_districts ADD COLUMN city_id BIGINT UNSIGNED AFTER city_code;
ALTER TABLE sys_districts ADD FOREIGN KEY fk_districts_city (city_id) REFERENCES sys_cities(id);

-- 4. 更新现有北京区域数据的城市ID
UPDATE sys_districts SET city_id = 1 WHERE city_code = 'BJ';

-- 5. 清空现有数据，重新插入完整的四城市区域数据
DELETE FROM sys_business_areas;
DELETE FROM sys_districts;
ALTER TABLE sys_business_areas AUTO_INCREMENT = 1;
ALTER TABLE sys_districts AUTO_INCREMENT = 1;

-- 6. 插入四个城市的区域数据
INSERT INTO sys_districts (id, code, name, city_code, city_id, sort, status) VALUES
-- 北京区域
(1, 'BJ001', '朝阳区', 'BJ', 1, 1, 'active'),
(2, 'BJ002', '海淀区', 'BJ', 1, 2, 'active'),
(3, 'BJ003', '西城区', 'BJ', 1, 3, 'active'),
(4, 'BJ004', '东城区', 'BJ', 1, 4, 'active'),
(5, 'BJ005', '丰台区', 'BJ', 1, 5, 'active'),
(6, 'BJ006', '石景山区', 'BJ', 1, 6, 'active'),
(7, 'BJ007', '通州区', 'BJ', 1, 7, 'active'),
(8, 'BJ008', '昌平区', 'BJ', 1, 8, 'active'),
(9, 'BJ009', '大兴区', 'BJ', 1, 9, 'active'),
(10, 'BJ010', '房山区', 'BJ', 1, 10, 'active'),

-- 上海区域
(11, 'SH001', '黄浦区', 'SH', 2, 1, 'active'),
(12, 'SH002', '徐汇区', 'SH', 2, 2, 'active'),
(13, 'SH003', '长宁区', 'SH', 2, 3, 'active'),
(14, 'SH004', '静安区', 'SH', 2, 4, 'active'),
(15, 'SH005', '普陀区', 'SH', 2, 5, 'active'),
(16, 'SH006', '虹口区', 'SH', 2, 6, 'active'),
(17, 'SH007', '杨浦区', 'SH', 2, 7, 'active'),
(18, 'SH008', '浦东新区', 'SH', 2, 8, 'active'),
(19, 'SH009', '闵行区', 'SH', 2, 9, 'active'),
(20, 'SH010', '宝山区', 'SH', 2, 10, 'active'),
(21, 'SH011', '嘉定区', 'SH', 2, 11, 'active'),
(22, 'SH012', '松江区', 'SH', 2, 12, 'active'),

-- 广州区域
(23, 'GZ001', '越秀区', 'GZ', 3, 1, 'active'),
(24, 'GZ002', '荔湾区', 'GZ', 3, 2, 'active'),
(25, 'GZ003', '海珠区', 'GZ', 3, 3, 'active'),
(26, 'GZ004', '天河区', 'GZ', 3, 4, 'active'),
(27, 'GZ005', '白云区', 'GZ', 3, 5, 'active'),
(28, 'GZ006', '黄埔区', 'GZ', 3, 6, 'active'),
(29, 'GZ007', '番禺区', 'GZ', 3, 7, 'active'),
(30, 'GZ008', '花都区', 'GZ', 3, 8, 'active'),
(31, 'GZ009', '南沙区', 'GZ', 3, 9, 'active'),
(32, 'GZ010', '从化区', 'GZ', 3, 10, 'active'),
(33, 'GZ011', '增城区', 'GZ', 3, 11, 'active'),

-- 深圳区域
(34, 'SZ001', '福田区', 'SZ', 4, 1, 'active'),
(35, 'SZ002', '罗湖区', 'SZ', 4, 2, 'active'),
(36, 'SZ003', '南山区', 'SZ', 4, 3, 'active'),
(37, 'SZ004', '盐田区', 'SZ', 4, 4, 'active'),
(38, 'SZ005', '宝安区', 'SZ', 4, 5, 'active'),
(39, 'SZ006', '龙岗区', 'SZ', 4, 6, 'active'),
(40, 'SZ007', '龙华区', 'SZ', 4, 7, 'active'),
(41, 'SZ008', '坪山区', 'SZ', 4, 8, 'active'),
(42, 'SZ009', '光明区', 'SZ', 4, 9, 'active'),
(43, 'SZ010', '大鹏新区', 'SZ', 4, 10, 'active');

-- 7. 插入四个城市的商圈数据
INSERT INTO sys_business_areas (id, code, name, district_id, city_code, sort, status) VALUES
-- 北京商圈 (朝阳区)
(1, 'BJ001001', '国贸商圈', 1, 'BJ', 1, 'active'),
(2, 'BJ001002', '三里屯商圈', 1, 'BJ', 2, 'active'),
(3, 'BJ001003', '望京商圈', 1, 'BJ', 3, 'active'),
(4, 'BJ001004', '亚运村商圈', 1, 'BJ', 4, 'active'),
(5, 'BJ001005', 'CBD商圈', 1, 'BJ', 5, 'active'),

-- 北京商圈 (海淀区)
(6, 'BJ002001', '中关村商圈', 2, 'BJ', 1, 'active'),
(7, 'BJ002002', '五道口商圈', 2, 'BJ', 2, 'active'),
(8, 'BJ002003', '西二旗商圈', 2, 'BJ', 3, 'active'),
(9, 'BJ002004', '上地商圈', 2, 'BJ', 4, 'active'),
(10, 'BJ002005', '万柳商圈', 2, 'BJ', 5, 'active'),

-- 北京商圈 (西城区)
(11, 'BJ003001', '金融街商圈', 3, 'BJ', 1, 'active'),
(12, 'BJ003002', '西单商圈', 3, 'BJ', 2, 'active'),
(13, 'BJ003003', '什刹海商圈', 3, 'BJ', 3, 'active'),
(14, 'BJ003004', '德胜门商圈', 3, 'BJ', 4, 'active'),

-- 北京商圈 (东城区)
(15, 'BJ004001', '王府井商圈', 4, 'BJ', 1, 'active'),
(16, 'BJ004002', '东单商圈', 4, 'BJ', 2, 'active'),
(17, 'BJ004003', '前门商圈', 4, 'BJ', 3, 'active'),
(18, 'BJ004004', '崇文门商圈', 4, 'BJ', 4, 'active'),

-- 上海商圈 (黄浦区)
(19, 'SH001001', '外滩商圈', 11, 'SH', 1, 'active'),
(20, 'SH001002', '南京路商圈', 11, 'SH', 2, 'active'),
(21, 'SH001003', '人民广场商圈', 11, 'SH', 3, 'active'),
(22, 'SH001004', '豫园商圈', 11, 'SH', 4, 'active'),

-- 上海商圈 (徐汇区)
(23, 'SH002001', '徐家汇商圈', 12, 'SH', 1, 'active'),
(24, 'SH002002', '衡山路商圈', 12, 'SH', 2, 'active'),
(25, 'SH002003', '田子坊商圈', 12, 'SH', 3, 'active'),

-- 上海商圈 (长宁区)
(26, 'SH003001', '中山公园商圈', 13, 'SH', 1, 'active'),
(27, 'SH003002', '古北商圈', 13, 'SH', 2, 'active'),

-- 上海商圈 (静安区)
(28, 'SH004001', '静安寺商圈', 14, 'SH', 1, 'active'),
(29, 'SH004002', '南京西路商圈', 14, 'SH', 2, 'active'),

-- 上海商圈 (浦东新区)
(30, 'SH008001', '陆家嘴商圈', 18, 'SH', 1, 'active'),
(31, 'SH008002', '张江商圈', 18, 'SH', 2, 'active'),
(32, 'SH008003', '金桥商圈', 18, 'SH', 3, 'active'),
(33, 'SH008004', '世纪公园商圈', 18, 'SH', 4, 'active'),

-- 广州商圈 (天河区)
(34, 'GZ004001', '天河城商圈', 26, 'GZ', 1, 'active'),
(35, 'GZ004002', '珠江新城商圈', 26, 'GZ', 2, 'active'),
(36, 'GZ004003', '体育中心商圈', 26, 'GZ', 3, 'active'),

-- 广州商圈 (越秀区)
(37, 'GZ001001', '北京路商圈', 23, 'GZ', 1, 'active'),
(38, 'GZ001002', '环市东商圈', 23, 'GZ', 2, 'active'),

-- 广州商圈 (海珠区)
(39, 'GZ003001', '江南西商圈', 25, 'GZ', 1, 'active'),
(40, 'GZ003002', '琶洲商圈', 25, 'GZ', 2, 'active'),

-- 广州商圈 (荔湾区)
(41, 'GZ002001', '上下九商圈', 24, 'GZ', 1, 'active'),
(42, 'GZ002002', '陈家祠商圈', 24, 'GZ', 2, 'active'),

-- 深圳商圈 (福田区)
(43, 'SZ001001', '华强北商圈', 34, 'SZ', 1, 'active'),
(44, 'SZ001002', '中心区商圈', 34, 'SZ', 2, 'active'),
(45, 'SZ001003', '车公庙商圈', 34, 'SZ', 3, 'active'),

-- 深圳商圈 (南山区)
(46, 'SZ003001', '科技园商圈', 36, 'SZ', 1, 'active'),
(47, 'SZ003002', '蛇口商圈', 36, 'SZ', 2, 'active'),
(48, 'SZ003003', '后海商圈', 36, 'SZ', 3, 'active'),

-- 深圳商圈 (罗湖区)
(49, 'SZ002001', '东门商圈', 35, 'SZ', 1, 'active'),
(50, 'SZ002002', '国贸商圈', 35, 'SZ', 2, 'active'),

-- 深圳商圈 (宝安区)
(51, 'SZ005001', '宝安中心商圈', 38, 'SZ', 1, 'active'),
(52, 'SZ005002', '西乡商圈', 38, 'SZ', 2, 'active');

-- 8. 更新楼盘表的城市字段为城市ID关联（可选，如果需要的话）
-- ALTER TABLE sys_buildings ADD COLUMN city_id BIGINT UNSIGNED AFTER city;
-- ALTER TABLE sys_buildings ADD FOREIGN KEY fk_buildings_city (city_id) REFERENCES sys_cities(id);

-- 完成提示
SELECT '✅ 城市、区域、商圈数据重构完成！' as message;
SELECT '📊 数据统计:' as message;
SELECT COUNT(*) as '城市数量' FROM sys_cities;
SELECT COUNT(*) as '区域数量' FROM sys_districts;
SELECT COUNT(*) as '商圈数量' FROM sys_business_areas;
