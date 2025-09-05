# 七牛云图片存储服务选购指南

**指南日期：** 2024年12月
**适用场景：** 户型图存储和管理
**推荐配置：** 对象存储 + 图片处理 + CDN

## 🎯 推荐购买的七牛云服务

### 1️⃣ **对象存储 (Kodo) - 核心必选**

#### **服务说明**
- **产品名称：** 对象存储 Kodo
- **功能：** 文件存储、管理、访问
- **适用：** 户型图片、文档等静态资源存储

#### **免费额度（每月）**
```
✅ 存储空间：10GB
✅ 下载流量：10GB  
✅ PUT/DELETE 请求：10万次
✅ GET 请求：100万次
```

#### **付费标准（超出免费额度后）**
| 项目 | 价格 | 说明 |
|------|------|------|
| 存储费用 | 0.148元/GB/月 | 标准存储 |
| 下载流量 | 0.29元/GB | 外网下行流量 |
| PUT/DELETE请求 | 0.01元/万次 | 上传/删除操作 |
| GET请求 | 0.01元/10万次 | 下载/访问操作 |

### 2️⃣ **智能多媒体服务 (DORA) - 强烈推荐**

#### **服务说明**
- **产品名称：** 智能多媒体服务 DORA
- **功能：** 图片处理、格式转换、压缩、水印
- **适用：** 户型图自动优化、缩略图生成

#### **免费额度（每月）**
```
✅ 图片处理：10万次
✅ 图片瘦身：1万次
✅ 内容审核：1000次
```

#### **付费标准**
| 功能 | 价格 | 说明 |
|------|------|------|
| 基础图片处理 | 0.025元/千次 | 缩放、裁剪、旋转等 |
| 图片格式转换 | 0.025元/千次 | JPG/PNG/WebP转换 |
| 图片瘦身 | 0.15元/千次 | 智能压缩 |

### 3️⃣ **CDN 加速服务 - 推荐配置**

#### **服务说明**
- **产品名称：** CDN 内容分发网络
- **功能：** 全球加速访问、降低延迟
- **适用：** 提升户型图加载速度

#### **免费额度（每月）**
```
✅ CDN流量：10GB
✅ HTTPS请求：100万次
```

#### **付费标准**
| 流量阶梯 | 价格 | 说明 |
|----------|------|------|
| 0-50GB | 0.24元/GB | 首档价格 |
| 50GB-500GB | 0.23元/GB | 第二档 |
| 500GB-1TB | 0.22元/GB | 第三档 |

## 💰 成本预估

### 📊 **小型项目（月活<1000用户）**

**预估使用量：**
- 存储空间：5GB（约5000张户型图）
- 月下载流量：8GB
- 图片处理：5万次/月
- CDN流量：8GB

**月费用预估：**
```
✅ 对象存储：0元（在免费额度内）
✅ 图片处理：0元（在免费额度内）
✅ CDN加速：0元（在免费额度内）

总计：0元/月 🎉
```

### 📈 **中型项目（月活1000-5000用户）**

**预估使用量：**
- 存储空间：25GB（约25000张户型图）
- 月下载流量：30GB
- 图片处理：20万次/月
- CDN流量：30GB

**月费用预估：**
```
📦 对象存储：
   - 存储费：(25-10) × 0.148 = 2.22元
   - 流量费：(30-10) × 0.29 = 5.8元

🎨 图片处理：
   - 处理费：(20-10) × 0.025 = 0.25元

🚀 CDN加速：
   - 流量费：(30-10) × 0.24 = 4.8元

总计：约13元/月 💰
```

### 🚀 **大型项目（月活>5000用户）**

**预估使用量：**
- 存储空间：100GB（约10万张户型图）
- 月下载流量：200GB
- 图片处理：100万次/月
- CDN流量：200GB

**月费用预估：**
```
📦 对象存储：
   - 存储费：(100-10) × 0.148 = 13.32元
   - 流量费：(200-10) × 0.29 = 55.1元

🎨 图片处理：
   - 处理费：(100-10) × 0.025 = 2.25元

🚀 CDN加速：
   - 流量费：200 × 0.24 = 48元

总计：约119元/月 💰
```

## 🛒 购买步骤指南

