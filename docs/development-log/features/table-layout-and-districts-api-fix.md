# 表格布局优化和区域接口修复

## 功能概述

完善了楼盘管理页面的两个重要功能：
1. 优化表格列布局：将在租数字段移动到更合理的位置
2. 修复区域数据获取问题：实现缺失的区域和商圈API接口

## 实现详情

### 1. 表格列布局优化

**目标**: 将在租数字段移动到序号字段后面，楼盘名称字段前面

**文件**: `rent-foren/src/views/rental/building/building-management.vue`

**原布局**:
```
序号 | 楼盘名称 | 区域 | 商圈 | 物业类型 | 详细地址 | 在租数 | 操作
```

**优化后布局**:
```
序号 | 在租数 | 楼盘名称 | 区域 | 商圈 | 物业类型 | 详细地址 | 操作
```

**代码修改**:
```vue
<!-- 原顺序 -->
<el-table-column label="序号" width="80" align="center">
  <!-- 序号模板 -->
</el-table-column>
<el-table-column label="楼盘名称" min-width="150">
  <!-- 楼盘名称模板 -->
</el-table-column>
<!-- ... 其他列 ... -->
<el-table-column prop="rent_count" label="在租数" width="80" align="center" />

<!-- 优化后 -->
<el-table-column label="序号" width="80" align="center">
  <!-- 序号模板 -->
</el-table-column>
<el-table-column prop="rent_count" label="在租数" width="80" align="center" />
<el-table-column label="楼盘名称" min-width="150">
  <!-- 楼盘名称模板 -->
</el-table-column>
<!-- ... 其他列 ... -->
```

**布局优势**:
- **逻辑顺序**: 序号 → 关键指标(在租数) → 主要信息(楼盘名称)
- **视觉效果**: 在租数作为重要业务指标更早展示
- **用户体验**: 便于快速识别热门楼盘

### 2. 区域和商圈API接口实现

**问题**: 前端请求 `/api/v1/districts` 接口返回404错误

**根因**: 后端缺少区域和商圈数据获取接口

**文件**: `cmd/api/routes/building_routes.go`

#### 区域列表接口

**接口**: `GET /api/v1/districts`

**实现**:
```go
// 获取区域列表
api.GET("/districts", func(c *gin.Context) {
    var districts []map[string]interface{}
    result := database.DB.Raw("SELECT id, code, name, city_code, sort, status FROM sys_districts WHERE status = 'active' ORDER BY sort ASC").Scan(&districts)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code":    500,
            "message": "获取区域列表失败",
            "error":   result.Error.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code":    200,
        "message": "获取区域列表成功",
        "data":    districts,
    })
})
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取区域列表成功",
  "data": [
    {
      "id": 1,
      "code": "BJ001",
      "name": "朝阳区",
      "city_code": "BJ",
      "sort": 1,
      "status": "active"
    },
    {
      "id": 2,
      "code": "BJ002",
      "name": "海淀区",
      "city_code": "BJ",
      "sort": 2,
      "status": "active"
    }
    // ... 更多区域
  ]
}
```

#### 商圈列表接口

**接口**: `GET /api/v1/business-areas[?districtId=区域ID]`

**实现**:
```go
// 获取商圈列表
api.GET("/business-areas", func(c *gin.Context) {
    districtId := c.Query("districtId")

    query := "SELECT id, code, name, district_id, city_code, sort, status FROM sys_business_areas WHERE status = 'active'"
    args := []interface{}{}

    if districtId != "" {
        query += " AND district_id = ?"
        args = append(args, districtId)
    }

    query += " ORDER BY sort ASC"

    var businessAreas []map[string]interface{}
    result := database.DB.Raw(query, args...).Scan(&businessAreas)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code":    500,
            "message": "获取商圈列表失败",
            "error":   result.Error.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code":    200,
        "message": "获取商圈列表成功",
        "data":    businessAreas,
    })
})
```

