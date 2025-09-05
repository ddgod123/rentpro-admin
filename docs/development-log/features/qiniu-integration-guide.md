# 七牛云集成使用指南

**创建日期：** 2024年12月
**适用版本：** rentpro-admin v1.0.0
**集成状态：** ✅ 配置完成，待启用

## 📋 文件清单

### 🔧 配置文件
- `config/qiniu.yml` - 七牛云主配置文件
- `common/config/qiniu.go` - 配置管理代码
- `common/utils/qiniu.go` - 七牛云服务封装
- `common/initialize/qiniu.go` - 初始化代码
- `examples/qiniu_upload_example.go` - 使用示例

### 📝 环境变量
需要在环境变量中设置以下值：
```bash
QINIU_ACCESS_KEY=your_access_key_here
QINIU_SECRET_KEY=your_secret_key_here
```

## 🚀 快速开始

### 第一步：获取七牛云密钥

1. **注册七牛云账号**
   - 访问 [七牛云官网](https://www.qiniu.com)
   - 完成实名认证

2. **获取Access Key和Secret Key**
   - 进入控制台 → 个人中心 → 密钥管理
   - 创建或查看Access Key和Secret Key

3. **创建存储空间**
   - 进入控制台 → 对象存储
   - 创建新的存储空间（Bucket）
   - 记录存储空间名称和域名

### 第二步：配置项目

1. **修改配置文件 `config/qiniu.yml`**
```yaml
qiniu:
  access_key: "your_access_key_here"          # 替换为真实的Access Key
  secret_key: "your_secret_key_here"          # 替换为真实的Secret Key
  bucket: "rentpro-floor-plans"               # 替换为你的存储空间名称
  domain: "your-domain.com"                   # 替换为你的访问域名
```

2. **设置环境变量**
```bash
# 开发环境
export QINIU_ACCESS_KEY="your_access_key_here"
export QINIU_SECRET_KEY="your_secret_key_here"

# 或者在 .env 文件中设置
QINIU_ACCESS_KEY=your_access_key_here
QINIU_SECRET_KEY=your_secret_key_here
```

3. **安装依赖**
```bash
go mod tidy
```

### 第三步：启用七牛云服务

在 `cmd/api/server.go` 中取消注释：
```go
// 初始化七牛云服务
fmt.Println("初始化七牛云服务...")
err = initialize.InitQiniu(config.Settings.Application.Mode)
if err != nil {
    log.Printf("⚠️  七牛云服务初始化失败: %v", err)
    log.Println("将使用本地文件存储")
}
```

### 第四步：测试服务

```bash
# 启动服务
go run main.go api --port 8002

# 测试健康检查
curl http://localhost:8002/api/v1/qiniu/health
```

## 🔧 配置说明

### 基础配置
```yaml
qiniu:
  access_key: "your_access_key"     # 七牛云Access Key
  secret_key: "your_secret_key"     # 七牛云Secret Key
  bucket: "your-bucket"             # 存储空间名称
  domain: "your-domain.com"         # 访问域名
  zone: "huadong"                   # 存储区域
  use_https: true                   # 是否使用HTTPS
  use_cdn_domains: true             # 是否使用CDN域名
```

### 上传配置
```yaml
upload:
  max_file_size: 5242880            # 最大文件大小 (5MB)
  allowed_types:                    # 允许的文件类型
    - "image/jpeg"
    - "image/png"
    - "image/gif"
  upload_dir: "floor-plans"         # 上传目录前缀
```

### 图片样式配置
```yaml
image_styles:
  thumbnail:                        # 缩略图
    name: "thumb"
    process: "imageView2/1/w/200/h/150/q/85/format/jpg"
  
  medium:                           # 中等尺寸
    name: "medium"
    process: "imageView2/1/w/800/h/600/q/85/format/jpg"
  
  large:                            # 大图
    name: "large"
    process: "imageView2/1/w/1200/h/900/q/90/format/jpg"
```

## 💻 代码使用示例

### 基础上传
```go
// 获取七牛云服务
qiniuService := utils.GetQiniuService()

// 上传文件
uploadResult, err := qiniuService.UploadFile(file, "custom_key.jpg")
if err != nil {
    return err
}

// 获取不同尺寸的URL
originalURL := uploadResult.OriginalURL
thumbnailURL := uploadResult.ThumbnailURL
mediumURL := uploadResult.MediumURL
```

### 删除文件
```go
// 从URL提取key
key := qiniuService.ExtractKeyFromURL(imageURL)

// 删除文件
err := qiniuService.DeleteFile(key)
if err != nil {
    return err
}
```

### 生成样式URL
```go
// 基础URL
baseURL := "https://your-domain.com/floor-plans/image.jpg"

// 生成缩略图URL
thumbnailURL := qiniuService.GetStyleURL(baseURL, "thumbnail")
// 结果: https://your-domain.com/floor-plans/image.jpg-thumb
```

## 🔄 替换现有上传逻辑

### 当前本地上传代码
```go
// 现有的本地上传逻辑 (cmd/api/server.go)
api.POST("/upload/floor-plan", func(c *gin.Context) {
    // ... 保存到本地文件 ...
    filePath := filepath.Join(uploadDir, fileName)
    c.SaveUploadedFile(file, filePath)
    
    // 生成本地URL
    fileURL := fmt.Sprintf("/uploads/floor-plans/%s", fileName)
    
    // 更新数据库
    database.DB.Model(&houseType).Update("floor_plan_url", fileURL)
})
```

### 替换为七牛云上传
```go
// 使用七牛云的上传逻辑
api.POST("/upload/floor-plan", func(c *gin.Context) {
    // ... 验证逻辑保持不变 ...
    
    // 获取七牛云服务
    qiniuService := utils.GetQiniuService()
    if qiniuService == nil {
        // 降级到本地存储
        // ... 原有本地上传逻辑 ...
        return
    }
    
    // 上传到七牛云
    customKey := fmt.Sprintf("floor_plan_%s_%d.jpg", houseTypeID, time.Now().Unix())
    uploadResult, err := qiniuService.UploadFile(file, customKey)
    if err != nil {
        // 降级到本地存储
        // ... 原有本地上传逻辑 ...
        return
    }
    
    // 更新数据库
    database.DB.Model(&houseType).Updates(map[string]interface{}{
        "floor_plan_url": uploadResult.OriginalURL,
        // 可选：存储多个尺寸的URL
        // "floor_plan_thumbnail_url": uploadResult.ThumbnailURL,
        // "floor_plan_medium_url": uploadResult.MediumURL,
    })
    
    // 返回结果
    c.JSON(http.StatusOK, gin.H{
        "code": 200,
        "message": "户型图上传成功",
        "data": gin.H{
            "original_url":  uploadResult.OriginalURL,
            "thumbnail_url": uploadResult.ThumbnailURL,
            "medium_url":    uploadResult.MediumURL,
            "large_url":     uploadResult.LargeURL,
        },
    })
})
```

## 🎨 前端配置

### 修改前端上传配置
```javascript
// 修改 FloorPlanForm.vue 中的上传配置
const uploadAction = ref('/api/v1/upload/floor-plan')  // 保持不变

// 处理上传成功，现在可以获得多个尺寸的URL
const handleUploadSuccess = (response) => {
  if (response.code === 200) {
    ElMessage.success('户型图上传成功')
    
    // 可以使用不同尺寸的图片
    console.log('原图:', response.data.original_url)
    console.log('缩略图:', response.data.thumbnail_url)
    console.log('中等尺寸:', response.data.medium_url)
    
    emit('success')
  }
}
```

### 使用不同尺寸的图片
```vue
<template>
  <!-- 列表页使用缩略图 -->
  <el-image :src="getThumbnailURL(row.floor_plan_url)" />
  
  <!-- 详情页使用中等尺寸 -->
  <el-image :src="getMediumURL(houseType.floor_plan_url)" />
  
  <!-- 预览使用大图 -->
  <el-image 
    :src="getMediumURL(houseType.floor_plan_url)"
    :preview-src-list="[getLargeURL(houseType.floor_plan_url)]"
  />
</template>

<script>
// 生成不同尺寸的URL
const getThumbnailURL = (originalURL) => {
  return originalURL + '-thumb'
}

const getMediumURL = (originalURL) => {
  return originalURL + '-medium'  
}

const getLargeURL = (originalURL) => {
  return originalURL + '-large'
}
</script>
```

## 🔍 故障排查

### 常见问题

1. **初始化失败**
```
错误: 七牛云配置验证失败: Access Key 未配置
解决: 检查配置文件中的access_key是否正确设置
```

2. **上传失败**
```
错误: 上传到七牛云失败: no such bucket
解决: 检查bucket名称是否正确，确保存储空间已创建
```

3. **域名访问失败**
```
错误: 图片无法访问
解决: 检查domain配置，确保域名已绑定到存储空间
```

### 调试模式
```go
// 在开发环境启用详细日志
log.SetLevel(log.DebugLevel)

// 检查七牛云服务状态
qiniuService := utils.GetQiniuService()
if qiniuService != nil {
    log.Println("七牛云服务已初始化")
} else {
    log.Println("七牛云服务未初始化")
}
```

## 📊 性能优化建议

### 1. 图片处理优化
```yaml
# 在七牛云控制台配置自动WebP转换
image_styles:
  webp_thumb:
    name: "webp-thumb"
    process: "imageView2/1/w/200/h/150/q/85/format/webp"
```

### 2. 缓存策略
```go
// 设置图片缓存头
c.Header("Cache-Control", "public, max-age=31536000")  // 1年
```

### 3. CDN配置
- 在七牛云控制台开启CDN加速
- 配置HTTPS证书
- 设置缓存规则

## 💰 成本控制

### 监控用量
```go
// 定期检查存储用量
files, err := qiniuService.ListFiles("", 1000)
if err == nil {
    log.Printf("当前文件数量: %d", len(files))
}
```

### 清理策略
```go
// 定期清理临时文件和无效文件
func CleanupOldFiles() {
    // 实现清理逻辑
}
```

## 🎉 部署检查清单

### 部署前检查
- [ ] 七牛云账号已注册并实名认证
- [ ] 存储空间已创建
- [ ] 域名已绑定
- [ ] Access Key和Secret Key已获取
- [ ] 配置文件已正确填写
- [ ] 环境变量已设置
- [ ] 依赖包已安装

### 部署后验证
- [ ] 服务启动成功
- [ ] 七牛云服务初始化成功
- [ ] 健康检查通过
- [ ] 上传功能正常
- [ ] 图片访问正常
- [ ] 不同尺寸图片正常生成

## 📞 技术支持

### 官方文档
- [七牛云对象存储文档](https://developer.qiniu.com/kodo)
- [Go SDK文档](https://developer.qiniu.com/kodo/1238/go)
- [图片处理文档](https://developer.qiniu.com/dora)

### 社区支持
- [七牛云开发者社区](https://segmentfault.com/t/%E4%B8%83%E7%89%9B%E4%BA%91)
- [GitHub Issues](https://github.com/qiniu/go-sdk)

现在您已经有了完整的七牛云集成方案，可以根据实际需要启用和配置！🎨✨