### **第一步：注册和实名认证**
1. 访问 [七牛云官网](https://www.qiniu.com)
2. 注册账号并完成实名认证
3. 绑定手机号和邮箱

### **第二步：开通对象存储服务**
1. 进入控制台 → 对象存储
2. 创建存储空间（Bucket）
   - **空间名称：** `rentpro-floor-plans`
   - **存储区域：** 华东-浙江（推荐，国内访问速度快）
   - **访问控制：** 公开空间（便于图片访问）
3. 绑定自定义域名（可选，推荐）

### **第三步：开通图片处理服务**
1. 在存储空间设置中
2. 开启"图片样式"功能
3. 配置常用的图片处理样式：
   - `thumb`: 缩略图 200x150
   - `medium`: 中等尺寸 800x600
   - `large`: 大图 1200x900

### **第四步：配置CDN加速（推荐）**
1. 进入控制台 → CDN
2. 添加加速域名
3. 配置CNAME解析
4. 开启HTTPS（如果需要）

## 🔧 技术集成

### **Go SDK安装**
```bash
go get github.com/qiniu/go-sdk/v7/storage
```

### **基础配置代码**
```go
package main

import (
    "context"
    "fmt"
    "github.com/qiniu/go-sdk/v7/auth/qbox"
    "github.com/qiniu/go-sdk/v7/storage"
)

// 七牛云配置
type QiniuConfig struct {
    AccessKey string
    SecretKey string
    Bucket    string
    Domain    string
}

// 初始化七牛云客户端
func NewQiniuClient(config QiniuConfig) *QiniuClient {
    mac := qbox.NewMac(config.AccessKey, config.SecretKey)
    
    cfg := storage.Config{
        Zone:          &storage.ZoneHuadong,    // 华东区域
        UseHTTPS:      true,                    // 使用HTTPS
        UseCdnDomains: true,                    // 使用CDN域名
    }
    
    return &QiniuClient{
        mac:    mac,
        config: cfg,
        bucket: config.Bucket,
        domain: config.Domain,
    }
}

type QiniuClient struct {
    mac    *qbox.Mac
    config storage.Config
    bucket string
    domain string
}

// 上传文件
func (q *QiniuClient) UploadFile(localFile string, key string) (string, error) {
    uploader := storage.NewFormUploader(&q.config)
    ret := storage.PutRet{}
    
    // 生成上传Token
    putPolicy := storage.PutPolicy{
        Scope: q.bucket,
    }
    upToken := putPolicy.UploadToken(q.mac)
    
    // 执行上传
    err := uploader.PutFile(context.Background(), &ret, upToken, key, localFile, nil)
    if err != nil {
        return "", err
    }
    
    // 返回访问URL
    return fmt.Sprintf("https://%s/%s", q.domain, key), nil
}

// 生成图片处理URL
func (q *QiniuClient) GetImageURL(key string, style string) string {
    baseURL := fmt.Sprintf("https://%s/%s", q.domain, key)
    if style != "" {
        return fmt.Sprintf("%s-%s", baseURL, style)
    }
    return baseURL
}
```

### **集成到现有项目**
```go
// 在现有的上传API中集成七牛云
api.POST("/upload/floor-plan", func(c *gin.Context) {
    // ... 现有的验证逻辑 ...
    
    // 保存到本地临时文件
    tempFile := filepath.Join("/tmp", fileName)
    c.SaveUploadedFile(file, tempFile)
    defer os.Remove(tempFile) // 清理临时文件
    
    // 上传到七牛云
    qiniuKey := fmt.Sprintf("floor-plans/%s", fileName)
    imageURL, err := qiniuClient.UploadFile(tempFile, qiniuKey)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code":    500,
            "message": "上传到七牛云失败",
            "error":   err.Error(),
        })
        return
    }
    
    // 生成不同尺寸的URL
    thumbnailURL := qiniuClient.GetImageURL(qiniuKey, "thumb")
    mediumURL := qiniuClient.GetImageURL(qiniuKey, "medium")
    
    // 更新数据库
    database.DB.Model(&houseType).Updates(map[string]interface{}{
        "floor_plan_url": imageURL,
        "floor_plan_thumbnail_url": thumbnailURL,
        "floor_plan_medium_url": mediumURL,
    })
    
    c.JSON(http.StatusOK, gin.H{
        "code":    200,
        "message": "户型图上传成功",
        "data": gin.H{
            "original_url":  imageURL,
            "thumbnail_url": thumbnailURL,
            "medium_url":    mediumURL,
        },
    })
})
```

## 🎨 图片处理样式配置

### **在七牛云控制台配置样式**

1. **缩略图样式 (thumb)**
   ```
   imageView2/1/w/200/h/150/q/85/format/jpg
   ```

2. **中等尺寸样式 (medium)**
   ```
   imageView2/1/w/800/h/600/q/85/format/jpg
   ```

3. **大图样式 (large)**
   ```
   imageView2/1/w/1200/h/900/q/90/format/jpg
   ```

### **前端使用示例**
```vue
<template>
  <!-- 列表页缩略图 -->
  <el-image :src="row.floor_plan_url + '-thumb'" />
  
  <!-- 详情页中等图片 -->
  <el-image :src="houseType.floor_plan_url + '-medium'" />
  
  <!-- 预览大图 -->
  <el-image 
    :src="houseType.floor_plan_url + '-medium'"
    :preview-src-list="[houseType.floor_plan_url + '-large']"
  />
</template>
```

## 💡 最佳实践建议

### **1. 成本优化**
- 🆓 **充分利用免费额度** - 小项目完全免费
- 📦 **合理规划存储** - 定期清理无用图片
- 🗜️ **使用图片压缩** - 减少存储和流量成本

### **2. 性能优化**
- 🚀 **启用CDN加速** - 提升全球访问速度
- 🎨 **配置图片样式** - 按需加载不同尺寸
- 📱 **WebP格式支持** - 现代浏览器优化

### **3. 安全建议**
- 🔐 **AccessKey安全** - 使用环境变量存储
- 🌐 **域名绑定** - 使用自己的域名
- 🔒 **HTTPS配置** - 确保传输安全

## 🎉 总结

### ✅ **推荐购买清单**
1. **对象存储 Kodo** - 必选，基础存储服务
2. **智能多媒体服务 DORA** - 推荐，图片处理功能
3. **CDN加速服务** - 推荐，提升访问速度

### 💰 **成本优势**
- **小项目完全免费** - 月免费额度足够使用
- **按需付费** - 用多少付多少，无浪费
- **价格透明** - 无隐藏费用

### 🚀 **技术优势**
- **Go SDK完善** - 集成简单
- **图片处理强大** - 实时处理，无需预生成
- **CDN全球加速** - 访问速度快

**对于您的户型图管理项目，建议从免费额度开始使用，随着业务增长再考虑付费！** 🎨✨
