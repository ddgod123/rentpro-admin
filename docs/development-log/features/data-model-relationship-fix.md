# 数据模型关联关系修复方案

**创建时间：** 2024年12月
**问题描述：** SysHouse 和 SysHouseType 之间缺少关联关系，存在数据冗余
**修复优先级：** 高

## 🔍 问题分析

### 当前问题
1. **缺失关联：** SysHouse 没有关联到 SysHouseType
2. **数据冗余：** 两个模型有大量重复字段
3. **业务逻辑混乱：** 户型和房屋实例的职责不清晰
4. **库存管理困难：** 无法准确统计户型库存

### 业务逻辑设计
```
楼盘 (SysBuildings)
├── 户型A (SysHouseType) - 模板/规格定义
│   ├── 房屋A1 (SysHouse) - 具体实例
│   ├── 房屋A2 (SysHouse) - 具体实例
│   └── 房屋A3 (SysHouse) - 具体实例
├── 户型B (SysHouseType) - 模板/规格定义
│   ├── 房屋B1 (SysHouse) - 具体实例
│   └── 房屋B2 (SysHouse) - 具体实例
```

## 🛠️ 修复方案

### 方案1：添加关联关系（推荐）

#### 1.1 修改 SysHouse 模型
```go
// SysHouse 房屋模型 - 具体房屋实例
type SysHouse struct {
    // 主键
    ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

    // 基础信息
    Name string `json:"name" gorm:"size:100;not null;index:idx_name"`
    Code string `json:"code" gorm:"size:50;not null;uniqueIndex:idx_code"`

    // 关联关系
    BuildingID  uint         `json:"buildingId" gorm:"not null;index:idx_building_id"`
    Building    SysBuildings `json:"building,omitempty" gorm:"foreignKey:BuildingID"`
    HouseTypeID uint         `json:"houseTypeId" gorm:"not null;index:idx_house_type_id"` // 新增
    HouseType   SysHouseType `json:"houseType,omitempty" gorm:"foreignKey:HouseTypeID"`  // 新增

    // 房屋位置信息（房屋特有）
    Floor      int    `json:"floor"`
    Unit       string `json:"unit" gorm:"size:20"`
    RoomNumber string `json:"roomNumber" gorm:"size:20"`

    // 个性化信息（可能与户型不同）
    ActualArea     float64 `json:"actualArea" gorm:"type:decimal(8,2)"` // 实际面积可能与户型略有差异
    ActualUsableArea float64 `json:"actualUsableArea" gorm:"type:decimal(8,2)"`
    Decoration     string  `json:"decoration" gorm:"size:50"` // 装修情况（房屋特有）
    ActualOrientation string `json:"actualOrientation" gorm:"size:50"` // 实际朝向可能与户型不同
    ActualView     string  `json:"actualView" gorm:"size:100"` // 实际景观

    // 价格信息（可能与户型基准价不同）
    ActualSalePrice    float64 `json:"actualSalePrice" gorm:"type:decimal(12,2);default:0"`
    ActualRentPrice    float64 `json:"actualRentPrice" gorm:"type:decimal(8,2);default:0"`
    PriceAdjustment    float64 `json:"priceAdjustment" gorm:"type:decimal(8,2);default:0"` // 价格调整
    PriceAdjustmentReason string `json:"priceAdjustmentReason" gorm:"size:200"` // 调价原因

    // 状态信息
    Status     string `json:"status" gorm:"size:20;not null;default:'available'"`
    SaleStatus string `json:"saleStatus" gorm:"size:20;default:'available'"`
    RentStatus string `json:"rentStatus" gorm:"size:20;default:'available'"`

    // 房屋特有信息
    MainImage   string   `json:"mainImage" gorm:"size:500"`
    ImageUrls   []string `json:"imageUrls" gorm:"type:json"`
    Facilities  []string `json:"facilities" gorm:"type:json"` // 配套设施
    Description string   `json:"description" gorm:"type:text"`
    Notes       string   `json:"notes" gorm:"type:text"`

    // 管理信息
    CreatedBy string     `json:"createdBy" gorm:"size:50"`
    UpdatedBy string     `json:"updatedBy" gorm:"size:50"`
    CreatedAt *time.Time `json:"createdAt" gorm:"autoCreateTime"`
    UpdatedAt *time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
    DeletedAt *time.Time `json:"deletedAt" gorm:"index"`
}
```

