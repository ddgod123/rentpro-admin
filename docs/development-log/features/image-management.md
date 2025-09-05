# 🖼️ 图片管理系统设计与实现

**创建日期：** 2024年12月
**适用版本：** rentpro-admin v1.0.0
**设计目标：** 基于七牛云的完整图片管理解决方案

## 📋 目录结构

```
rentpro-admin-main/
├── common/models/image/           # 图片数据模型
│   └── sys_image.go              # 图片和分类模型定义
├── common/utils/
│   ├── qiniu.go                  # 七牛云基础服务
│   └── image_manager.go          # 图片管理器
├── config/sql/migrations/
│   └── create_images_table.sql   # 数据库迁移脚本
├── cmd/api/server.go             # API接口实现
├── examples/
│   └── image_management_test.go  # 测试示例
└── docs/development-log/features/
    └── image-management.md       # 本文档
```

## 🏗️ 系统架构设计

### 1. 核心组件

#### 🖼️ 图片管理器 (ImageManager)
```go
type ImageManager struct {
    qiniuService *QiniuService
    db           *gorm.DB
}
```

**主要功能：**
- 文件上传到七牛云
- 自动生成多种尺寸图片
- 数据库记录管理
- 图片分类和模块关联
- 主图设置和排序

#### 📊 数据模型设计

**图片主表 (sys_images)：**
```sql
CREATE TABLE sys_images (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    name            VARCHAR(200) NOT NULL,           -- 图片名称
    file_name       VARCHAR(255) NOT NULL,           -- 原始文件名
    file_size       BIGINT NOT NULL,                 -- 文件大小
    mime_type       VARCHAR(100) NOT NULL,           -- MIME类型
    extension       VARCHAR(10) NOT NULL,            -- 文件扩展名

    -- 存储路径
    key             VARCHAR(500) NOT NULL,           -- 七牛云Key
    url             VARCHAR(1000) NOT NULL,          -- 原始URL
    thumbnail_url   VARCHAR(1000),                   -- 缩略图URL
    medium_url      VARCHAR(1000),                   -- 中等尺寸URL
    large_url       VARCHAR(1000),                   -- 大图URL

    -- 分类信息
    category        VARCHAR(50) NOT NULL,            -- 分类
    module          VARCHAR(50) NOT NULL,            -- 所属模块
    module_id       BIGINT,                          -- 模块ID

    -- 状态控制
    is_public       BOOLEAN DEFAULT TRUE,            -- 是否公开
    is_main         BOOLEAN DEFAULT FALSE,           -- 是否主图
    sort_order      INT DEFAULT 0,                   -- 排序
    status          VARCHAR(20) DEFAULT 'active',    -- 状态

    -- 审计字段
    created_by      BIGINT,
    updated_by      BIGINT,
    created_at      DATETIME,
    updated_at      DATETIME,
    deleted_at      DATETIME
);
```

**分类配置表 (sys_image_categories)：**
```sql
CREATE TABLE sys_image_categories (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    code            VARCHAR(50) UNIQUE NOT NULL,     -- 分类编码
    name            VARCHAR(100) NOT NULL,           -- 分类名称
    description     VARCHAR(200),                    -- 描述
    max_size        BIGINT DEFAULT 5242880,          -- 最大文件大小
    allowed_types   JSON,                            -- 允许的文件类型
    max_count       INT DEFAULT 10,                  -- 最大上传数量
    is_required     BOOLEAN DEFAULT FALSE,           -- 是否必填
    status          VARCHAR(20) DEFAULT 'active'     -- 状态
);
```

### 2. 设计模式

#### 🎯 策略模式 - 图片处理
```go
// 不同分类的处理策略
type ImageProcessor interface {
    Process(file *multipart.FileHeader) (*UploadResult, error)
    Validate(file *multipart.FileHeader) error
}

type BuildingImageProcessor struct{}  // 楼盘图片处理器
type AvatarImageProcessor struct{}    // 头像图片处理器
```

#### 🔧 工厂模式 - 处理器创建
```go
func CreateImageProcessor(category string) ImageProcessor {
    switch category {
    case "building":
        return &BuildingImageProcessor{}
    case "avatar":
        return &AvatarImageProcessor{}
    default:
        return &DefaultImageProcessor{}
    }
}
```

## 🚀 API 接口设计

### 1. 基础 CRUD 接口

