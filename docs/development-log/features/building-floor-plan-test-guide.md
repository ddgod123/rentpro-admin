# 🏗️ 楼盘-户型图上传功能测试指南

**创建日期：** 2024年12月
**测试版本：** rentpro-admin v1.0.0
**测试目标：** 验证楼盘文件夹管理和户型图上传功能

## 📋 测试概述

### 测试场景
1. **创建新楼盘** → 自动生成楼盘文件夹结构
2. **创建户型** → 关联到对应楼盘
3. **上传户型图** → 保存到楼盘的floor-plans文件夹
4. **图片管理** → 查看、更新、删除图片
5. **统计查询** → 获取楼盘图片统计信息

### 文件夹结构
```
七牛云存储空间/
├── buildings/
│   ├── {buildingId}/           # 楼盘文件夹
│   │   ├── floor-plans/        # 户型图文件夹
│   │   │   ├── floor_plan_{timestamp}_{filename}
│   │   │   └── ...
│   │   ├── images/             # 楼盘图片文件夹
│   │   └── documents/          # 相关文档文件夹
│   └── common/                 # 通用文件夹（无楼盘ID时使用）
```

## 🚀 快速测试

### 1. 启动服务
```bash
# 启动后端服务
cd /Users/mac/go/src/rentPro/houduan/rentpro-admin-main
go run main.go api -c config/settings.yml -p 8002

# 启动前端服务（可选）
cd /Users/mac/go/src/rentPro/houduan/rent-foren
npm run dev
```

### 2. 运行测试脚本
```bash
# 给测试脚本执行权限
chmod +x test_building_floor_plan.sh

# 运行完整测试
./test_building_floor_plan.sh
```

### 3. 手动测试步骤

#### 第一步：创建楼盘
```bash
curl -X POST "http://localhost:8002/api/v1/buildings" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试楼盘_$(date +%s)",
    "district": "朝阳区",
    "businessArea": "三里屯",
    "propertyType": "住宅",
    "status": "available",
    "description": "测试楼盘"
  }'
```

**预期响应：**
```json
{
    "code": 201,
    "message": "楼盘创建成功",
    "data": {
        "id": 123,
        "name": "测试楼盘_1734567890",
        "district": "朝阳区",
        "businessArea": "三里屯",
        "folderStructure": "buildings/123/"
    }
}
```

#### 第二步：创建户型
```bash
curl -X POST "http://localhost:8002/api/v1/house-types" \
  -H "Content-Type: application/json" \
  -d '{
    "buildingId": 123,
    "name": "一室一厅户型",
    "type": "1室1厅",
    "area": 65.5,
    "price": 3500000,
    "status": "available"
  }'
```

#### 第三步：上传户型图
```bash
# 准备测试图片文件
echo "创建测试图片文件..."
convert -size 800x600 xc:blue test_floor_plan.jpg

# 上传户型图
curl -X POST "http://localhost:8002/api/v1/upload/floor-plan" \
  -F "house_type_id=456" \
  -F "file=@test_floor_plan.jpg"
```

**预期响应：**
```json
{
    "code": 200,
    "message": "户型图上传成功",
    "data": {
        "image_id": 789,
        "original_url": "https://cdn.domain.com/buildings/123/floor-plans/floor_plan_1734567890_test.jpg",
        "thumbnail_url": "https://cdn.domain.com/buildings/123/floor-plans/floor_plan_1734567890_test.jpg-thumb",
        "building_id": 123,
        "house_type_id": 456
    }
}
```

## 🔍 详细测试用例

### 1. 楼盘文件夹初始化测试

**测试目的：** 验证楼盘创建时是否正确初始化文件夹结构

**测试步骤：**
1. 创建新楼盘
2. 检查响应中是否包含文件夹信息
3. 验证七牛云存储中是否创建了对应文件夹

**预期结果：**
- ✅ 楼盘创建成功
- ✅ 返回楼盘ID和文件夹结构信息
- ✅ 七牛云中创建了 `buildings/{id}/` 文件夹

### 2. 户型图上传测试

**测试目的：** 验证户型图是否上传到正确的楼盘文件夹

**测试步骤：**
1. 上传户型图到指定户型
2. 检查图片URL是否包含正确的楼盘路径
3. 验证数据库中是否正确记录了图片信息

**预期结果：**
- ✅ 图片上传成功
- ✅ URL格式：`buildings/{buildingId}/floor-plans/{filename}`
- ✅ 数据库记录完整（包含楼盘ID、户型ID等信息）

