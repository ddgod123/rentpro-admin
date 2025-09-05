package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"

	"rentPro/rentpro-admin/common/config"
)

// QiniuService 七牛云服务
type QiniuService struct {
	mac           *qbox.Mac
	config        storage.Config
	bucket        string
	domain        string
	qiniuConfig   *config.QiniuConfig
	configManager *config.QiniuConfigManager
}

// UploadResult 上传结果
type UploadResult struct {
	Key          string            `json:"key"`           // 存储key
	Hash         string            `json:"hash"`          // 文件hash
	Size         int64             `json:"size"`          // 文件大小
	OriginalURL  string            `json:"original_url"`  // 原图URL
	ThumbnailURL string            `json:"thumbnail_url"` // 缩略图URL
	MediumURL    string            `json:"medium_url"`    // 中等尺寸URL
	LargeURL     string            `json:"large_url"`     // 大图URL
	Styles       map[string]string `json:"styles"`        // 所有样式URL
}

// NewQiniuService 创建七牛云服务实例
func NewQiniuService() (*QiniuService, error) {
	qiniuConfig := config.GetQiniuConfig()
	if qiniuConfig == nil {
		return nil, fmt.Errorf("七牛云配置未初始化")
	}

	// 创建认证对象
	mac := qbox.NewMac(qiniuConfig.AccessKey, qiniuConfig.SecretKey)

	// 配置存储区域
	var zone *storage.Region
	switch qiniuConfig.Zone {
	case "huadong":
		zone = &storage.ZoneHuadong
	case "huabei":
		zone = &storage.ZoneHuabei
	case "huanan":
		zone = &storage.ZoneHuanan
	case "beimei":
		zone = &storage.ZoneBeimei
	default:
		zone = &storage.ZoneHuadong // 默认华东
	}

	// 存储配置
	cfg := storage.Config{
		Zone:          zone,
		UseHTTPS:      qiniuConfig.UseHTTPS,
		UseCdnDomains: qiniuConfig.UseCdnDomains,
	}

	return &QiniuService{
		mac:           mac,
		config:        cfg,
		bucket:        qiniuConfig.Bucket,
		domain:        qiniuConfig.Domain,
		qiniuConfig:   qiniuConfig,
		configManager: config.QiniuConfigInstance,
	}, nil
}

// UploadFile 上传文件
func (q *QiniuService) UploadFile(file *multipart.FileHeader, customKey string) (*UploadResult, error) {
	// 生成上传key
	var key string
	if customKey != "" {
		key = q.configManager.GetUploadKey(customKey)
	} else {
		// 自动生成key
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("floor_plan_%d%s", time.Now().UnixNano(), ext)
		key = q.configManager.GetUploadKey(filename)
	}

	// 创建临时文件
	tempFile, err := q.saveToTempFile(file)
	if err != nil {
		return nil, fmt.Errorf("保存临时文件失败: %v", err)
	}
	defer os.Remove(tempFile) // 清理临时文件

	// 执行上传
	uploader := storage.NewFormUploader(&q.config)
	ret := storage.PutRet{}

	// 生成上传Token
	putPolicy := storage.PutPolicy{
		Scope: q.bucket,
	}
	upToken := putPolicy.UploadToken(q.mac)

	// 上传文件
	err = uploader.PutFile(context.Background(), &ret, upToken, key, tempFile, nil)
	if err != nil {
		return nil, fmt.Errorf("上传到七牛云失败: %v", err)
	}

	// 生成访问URL
	originalURL := q.configManager.GetPublicURL(key)

	// 生成各种样式的URL
	styles := make(map[string]string)
	for styleName := range q.qiniuConfig.ImageStyles {
		styles[styleName] = q.configManager.GetImageStyleURL(originalURL, styleName)
	}

	return &UploadResult{
		Key:          key,
		Hash:         ret.Hash,
		Size:         file.Size,
		OriginalURL:  originalURL,
		ThumbnailURL: styles["thumbnail"],
		MediumURL:    styles["medium"],
		LargeURL:     styles["large"],
		Styles:       styles,
	}, nil
}

// UploadFromPath 从本地路径上传文件
func (q *QiniuService) UploadFromPath(localPath string, customKey string) (*UploadResult, error) {
	// 检查文件是否存在
	fileInfo, err := os.Stat(localPath)
	if err != nil {
		return nil, fmt.Errorf("文件不存在: %v", err)
	}

	// 生成上传key
	var key string
	if customKey != "" {
		key = q.configManager.GetUploadKey(customKey)
	} else {
		// 使用文件名生成key
		filename := filepath.Base(localPath)
		key = q.configManager.GetUploadKey(filename)
	}

	// 执行上传
	uploader := storage.NewFormUploader(&q.config)
	ret := storage.PutRet{}

	// 生成上传Token
	putPolicy := storage.PutPolicy{
		Scope: q.bucket,
	}
	upToken := putPolicy.UploadToken(q.mac)

	// 上传文件
	err = uploader.PutFile(context.Background(), &ret, upToken, key, localPath, nil)
	if err != nil {
		return nil, fmt.Errorf("上传到七牛云失败: %v", err)
	}

	// 生成访问URL
	originalURL := q.configManager.GetPublicURL(key)

	// 生成各种样式的URL
	styles := make(map[string]string)
	for styleName := range q.qiniuConfig.ImageStyles {
		styles[styleName] = q.configManager.GetImageStyleURL(originalURL, styleName)
	}

	return &UploadResult{
		Key:          key,
		Hash:         ret.Hash,
		Size:         fileInfo.Size(),
		OriginalURL:  originalURL,
		ThumbnailURL: styles["thumbnail"],
		MediumURL:    styles["medium"],
		LargeURL:     styles["large"],
		Styles:       styles,
	}, nil
}