#### 1.2 优化 SysHouseType 模型
```go
// SysHouseType 户型模型 - 户型模板/规格定义
type SysHouseType struct {
    // 主键
    ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

    // 基础信息
    Name        string `json:"name" gorm:"size:100;not null;index:idx_name"`
    Code        string `json:"code" gorm:"size:50;not null;uniqueIndex:idx_code"`
    Description string `json:"description" gorm:"type:text"`

    // 楼盘关联
    BuildingID uint         `json:"building_id" gorm:"not null;index:idx_building_id"`
    Building   SysBuildings `json:"building,omitempty" gorm:"foreignKey:BuildingID"`

    // 户型规格（标准规格）
    StandardArea    float64 `json:"standardArea" gorm:"type:decimal(8,2);not null"`
    Rooms          int     `json:"rooms" gorm:"not null;default:1"`
    Halls          int     `json:"halls" gorm:"not null;default:1"`
    Bathrooms      int     `json:"bathrooms" gorm:"not null;default:1"`
    Balconies      int     `json:"balconies" gorm:"default:0"`
    FloorHeight    float64 `json:"floorHeight" gorm:"type:decimal(4,2)"`

    // 标准朝向和景观
    StandardOrientation string `json:"standardOrientation" gorm:"size:50"`
    StandardView        string `json:"standardView" gorm:"size:100"`

    // 基准价格
    BaseSalePrice    float64 `json:"baseSalePrice" gorm:"type:decimal(12,2);default:0"`
    BaseRentPrice    float64 `json:"baseRentPrice" gorm:"type:decimal(8,2);default:0"`
    BaseSalePricePer float64 `json:"baseSalePricePer" gorm:"type:decimal(8,2);default:0"`
    BaseRentPricePer float64 `json:"baseRentPricePer" gorm:"type:decimal(6,2);default:0"`

    // 库存统计（自动计算）
    TotalStock    int `json:"totalStock" gorm:"default:0"`
    AvailableStock int `json:"availableStock" gorm:"default:0"`
    SoldStock     int `json:"soldStock" gorm:"default:0"`
    RentedStock   int `json:"rentedStock" gorm:"default:0"`

    // 户型状态
    Status string `json:"status" gorm:"size:20;not null;default:'active'"`
    IsHot  bool   `json:"isHot" gorm:"default:false"`

    // 户型展示
    MainImage    string   `json:"mainImage" gorm:"size:500"`
    FloorPlanUrl string   `json:"floorPlanUrl" gorm:"size:500"` // 户型图
    ImageUrls    []string `json:"imageUrls" gorm:"type:json"`
    Tags         []string `json:"tags" gorm:"type:json"`

    // 管理信息
    CreatedBy string     `json:"createdBy" gorm:"size:50"`
    UpdatedBy string     `json:"updatedBy" gorm:"size:50"`
    CreatedAt *time.Time `json:"createdAt" gorm:"autoCreateTime"`
    UpdatedAt *time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
    DeletedAt *time.Time `json:"deletedAt" gorm:"index"`

    // 关联的房屋列表
    Houses []SysHouse `json:"houses,omitempty" gorm:"foreignKey:HouseTypeID"`
}
```

### 方案2：数据库迁移脚本

#### 2.1 添加外键约束
```sql
-- 为 sys_houses 表添加 house_type_id 字段
ALTER TABLE sys_houses ADD COLUMN house_type_id INT UNSIGNED;
ALTER TABLE sys_houses ADD INDEX idx_house_type_id (house_type_id);
ALTER TABLE sys_houses ADD FOREIGN KEY (house_type_id) REFERENCES sys_house_types(id);

-- 更新现有数据（示例逻辑）
-- 根据房屋规格匹配对应的户型
UPDATE sys_houses h 
SET house_type_id = (
    SELECT ht.id 
    FROM sys_house_types ht 
    WHERE ht.building_id = h.building_id 
    AND ht.rooms = h.rooms 
    AND ht.halls = h.halls 
    AND ht.bathrooms = h.bathrooms
    AND ABS(ht.standard_area - h.area) < 5  -- 面积差异小于5平米
    LIMIT 1
);
```