#### 📤 上传图片
```http
POST /api/v1/images/upload
Content-Type: multipart/form-data

Form Data:
- file: <图片文件>
- category: building|house|avatar|banner
- module: rental|user|system
- moduleId: <关联ID>
- isMain: true|false
- isPublic: true|false
```

**响应示例：**
```json
{
    "code": 200,
    "message": "图片上传成功",
    "data": {
        "id": 1,
        "name": "building_001.jpg",
        "url": "https://cdn.domain.com/images/building_001.jpg",
        "thumbnailUrl": "https://cdn.domain.com/images/building_001.jpg-thumb",
        "mediumUrl": "https://cdn.domain.com/images/building_001.jpg-medium",
        "largeUrl": "https://cdn.domain.com/images/building_001.jpg-large",
        "category": "building",
        "module": "rental",
        "isMain": false,
        "createdAt": "2024-12-01T10:00:00Z"
    }
}
```

#### 📋 获取图片列表
```http
GET /api/v1/images?page=1&pageSize=10&category=building&module=rental&moduleId=123
```

#### 📖 获取图片详情
```http
GET /api/v1/images/1
```

#### ✏️ 更新图片信息
```http
PUT /api/v1/images/1
Content-Type: application/json

{
    "name": "新图片名称",
    "description": "图片描述",
    "category": "building",
    "isMain": true,
    "isPublic": true,
    "sortOrder": 1
}
```

#### 🗑️ 删除图片
```http
DELETE /api/v1/images/1
```

#### 📦 批量删除
```http
DELETE /api/v1/images/batch
Content-Type: application/json

{
    "ids": [1, 2, 3, 4, 5]
}
```

### 2. 高级功能接口

#### 🎯 设置主图
```http
PUT /api/v1/images/1/set-main
Content-Type: application/json

{
    "module": "rental",
    "moduleId": 123
}
```

#### 📊 获取模块图片
```http
GET /api/v1/images/module/rental/123?category=building
```

#### 📈 获取统计信息
```http
GET /api/v1/images/stats
```

**统计响应：**
```json
{
    "code": 200,
    "data": {
        "totalImages": 1250,
        "totalSize": 524288000,
        "categoryStats": {
            "building": 450,
            "house": 380,
            "avatar": 120,
            "banner": 50,
            "floor_plan": 200,
            "certificate": 50
        },
        "moduleStats": {
            "rental": 1030,
            "user": 120,
            "system": 100
        },
        "todayUploads": 15,
        "storageUsed": 524288000
    }
}
```

## 🎨 图片处理策略

### 1. 七牛云图片样式配置

**缩略图样式 (thumbnail):**
```yaml
image_styles:
  thumbnail:
    name: "thumb"
    process: "imageView2/1/w/200/h/150/q/85/format/jpg"
    description: "缩略图 200x150"
```

**中等尺寸样式 (medium):**
```yaml
medium:
  name: "medium"
  process: "imageView2/1/w/800/h/600/q/85/format/jpg"
  description: "中等尺寸 800x600"
```

**大图样式 (large):**
```yaml
large:
  name: "large"
  process: "imageView2/1/w/1200/h/900/q/90/format/jpg"
  description: "大图 1200x900"
```

### 2. 分类处理策略

**楼盘图片 (building):**
- 支持格式：JPEG, PNG, GIF, WebP
- 最大尺寸：5MB
- 建议尺寸：1200x900
- 生成样式：缩略图、中等、大图

**头像图片 (avatar):**
- 支持格式：JPEG, PNG
- 最大尺寸：2MB
- 建议尺寸：200x200
- 生成样式：缩略图、中等

**横幅图片 (banner):**
- 支持格式：JPEG, PNG
- 最大尺寸：3MB
- 建议尺寸：1920x600
- 生成样式：中等、大图

## 🔐 权限控制设计

### 1. 基于角色的访问控制

```go
// 权限检查中间件
func ImagePermissionMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetUint64("user_id")
        roleID := c.GetUint64("role_id")

        // 检查用户是否有图片管理权限
        if !hasImagePermission(userID, roleID) {
            c.JSON(403, gin.H{
                "code": 403,
                "message": "没有图片管理权限",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### 2. 私有图片访问控制

```go
// 生成私有图片访问URL
func (q *QiniuService) GeneratePrivateURL(key string, expires int64) string {
    deadline := time.Now().Add(time.Duration(expires) * time.Second).Unix()
    privateURL := storage.MakePrivateURL(q.mac, q.domain, key, deadline)
    return privateURL
}
```

## 📊 监控和统计

### 1. 存储使用量监控

```go
type StorageMonitor struct {
    qiniuService *QiniuService
    db           *gorm.DB
}

