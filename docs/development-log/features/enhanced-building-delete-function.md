# 完善楼盘删除功能实现

## 功能概述

实现了两个重要的楼盘删除功能增强：
1. 删除楼盘后自动刷新页面数据
2. 删除楼盘前检查是否有关联户型数据，有则阻止删除

## 实现详情

### 1. 后端API增强

**文件**: `cmd/api/routes/building_routes.go`

**删除接口逻辑增强**:
```go
// 删除楼盘
api.DELETE("/buildings/:id", func(c *gin.Context) {
    id := c.Param("id")

    // 检查楼盘是否存在
    var buildingExists int64
    database.DB.Raw("SELECT COUNT(*) FROM sys_buildings WHERE id = ? AND deleted_at IS NULL", id).Scan(&buildingExists)
    if buildingExists == 0 {
        c.JSON(http.StatusNotFound, gin.H{
            "code":    404,
            "message": "楼盘不存在",
        })
        return
    }

    // 检查是否有关联的户型数据
    var houseTypeCount int64
    database.DB.Raw("SELECT COUNT(*) FROM sys_house_types WHERE building_id = ? AND deleted_at IS NULL", id).Scan(&houseTypeCount)
    
    if houseTypeCount > 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "code":    400,
            "message": "该楼盘下还有户型数据，无法删除",
            "data": gin.H{
                "house_type_count": houseTypeCount,
            },
        })
        return
    }

    // 删除数据库记录（软删除）
    result := database.DB.Exec("UPDATE sys_buildings SET deleted_at = NOW() WHERE id = ?", id)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code":    500,
            "message": "删除楼盘失败",
            "error":   result.Error.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code":    200,
        "message": "删除楼盘成功",
    })
})
```

**关键改进**:
1. **存在性检查**: 删除前检查楼盘是否存在且未被删除
2. **关联数据检查**: 查询该楼盘下是否有未删除的户型数据
3. **详细错误信息**: 返回具体的错误信息和户型数量
4. **状态码区分**: 使用不同的HTTP状态码区分不同的错误类型

### 2. 前端删除逻辑优化

**文件**: `rent-foren/src/views/rental/building/building-management.vue`

**删除确认对话框增强**:
```typescript
const handleDelete = (row: any) => {
  ElMessageBox.confirm(
    `确定要删除楼盘 "${row.name}" 吗？`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
      dangerouslyUseHTMLString: true,
      message: `
        <div>
          <p>确定要删除楼盘 "<strong>${row.name}</strong>" 吗？</p>
          <p style="color: #E6A23C; font-size: 12px; margin-top: 8px;">
            <i class="el-icon-warning"></i> 注意：如果楼盘下有户型数据，将无法删除
          </p>
        </div>
      `
    }
  ).then(async () => {
    try {
      await deleteBuilding(row.id)
      ElMessage.success('删除成功')
      // 刷新页面数据
      await fetchBuildings()
    } catch (error: any) {
      console.error('删除楼盘失败:', error)
      
      // 处理具体的错误信息
      if (error.message && error.message.includes('户型数据')) {
        ElMessageBox.alert(
          '该楼盘下还有户型数据，请先删除所有户型后再删除楼盘。',
          '无法删除',
          {
            confirmButtonText: '知道了',
            type: 'warning'
          }
        )
      } else {
        ElMessage.error(error.message || '删除楼盘失败')
      }
    }
  }).catch(() => {
    // 用户取消删除
  })
}
```

**关键改进**:
1. **增强的确认对话框**: 添加了警告信息，提醒用户删除限制
2. **自动刷新数据**: 删除成功后自动调用 `fetchBuildings()` 刷新列表
3. **专门的错误处理**: 对户型数据冲突错误显示专门的对话框
4. **用户友好的提示**: 清晰告知用户如何解决问题

## 测试结果

### 1. 删除有户型数据的楼盘
```bash
curl -X DELETE "http://localhost:8002/api/v1/buildings/9"
```
**响应**:
```json
{
  "code": 400,
  "data": {
    "house_type_count": 1
  },
  "message": "该楼盘下还有户型数据，无法删除"
}
```

### 2. 删除无户型数据的楼盘
```bash
curl -X DELETE "http://localhost:8002/api/v1/buildings/5"
```
**响应**:
```json
{
  "code": 200,
  "message": "删除楼盘成功"
}
```

### 3. 数据库状态验证
删除前楼盘总数: 7
删除后楼盘总数: 5
删除的楼盘: ID 4, 5 (无户型数据)
保留的楼盘: ID 1 (7个户型), ID 9 (1个户型) 等

## 功能特点

### ✅ 已实现功能
1. **智能删除保护**: 防止误删有业务数据的楼盘
2. **实时数据刷新**: 删除后立即更新列表显示
3. **友好的用户交互**: 清晰的错误提示和操作引导
4. **数据完整性保护**: 确保删除操作不会破坏业务数据关联

### 🔄 业务流程
1. 用户点击删除按钮
2. 显示确认对话框，提醒可能的限制
3. 后端检查楼盘是否存在
4. 后端检查是否有关联的户型数据
5. 如有户型数据，返回错误信息
6. 如无户型数据，执行软删除
7. 前端根据结果显示成功或错误信息
8. 成功删除后自动刷新页面数据

### 📋 注意事项
- 使用软删除机制，数据仍保留在数据库中
- 需要先删除所有户型数据才能删除楼盘
- 删除操作不可逆，需谨慎操作
- 前端会自动刷新数据，确保显示最新状态

## 相关文件

**后端文件**:
- `cmd/api/routes/building_routes.go` - 删除接口逻辑

**前端文件**:
- `rent-foren/src/views/rental/building/building-management.vue` - 删除UI逻辑
- `rent-foren/src/api/building.ts` - 删除API调用

**数据库表**:
- `sys_buildings` - 楼盘主表
- `sys_house_types` - 户型表，关联building_id
