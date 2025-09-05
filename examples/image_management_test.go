package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"rentPro/rentpro-admin/common/config"
	"rentPro/rentpro-admin/common/initialize"
	"rentPro/rentpro-admin/common/utils"
)

// 图片管理API测试示例
func main() {
	fmt.Println("🖼️  七牛云图片管理API测试示例")
	fmt.Println("=" * 50)

	// 初始化配置
	fmt.Println("1. 初始化配置...")
	err := config.InitQiniuConfig("config/qiniu.yml", "development")
	if err != nil {
		fmt.Printf("❌ 配置初始化失败: %v\n", err)
		return
	}

	// 初始化七牛云服务
	fmt.Println("2. 初始化七牛云服务...")
	err = initialize.InitQiniu("development")
	if err != nil {
		fmt.Printf("❌ 七牛云服务初始化失败: %v\n", err)
		return
	}

	// 初始化图片管理器
	fmt.Println("3. 初始化图片管理器...")
	err = utils.InitImageManager()
	if err != nil {
		fmt.Printf("❌ 图片管理器初始化失败: %v\n", err)
		return
	}

	fmt.Println("✅ 初始化完成！")
	fmt.Println()

	// 演示API调用
	demonstrateAPIUsage()

	// 演示图片管理器功能
	demonstrateImageManager()
}

// 演示API调用示例
func demonstrateAPIUsage() {
	fmt.Println("📡 API调用示例:")
	fmt.Println("-" * 30)

	baseURL := "http://localhost:8002/api/v1"

	// 1. 获取图片统计
	fmt.Println("1. 获取图片统计信息:")
	fmt.Printf("   GET %s/images/stats\n", baseURL)

	// 2. 上传图片
	fmt.Println("2. 上传图片:")
	fmt.Printf("   POST %s/images/upload\n", baseURL)
	fmt.Printf("   Form Data:\n")
	fmt.Printf("     - file: <图片文件>\n")
	fmt.Printf("     - category: building\n")
	fmt.Printf("     - module: rental\n")
	fmt.Printf("     - moduleId: 123\n")
	fmt.Printf("     - isMain: false\n")
	fmt.Printf("     - isPublic: true\n")

	// 3. 获取图片列表
	fmt.Println("3. 获取图片列表:")
	fmt.Printf("   GET %s/images?page=1&pageSize=10&category=building\n", baseURL)

	// 4. 获取图片详情
	fmt.Println("4. 获取图片详情:")
	fmt.Printf("   GET %s/images/1\n", baseURL)

	// 5. 更新图片信息
	fmt.Println("5. 更新图片信息:")
	fmt.Printf("   PUT %s/images/1\n", baseURL)
	fmt.Printf(`   Body: {"name": "新图片名称", "description": "新描述"}\n`)

	// 6. 删除图片
	fmt.Println("6. 删除图片:")
	fmt.Printf("   DELETE %s/images/1\n", baseURL)

	// 7. 批量删除图片
	fmt.Println("7. 批量删除图片:")
	fmt.Printf("   DELETE %s/images/batch\n", baseURL)
	fmt.Printf(`   Body: {"ids": [1, 2, 3]}\n`)

	// 8. 获取模块图片
	fmt.Println("8. 获取模块图片:")
	fmt.Printf("   GET %s/images/module/rental/123?category=building\n", baseURL)

	// 9. 设置主图
	fmt.Println("9. 设置主图:")
	fmt.Printf("   PUT %s/images/1/set-main\n", baseURL)
	fmt.Printf(`   Body: {"module": "rental", "moduleId": 123}\n`)

	fmt.Println()
}