// 获取存储统计
func (sm *StorageMonitor) GetStorageStats() (*StorageStats, error) {
    // 查询数据库统计
    var stats StorageStats
    sm.db.Model(&SysImage{}).Select(
        "COUNT(*) as total_files",
        "SUM(file_size) as total_size",
        "COUNT(CASE WHEN created_at >= CURDATE() THEN 1 END) as today_uploads",
    ).Scan(&stats)

    return &stats, nil
}
```

### 2. 性能监控

```go
type PerformanceMonitor struct {
    uploadCount    int64
    uploadTime     time.Duration
    errorCount     int64
    lastUploadTime time.Time
}

// 记录上传性能
func (pm *PerformanceMonitor) RecordUpload(duration time.Duration, success bool) {
    pm.uploadCount++
    pm.uploadTime += duration
    pm.lastUploadTime = time.Now()

    if !success {
        pm.errorCount++
    }
}
```

## 🧪 测试策略

### 1. 单元测试

```go
func TestImageManager_UploadImage(t *testing.T) {
    // 创建模拟文件
    file := createMockMultipartFile("test.jpg", 1024)

    // 创建上传请求
    req := &ImageUploadRequest{
        Category: "building",
        Module:   "rental",
        ModuleID: 123,
        IsMain:   false,
        IsPublic: true,
    }

    // 执行上传
    imageManager := NewImageManager()
    result, err := imageManager.UploadImage(file, req, 1)

    // 断言结果
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "building", result.Category)
}
```

### 2. 集成测试

```go
func TestImageAPI_UploadFlow(t *testing.T) {
    // 启动测试服务器
    router := setupTestRouter()

    // 创建测试文件
    file := createTestImageFile()

    // 发送上传请求
    req := createUploadRequest(file)
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    // 验证响应
    assert.Equal(t, 200, w.Code)

    var response ImageUploadResponse
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.True(t, response.Success)
}
```

## 🚀 部署和运维

### 1. 数据库迁移

```bash
# 执行数据库迁移
mysql -u root -p rentpro < config/sql/migrations/create_images_table.sql
```

### 2. 配置检查清单

- [ ] 七牛云账号配置正确
- [ ] Access Key 和 Secret Key 已设置
- [ ] 存储空间已创建
- [ ] 域名已绑定
- [ ] 图片样式已配置
- [ ] 数据库表已创建
- [ ] API 接口已注册

### 3. 监控指标

**业务指标：**
- 图片上传成功率
- 平均上传耗时
- 存储使用量
- 分类分布统计

**系统指标：**
- API 响应时间
- 错误率统计
- 七牛云服务可用性
- 数据库连接池状态

## 📈 扩展计划

### 1. 功能扩展

- [ ] 图片压缩优化
- [ ] 批量水印处理
- [ ] 智能裁剪功能
- [ ] 图片审核服务
- [ ] CDN 分发优化

### 2. 性能优化

- [ ] 图片懒加载
- [ ] WebP 格式支持
- [ ] 缓存策略优化
- [ ] 并发上传支持

### 3. 安全增强

- [ ] 图片内容审核
- [ ] 上传频率限制
- [ ] 文件类型深度检测
- [ ] 私有空间加密

## 🎯 使用示例

### 前端调用示例

```javascript
// 上传图片
const uploadImage = async (file, category, moduleId) => {
  const formData = new FormData();
  formData.append('file', file);
  formData.append('category', category);
  formData.append('module', 'rental');
  formData.append('moduleId', moduleId);
  formData.append('isMain', false);
  formData.append('isPublic', true);

  const response = await fetch('/api/v1/images/upload', {
    method: 'POST',
    body: formData,
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });

  return response.json();
};

// 获取图片列表
const getImages = async (category, moduleId) => {
  const response = await fetch(
    `/api/v1/images?category=${category}&module=rental&moduleId=${moduleId}`
  );
  return response.json();
};
```

这个图片管理系统提供了完整的图片管理解决方案，集成了七牛云的高性能存储服务，支持多种业务场景的图片管理需求。通过合理的架构设计和丰富的功能特性，可以有效提升应用的图片处理能力和用户体验。🎨✨
