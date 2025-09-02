package base

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// SQLFileLoader SQL文件加载器
type SQLFileLoader struct {
	ConfigPath string // 配置文件根路径，默认为 "config/sql/data"
}

// NewSQLFileLoader 创建SQL文件加载器
func NewSQLFileLoader(configPath string) *SQLFileLoader {
	if configPath == "" {
		configPath = "config/sql/data"
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

		// 检查是否是INSERT语句
		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(sqlStmt)), "INSERT") {
			// 对于INSERT语句，先检查是否可以跳过（避免重复插入）
			if shouldSkipInsert(db, sqlStmt) {
				fmt.Printf("信息: 跳过已存在的数据插入: %s\n", sqlStmt)
				continue
			}
		}

		if err := db.Exec(sqlStmt).Error; err != nil {
			// 如果是重复键错误，跳过（防止重复插入）
			if strings.Contains(err.Error(), "Duplicate entry") ||
				strings.Contains(err.Error(), "duplicate key") {
				fmt.Printf("警告: 跳过重复数据插入: %s\n", sqlStmt)
				continue
			}
			return fmt.Errorf("执行SQL语句失败: %v, SQL: %s", err, sqlStmt)
		}
	}

	return nil
}

// shouldSkipInsert 检查是否应该跳过插入语句（避免重复插入数据）
func shouldSkipInsert(db *gorm.DB, sqlStmt string) bool {
	// 解析INSERT语句，提取表名和主键值
	// 检查是否是sys_dept表的插入
	if strings.Contains(sqlStmt, "sys_dept") && strings.Contains(sqlStmt, "VALUES") {
		return shouldSkipDeptInsert(db, sqlStmt)
	}

	// 检查是否是sys_post表的插入
	if strings.Contains(sqlStmt, "sys_post") && strings.Contains(sqlStmt, "VALUES") {
		return shouldSkipPostInsert(db, sqlStmt)
	}

	// 检查是否是sys_role表的插入
	if strings.Contains(sqlStmt, "sys_role") && strings.Contains(sqlStmt, "VALUES") {
		return shouldSkipRoleInsert(db, sqlStmt)
	}

	// 检查是否是sys_menu表的插入
	if strings.Contains(sqlStmt, "sys_menu") && strings.Contains(sqlStmt, "VALUES") {
		return shouldSkipMenuInsert(db, sqlStmt)
	}

	// 检查是否是sys_user表的插入
	if strings.Contains(sqlStmt, "sys_user") && strings.Contains(sqlStmt, "VALUES") {
		return shouldSkipUserInsert(db, sqlStmt)
	}

	// 对于其他表或无法解析的情况，不跳过
	return false
}

// shouldSkipDeptInsert 检查是否应该跳过部门数据插入
func shouldSkipDeptInsert(db *gorm.DB, sqlStmt string) bool {
	// 提取VALUES部分
	re := regexp.MustCompile(`VALUES\s*$$[^)]*$$`)
	matches := re.FindStringSubmatch(sqlStmt)
	if len(matches) > 0 {
		valuesPart := matches[0]
		// 提取每个值组
		re = regexp.MustCompile(`$$(.*?)$$`)
		valueGroups := re.FindAllString(valuesPart, -1)

		for _, group := range valueGroups {
			// 提取第一个值作为ID
			group = strings.Trim(group, "()")
			parts := strings.Split(group, ",")
			if len(parts) > 0 {
				// 清理ID值（去除空格和引号）
				idStr := strings.TrimSpace(parts[0])
				idStr = strings.Trim(idStr, "'\"")
				id, err := strconv.Atoi(idStr)
				if err == nil {
					// 检查该ID是否已存在
					var count int64
					db.Table("sys_dept").Where("id = ?", id).Count(&count)
					if count > 0 {
						return true // 如果任何一个ID已存在，则跳过整个插入
					}
				}
			}
		}
	}
	return false
}

// shouldSkipPostInsert 检查是否应该跳过岗位数据插入
func shouldSkipPostInsert(db *gorm.DB, sqlStmt string) bool {
	// 提取VALUES部分
	re := regexp.MustCompile(`VALUES\s*$$[^)]*$$`)
	matches := re.FindStringSubmatch(sqlStmt)
	if len(matches) > 0 {
		valuesPart := matches[0]
		// 提取每个值组
		re = regexp.MustCompile(`$$(.*?)$$`)
		valueGroups := re.FindAllString(valuesPart, -1)

		for _, group := range valueGroups {
			// 提取第一个值作为ID
			group = strings.Trim(group, "()")
			parts := strings.Split(group, ",")
			if len(parts) > 0 {
				// 清理ID值（去除空格和引号）
				idStr := strings.TrimSpace(parts[0])
				idStr = strings.Trim(idStr, "'\"")
				id, err := strconv.Atoi(idStr)
				if err == nil {
					// 检查该ID是否已存在
					var count int64
					db.Table("sys_post").Where("id = ?", id).Count(&count)
					if count > 0 {
						return true // 如果任何一个ID已存在，则跳过整个插入
					}
				}
			}
		}
	}
	return false
}

