package base

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

// SQLFileLoader SQL文件加载器
type SQLFileLoader struct {
	ConfigPath string // 配置文件根路径，默认为 "config/data"
}

// NewSQLFileLoader 创建SQL文件加载器
func NewSQLFileLoader(configPath string) *SQLFileLoader {
	if configPath == "" {
		configPath = "config/data"
	}
	return &SQLFileLoader{
		ConfigPath: configPath,
	}
}

// LoadAndExecuteSQL 加载并执行SQL文件
func (loader *SQLFileLoader) LoadAndExecuteSQL(db *gorm.DB, filename string) error {
	sqlStatements, err := loader.ReadSQLFromFile(filename)
	if err != nil {
		return fmt.Errorf("读取SQL文件失败: %v", err)
	}

	// 执行SQL语句
	for _, sqlStmt := range sqlStatements {
		if strings.TrimSpace(sqlStmt) == "" {
			continue
		}

		if err := db.Exec(sqlStmt).Error; err != nil {
			// 如果是重复键错误，跳过（防止重复插入）
			if strings.Contains(err.Error(), "Duplicate entry") ||
				strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			return fmt.Errorf("执行SQL语句失败: %v, SQL: %s", err, sqlStmt)
		}
	}

	return nil
}

// ReadSQLFromFile 从SQL文件中读取INSERT语句
func (loader *SQLFileLoader) ReadSQLFromFile(filename string) ([]string, error) {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// 构造完整的文件路径
	fullPath := filepath.Join(wd, loader.ConfigPath, filename)

	// 检查文件是否存在
	if !loader.fileExists(fullPath) {
		return nil, fmt.Errorf("SQL文件不存在: %s", fullPath)
	}

	// 读取文件内容
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法读取文件 %s: %v", fullPath, err)
	}

	// 解析SQL语句
	sqlStatements := loader.ParseSQLStatements(string(content))
	return sqlStatements, nil
}

// ParseSQLStatements 解析SQL文件内容，提取INSERT语句
func (loader *SQLFileLoader) ParseSQLStatements(content string) []string {
	var statements []string

	// 按行分割内容
	scanner := bufio.NewScanner(strings.NewReader(content))
	var currentStatement strings.Builder

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}

		// 检查是否是INSERT语句
		insertRegex := regexp.MustCompile(`^INSERT\s+INTO`)
		if insertRegex.MatchString(strings.ToUpper(line)) {
			// 开始新的INSERT语句
			currentStatement.Reset()
			currentStatement.WriteString(line)

			// 检查是否是完整语句（以分号结尾）
			if strings.HasSuffix(line, ";") {
				// 移除分号并添加到语句列表
				stmt := strings.TrimSuffix(currentStatement.String(), ";")
				statements = append(statements, stmt)
				currentStatement.Reset()
			}
		} else if currentStatement.Len() > 0 {
			// 继续当前语句
			currentStatement.WriteString(" ")
			currentStatement.WriteString(line)

			// 检查是否语句结束
			if strings.HasSuffix(line, ";") {
				// 移除分号并添加到语句列表
				stmt := strings.TrimSuffix(currentStatement.String(), ";")
				statements = append(statements, stmt)
				currentStatement.Reset()
			}
		}
	}

	return statements
}

// fileExists 检查文件是否存在
func (loader *SQLFileLoader) fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetAvailableFiles 获取可用的SQL文件列表
func (loader *SQLFileLoader) GetAvailableFiles() ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	dirPath := filepath.Join(wd, loader.ConfigPath)
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}

	return sqlFiles, nil
}
