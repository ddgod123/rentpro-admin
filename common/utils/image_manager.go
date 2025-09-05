package utils

import (
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"rentPro/rentpro-admin/common/database"
	"rentPro/rentpro-admin/common/models/image"

	"gorm.io/gorm"
)

// ImageManager 图片管理器
type ImageManager struct {
	qiniuService *QiniuService
	db           *gorm.DB
}

// NewImageManager 创建图片管理器
func NewImageManager() (*ImageManager, error) {
	qiniuService := GetQiniuService()
	if qiniuService == nil {
		return nil, fmt.Errorf("七牛云服务未初始化")
	}

	return &ImageManager{
		qiniuService: qiniuService,
		db:           database.DB,
	}, nil
}

// UploadImage 上传图片
func (im *ImageManager) UploadImage(file *multipart.FileHeader, req *image.ImageUploadRequest, userID uint64) (*image.SysImage, error) {
	// 验证文件
	if err := im.qiniuService.ValidateFile(file); err != nil {
		return nil, err
	}

	// 生成存储Key，支持楼盘文件夹结构
	fileName := fmt.Sprintf("%s_%d_%s", req.Category, time.Now().UnixNano(), file.Filename)
	var customKey string

	// 如果是楼盘相关的图片，使用楼盘文件夹结构
	if req.Module == "building" || req.Module == "house" {
		if req.ModuleID > 0 {
			// 格式: buildings/{buildingId}/{category}/{timestamp}_{filename}
			customKey = fmt.Sprintf("buildings/%d/%s/%s", req.ModuleID, req.Category, fileName)
		} else {
			// 如果没有指定楼盘ID，使用通用楼盘文件夹
			customKey = fmt.Sprintf("buildings/common/%s/%s", req.Category, fileName)
		}
	} else {
		// 其他模块使用原有逻辑
		customKey = im.qiniuService.configManager.GetUploadKey(fileName)
	}

	// 上传到七牛云
	uploadResult, err := im.qiniuService.UploadFile(file, customKey)
	if err != nil {
		return nil, fmt.Errorf("上传到七牛云失败: %v", err)
	}

	// 保存到数据库
	img := &image.SysImage{
		Name:         file.Filename,
		Description:  "",
		FileName:     file.Filename,
		FileSize:     file.Size,
		MimeType:     file.Header.Get("Content-Type"),
		Extension:    im.getFileExtension(file.Filename),
		Key:          uploadResult.Key,
		URL:          uploadResult.OriginalURL,
		ThumbnailURL: uploadResult.ThumbnailURL,
		MediumURL:    uploadResult.MediumURL,
		LargeURL:     uploadResult.LargeURL,
		Category:     req.Category,
		Module:       req.Module,
		ModuleID:     req.ModuleID,
		IsPublic:     req.IsPublic,
		IsMain:       req.IsMain,
		SortOrder:    0,
		Status:       "active",
		CreatedBy:    userID,
		UpdatedBy:    userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := im.db.Create(img).Error; err != nil {
		// 如果数据库保存失败，删除已上传的文件
		im.qiniuService.DeleteFile(uploadResult.Key)
		return nil, fmt.Errorf("保存到数据库失败: %v", err)
	}

	return img, nil
}

// GetImage 获取图片信息
func (im *ImageManager) GetImage(id uint64) (*image.SysImage, error) {
	var img image.SysImage
	if err := im.db.Where("id = ?", id).First(&img).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("图片不存在")
		}
		return nil, fmt.Errorf("查询图片失败: %v", err)
	}
	return &img, nil
}

// GetImagesByModule 根据模块获取图片列表
func (im *ImageManager) GetImagesByModule(module string, moduleID uint64, category string) ([]*image.SysImage, error) {
	var images []*image.SysImage
	query := im.db.Where("module = ? AND module_id = ? AND status = 'active'", module, moduleID)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Order("is_main DESC, sort_order ASC, created_at DESC").Find(&images).Error; err != nil {
		return nil, fmt.Errorf("查询图片列表失败: %v", err)
	}

	return images, nil
}

