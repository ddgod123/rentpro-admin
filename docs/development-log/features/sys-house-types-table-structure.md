# sys_house_types 户型表结构分析

**分析日期：** 2024年12月
**表名：** `sys_house_types`
**用途：** 存储楼盘户型模板信息
**引擎：** InnoDB
**字符集：** utf8mb4

## 📋 表结构概览

### 基本信息
- **表名：** sys_house_types
- **引擎：** InnoDB
- **自增ID：** 当前最大值 12
- **字符集：** utf8mb4_0900_ai_ci
- **外键约束：** 与 sys_buildings 表关联

## 🏗️ 字段详细说明

### 1. 主键和标识字段

| 字段名 | 类型 | 约束 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `id` | bigint unsigned | NOT NULL, AUTO_INCREMENT, PRIMARY KEY | - | 户型唯一标识符 |
| `name` | varchar(100) | NOT NULL, INDEX | - | 户型名称，如"经典一居"、"舒适两居" |
| `code` | varchar(50) | NOT NULL, UNIQUE | - | 户型编码，如"A1"、"B2"，全局唯一 |
| `description` | text | NULL | NULL | 户型描述信息 |

### 2. 关联关系字段

| 字段名 | 类型 | 约束 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `building_id` | bigint unsigned | NOT NULL, INDEX, FOREIGN KEY | - | 所属楼盘ID，关联 sys_buildings 表 |

### 3. 户型规格字段

| 字段名 | 类型 | 约束 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `standard_area` | decimal(8,2) | NOT NULL, INDEX | - | 标准建筑面积(平方米)，精确到小数点后2位 |
| `rooms` | bigint | NOT NULL | 1 | 房间数量（几室） |
| `halls` | bigint | NOT NULL | 1 | 客厅数量（几厅） |
| `bathrooms` | bigint | NOT NULL | 1 | 卫生间数量（几卫） |
| `balconies` | bigint | NULL | 0 | 阳台数量 |
| `floor_height` | decimal(4,2) | NULL | NULL | 层高(米)，精确到小数点后2位 |

### 4. 户型特征字段

| 字段名 | 类型 | 约束 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `standard_orientation` | varchar(50) | NULL | NULL | 标准朝向，如"南北"、"东西" |
| `standard_view` | varchar(100) | NULL | NULL | 标准景观，如"园景"、"海景" |

### 5. 价格相关字段

| 字段名 | 类型 | 约束 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `base_sale_price` | decimal(12,2) | NULL, INDEX | 0.00 | 基准售价(元)，总价 |
| `base_rent_price` | decimal(8,2) | NULL, INDEX | 0.00 | 基准月租金(元) |
| `base_sale_price_per` | decimal(8,2) | NULL | 0.00 | 基准单价(元/平方米) |
| `base_rent_price_per` | decimal(6,2) | NULL | 0.00 | 基准租金单价(元/平方米/月) |

### 6. 库存管理字段

| 字段名 | 类型 | 约束 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `total_stock` | bigint | NOT NULL | 0 | 总库存数量 |
| `available_stock` | int | NULL | 0 | 可用库存数量 |
| `sold_stock` | int | NULL | 0 | 已售库存数量 |
| `rented_stock` | int | NULL | 0 | 已租库存数量 |
| `reserved_stock` | bigint | NOT NULL | 0 | 预留库存数量 |

### 7. 状态和标签字段

| 字段名 | 类型 | 约束 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `status` | varchar(20) | NOT NULL, INDEX | 'active' | 户型状态：active(正常)、inactive(停用) |
| `is_hot` | tinyint(1) | NULL, INDEX | 0 | 是否热门户型：1(是)、0(否) |

### 8. 图片和媒体字段

| 字段名 | 类型 | 约束 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `main_image` | varchar(500) | NULL | NULL | 主图片URL |
| `floor_plan_url` | varchar(500) | NULL | NULL | 户型图URL |
| `image_urls` | json | NULL | NULL | 图片URL数组，JSON格式存储 |
| `tags` | json | NULL | NULL | 户型标签数组，JSON格式存储 |

### 9. 审计字段

| 字段名 | 类型 | 约束 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `created_by` | varchar(50) | NULL | NULL | 创建人 |
| `updated_by` | varchar(50) | NULL | NULL | 最后更新人 |
| `created_at` | datetime(3) | NULL | NULL | 创建时间，精确到毫秒 |
| `updated_at` | datetime(3) | NULL | NULL | 更新时间，精确到毫秒 |
| `deleted_at` | datetime(3) | NULL, INDEX | NULL | 软删除时间，GORM软删除支持 |

## 🔗 索引结构

### 1. 主键索引
```sql
PRIMARY KEY (`id`)
```
- **类型：** BTREE
- **唯一性：** 是
- **用途：** 主键约束，保证记录唯一性

### 2. 唯一索引
```sql
UNIQUE KEY `idx_code` (`code`)
```
- **字段：** code
- **用途：** 确保户型编码全局唯一

### 3. 普通索引

#### 名称索引
```sql
KEY `idx_name` (`name`)
```
- **字段：** name
- **用途：** 按户型名称查询优化

#### 楼盘关联索引
```sql
KEY `idx_building_id` (`building_id`)
```
- **字段：** building_id
- **用途：** 按楼盘查询户型列表优化

#### 面积价格复合索引
```sql
KEY `idx_area` (`standard_area`, `base_sale_price`)
```
- **字段：** standard_area, base_sale_price
- **用途：** 按面积和价格范围查询优化