### 3. 批量图片管理测试

**测试目的：** 验证楼盘图片的批量查询和管理功能

**测试步骤：**
1. 为同一楼盘上传多张不同类型的图片
2. 查询楼盘的所有图片
3. 查询楼盘的户型图
4. 按分类查询图片

**测试命令：**
```bash
# 获取楼盘所有图片
curl "http://localhost:8002/api/v1/buildings/123/images"

# 获取楼盘户型图
curl "http://localhost:8002/api/v1/buildings/123/floor-plans"

# 按分类获取图片
curl "http://localhost:8002/api/v1/buildings/123/images?category=floor_plan"
```

### 4. 图片更新和删除测试

**测试目的：** 验证图片信息的更新和删除功能

**测试步骤：**
1. 更新图片名称和描述
2. 设置/取消主图
3. 删除图片
4. 验证七牛云文件是否被正确删除

**测试命令：**
```bash
# 更新图片信息
curl -X PUT "http://localhost:8002/api/v1/images/789" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "更新后的户型图",
    "description": "更新描述"
  }'

# 删除图片
curl -X DELETE "http://localhost:8002/api/v1/images/789"
```

### 5. 统计信息测试

**测试目的：** 验证图片统计功能是否正常工作

**测试步骤：**
1. 查看总体图片统计
2. 查看分类统计
3. 查看今日上传统计

**测试命令：**
```bash
curl "http://localhost:8002/api/v1/images/stats"
```

**预期响应：**
```json
{
    "code": 200,
    "data": {
        "totalImages": 15,
        "totalSize": 52428800,
        "categoryStats": {
            "building": 5,
            "house": 8,
            "floor_plan": 2
        },
        "moduleStats": {
            "rental": 15
        },
        "todayUploads": 3
    }
}
```

## 🐛 常见问题排查

### 1. 文件夹创建失败
```
问题：楼盘创建后没有生成文件夹结构
解决：
1. 检查图片管理器是否正确初始化
2. 查看控制台日志是否有错误信息
3. 验证七牛云配置是否正确
```

### 2. 图片上传失败
```
问题：户型图上传失败
解决：
1. 检查户型ID是否存在
2. 验证图片文件格式和大小
3. 查看七牛云服务状态
4. 检查网络连接
```

### 3. 图片URL错误
```
问题：上传后的图片URL不正确
解决：
1. 检查楼盘ID是否正确传递
2. 验证图片管理器的文件夹生成逻辑
3. 查看七牛云配置的域名设置
```

### 4. 数据库记录不完整
```
问题：图片上传成功但数据库记录缺失
解决：
1. 检查数据库连接状态
2. 查看事务回滚逻辑
3. 验证图片管理器的数据库操作
```

## 📊 测试覆盖范围

### ✅ 已测试功能
- [x] 楼盘创建和文件夹初始化
- [x] 户型创建
- [x] 户型图上传到指定文件夹
- [x] 图片URL格式验证
- [x] 数据库记录完整性
- [x] 图片列表查询
- [x] 图片详情获取
- [x] 图片信息更新
- [x] 图片删除
- [x] 统计信息获取

### 🔄 待测试功能
- [ ] 批量图片上传
- [ ] 图片压缩优化
- [ ] CDN分发验证
- [ ] 高并发上传测试
- [ ] 大文件上传测试

## 🎯 测试结果评估

### 成功标准
1. **功能完整性**：所有API接口正常工作
2. **数据正确性**：图片URL格式正确，数据库记录完整
3. **文件夹结构**：按楼盘组织文件结构清晰
4. **错误处理**：异常情况处理正确
5. **性能表现**：上传速度和查询效率良好

### 性能指标
- **上传时间**：< 3秒（小文件）
- **查询响应**：< 500ms
- **并发处理**：支持10+并发上传
- **存储效率**：图片自动压缩优化

## 📞 技术支持

### 联系方式
- **开发者**：系统管理员
- **文档位置**：`docs/development-log/features/`
- **测试脚本**：`test_building_floor_plan.sh`

### 相关文档
- [七牛云集成指南](qiniu-integration-guide.md)
- [图片管理系统设计](image-management.md)
- [API接口文档](api-documentation.md)

---

**测试完成标记：** ⬜ 未开始 🔄 进行中 ✅ 已完成

**测试日期：** ________

**测试人员：** ________

**测试结果：** ________
