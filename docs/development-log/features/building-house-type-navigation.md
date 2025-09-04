# 楼盘名称按钮化及户型页面跳转实现方案

**创建时间：** 2024年12月
**需求：** 将楼盘名称字段改为按钮，点击跳转到楼盘户型页面
**优先级：** 中等

## 📋 需求分析

### 功能目标
1. **楼盘名称按钮化：** 将表格中的楼盘名称从文本改为可点击的按钮
2. **页面跳转：** 点击楼盘名称后跳转到该楼盘的户型管理页面
3. **数据传递：** 将楼盘ID和相关信息传递到户型页面
4. **面包屑导航：** 在户型页面显示正确的导航路径

## 🎯 实现方案

### 1. 前端路由设计

#### 1.1 路由结构调整
**文件位置：** `/rent-foren/src/router/index.ts`

```typescript
// 租赁管理模块 - 包含六个子模块
{
  path: '/rental',
  component: Layout,
  redirect: '/rental/building',
  meta: { title: '租赁管理', icon: 'OfficeBuilding' },
  children: [
    // 楼盘管理
    {
      path: 'building',
      name: 'Building',
      component: () => import('@/views/rental/building/building-management.vue'),
      meta: { title: '楼盘管理' }
    },
    // 楼盘户型管理 - 新增路由
    {
      path: 'building/:buildingId/house-types',
      name: 'BuildingHouseTypes',
      component: () => import('@/views/rental/building/house-types.vue'),
      meta: { 
        title: '户型管理',
        hidden: true, // 不在菜单中显示
        breadcrumb: true // 显示面包屑
      },
      props: true // 将路由参数作为props传递给组件
    },
    // 房屋管理
    {
      path: 'house',
      name: 'House',
      component: () => import('@/views/rental/house/index.vue'),
      meta: { title: '房屋管理' }
    },
    // 其他路由...
  ]
}
```

#### 1.2 动态路由参数
- **路由路径：** `/rental/building/:buildingId/house-types`
- **参数传递：** `buildingId` 作为路由参数
- **示例URL：** `/rental/building/1/house-types`

### 2. 前端页面实现

#### 2.1 修改楼盘管理页面
**文件位置：** `/rent-foren/src/views/rental/building/building-management.vue`

##### 表格列修改 (第66行)：
```vue
<!-- 原来的文本列 -->
<!-- <el-table-column prop="name" label="楼盘名称" min-width="150" /> -->

<!-- 修改为按钮列 -->
<el-table-column label="楼盘名称" min-width="150">
  <template #default="{ row }">
    <el-button 
      type="primary" 
      link 
      @click="handleViewHouseTypes(row)"
      class="building-name-btn"
    >
      {{ row.name }}
    </el-button>
  </template>
</el-table-column>
```

##### 添加跳转方法 (script部分)：
```typescript
import { useRouter } from 'vue-router'

const router = useRouter()

// 跳转到楼盘户型页面
const handleViewHouseTypes = (building: any) => {
  router.push({
    name: 'BuildingHouseTypes',
    params: {
      buildingId: building.id
    },
    query: {
      buildingName: building.name, // 传递楼盘名称用于面包屑显示
      district: building.district,
      businessArea: building.business_area
    }
  })
}
```

##### 样式添加：
```scss
.building-name-btn {
  padding: 0;
  font-weight: 500;
  
  &:hover {
    color: var(--el-color-primary-light-3);
  }
}
```

#### 2.2 创建楼盘户型页面
**文件位置：** `/rent-foren/src/views/rental/building/house-types.vue`