#### 2.2 库存统计触发器
```sql
-- 创建触发器自动更新户型库存统计
DELIMITER $$

CREATE TRIGGER update_house_type_stock_after_house_insert
AFTER INSERT ON sys_houses
FOR EACH ROW
BEGIN
    UPDATE sys_house_types 
    SET total_stock = (
        SELECT COUNT(*) FROM sys_houses 
        WHERE house_type_id = NEW.house_type_id AND deleted_at IS NULL
    ),
    available_stock = (
        SELECT COUNT(*) FROM sys_houses 
        WHERE house_type_id = NEW.house_type_id 
        AND status = 'available' AND deleted_at IS NULL
    ),
    sold_stock = (
        SELECT COUNT(*) FROM sys_houses 
        WHERE house_type_id = NEW.house_type_id 
        AND sale_status = 'sold' AND deleted_at IS NULL
    ),
    rented_stock = (
        SELECT COUNT(*) FROM sys_houses 
        WHERE house_type_id = NEW.house_type_id 
        AND rent_status = 'rented' AND deleted_at IS NULL
    )
    WHERE id = NEW.house_type_id;
END$$

CREATE TRIGGER update_house_type_stock_after_house_update
AFTER UPDATE ON sys_houses
FOR EACH ROW
BEGIN
    -- 更新新户型的库存统计
    IF NEW.house_type_id IS NOT NULL THEN
        UPDATE sys_house_types 
        SET total_stock = (SELECT COUNT(*) FROM sys_houses WHERE house_type_id = NEW.house_type_id AND deleted_at IS NULL),
            available_stock = (SELECT COUNT(*) FROM sys_houses WHERE house_type_id = NEW.house_type_id AND status = 'available' AND deleted_at IS NULL),
            sold_stock = (SELECT COUNT(*) FROM sys_houses WHERE house_type_id = NEW.house_type_id AND sale_status = 'sold' AND deleted_at IS NULL),
            rented_stock = (SELECT COUNT(*) FROM sys_houses WHERE house_type_id = NEW.house_type_id AND rent_status = 'rented' AND deleted_at IS NULL)
        WHERE id = NEW.house_type_id;
    END IF;
    
    -- 如果户型ID发生变化，也要更新旧户型的统计
    IF OLD.house_type_id IS NOT NULL AND OLD.house_type_id != NEW.house_type_id THEN
        UPDATE sys_house_types 
        SET total_stock = (SELECT COUNT(*) FROM sys_houses WHERE house_type_id = OLD.house_type_id AND deleted_at IS NULL),
            available_stock = (SELECT COUNT(*) FROM sys_houses WHERE house_type_id = OLD.house_type_id AND status = 'available' AND deleted_at IS NULL),
            sold_stock = (SELECT COUNT(*) FROM sys_houses WHERE house_type_id = OLD.house_type_id AND sale_status = 'sold' AND deleted_at IS NULL),
            rented_stock = (SELECT COUNT(*) FROM sys_houses WHERE house_type_id = OLD.house_type_id AND rent_status = 'rented' AND deleted_at IS NULL)
        WHERE id = OLD.house_type_id;
    END IF;
END$$

DELIMITER ;
```

### 方案3：API 接口调整

#### 3.1 户型相关接口
```go
// 获取户型列表（包含库存统计）
GET /api/v1/house-types?building_id=1&status=active

// 获取户型详情（包含关联房屋）
GET /api/v1/house-types/1?include=houses

// 获取户型的可用房屋列表
GET /api/v1/house-types/1/available-houses
```

#### 3.2 房屋相关接口
```go
// 获取房屋列表（包含户型信息）
GET /api/v1/houses?building_id=1&house_type_id=2&status=available

// 创建房屋时必须指定户型
POST /api/v1/houses
{
    "name": "A座1001室",
    "building_id": 1,
    "house_type_id": 2,  // 必需字段
    "floor": 10,
    "unit": "A",
    "room_number": "1001"
}
```

## 📋 实施步骤

### 阶段1：模型修改（1-2天）
1. ✅ 修改 SysHouse 模型添加 HouseTypeID 字段
2. ✅ 优化 SysHouseType 模型字段命名
3. ✅ 添加关联方法和业务方法
4. ✅ 更新数据库迁移脚本

### 阶段2：数据迁移（1天）
1. ✅ 执行数据库结构变更
2. ✅ 创建数据迁移脚本
3. ✅ 建立现有数据的关联关系
4. ✅ 验证数据完整性

### 阶段3：API调整（2-3天）
1. ✅ 更新房屋管理API
2. ✅ 更新户型管理API
3. ✅ 添加关联查询接口
4. ✅ 更新库存统计逻辑

### 阶段4：前端适配（2-3天）
1. ✅ 更新房屋管理页面
2. ✅ 更新户型管理页面
3. ✅ 添加户型选择组件
4. ✅ 更新库存展示逻辑

## 🎯 预期收益

### 业务价值
- **数据一致性：** 消除冗余，确保数据准确性
- **库存管理：** 精确的户型库存统计和管理
- **业务清晰：** 明确的户型模板和房屋实例关系
- **扩展性：** 支持更复杂的业务场景

### 技术价值
- **性能优化：** 减少数据冗余，优化查询性能
- **维护性：** 清晰的数据模型，便于后续维护
- **可扩展：** 支持户型变体、动态定价等高级功能

## ⚠️ 风险评估

### 技术风险
- **数据迁移风险：** 现有数据可能无法完美匹配户型
- **性能影响：** 增加关联查询可能影响性能
- **兼容性：** 需要同步更新前端代码

### 缓解措施
- **备份数据：** 迁移前完整备份数据库
- **分步实施：** 先在测试环境验证
- **回滚方案：** 准备数据回滚脚本
- **性能监控：** 监控API响应时间