**功能特性**:
1. **全量查询**: 不传 `districtId` 参数，返回所有商圈
2. **按区域筛选**: 传入 `districtId` 参数，返回指定区域下的商圈
3. **排序**: 按 `sort` 字段升序排列
4. **状态过滤**: 只返回 `status = 'active'` 的数据

## 测试验证

### 1. 表格布局测试

**验证方式**: 前端页面查看表格列顺序

**期望结果**: 序号 → 在租数 → 楼盘名称 → 其他字段

✅ **测试通过**: 表格列顺序符合预期

### 2. 区域接口测试

**测试命令**:
```bash
curl -X GET "http://localhost:8002/api/v1/districts"
```

**测试结果**:
```json
{
  "code": 200,
  "message": "获取区域列表成功",
  "data": [
    {"id": 1, "name": "朝阳区", "sort": 1},
    {"id": 2, "name": "海淀区", "sort": 2},
    // ... 共10个区域
  ]
}
```

✅ **测试通过**: 返回10个北京区域数据

### 3. 商圈接口测试

**全量查询**:
```bash
curl -X GET "http://localhost:8002/api/v1/business-areas"
```

**按区域筛选**:
```bash
curl -X GET "http://localhost:8002/api/v1/business-areas?districtId=1"
```

**测试结果**:
- 全量查询: 返回所有商圈数据
- 按区域筛选: 返回朝阳区(ID:1)下的5个商圈
  - 国贸商圈
  - 三里屯商圈
  - 望京商圈
  - 亚运村商圈
  - CBD商圈

✅ **测试通过**: 接口功能正常，数据准确

## 数据库表结构

### sys_districts (区域表)
```sql
CREATE TABLE `sys_districts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(20) NOT NULL COMMENT '区域代码',
  `name` varchar(100) NOT NULL COMMENT '区域名称',
  `city_code` varchar(10) NOT NULL COMMENT '城市代码',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
  `status` varchar(20) NOT NULL DEFAULT 'active' COMMENT '状态',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_sys_districts_code` (`code`),
  KEY `idx_sys_districts_city_code` (`city_code`)
);
```

### sys_business_areas (商圈表)
```sql
CREATE TABLE `sys_business_areas` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(20) NOT NULL COMMENT '商圈代码',
  `name` varchar(100) NOT NULL COMMENT '商圈名称',
  `district_id` int NOT NULL COMMENT '所属区域ID',
  `city_code` varchar(10) NOT NULL COMMENT '城市代码',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
  `status` varchar(20) NOT NULL DEFAULT 'active' COMMENT '状态',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_sys_business_areas_code` (`code`),
  KEY `idx_sys_business_areas_district_id` (`district_id`)
);
```

## 业务价值

### 📊 表格布局优化
1. **业务优先**: 在租数作为关键指标优先展示
2. **视觉层次**: 重要信息前置，提升用户体验
3. **逻辑清晰**: 序号 → 指标 → 名称的逻辑顺序

### 🔧 接口功能完善
1. **功能完整**: 补齐缺失的基础数据接口
2. **级联筛选**: 支持区域-商圈的级联选择
3. **数据准确**: 基于数据库真实数据，确保准确性
4. **扩展性好**: 接口设计灵活，支持后续扩展

## 相关文件

**前端文件**:
- `rent-foren/src/views/rental/building/building-management.vue` - 表格布局优化
- `rent-foren/src/api/building.ts` - 区域商圈接口类型定义

**后端文件**:
- `cmd/api/routes/building_routes.go` - 区域商圈接口实现

**数据库表**:
- `sys_districts` - 区域数据表
- `sys_business_areas` - 商圈数据表

## 注意事项

1. **接口性能**: 商圈数据较多时可考虑增加缓存
2. **数据一致性**: 确保前端显示的区域商圈与数据库数据一致
3. **错误处理**: 接口已包含完整的错误处理逻辑
4. **扩展性**: 接口设计支持后续增加更多筛选条件