// UpdateImage 更新图片信息
func (im *ImageManager) UpdateImage(id uint64, req *image.ImageUpdateRequest, userID uint64) error {
	updateData := map[string]interface{}{
		"updated_by": userID,
		"updated_at": time.Now(),
	}

	if req.Name != "" {
		updateData["name"] = req.Name
	}
	if req.Description != "" {
		updateData["description"] = req.Description
	}
	if req.Category != "" {
		updateData["category"] = req.Category
	}
	if req.Status != "" {
		updateData["status"] = req.Status
	}
	updateData["is_main"] = req.IsMain
	updateData["is_public"] = req.IsPublic
	updateData["sort_order"] = req.SortOrder

	if err := im.db.Model(&image.SysImage{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
		return fmt.Errorf("更新图片信息失败: %v", err)
	}

	return nil
}

// DeleteImage 删除图片
func (im *ImageManager) DeleteImage(id uint64, userID uint64) error {
	var img image.SysImage
	if err := im.db.Where("id = ?", id).First(&img).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("图片不存在")
		}
		return fmt.Errorf("查询图片失败: %v", err)
	}

	// 从七牛云删除文件
	if err := im.qiniuService.DeleteFile(img.Key); err != nil {
		// 记录错误，但不阻止数据库删除
		fmt.Printf("删除七牛云文件失败: %v\n", err)
	}

	// 从数据库删除记录
	if err := im.db.Delete(&img).Error; err != nil {
		return fmt.Errorf("删除图片记录失败: %v", err)
	}

	return nil
}

// BatchDeleteImages 批量删除图片
func (im *ImageManager) BatchDeleteImages(ids []uint64, userID uint64) error {
	var images []image.SysImage
	if err := im.db.Where("id IN (?)", ids).Find(&images).Error; err != nil {
		return fmt.Errorf("查询图片失败: %v", err)
	}

	// 删除七牛云文件
	for _, img := range images {
		if err := im.qiniuService.DeleteFile(img.Key); err != nil {
			fmt.Printf("删除七牛云文件失败 [%s]: %v\n", img.Key, err)
		}
	}

	// 批量删除数据库记录
	if err := im.db.Where("id IN (?)", ids).Delete(&image.SysImage{}).Error; err != nil {
		return fmt.Errorf("批量删除图片记录失败: %v", err)
	}

	return nil
}

// ListImages 获取图片列表
func (im *ImageManager) ListImages(req *image.ImageListRequest) (*image.ImageListResponse, error) {
	var images []*image.SysImage
	var total int64

	query := im.db.Model(&image.SysImage{})

	// 构建查询条件
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	if req.Module != "" {
		query = query.Where("module = ?", req.Module)
	}
	if req.ModuleID > 0 {
		query = query.Where("module_id = ?", req.ModuleID)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.IsMain != nil {
		query = query.Where("is_main = ?", *req.IsMain)
	}
	if req.IsPublic != nil {
		query = query.Where("is_public = ?", *req.IsPublic)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("获取总数失败: %v", err)
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	orderBy := req.OrderBy
	if req.OrderDir == "asc" {
		orderBy += " ASC"
	} else {
		orderBy += " DESC"
	}

	if err := query.Offset(offset).Limit(req.PageSize).Order(orderBy).Find(&images).Error; err != nil {
		return nil, fmt.Errorf("查询图片列表失败: %v", err)
	}

	return &image.ImageListResponse{
		Total: total,
		List:  images,
	}, nil
}

// GetImageStats 获取图片统计信息
func (im *ImageManager) GetImageStats() (*image.ImageStats, error) {
	stats := &image.ImageStats{
		CategoryStats: make(map[string]int64),
		ModuleStats:   make(map[string]int64),
	}

	// 总图片数和总大小
	var result struct {
		TotalImages int64
		TotalSize   int64
	}
	if err := im.db.Model(&image.SysImage{}).Where("status = 'active'").Select("COUNT(*) as total_images, SUM(file_size) as total_size").Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("获取统计信息失败: %v", err)
	}
	stats.TotalImages = result.TotalImages
	stats.TotalSize = result.TotalSize
	stats.StorageUsed = result.TotalSize

	// 分类统计
	var categoryStats []struct {
		Category string
		Count    int64
	}
	if err := im.db.Model(&image.SysImage{}).Where("status = 'active'").Select("category, COUNT(*) as count").Group("category").Scan(&categoryStats).Error; err != nil {
		return nil, fmt.Errorf("获取分类统计失败: %v", err)
	}
	for _, stat := range categoryStats {
		stats.CategoryStats[stat.Category] = stat.Count
	}

	// 模块统计
	var moduleStats []struct {
		Module string
		Count  int64
	}
	if err := im.db.Model(&image.SysImage{}).Where("status = 'active'").Select("module, COUNT(*) as count").Group("module").Scan(&moduleStats).Error; err != nil {
		return nil, fmt.Errorf("获取模块统计失败: %v", err)
	}
	for _, stat := range moduleStats {
		stats.ModuleStats[stat.Module] = stat.Count
	}

	// 今日上传数
	today := time.Now().Format("2006-01-02")
	var todayUploads int64
	if err := im.db.Model(&image.SysImage{}).Where("DATE(created_at) = ? AND status = 'active'", today).Count(&todayUploads).Error; err != nil {
		return nil, fmt.Errorf("获取今日上传数失败: %v", err)
	}
	stats.TodayUploads = todayUploads

	return stats, nil
}

// SetMainImage 设置主图
func (im *ImageManager) SetMainImage(module string, moduleID uint64, imageID uint64, userID uint64) error {
	// 取消其他主图
	if err := im.db.Model(&image.SysImage{}).Where("module = ? AND module_id = ? AND id != ?", module, moduleID, imageID).Update("is_main", false).Error; err != nil {
		return fmt.Errorf("取消其他主图失败: %v", err)
	}

	// 设置新主图
	if err := im.db.Model(&image.SysImage{}).Where("id = ?", imageID).Updates(map[string]interface{}{
		"is_main":    true,
		"updated_by": userID,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("设置主图失败: %v", err)
	}

	return nil
}

// getFileExtension 获取文件扩展名
func (im *ImageManager) getFileExtension(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i+1:]
		}
	}
	return ""
}

