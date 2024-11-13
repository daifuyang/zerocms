package utils

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

// 生成一个随机盐值

func GenerateSalt(size int) (string, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// 生成hash密码

func GenerateHash(content string, salt string) (string, error) {
	// bcrypt.DefaultCost 默认使用 10 次迭代
	hash, err := bcrypt.GenerateFromPassword([]byte(content+salt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// 密码比对

func ComparePassword(hashedPassword, plainPassword string) bool {
	// 使用 bcrypt.CompareHashAndPassword 比较哈希值和明文密码
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
