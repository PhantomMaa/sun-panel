package util

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"sort"
	"time"
)

const (
	TimeFormatMode1 = "2006-01-02 15:04:05" // 标准格式
)

func GetTime() string {
	return time.Unix(time.Now().Unix(), 0).Format(TimeFormatMode1)
}

func Md5(str string) string {
	md5Byte := md5.Sum([]byte(str))
	return hex.EncodeToString(md5Byte[:])
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func InArray[T uint | int | int8 | int64 | float32 | float64 | string](arr []T, item T) bool {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})

	index := sort.Search(len(arr), func(i int) bool {
		return arr[i] >= item
	})

	return index < len(arr) && arr[index] == item
}

func PasswordEncryption(password string) string {
	return Md5(Md5(Md5(password)))
}

// ToJSONString toJSON 将对象转换为JSON字符串，如果出错则返回"{}"
func ToJSONString(v any) string {
	if v == nil {
		return "{}"
	}
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func RedirectURL(rootUrl, provider string) string {
	return rootUrl + "/api/oauth/" + provider + "/callback"
}

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:length]
}