```vue
<template>
  <div class="house-types-container">
    <!-- 面包屑导航 -->
    <el-breadcrumb class="breadcrumb mb-20">
      <el-breadcrumb-item :to="{ name: 'Building' }">楼盘管理</el-breadcrumb-item>
      <el-breadcrumb-item>{{ buildingInfo.name || '楼盘户型' }}</el-breadcrumb-item>
    </el-breadcrumb>

    <el-card>
      <template #header>
        <div class="card-header">
          <div class="header-info">
            <h3>{{ buildingInfo.name }} - 户型管理</h3>
            <div class="building-meta">
              <el-tag>{{ buildingInfo.district }}</el-tag>
              <el-tag type="info" v-if="buildingInfo.businessArea">{{ buildingInfo.businessArea }}</el-tag>
            </div>
          </div>
          <div class="header-actions">
            <el-button @click="handleBack">返回楼盘列表</el-button>
            <el-button type="primary" @click="handleAddHouseType">新增户型</el-button>
          </div>
        </div>
      </template>
      
      <!-- 户型统计卡片 -->
      <div class="stats-cards mb-20">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-card class="stats-card">
              <div class="stats-content">
                <div class="stats-number">{{ stats.totalTypes }}</div>
                <div class="stats-label">户型总数</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stats-card">
              <div class="stats-content">
                <div class="stats-number">{{ stats.totalHouses }}</div>
                <div class="stats-label">房屋总数</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stats-card">
              <div class="stats-content">
                <div class="stats-number">{{ stats.availableHouses }}</div>
                <div class="stats-label">可用房屋</div>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stats-card">
              <div class="stats-content">
                <div class="stats-number">{{ stats.soldRentedHouses }}</div>
                <div class="stats-label">已售/已租</div>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>
      
      <!-- 户型列表表格 -->
      <el-table :data="houseTypesData" border style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="户型名称" min-width="120" />
        <el-table-column prop="code" label="户型编码" width="100" />
        <el-table-column label="户型规格" width="120">
          <template #default="{ row }">
            {{ row.rooms }}室{{ row.halls }}厅{{ row.bathrooms }}卫
          </template>
        </el-table-column>
        <el-table-column prop="standard_area" label="标准面积" width="100" align="right">
          <template #default="{ row }">
            {{ row.standard_area }}㎡
          </template>
        </el-table-column>
        <el-table-column label="基准价格" width="150" align="right">
          <template #default="{ row }">
            <div v-if="row.base_sale_price > 0">
              售: {{ formatPrice(row.base_sale_price) }}万
            </div>
            <div v-if="row.base_rent_price > 0">
              租: {{ row.base_rent_price }}元/月
            </div>
          </template>
        </el-table-column>
        <el-table-column label="库存统计" width="120" align="center">
          <template #default="{ row }">
            <div class="stock-info">
              <div>总数: {{ row.total_stock }}</div>
              <div>可用: {{ row.available_stock }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'warning'">
              {{ row.status === 'active' ? '正常' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleViewHouses(row)">查看房屋</el-button>
            <el-button type="info" size="small" @click="handleEditHouseType(row)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDeleteHouseType(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination-area mt-20">
        <el-pagination
          v-model:current-page="pagination.currentPage"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getHouseTypesByBuilding, getBuildingInfo, deleteHouseType } from '@/api/building'

// 路由相关
const router = useRouter()
const route = useRoute()

// Props
const props = defineProps<{
  buildingId: string
}>()

// 楼盘信息
const buildingInfo = reactive({
  id: '',
  name: '',
  district: '',
  businessArea: ''
})

// 表格数据
const houseTypesData = ref<any[]>([])
const loading = ref(false)

// 分页
const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
})

// 统计数据
const stats = computed(() => {
  return {
    totalTypes: houseTypesData.value.length,
    totalHouses: houseTypesData.value.reduce((sum, item) => sum + (item.total_stock || 0), 0),
    availableHouses: houseTypesData.value.reduce((sum, item) => sum + (item.available_stock || 0), 0),
    soldRentedHouses: houseTypesData.value.reduce((sum, item) => sum + (item.sold_stock || 0) + (item.rented_stock || 0), 0)
  }
})

// 获取楼盘信息
const fetchBuildingInfo = async () => {
  try {
    // 从query参数获取楼盘信息
    if (route.query.buildingName) {
      buildingInfo.name = route.query.buildingName as string
      buildingInfo.district = route.query.district as string || ''
      buildingInfo.businessArea = route.query.businessArea as string || ''
    } else {
      // 从API获取楼盘信息
      const response = await getBuildingInfo(props.buildingId)
      if (response && response.data) {
        Object.assign(buildingInfo, response.data)
      }
    }
  } catch (error) {
    console.error('获取楼盘信息失败:', error)
    ElMessage.error('获取楼盘信息失败')
  }
}

// 获取户型列表
const fetchHouseTypes = async () => {
  loading.value = true
  try {
    const params = {
      buildingId: props.buildingId,
      page: pagination.currentPage,
      pageSize: pagination.pageSize
    }
    
    const response = await getHouseTypesByBuilding(params)
    if (response && response.data) {
      houseTypesData.value = response.data
      pagination.total = response.total || response.data.length
    }
  } catch (error) {
    console.error('获取户型列表失败:', error)
    ElMessage.error('获取户型列表失败')
  } finally {
    loading.value = false
  }
}

// 格式化价格
const formatPrice = (price: number) => {
  return (price / 10000).toFixed(0)
}

// 返回楼盘列表
const handleBack = () => {
  router.push({ name: 'Building' })
}

// 新增户型
const handleAddHouseType = () => {
  // TODO: 实现新增户型功能
  ElMessage.info('新增户型功能开发中...')
}

// 查看房屋
const handleViewHouses = (houseType: any) => {
  router.push({
    name: 'House',
    query: {
      buildingId: props.buildingId,
      houseTypeId: houseType.id,
      buildingName: buildingInfo.name,
      houseTypeName: houseType.name
    }
  })
}

// 编辑户型
const handleEditHouseType = (houseType: any) => {
  // TODO: 实现编辑户型功能
  ElMessage.info('编辑户型功能开发中...')
}

// 删除户型
const handleDeleteHouseType = (houseType: any) => {
  ElMessageBox.confirm(
    `确定要删除户型 "${houseType.name}" 吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await deleteHouseType(houseType.id)
      ElMessage.success('删除成功')
      fetchHouseTypes()
    } catch (error: any) {
      console.error('删除户型失败:', error)
      ElMessage.error(error.message || '删除户型失败')
    }
  }).catch(() => {
    // 用户取消删除
  })
}

