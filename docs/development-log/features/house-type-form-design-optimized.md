# 新增户型功能 - 优化表单设计方案

**创建时间：** 2024年12月
**功能模块：** 户型管理
**设计理念：** 分步骤渐进式表单，核心信息优先

## 🎯 **表单设计理念**

### **设计原则**
1. **核心优先：** 最重要的信息放在第一步
2. **渐进式：** 从必填到选填，降低用户心理负担
3. **直观易懂：** 每步聚焦特定类型的信息
4. **灵活提交：** 完成第一步即可保存，后续可补充

## 📋 **分步骤表单设计**

### 🔴 **第一步：基础信息（必填）**
**目标：** 创建户型的最小可用信息

```vue
┌─ 第1步：基础信息 ─────────────────────┐
│                                      │
│ 户型名称: [____________________] *   │
│ 户型编码: [____________________] *   │
│ 建筑面积: [____________________] * ㎡│
│ 所属楼盘: [下拉选择_____________] *   │
│                                      │
│ [上一步]              [保存并继续] │
└──────────────────────────────────────┘
```

**字段详情：**
- **`name`** - 户型名称 *(varchar 100)*
- **`code`** - 户型编码 *(varchar 50)*
- **`standard_area`** - 建筑面积 *(decimal 8,2)*
- **`building_id`** - 所属楼盘 *(关联选择)*

**验证规则：**
```typescript
const step1Rules = {
  name: [
    { required: true, message: '请输入户型名称', trigger: 'blur' },
    { min: 1, max: 100, message: '户型名称长度为1-100个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入户型编码', trigger: 'blur' },
    { pattern: /^[A-Z0-9]+$/, message: '户型编码只能包含大写字母和数字', trigger: 'blur' }
  ],
  standard_area: [
    { required: true, message: '请输入建筑面积', trigger: 'blur' },
    { type: 'number', min: 0.01, max: 9999.99, message: '面积范围0.01-9999.99㎡', trigger: 'blur' }
  ],
  building_id: [
    { required: true, message: '请选择所属楼盘', trigger: 'change' }
  ]
}
```

### 🟡 **第二步：户型图片（选填）**
**目标：** 上传户型相关图片

```vue
┌─ 第2步：户型图片 ─────────────────────┐
│                                      │
│ 户型图片:                            │
│ ┌─────────────┐ ┌─────────────┐      │
│ │   主图上传   │ │  户型图上传  │      │
│ │     +       │ │     +       │      │
│ │  点击上传    │ │  点击上传    │      │
│ └─────────────┘ └─────────────┘      │
│                                      │
│ 更多图片:                            │
│ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐      │
│ │  +  │ │  +  │ │  +  │ │  +  │      │
│ └─────┘ └─────┘ └─────┘ └─────┘      │
│                                      │
│ [上一步]              [保存并继续] │
└──────────────────────────────────────┘
```

**字段详情：**
- **`main_image`** - 主图URL *(varchar 500)*
- **`floor_plan_url`** - 户型图URL *(varchar 500)*
- **`image_urls`** - 更多图片URL列表 *(JSON)*

**功能特性：**
- 图片预览功能
- 拖拽上传支持
- 图片压缩和格式验证
- 支持批量上传

### 🟢 **第三步：户型布局（选填）**
**目标：** 详细的户型规格信息

```vue
┌─ 第3步：户型布局 ─────────────────────┐
│                                      │
│ 房间布局:                            │
│ 房间数: [▼] 室  客厅数: [▼] 厅       │
│ 卫生间: [▼] 卫  阳台数: [▼] 个       │
│                                      │
│ 规格信息:                            │
│ 层    高: [__________] 米             │
│ 朝    向: [__________]               │
│ 景    观: [__________]               │
│                                      │
│ [上一步]              [保存并继续] │
└──────────────────────────────────────┘
```

**字段详情：**
- **`rooms`** - 房间数 *(默认1)*
- **`halls`** - 客厅数 *(默认1)*
- **`bathrooms`** - 卫生间数 *(默认1)*
- **`balconies`** - 阳台数 *(默认0)*
- **`floor_height`** - 层高 *(decimal 4,2)*
- **`standard_orientation`** - 朝向 *(varchar 50)*
- **`standard_view`** - 景观 *(varchar 100)*

**交互特性：**
- 数字选择器（1-9）
- 朝向预设选项：南向、北向、东向、西向、南北、东西
- 景观预设选项：小区景观、城市景观、海景、山景、园景

### 🟢 **第四步：其他信息（选填）**
**目标：** 价格、描述、状态等扩展信息

```vue
┌─ 第4步：其他信息 ─────────────────────┐
│                                      │
│ 价格信息:                            │
│ 基准售价: [__________________] 元    │
│ 基准租金: [__________________] 元/月 │
│ □ 自动计算单价                       │
│                                      │
│ 户型描述:                            │
│ ┌────────────────────────────────────┐│
│ │                                    ││
│ │                                    ││
│ │                                    ││
│ └────────────────────────────────────┘│
│                                      │
│ 状态设置:                            │
│ 状    态: [▼ 正常] □ 热门户型         │
│                                      │
│ 户型标签:                            │
│ [精装修] [南北通透] [采光好] [+添加]  │
│                                      │
│ [上一步]                    [完成] │
└──────────────────────────────────────┘
```

