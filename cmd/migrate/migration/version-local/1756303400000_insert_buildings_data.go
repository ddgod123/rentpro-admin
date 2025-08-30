package version_local

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"rentPro/rentpro-admin/cmd/migrate/migration"
	commonModels "rentPro/rentpro-admin/common/models/base"
	"strings"

	"gorm.io/gorm"
)

func init() {
	migration.Migrate.SetVersion("1756303400000", insertBuildingsData)
}

// insertBuildingsData 从db.sql文件读取并插入楼盘初始化数据
func insertBuildingsData(db *gorm.DB, version string) error {
	// 读取SQL文件
	sqlStatements, err := readSQLFromFile("config/db.sql")
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
			if strings.Contains(err.Error(), "Duplicate entry") {
				continue
			}
			return fmt.Errorf("执行SQL语句失败: %v, SQL: %s", err, sqlStmt)
		}
	}

	// 记录迁移完成
	return db.Create(&commonModels.Migration{
		Version: version,
		Name:    "从db.sql插入楼盘初始化数据",
		Status:  "completed",
	}).Error
}

// readSQLFromFile 从SQL文件中读取INSERT语句
func readSQLFromFile(filePath string) ([]string, error) {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// 构造完整的文件路径
	fullPath := filepath.Join(wd, filePath)

	// 读取文件内容
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法读取文件 %s: %v", fullPath, err)
	}

	// 解析SQL语句
	sqlStatements := parseSQLStatements(string(content))
	return sqlStatements, nil
}

// parseSQLStatements 解析SQL文件内容，提取INSERT语句
func parseSQLStatements(content string) []string {
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
		insertRegex := regexp.MustCompile(`^INSERT\s+INTO\s+sys_buildings`)
		if insertRegex.MatchString(line) {
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
