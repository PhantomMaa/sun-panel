package common

import (
	"crypto/md5"
	"encoding/hex"
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

var Version string

func GetVersion() string {
	if Version == "" {
		Version = "v1.0.0"
	}
	return Version
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