**字段详情：**
- **`base_sale_price`** - 基准售价 *(decimal 12,2)*
- **`base_rent_price`** - 基准租金 *(decimal 8,2)*
- **`base_sale_price_per`** - 售价单价 *(自动计算)*
- **`base_rent_price_per`** - 租金单价 *(自动计算)*
- **`description`** - 户型描述 *(text)*
- **`status`** - 状态 *(默认active)*
- **`is_hot`** - 是否热门 *(默认false)*
- **`tags`** - 户型标签 *(JSON)*

## 🎨 **前端实现设计**

### **Vue组件结构**
```
HouseTypeForm.vue
├── HouseTypeStep1.vue (基础信息)
├── HouseTypeStep2.vue (户型图片)
├── HouseTypeStep3.vue (户型布局)
└── HouseTypeStep4.vue (其他信息)
```

### **状态管理**
```typescript
interface HouseTypeFormData {
  // 第一步：基础信息
  name: string
  code: string
  standard_area: number
  building_id: number
  
  // 第二步：图片信息
  main_image?: string
  floor_plan_url?: string
  image_urls?: string[]
  
  // 第三步：布局信息
  rooms: number
  halls: number
  bathrooms: number
  balconies: number
  floor_height?: number
  standard_orientation?: string
  standard_view?: string
  
  // 第四步：其他信息
  base_sale_price?: number
  base_rent_price?: number
  description?: string
  status: string
  is_hot: boolean
  tags?: string[]
}
```

### **分步提交逻辑**
```typescript
const formSteps = [
  { step: 1, title: '基础信息', required: true },
  { step: 2, title: '户型图片', required: false },
  { step: 3, title: '户型布局', required: false },
  { step: 4, title: '其他信息', required: false }
]

// 每步都可以保存草稿
const saveDraft = async (stepData: Partial<HouseTypeFormData>) => {
  // 保存到本地存储或发送到服务器
}

// 最终提交
const submitForm = async (formData: HouseTypeFormData) => {
  // 验证必填字段
  // 提交到服务器
  // 清理草稿数据
}
```

## 🔧 **后端API设计**

### **分步提交支持**
```go
// POST /api/v1/house-types (创建)
// PUT /api/v1/house-types/:id (更新)
type CreateHouseTypeRequest struct {
    // 第一步：必填字段
    Name         string  `json:"name" binding:"required,min=1,max=100"`
    Code         string  `json:"code" binding:"required,max=50"`
    StandardArea float64 `json:"standard_area" binding:"required,gt=0"`
    BuildingID   uint    `json:"building_id" binding:"required,gt=0"`
    
    // 第二步：图片字段
    MainImage    *string   `json:"main_image,omitempty"`
    FloorPlanUrl *string   `json:"floor_plan_url,omitempty"`
    ImageUrls    *[]string `json:"image_urls,omitempty"`
    
    // 第三步：布局字段
    Rooms               *int     `json:"rooms,omitempty" binding:"omitempty,gte=1"`
    Halls               *int     `json:"halls,omitempty" binding:"omitempty,gte=1"`
    Bathrooms           *int     `json:"bathrooms,omitempty" binding:"omitempty,gte=1"`
    Balconies           *int     `json:"balconies,omitempty" binding:"omitempty,gte=0"`
    FloorHeight         *float64 `json:"floor_height,omitempty"`
    StandardOrientation *string  `json:"standard_orientation,omitempty"`
    StandardView        *string  `json:"standard_view,omitempty"`
    
    // 第四步：其他字段
    BaseSalePrice *float64 `json:"base_sale_price,omitempty" binding:"omitempty,gte=0"`
    BaseRentPrice *float64 `json:"base_rent_price,omitempty" binding:"omitempty,gte=0"`
    Description   *string  `json:"description,omitempty"`
    Status        *string  `json:"status,omitempty"`
    IsHot         *bool    `json:"is_hot,omitempty"`
    Tags          *[]string `json:"tags,omitempty"`
}
```

## 💡 **用户体验优化**

### **1. 进度指示**
- 步骤导航条显示当前进度
- 已完成步骤显示绿色勾选
- 当前步骤高亮显示

### **2. 数据保存**
- 每步完成后自动保存草稿
- 支持随时退出和继续编辑
- 浏览器关闭前提醒保存

### **3. 智能提示**
- 户型编码格式提示
- 面积合理范围建议
- 价格单价自动计算显示

### **4. 快捷操作**
- 支持键盘导航（Tab键）
- 回车键快速下一步
- 支持跳步编辑（已完成的步骤）

## 🎯 **实施优先级**

### **第一阶段（核心功能）**
1. ✅ 实现第一步基础信息表单
2. ✅ 后端API创建和验证
3. ✅ 基础的表单验证和提交

### **第二阶段（图片功能）**
1. 🔄 图片上传组件
2. 🔄 图片预览和管理
3. 🔄 图片存储和访问

### **第三阶段（完整功能）**
1. 🔄 第三步布局信息表单
2. 🔄 第四步扩展信息表单
3. 🔄 分步导航和状态管理

### **第四阶段（体验优化）**
1. 🔄 草稿保存功能
2. 🔄 智能提示和计算
3. 🔄 表单验证优化

这样的设计既保证了核心功能的快速实现，又为后续的功能扩展留下了充分的空间。您觉得这个优化后的表单设计如何？
