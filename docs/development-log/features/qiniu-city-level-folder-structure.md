# 七牛云城市级别文件夹结构优化

## 📋 功能概述

为了更好地组织和管理七牛云存储中的楼盘文件，我们优化了文件夹结构，引入了城市级别的分层管理。

## 🔄 改进前后对比

### 改进前的文件夹结构
```
buildings/
  ├── 1-楼盘A/
  │   ├── floor-plans/
  │   ├── site-plans/
  │   └── ...
  ├── 2-楼盘B/
  │   ├── floor-plans/
  │   ├── site-plans/
  │   └── ...
  └── ...
```

### 改进后的文件夹结构
```
buildings/
  ├── 北京/
  │   ├── 1-楼盘A/
  │   │   ├── floor-plans/
  │   │   ├── site-plans/
  │   │   ├── environment/
  │   │   ├── building-images/
  │   │   ├── interior/
  │   │   ├── facilities/
  │   │   └── documents/
  │   └── 2-楼盘B/
  │       └── ...
  ├── 上海/
  │   ├── 3-楼盘C/
  │   │   └── ...
  │   └── ...
  ├── 广州/
  │   └── ...
  └── 深圳/
      └── ...
```

## 🛠️ 技术实现

### 1. 数据模型更新

修改了 `CreateBuildingFolder` 函数，从数据库获取楼盘的城市信息：

```go
// 获取楼盘信息包括城市
var building struct {
    ID   uint64 `json:"id"`
    Name string `json:"name"`
    City string `json:"city"`
}

if err := im.db.Table("sys_buildings").
    Select("id, name, city").
    Where("id = ?", buildingID).
    First(&building).Error; err != nil {
    return fmt.Errorf("获取楼盘信息失败: %v", err)
}
```

### 2. 文件夹路径生成

更新了 `createQiniuFolderStructure` 函数签名和实现：

```go
// 修改前
func (im *ImageManager) createQiniuFolderStructure(buildingID uint64, buildingName string, folders map[string]string) error

// 修改后
func (im *ImageManager) createQiniuFolderStructure(buildingID uint64, buildingName, cityName string, folders map[string]string) error
```

### 3. 七牛云路径构建

新的文件夹路径构建逻辑：

```go
// 处理城市名称和楼盘名称，确保适合作为文件夹名称
safeCityName := im.sanitizeFolderName(cityName)
safeBuildingName := im.sanitizeFolderName(buildingName)
buildingFolderName := fmt.Sprintf("%d-%s", buildingID, safeBuildingName)

// 创建文件夹标记文件的key，使用城市/楼盘/子文件夹的层级结构
folderKey := fmt.Sprintf("buildings/%s/%s/%s/.folder", safeCityName, buildingFolderName, folder)
```

### 4. 标记文件内容增强

标记文件现在包含更多信息：

```json
{
  "building_id": 18,
  "building_name": "深圳测试楼盘",
  "city_name": "深圳",
  "building_folder_name": "18-深圳测试楼盘",
  "folder_type": "floor-plans",
  "description": "户型图",
  "folder_path": "buildings/深圳/18-深圳测试楼盘/floor-plans/",
  "created_at": "2025-09-05 17:30:15",
  "purpose": "This file marks the existence of this folder structure with city-level organization"
}
```

## ✅ 功能优势

### 1. 更好的文件组织
- **城市级别分类**：按城市对楼盘进行分组，便于管理
- **清晰的层级结构**：城市 → 楼盘 → 文件类型的三级结构
- **可扩展性**：支持更多城市的无缝扩展

### 2. 提升管理效率
- **快速定位**：可以快速定位到特定城市的楼盘文件
- **批量操作**：可以按城市进行批量文件操作
- **权限管理**：未来可以基于城市进行权限分配

### 3. 存储优化
- **逻辑分区**：不同城市的文件逻辑分离
- **减少冲突**：降低同名楼盘的文件夹冲突概率
- **便于备份**：可以按城市进行分别备份

## 🧪 测试验证

### 测试用例

1. **深圳楼盘测试**
   ```bash
   curl -X POST "http://localhost:8002/api/v1/buildings" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "深圳测试楼盘",
       "city": "深圳",
       "district": "福田区",
       "businessArea": "华强北商圈",
       "propertyType": "住宅",
       "description": "测试城市级别文件夹结构",
       "status": "active"
     }'
   ```

2. **广州楼盘测试**
   ```bash
   curl -X POST "http://localhost:8002/api/v1/buildings" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "广州测试楼盘",
       "city": "广州",
       "district": "天河区",
       "businessArea": "天河城商圈",
       "propertyType": "办公",
       "description": "测试广州城市文件夹结构",
       "status": "active"
     }'
   ```

### 验证结果

✅ 成功创建城市级别的文件夹结构：
- `buildings/深圳/18-深圳测试楼盘/`
- `buildings/广州/19-广州测试楼盘/`

每个楼盘文件夹下包含完整的子文件夹：
- `floor-plans/` (户型图)
- `site-plans/` (小区平面图)
- `environment/` (小区环境图)
- `building-images/` (楼盘外观图)
- `interior/` (室内样板图)
- `facilities/` (配套设施图)
- `documents/` (相关文档)

## 🔮 未来扩展

### 1. 区域级别细分
可以进一步细分到区域级别：
```
buildings/
  ├── 北京/
  │   ├── 朝阳区/
  │   ├── 海淀区/
  │   └── ...
  └── ...
```

### 2. 文件类型标签
为不同类型的文件添加更详细的标签和分类。

### 3. 自动清理机制
实现过期或无用文件的自动清理机制。

## 📅 实施日期
2025年9月5日

## 👥 相关人员
- 开发者：AI Assistant
- 测试验证：完成

## 🏷️ 标签
`七牛云` `文件管理` `城市分级` `存储优化` `楼盘管理`
