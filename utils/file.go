package utils

import (
	"fmt"
	"io"
	"os"
)

// 检查文件是否存在

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// 如果文件不存在，则创建它

func CreateFile(filename string) error {
	if FileExists(filename) {
		return nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()
	return nil
}

func ReadFile(filePath string) ([]byte, error) {
	// 打开 JSON 文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 读取文件内容
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}
