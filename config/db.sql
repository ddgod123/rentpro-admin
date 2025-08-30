-- 楼盘数据初始化
-- 创建时间: 2025-08-28
-- 说明: 系统楼盘基础数据

INSERT INTO sys_buildings (id, name, developer, address, city, district, sub_district, property_type, property_company, description, status, is_hot, created_by, updated_by, created_at, updated_at, deleted_at) VALUES 
(1, '滨江一号', '滨江地产', '上海市浦东新区张江高科技园区博云路2号', '上海市', '浦东新区', '张江镇', '住宅', '滨江物业', '滨江一号是滨江地产打造的高品质住宅项目，位于浦东新区核心位置，配套设施完善。', 'active', 1, 'admin', 'admin', NOW(), NOW(), NULL),
(2, '城市之光', '城市发展', '北京市朝阳区建国路88号', '北京市', '朝阳区', '建外街道', '商业', '城市物业', '城市之光是集购物、餐饮、娱乐于一体的综合性商业中心，地处北京市中心繁华地段。', 'active', 0, 'admin', 'admin', NOW(), NOW(), NULL),
(3, '科技园大厦', '科技地产', '深圳市南山区科技园南区', '深圳市', '南山区', '科技园', '办公', '科技物业', '科技园大厦是专为高科技企业打造的现代化办公大楼，配备高速网络和智能办公系统。', 'active', 0, 'admin', 'admin', NOW(), NOW(), NULL),
(4, '湖景花园', '湖景地产', '杭州市西湖区文三路478号', '杭州市', '西湖区', '文三路', '住宅', '湖景物业', '湖景花园依湖而建，环境优美，空气清新，是理想的居住选择。', 'active', 1, 'admin', 'admin', NOW(), NOW(), NULL),
(5, '金融中心', '金融地产', '广州市天河区珠江新城冼村路', '广州市', '天河区', '珠江新城', '商业', '金融物业', '金融中心是广州市标志性建筑，吸引了众多金融机构和企业入驻。', 'pending', 0, 'admin', 'admin', NOW(), NOW(), NULL);
