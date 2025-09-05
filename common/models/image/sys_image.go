package image

import (
	"time"

	"gorm.io/gorm"
)

// SysImage 图片管理模型
type SysImage struct {
	ID          uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"size:200;not null;comment:图片名称"`
	Description string `json:"description" gorm:"size:500;comment:图片描述"`
	FileName    string `json:"fileName" gorm:"size:255;not null;comment:原始文件名"`
	FileSize    int64  `json:"fileSize" gorm:"not null;comment:文件大小(字节)"`
	MimeType    string `json:"mimeType" gorm:"size:100;not null;comment:MIME类型"`
	Extension   string `json:"extension" gorm:"size:10;not null;comment:文件扩展名"`

	// 存储路径信息
	Key          string `json:"key" gorm:"size:500;not null;comment:七牛云存储Key"`
	URL          string `json:"url" gorm:"size:1000;not null;comment:原始图片URL"`
	ThumbnailURL string `json:"thumbnailUrl" gorm:"size:1000;comment:缩略图URL"`
	MediumURL    string `json:"mediumUrl" gorm:"size:1000;comment:中等尺寸URL"`
	LargeURL     string `json:"largeUrl" gorm:"size:1000;comment:大图URL"`

	// 分类信息
	Category string `json:"category" gorm:"size:50;not null;default:'default';comment:图片分类(building/house/avatar/banner等)"`
	Module   string `json:"module" gorm:"size:50;not null;default:'common';comment:所属模块"`
	ModuleID uint64 `json:"moduleId" gorm:"default:0;comment:模块关联ID"`

	// 图片属性
	Width  int    `json:"width" gorm:"comment:图片宽度"`
	Height int    `json:"height" gorm:"comment:图片高度"`
	Hash   string `json:"hash" gorm:"size:100;comment:文件Hash"`

	// 状态控制
	IsPublic  bool   `json:"isPublic" gorm:"default:true;comment:是否公开访问"`
	IsMain    bool   `json:"isMain" gorm:"default:false;comment:是否为主图"`
	SortOrder int    `json:"sortOrder" gorm:"default:0;comment:排序序号"`
	Status    string `json:"status" gorm:"size:20;default:'active';comment:状态(active/inactive/deleted)"`

	// 审计字段
	CreatedBy uint64         `json:"createdBy" gorm:"comment:创建者ID"`
	UpdatedBy uint64         `json:"updatedBy" gorm:"comment:更新者ID"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

// TableName 指定表名
func (SysImage) TableName() string {
	return "sys_images"
}

// SysImageCategory 图片分类配置
type SysImageCategory struct {
	ID           uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Code         string    `json:"code" gorm:"size:50;unique;not null;comment:分类编码"`
	Name         string    `json:"name" gorm:"size:100;not null;comment:分类名称"`
	Description  string    `json:"description" gorm:"size:200;comment:分类描述"`
	MaxSize      int64     `json:"maxSize" gorm:"comment:最大文件大小"`
	AllowedTypes []string  `json:"allowedTypes" gorm:"type:json;comment:允许的文件类型"`
	MaxCount     int       `json:"maxCount" gorm:"default:10;comment:最大上传数量"`
	IsRequired   bool      `json:"isRequired" gorm:"default:false;comment:是否必填"`
	Status       string    `json:"status" gorm:"size:20;default:'active';comment:状态"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// TableName 指定表名
func (SysImageCategory) TableName() string {
	return "sys_image_categories"
}

// ImageUploadRequest 图片上传请求
type ImageUploadRequest struct {
	Category string `json:"category" binding:"required"` // 图片分类
	Module   string `json:"module" binding:"required"`   // 所属模块
	ModuleID uint64 `json:"moduleId"`                    // 模块关联ID
	IsMain   bool   `json:"isMain"`                      // 是否为主图
	IsPublic bool   `json:"isPublic"`                    // 是否公开访问
}

// ImageUpdateRequest 图片更新请求
type ImageUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	IsMain      bool   `json:"isMain"`
	IsPublic    bool   `json:"isPublic"`
	SortOrder   int    `json:"sortOrder"`
	Status      string `json:"status"`
}

// ImageBatchDeleteRequest 批量删除请求
type ImageBatchDeleteRequest struct {
	IDs []uint64 `json:"ids" binding:"required"`
}

// ImageListRequest 图片列表请求
type ImageListRequest struct {
	Category string `form:"category"`
	Module   string `form:"module"`
	ModuleID uint64 `form:"moduleId"`
	Status   string `form:"status"`
	IsMain   *bool  `form:"isMain"`
	IsPublic *bool  `form:"isPublic"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=10"`
	OrderBy  string `form:"orderBy,default=created_at"`
	OrderDir string `form:"orderDir,default=desc"`
}

// ImageListResponse 图片列表响应
type ImageListResponse struct {
	Total int64       `json:"total"`
	List  []*SysImage `json:"list"`
}

// ImageStats 图片统计信息
type ImageStats struct {
	TotalImages   int64            `json:"totalImages"`
	TotalSize     int64            `json:"totalSize"`
	CategoryStats map[string]int64 `json:"categoryStats"`
	ModuleStats   map[string]int64 `json:"moduleStats"`
	TodayUploads  int64            `json:"todayUploads"`
	StorageUsed   int64            `json:"storageUsed"`
}