// 分页处理
const handleSizeChange = (val: number) => {
  pagination.pageSize = val
  pagination.currentPage = 1
  fetchHouseTypes()
}

const handleCurrentChange = (val: number) => {
  pagination.currentPage = val
  fetchHouseTypes()
}

// 页面初始化
onMounted(() => {
  fetchBuildingInfo()
  fetchHouseTypes()
})
</script>

<style scoped lang="scss">
.house-types-container {
  .breadcrumb {
    background: #f5f7fa;
    padding: 12px 16px;
    border-radius: 4px;
  }
  
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    
    .header-info {
      h3 {
        margin: 0 0 8px 0;
        color: #303133;
      }
      
      .building-meta {
        display: flex;
        gap: 8px;
      }
    }
    
    .header-actions {
      display: flex;
      gap: 12px;
    }
  }
  
  .stats-cards {
    .stats-card {
      text-align: center;
      
      .stats-content {
        .stats-number {
          font-size: 28px;
          font-weight: bold;
          color: #409eff;
          margin-bottom: 8px;
        }
        
        .stats-label {
          font-size: 14px;
          color: #909399;
        }
      }
    }
  }
  
  .stock-info {
    font-size: 12px;
    line-height: 1.5;
    
    div {
      white-space: nowrap;
    }
  }
  
  .pagination-area {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }
}
</style>
```

### 3. 后端API接口实现

#### 3.1 户型相关API
**文件位置：** `/rentpro-admin-main/cmd/api/server.go`

##### 添加户型管理API (在现有API组中添加)：
```go
// 获取楼盘的户型列表
api.GET("/buildings/:buildingId/house-types", func(c *gin.Context) {
    buildingID := c.Param("buildingId")
    page := c.DefaultQuery("page", "1")
    pageSize := c.DefaultQuery("pageSize", "10")
    
    // 转换分页参数
    pageNum, _ := strconv.Atoi(page)
    size, _ := strconv.Atoi(pageSize)
    
    if pageNum < 1 {
        pageNum = 1
    }
    if size < 1 {
        size = 10
    }
    
    // 构造查询
    offset := (pageNum - 1) * size
    
    query := `
        SELECT 
            id, name, code, description, building_id,
            standard_area, rooms, halls, bathrooms, balconies, floor_height,
            standard_orientation, standard_view,
            base_sale_price, base_rent_price, base_sale_price_per, base_rent_price_per,
            total_stock, available_stock, sold_stock, rented_stock, reserved_stock,
            status, is_hot, main_image, floor_plan_url,
            created_at, updated_at
        FROM sys_house_types 
        WHERE building_id = ? AND deleted_at IS NULL
        ORDER BY id DESC 
        LIMIT ? OFFSET ?`
    
    countQuery := `
        SELECT COUNT(*) 
        FROM sys_house_types 
        WHERE building_id = ? AND deleted_at IS NULL`
    
    // 执行查询
    var houseTypes []map[string]interface{}
    database.DB.Raw(query, buildingID, size, offset).Scan(&houseTypes)
    
    // 查询总数
    var total int64
    database.DB.Raw(countQuery, buildingID).Scan(&total)
    
    c.JSON(http.StatusOK, gin.H{
        "code":    200,
        "message": "户型列表获取成功",
        "data":    houseTypes,
        "total":   total,
        "page":    pageNum,
        "size":    size,
    })
})

// 获取楼盘基础信息
api.GET("/buildings/:id/info", func(c *gin.Context) {
    id := c.Param("id")
    
    var building map[string]interface{}
    result := database.DB.Raw(`
        SELECT id, name, district, business_area, property_type, 
               detailed_address, property_company, status, is_hot
        FROM sys_buildings 
        WHERE id = ? AND deleted_at IS NULL`, id).Scan(&building)
    
    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "code":    404,
            "message": "楼盘不存在",
        })
        return
    }
    
    if len(building) == 0 {
        c.JSON(http.StatusNotFound, gin.H{
            "code":    404,
            "message": "楼盘不存在",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "code":    200,
        "message": "获取成功",
        "data":    building,
    })
})