#### 租金索引
```sql
KEY `idx_rent_price` (`base_rent_price`)
```
- **字段：** base_rent_price
- **用途：** 按租金查询优化

#### 状态索引
```sql
KEY `idx_status` (`status`)
```
- **字段：** status
- **用途：** 按状态筛选优化

#### 热门户型索引
```sql
KEY `idx_is_hot` (`is_hot`)
```
- **字段：** is_hot
- **用途：** 热门户型筛选优化

#### 软删除索引
```sql
KEY `idx_sys_house_types_deleted_at` (`deleted_at`)
```
- **字段：** deleted_at
- **用途：** GORM软删除查询优化

## 🔒 约束关系

### 外键约束
```sql
CONSTRAINT `fk_sys_house_types_building` 
FOREIGN KEY (`building_id`) REFERENCES `sys_buildings` (`id`)
```
- **关联表：** sys_buildings
- **关联字段：** id
- **约束类型：** 外键约束
- **作用：** 确保户型必须属于一个有效的楼盘

## 📊 数据类型分析

### 数值类型使用
- **整数类型：** 
  - `bigint unsigned`: 主键和大整数字段
  - `bigint`: 房间数量等计数字段
  - `int`: 库存数量字段
  - `tinyint(1)`: 布尔值字段
  
- **小数类型：**
  - `decimal(8,2)`: 面积、单价等精度要求高的字段
  - `decimal(12,2)`: 总价等大数值字段
  - `decimal(6,2)`: 租金单价等小数值字段
  - `decimal(4,2)`: 层高等小范围精确数值

### 字符串类型使用
- **varchar(100)**: 户型名称
- **varchar(50)**: 户型编码、朝向、创建人等
- **varchar(500)**: 图片URL等长字符串
- **text**: 描述信息等长文本

### JSON类型使用
- **image_urls**: 存储多个图片URL
- **tags**: 存储户型标签数组

### 时间类型使用
- **datetime(3)**: 精确到毫秒的时间戳

## 🎯 业务逻辑分析

### 1. 户型模板概念
`sys_house_types` 表设计为**户型模板**，定义了某个楼盘中某种户型的标准规格和基础信息：

- **标准规格：** 面积、房间数、层高等
- **基准价格：** 售价和租金的基准值
- **库存管理：** 该户型的房源数量统计

### 2. 价格体系
表中包含4个价格相关字段，形成完整的价格体系：

```
总价 = 单价 × 面积
base_sale_price = base_sale_price_per × standard_area
base_rent_price = base_rent_price_per × standard_area
```

### 3. 库存体系
库存字段之间的关系：

```
total_stock = available_stock + sold_stock + rented_stock + reserved_stock
```

- **total_stock**: 该户型的房源总数
- **available_stock**: 可售/可租的房源数
- **sold_stock**: 已售出的房源数
- **rented_stock**: 已出租的房源数
- **reserved_stock**: 预留/锁定的房源数

### 4. 状态管理
- **status**: 户型模板的启用状态
  - `active`: 正常使用
  - `inactive`: 停用（不再销售/出租）
- **is_hot**: 热门户型标记，用于推荐展示

### 5. 多媒体支持
- **main_image**: 户型主图
- **floor_plan_url**: 户型图（平面图）
- **image_urls**: 多张户型相关图片（JSON数组）

## 🔧 优化建议

### 1. 索引优化
当前索引设计合理，覆盖了主要查询场景：
- ✅ 按楼盘查询户型列表
- ✅ 按面积和价格范围筛选
- ✅ 按状态和热门标记筛选
- ✅ 软删除查询优化

### 2. 数据类型优化
- ✅ 价格使用 decimal 类型，避免浮点数精度问题
- ✅ 时间字段精确到毫秒，满足高并发场景
- ✅ JSON 字段用于存储数组数据，灵活性好

### 3. 业务逻辑优化建议
- **价格一致性检查：** 可以添加触发器确保总价和单价的计算一致性
- **库存一致性检查：** 可以添加约束确保库存数量的逻辑正确性
- **状态变更日志：** 考虑添加状态变更历史记录

## 📈 使用场景

### 1. 前端展示场景
- **户型列表页：** 按楼盘显示所有户型
- **户型详情页：** 显示户型完整信息
- **价格筛选：** 按价格范围筛选户型
- **热门推荐：** 显示热门户型

### 2. 后台管理场景
- **户型管理：** CRUD操作
- **价格管理：** 批量调整价格
- **库存管理：** 实时更新库存状态
- **图片管理：** 上传和管理户型图片

### 3. 业务统计场景
- **销售统计：** 按户型统计销售情况
- **库存报表：** 各户型库存状况
- **价格分析：** 户型价格分布和趋势

## 🎉 总结

`sys_house_types` 表设计完善，具有以下特点：

### 优点
- ✅ **数据结构清晰：** 字段命名规范，类型选择合适
- ✅ **索引设计合理：** 覆盖主要查询场景，性能优化到位
- ✅ **业务逻辑完整：** 支持价格、库存、状态等完整业务流程
- ✅ **扩展性良好：** JSON字段支持灵活的标签和图片存储
- ✅ **数据完整性：** 外键约束和唯一约束保证数据一致性

### 设计亮点
- **户型模板概念：** 将户型作为模板，支持批量房源管理
- **完整价格体系：** 总价和单价双重价格体系
- **灵活库存管理：** 多维度库存统计支持
- **多媒体支持：** 完善的图片和户型图管理
- **软删除支持：** GORM集成的软删除功能

这个表结构能够很好地支持房产租赁系统中的户型管理需求！🏠
