package utils

import (
	"fmt"
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