// 全局图片管理器实例
var ImageManagerInstance *ImageManager

// InitImageManager 初始化图片管理器
func InitImageManager() error {
	var err error
	ImageManagerInstance, err = NewImageManager()
	if err != nil {
		return fmt.Errorf("初始化图片管理器失败: %v", err)
	}
	return nil
}

// UploadBuildingFloorPlan 上传楼盘户型图
func (im *ImageManager) UploadBuildingFloorPlan(file *multipart.FileHeader, buildingID uint64, houseTypeID uint64, userID uint64) (*image.SysImage, error) {
	// 验证文件
	if err := im.qiniuService.ValidateFile(file); err != nil {
		return nil, err
	}

	// 生成存储Key，使用楼盘文件夹结构
	fileName := fmt.Sprintf("floor_plan_%d_%s", time.Now().UnixNano(), file.Filename)
	customKey := fmt.Sprintf("buildings/%d/floor-plans/%s", buildingID, fileName)

	// 上传到七牛云
	uploadResult, err := im.qiniuService.UploadFile(file, customKey)
	if err != nil {
		return nil, fmt.Errorf("上传到七牛云失败: %v", err)
	}

	// 保存到数据库
	img := &image.SysImage{
		Name:         file.Filename,
		Description:  fmt.Sprintf("楼盘%d的户型图", buildingID),
		FileName:     file.Filename,
		FileSize:     file.Size,
		MimeType:     file.Header.Get("Content-Type"),
		Extension:    im.getFileExtension(file.Filename),
		Key:          uploadResult.Key,
		URL:          uploadResult.OriginalURL,
		ThumbnailURL: uploadResult.ThumbnailURL,
		MediumURL:    uploadResult.MediumURL,
		LargeURL:     uploadResult.LargeURL,
		Category:     "floor_plan",
		Module:       "house",
		ModuleID:     houseTypeID, // 使用户型ID作为模块ID
		IsPublic:     true,
		IsMain:       false,
		SortOrder:    0,
		Status:       "active",
		CreatedBy:    userID,
		UpdatedBy:    userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := im.db.Create(img).Error; err != nil {
		// 如果数据库保存失败，删除已上传的文件
		im.qiniuService.DeleteFile(uploadResult.Key)
		return nil, fmt.Errorf("保存到数据库失败: %v", err)
	}

	return img, nil
}

// GetBuildingImages 获取楼盘的所有图片
func (im *ImageManager) GetBuildingImages(buildingID uint64, category string) ([]*image.SysImage, error) {
	var images []*image.SysImage

	// 构建查询条件：楼盘ID匹配 或 户型属于该楼盘
	query := im.db.Where("(module = 'building' AND module_id = ?) OR (module = 'house' AND module_id IN (SELECT id FROM sys_house_types WHERE building_id = ?))",
		buildingID, buildingID)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Where("status = 'active'").Order("is_main DESC, sort_order ASC, created_at DESC").Find(&images).Error; err != nil {
		return nil, fmt.Errorf("查询楼盘图片失败: %v", err)
	}

	return images, nil
}

// GetBuildingFloorPlans 获取楼盘的所有户型图
func (im *ImageManager) GetBuildingFloorPlans(buildingID uint64) ([]*image.SysImage, error) {
	return im.GetBuildingImages(buildingID, "floor_plan")
}

// CreateBuildingFolder 创建楼盘文件夹结构并在七牛云上创建相关目录
func (im *ImageManager) CreateBuildingFolder(buildingID uint64, buildingName string) error {
	// 验证楼盘是否存在
	var count int64
	if err := im.db.Table("sys_buildings").Where("id = ?", buildingID).Count(&count).Error; err != nil {
		return fmt.Errorf("验证楼盘存在性失败: %v", err)
	}

	if count == 0 {
		return fmt.Errorf("楼盘不存在: %d", buildingID)
	}

	// 定义楼盘文件夹结构
	folderStructure := map[string]string{
		"floor-plans":     "户型图",
		"site-plans":      "小区平面图",
		"environment":     "小区环境图",
		"building-images": "楼盘外观图",
		"interior":        "室内样板图",
		"facilities":      "配套设施图",
		"documents":       "相关文档",
	}

	// 在七牛云上创建文件夹标记文件
	if err := im.createQiniuFolderStructure(buildingID, buildingName, folderStructure); err != nil {
		fmt.Printf("⚠️  七牛云文件夹创建失败: %v\n", err)
		// 不阻止楼盘创建，只记录错误
	}

	// 在数据库中记录楼盘文件夹信息
	if err := im.recordBuildingFolderInfo(buildingID, buildingName, folderStructure); err != nil {
		fmt.Printf("⚠️  数据库文件夹信息记录失败: %v\n", err)
	}

	// 处理楼盘名称用于显示
	safeBuildingName := im.sanitizeFolderName(buildingName)
	buildingFolderName := fmt.Sprintf("%d-%s", buildingID, safeBuildingName)

	fmt.Printf("✅ 楼盘文件夹结构创建完成: 楼盘ID=%d, 名称=%s\n", buildingID, buildingName)
	fmt.Printf("📁 文件夹结构: buildings/%s/\n", buildingFolderName)
	for folder, desc := range folderStructure {
		fmt.Printf("   ├── %s/     (%s)\n", folder, desc)
	}

	return nil
}

// createQiniuFolderStructure 在七牛云上创建文件夹结构
func (im *ImageManager) createQiniuFolderStructure(buildingID uint64, buildingName string, folders map[string]string) error {
	if im.qiniuService == nil {
		return fmt.Errorf("七牛云服务未初始化")
	}

	// 处理楼盘名称，确保适合作为文件夹名称
	safeBuildingName := im.sanitizeFolderName(buildingName)
	buildingFolderName := fmt.Sprintf("%d-%s", buildingID, safeBuildingName)

	// 为每个文件夹创建一个标记文件（因为七牛云不支持空文件夹）
	for folder, desc := range folders {
		// 创建文件夹标记文件的key，使用新的命名格式
		folderKey := fmt.Sprintf("buildings/%s/%s/.folder", buildingFolderName, folder)

		// 创建标记文件内容
		content := fmt.Sprintf(`{
  "building_id": %d,
  "building_name": "%s",
  "building_folder_name": "%s",
  "folder_type": "%s",
  "description": "%s",
  "created_at": "%s",
  "purpose": "This file marks the existence of this folder structure"
}`, buildingID, buildingName, buildingFolderName, folder, desc, time.Now().Format("2006-01-02 15:04:05"))

		// 上传标记文件到七牛云
		if err := im.qiniuService.UploadText(folderKey, content); err != nil {
			fmt.Printf("⚠️  创建文件夹标记失败 %s: %v\n", folder, err)
			continue
		}

		fmt.Printf("📁 创建七牛云文件夹: buildings/%s/%s/\n", buildingFolderName, folder)
	}

	return nil
}

// sanitizeFolderName 清理楼盘名称，确保适合作为文件夹名称
func (im *ImageManager) sanitizeFolderName(name string) string {
	// 替换不适合文件夹名称的字符
	replacements := map[string]string{
		" ":  "-", // 空格替换为横线
		"/":  "-", // 斜杠替换为横线
		"\\": "-", // 反斜杠替换为横线
		":":  "-", // 冒号替换为横线
		"*":  "-", // 星号替换为横线
		"?":  "-", // 问号替换为横线
		"\"": "-", // 双引号替换为横线
		"<":  "-", // 小于号替换为横线
		">":  "-", // 大于号替换为横线
		"|":  "-", // 竖线替换为横线
		"（":  "(", // 中文括号替换为英文括号
		"）":  ")", // 中文括号替换为英文括号
		"【":  "[", // 中文方括号替换为英文方括号
		"】":  "]", // 中文方括号替换为英文方括号
	}

	result := name
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}

	// 限制长度，避免文件夹名称过长
	if len(result) > 50 {
		// 截取前50个字符，但确保不会截断中文字符
		runes := []rune(result)
		if len(runes) > 50 {
			result = string(runes[:50])
		}
	}

	return result
}

// recordBuildingFolderInfo 在数据库中记录楼盘文件夹信息
func (im *ImageManager) recordBuildingFolderInfo(buildingID uint64, buildingName string, folders map[string]string) error {
	// 可以在这里创建一个楼盘文件夹配置表来记录文件夹结构信息
	// 暂时只在日志中记录
	fmt.Printf("📝 记录楼盘文件夹信息: ID=%d, 名称=%s, 文件夹数量=%d\n", buildingID, buildingName, len(folders))
	return nil
}

// GetImageManager 获取图片管理器实例
func GetImageManager() *ImageManager {
	return ImageManagerInstance
}