// shouldSkipRoleInsert 检查是否应该跳过角色数据插入
func shouldSkipRoleInsert(db *gorm.DB, sqlStmt string) bool {
	// 提取VALUES部分
	re := regexp.MustCompile(`VALUES\s*$$[^)]*$$`)
	matches := re.FindStringSubmatch(sqlStmt)
	if len(matches) > 0 {
		valuesPart := matches[0]
		// 提取每个值组
		re = regexp.MustCompile(`$$(.*?)$$`)
		valueGroups := re.FindAllString(valuesPart, -1)

		for _, group := range valueGroups {
			// 提取第一个值作为ID
			group = strings.Trim(group, "()")
			parts := strings.Split(group, ",")
			if len(parts) > 0 {
				// 清理ID值（去除空格和引号）
				idStr := strings.TrimSpace(parts[0])
				idStr = strings.Trim(idStr, "'\"")
				id, err := strconv.Atoi(idStr)
				if err == nil {
					// 检查该ID是否已存在
					var count int64
					db.Table("sys_role").Where("id = ?", id).Count(&count)
					if count > 0 {
						return true // 如果任何一个ID已存在，则跳过整个插入
					}
				}
			}
		}
	}
	return false
}

// shouldSkipMenuInsert 检查是否应该跳过菜单数据插入
func shouldSkipMenuInsert(db *gorm.DB, sqlStmt string) bool {
	// 提取VALUES部分
	re := regexp.MustCompile(`VALUES\s*$$[^)]*$$`)
	matches := re.FindStringSubmatch(sqlStmt)
	if len(matches) > 0 {
		valuesPart := matches[0]
		// 提取每个值组
		re = regexp.MustCompile(`$$(.*?)$$`)
		valueGroups := re.FindAllString(valuesPart, -1)

		for _, group := range valueGroups {
			// 提取第一个值作为ID
			group = strings.Trim(group, "()")
			parts := strings.Split(group, ",")
			if len(parts) > 0 {
				// 清理ID值（去除空格和引号）
				idStr := strings.TrimSpace(parts[0])
				idStr = strings.Trim(idStr, "'\"")
				id, err := strconv.Atoi(idStr)
				if err == nil {
					// 检查该ID是否已存在
					var count int64
					db.Table("sys_menu").Where("id = ?", id).Count(&count)
					if count > 0 {
						return true // 如果任何一个ID已存在，则跳过整个插入
					}
				}
			}
		}
	}
	return false
}

// shouldSkipUserInsert 检查是否应该跳过用户数据插入
func shouldSkipUserInsert(db *gorm.DB, sqlStmt string) bool {
	// 提取VALUES部分
	re := regexp.MustCompile(`VALUES\s*$$[^)]*$$`)
	matches := re.FindStringSubmatch(sqlStmt)
	if len(matches) > 0 {
		valuesPart := matches[0]
		// 提取每个值组
		re = regexp.MustCompile(`$$(.*?)$$`)
		valueGroups := re.FindAllString(valuesPart, -1)

		for _, group := range valueGroups {
			// 提取第一个值作为ID
			group = strings.Trim(group, "()")
			parts := strings.Split(group, ",")
			if len(parts) > 0 {
				// 清理ID值（去除空格和引号）
				idStr := strings.TrimSpace(parts[0])
				idStr = strings.Trim(idStr, "'\"")
				id, err := strconv.Atoi(idStr)
				if err == nil {
					// 检查该ID是否已存在
					var count int64
					db.Table("sys_user").Where("id = ?", id).Count(&count)
					if count > 0 {
						return true // 如果任何一个ID已存在，则跳过整个插入
					}
				}
			}
		}
	}
	return false
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

	// 如果文件不存在，尝试使用ConfigPath作为绝对路径
	if !loader.fileExists(fullPath) {
		fullPath = filepath.Join(loader.ConfigPath, filename)
	}

	// 检查文件是否存在
	if !loader.fileExists(fullPath) {
		return nil, fmt.Errorf("SQL文件不存在: %s (尝试路径: %s)", filename, fullPath)
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