// DeleteFile 删除文件
func (q *QiniuService) DeleteFile(key string) error {
	bucketManager := storage.NewBucketManager(q.mac, &q.config)
	err := bucketManager.Delete(q.bucket, key)
	if err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}
	return nil
}

// GetFileInfo 获取文件信息
func (q *QiniuService) GetFileInfo(key string) (*storage.FileInfo, error) {
	bucketManager := storage.NewBucketManager(q.mac, &q.config)
	fileInfo, err := bucketManager.Stat(q.bucket, key)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}
	return &fileInfo, nil
}

// UploadText 上传文本内容到七牛云
func (q *QiniuService) UploadText(key string, content string) error {
	upToken := q.GenerateUploadToken(key, 3600) // 1小时有效期

	formUploader := storage.NewFormUploader(&q.config)
	ret := storage.PutRet{}

	// 使用字符串内容直接上传
	putExtra := storage.PutExtra{}

	err := formUploader.Put(context.Background(), &ret, upToken, key, strings.NewReader(content), int64(len(content)), &putExtra)
	if err != nil {
		return fmt.Errorf("上传文本内容失败: %v", err)
	}

	return nil
}

// GenerateUploadToken 生成上传Token（供前端直传使用）
func (q *QiniuService) GenerateUploadToken(key string, expires int64) string {
	putPolicy := storage.PutPolicy{
		Scope:   q.bucket + ":" + key,
		Expires: uint64(expires),
	}
	return putPolicy.UploadToken(q.mac)
}

// GenerateDownloadURL 生成私有空间下载URL
func (q *QiniuService) GenerateDownloadURL(key string, expires int64) string {
	deadline := time.Now().Add(time.Duration(expires) * time.Second).Unix()
	privateAccessURL := storage.MakePrivateURL(q.mac, q.domain, key, deadline)
	return privateAccessURL
}

// ValidateFile 验证文件
func (q *QiniuService) ValidateFile(file *multipart.FileHeader) error {
	// 检查文件大小
	if file.Size > q.configManager.GetMaxFileSize() {
		return fmt.Errorf("文件大小超过限制: %d bytes", q.configManager.GetMaxFileSize())
	}

	// 检查文件类型
	contentType := file.Header.Get("Content-Type")
	if !q.configManager.IsAllowedFileType(contentType) {
		return fmt.Errorf("不支持的文件类型: %s", contentType)
	}

	return nil
}

// GetStyleURL 获取指定样式的URL
func (q *QiniuService) GetStyleURL(originalURL string, styleName string) string {
	return q.configManager.GetImageStyleURL(originalURL, styleName)
}

// ListFiles 列出文件
func (q *QiniuService) ListFiles(prefix string, limit int) ([]storage.ListItem, error) {
	bucketManager := storage.NewBucketManager(q.mac, &q.config)

	entries, _, _, hasNext, err := bucketManager.ListFiles(q.bucket, prefix, "", "", limit)
	if err != nil {
		return nil, fmt.Errorf("列出文件失败: %v", err)
	}

	_ = hasNext // 可以用于分页
	return entries, nil
}

// saveToTempFile 保存到临时文件
func (q *QiniuService) saveToTempFile(file *multipart.FileHeader) (string, error) {
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "qiniu_upload_*"+filepath.Ext(file.Filename))
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	// 复制文件内容
	_, err = tempFile.ReadFrom(src)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", err
	}

	return tempFile.Name(), nil
}

// ExtractKeyFromURL 从URL中提取key
func (q *QiniuService) ExtractKeyFromURL(url string) string {
	// 移除域名和协议部分
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, q.domain+"/")

	// 移除样式后缀
	for _, style := range q.qiniuConfig.ImageStyles {
		styleSuffix := "-" + style.Name
		if strings.HasSuffix(url, styleSuffix) {
			url = strings.TrimSuffix(url, styleSuffix)
			break
		}
	}

	return url
}

// 全局服务实例
var QiniuServiceInstance *QiniuService

// InitQiniuService 初始化七牛云服务
func InitQiniuService() error {
	var err error
	QiniuServiceInstance, err = NewQiniuService()
	if err != nil {
		return fmt.Errorf("初始化七牛云服务失败: %v", err)
	}
	return nil
}

// GetQiniuService 获取七牛云服务实例
func GetQiniuService() *QiniuService {
	return QiniuServiceInstance
}