// 演示图片管理器功能
func demonstrateImageManager() {
	fmt.Println("🔧 图片管理器功能演示:")
	fmt.Println("-" * 30)

	imageManager := utils.GetImageManager()
	if imageManager == nil {
		fmt.Println("❌ 图片管理器未初始化")
		return
	}

	// 1. 获取统计信息
	fmt.Println("1. 获取图片统计信息:")
	stats, err := imageManager.GetImageStats()
	if err != nil {
		fmt.Printf("❌ 获取统计信息失败: %v\n", err)
	} else {
		fmt.Printf("   总图片数: %d\n", stats.TotalImages)
		fmt.Printf("   总存储大小: %d bytes\n", stats.TotalSize)
		fmt.Printf("   今日上传: %d\n", stats.TodayUploads)
		fmt.Printf("   分类统计: %+v\n", stats.CategoryStats)
		fmt.Printf("   模块统计: %+v\n", stats.ModuleStats)
	}

	fmt.Println()

	// 2. 模拟上传文件
	fmt.Println("2. 模拟文件上传:")
	// 这里可以添加实际的文件上传测试
	fmt.Println("   💡 提示: 运行实际服务器后，可以使用curl或Postman测试文件上传")

	fmt.Println()

	// 3. 显示支持的功能
	fmt.Println("3. 支持的功能:")
	fmt.Println("   ✅ 文件上传到七牛云")
	fmt.Println("   ✅ 自动生成多种尺寸图片")
	fmt.Println("   ✅ 图片分类管理")
	fmt.Println("   ✅ 模块关联管理")
	fmt.Println("   ✅ 主图设置")
	fmt.Println("   ✅ 批量操作")
	fmt.Println("   ✅ 统计信息")
	fmt.Println("   ✅ 权限控制")

	fmt.Println()

	// 4. 使用建议
	fmt.Println("4. 使用建议:")
	fmt.Println("   📁 分类管理: 为不同业务场景创建分类")
	fmt.Println("   🏷️  标签系统: 使用模块+模块ID关联业务数据")
	fmt.Println("   🖼️  多尺寸: 利用七牛云的图片处理功能")
	fmt.Println("   🔒 权限控制: 根据用户权限控制图片访问")
	fmt.Println("   📊 监控统计: 定期查看存储使用情况")

	fmt.Println()
	fmt.Println("🎉 图片管理系统已准备就绪！")
}

// 实际的文件上传测试示例
func uploadTestFile() {
	fmt.Println("🧪 文件上传测试:")

	// 创建测试文件
	testFilePath := "test_image.jpg"
	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		fmt.Printf("⚠️  测试文件不存在: %s\n", testFilePath)
		fmt.Println("💡 请准备一个测试图片文件")
		return
	}

	// 打开测试文件
	file, err := os.Open(testFilePath)
	if err != nil {
		fmt.Printf("❌ 打开测试文件失败: %v\n", err)
		return
	}
	defer file.Close()

	// 创建multipart表单
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// 添加文件
	fw, err := w.CreateFormFile("file", filepath.Base(testFilePath))
	if err != nil {
		fmt.Printf("❌ 创建表单文件失败: %v\n", err)
		return
	}
	if _, err = io.Copy(fw, file); err != nil {
		fmt.Printf("❌ 复制文件内容失败: %v\n", err)
		return
	}

	// 添加表单字段
	w.WriteField("category", "building")
	w.WriteField("module", "test")
	w.WriteField("moduleId", "1")
	w.WriteField("isMain", "true")
	w.WriteField("isPublic", "true")
	w.Close()

	// 发送HTTP请求
	req, err := http.NewRequest("POST", "http://localhost:8002/api/v1/images/upload", &b)
	if err != nil {
		fmt.Printf("❌ 创建请求失败: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer YOUR_TOKEN_HERE") // 需要实际的token

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ 发送请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("📡 响应状态: %s\n", resp.Status)
	if resp.StatusCode == 200 {
		fmt.Println("✅ 文件上传成功！")
	} else {
		fmt.Println("❌ 文件上传失败")
	}
}

func init() {
	// 设置Go模块路径
	os.Setenv("GO111MODULE", "on")
}
