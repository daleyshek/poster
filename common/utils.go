package common

import (
	"math/rand"
	"time"
)

// GetRandomName 获取随机的文件名
func GetRandomName(length int) (name string) {
	dic := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	maxL := len([]byte(dic))
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		name = name + string(dic[rand.Intn(maxL)])
	}
	return name
}