// 删除户型
api.DELETE("/house-types/:id", func(c *gin.Context) {
    id := c.Param("id")
    
    // 检查是否有关联的房屋
    var houseCount int64
    database.DB.Raw("SELECT COUNT(*) FROM sys_houses WHERE house_type_id = ? AND deleted_at IS NULL", id).Scan(&houseCount)
    
    if houseCount > 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "code":    400,
            "message": fmt.Sprintf("该户型下还有 %d 套房屋，无法删除", houseCount),
        })
        return
    }
    
    // 软删除户型
    result := database.DB.Exec("UPDATE sys_house_types SET deleted_at = NOW() WHERE id = ?", id)
    
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code":    500,
            "message": "删除户型失败",
            "error":   result.Error.Error(),
        })
        return
    }
    
    if result.RowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{
            "code":    404,
            "message": "户型不存在",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "code":    200,
        "message": "删除成功",
    })
})
```

#### 3.2 API接口封装
**文件位置：** `/rent-foren/src/api/building.ts`

##### 添加户型相关接口：
```typescript
// 户型相关接口类型定义
export interface HouseType {
  id: number
  name: string
  code: string
  description?: string
  building_id: number
  standard_area: number
  rooms: number
  halls: number
  bathrooms: number
  balconies: number
  floor_height?: number
  standard_orientation?: string
  standard_view?: string
  base_sale_price: number
  base_rent_price: number
  base_sale_price_per: number
  base_rent_price_per: number
  total_stock: number
  available_stock: number
  sold_stock: number
  rented_stock: number
  reserved_stock: number
  status: string
  is_hot: boolean
  main_image?: string
  floor_plan_url?: string
  created_at: string
  updated_at: string
}

export interface HouseTypeQuery {
  buildingId: string | number
  page?: number
  pageSize?: number
}

// 获取楼盘的户型列表
export function getHouseTypesByBuilding(params: HouseTypeQuery) {
  return request<{
    data: HouseType[]
    total: number
    page: number
    size: number
  }>({
    url: `/buildings/${params.buildingId}/house-types`,
    method: 'get',
    params: {
      page: params.page,
      pageSize: params.pageSize
    }
  })
}

// 获取楼盘基础信息
export function getBuildingInfo(id: string | number) {
  return request<Building>({
    url: `/buildings/${id}/info`,
    method: 'get'
  })
}

// 删除户型
export function deleteHouseType(id: number) {
  return request({
    url: `/house-types/${id}`,
    method: 'delete'
  })
}
```

## 🔄 实施步骤

### 阶段1：路由和基础页面 (1-2天)
1. ✅ 修改路由配置，添加户型管理路由
2. ✅ 创建户型管理页面基础结构
3. ✅ 实现面包屑导航

### 阶段2：前端交互 (1天)
1. ✅ 修改楼盘列表页面，将名称改为按钮
2. ✅ 实现跳转逻辑和参数传递
3. ✅ 完善户型页面的数据展示

### 阶段3：后端API (1-2天)
1. ✅ 实现户型列表查询API
2. ✅ 实现楼盘信息查询API
3. ✅ 实现户型删除API
4. ✅ 测试API接口功能

### 阶段4：功能完善 (1天)
1. ✅ 添加统计卡片数据
2. ✅ 完善错误处理
3. ✅ 优化用户体验
4. ✅ 添加样式美化

## 📊 用户体验设计

### 交互流程
1. **楼盘列表页面：** 用户看到楼盘名称为蓝色链接按钮
2. **点击跳转：** 点击楼盘名称，页面跳转到户型管理页面
3. **户型页面：** 显示面包屑导航，楼盘信息，户型列表和统计
4. **返回导航：** 用户可以通过面包屑或返回按钮回到楼盘列表

### 视觉设计
- **按钮样式：** 使用 `type="primary" link` 的样式，看起来像链接但有按钮的交互
- **面包屑：** 清晰的导航路径，支持点击跳转
- **统计卡片：** 直观展示户型和房屋统计数据
- **表格布局：** 合理的列宽和数据展示

## ⚠️ 注意事项

### 技术考虑
1. **路由参数验证：** 确保buildingId参数的有效性
2. **权限控制：** 后续可能需要添加页面访问权限
3. **数据缓存：** 考虑添加楼盘信息的缓存机制
4. **错误处理：** 完善各种异常情况的处理

### 扩展性
1. **户型详情：** 后续可以添加户型详情页面
2. **房屋列表：** 从户型页面可以跳转到房屋列表
3. **数据筛选：** 户型页面可以添加搜索和筛选功能
4. **批量操作：** 支持批量管理户型

这个方案提供了完整的从楼盘名称按钮化到户型页面展示的实现路径，包含了前后端的完整代码和详细的实施步骤。
