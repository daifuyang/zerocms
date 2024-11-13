package utils

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// 去除 DSN 中的数据库名称
func RemoveDatabaseName(dsn string) string {
	// 正则表达式匹配 DSN 中的 /database 部分
	parts := strings.SplitN(dsn, "/", 2)
	dsnWithoutDB := parts[0] + "/"
	return dsnWithoutDB
}

// 解析 DSN 提取数据库名称

func ParseDatabaseName(dsn string) (string, error) {
	re := regexp.MustCompile(`/([a-zA-Z0-9_]+)\\?`)
	match := re.FindStringSubmatch(dsn)
	if len(match) < 2 {
		return "", fmt.Errorf("database name not found")
	}
	return match[1], nil
}

// 创建数据库

func CreateDatabase(db *sql.DB, dbName string) error {
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)
	_, err := db.Exec(query)
	if err != nil {
		fmt.Printf("Error creating database: %v\n", err)
	}
	return err
}

// 读取并执行 SQL 文件

func ExecuteSQLFile(db *sql.DB, filePath string) error {
	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %v", err)
	}
	// 按分号分割为多个 SQL 语句
	queries := strings.Split(string(content), ";")

	for _, query := range queries {
		// 去除每条语句的前后空格
		query = strings.TrimSpace(query)
		if query == "" {
			continue // 跳过空语句
		}

		// 执行每条语句
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to execute query: %v", err)
		}
	}

	return nil
}

// 执行文件夹下所有的 SQL 文件

func ExecuteSQLFilesInDirectory(db *sql.DB, dirPath string) error {
	// 遍历文件夹下的所有 .sql 文件
	return filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 只处理 .sql 文件
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".sql") {
			if err := ExecuteSQLFile(db, path); err != nil {
				return fmt.Errorf("failed to execute SQL file %s: %v", path, err)
			}
		}
		return nil
	})
}
