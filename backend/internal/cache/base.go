package cache

import (
	"time"
)

// Cacher 缓存接口-基于内存实现
type Cacher[T any] interface {
	Set(k string, v T, d time.Duration)

	Get(k string) (T, bool)

	// SetDefault 设置-过期时间采用默认值
	SetDefault(k string, v T)

	Delete(k string)

	// SetKeepExpiration 设置值，但不重置过期时间
	SetKeepExpiration(k string, v T)
}
