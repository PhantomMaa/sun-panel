package kvcache

import (
	"encoding/json"
	"time"

	"github.com/patrickmn/go-cache"
)

// 参考：https://blog.csdn.net/u014459543/article/details/108429469
type LocalCache[T any] struct {
	gocahce *cache.Cache
	Result  T
}

type GoCacheValue[T any] struct {
	Value T
}

// NewLocalCache 创建一个goCache结构体
// cache.New(5*time.Minute, 60*time.Second)，清理过期的item间隔 0.不清理
func NewLocalCache[T any](defaultExpiration time.Duration, cleanupInterval time.Duration) *LocalCache[T] {
	cacheAdapter := cache.New(defaultExpiration, cleanupInterval)
	return &LocalCache[T]{
		gocahce: cacheAdapter,
	}
}

func (c *LocalCache[T]) Set(k string, x T, d time.Duration) {
	c.gocahce.Set(k, GoCacheValue[T]{Value: x}, d)
}

func (c *LocalCache[T]) Get(k string) (T, bool) {
	if v, ok := c.gocahce.Get(k); ok {
		if value, okv := v.(GoCacheValue[T]); okv {
			return value.Value, true
		}
	}
	return c.Result, false
}

// SetDefault 设置cache 无时间参数
func (c *LocalCache[T]) SetDefault(k string, v T) {
	c.gocahce.SetDefault(k, GoCacheValue[T]{Value: v})
}

// SetKeepExpiration 设置并保持原始的过期时间
func (c *LocalCache[T]) SetKeepExpiration(k string, v T) {
	_, expirationTime, ok := c.gocahce.GetWithExpiration(k)

	now := time.Now()
	differ := expirationTime.Sub(now)
	// 如果 过期值不为零值 && 未过期 && 过期时间大于现在的时间
	// 将保持不变原始的过期时间来计算时间
	if !expirationTime.IsZero() && ok && differ > 0 {
		// newExpiration := now.Unix() + int64(math.Round(differ.Seconds()))
		// fmt.Println("旧的过期时间", expirationTime.Unix())
		// fmt.Println("时间限制差", math.Round(differ.Seconds()))
		// fmt.Println("新的过期时间", newExpiration)
		c.gocahce.Set(k, GoCacheValue[T]{Value: v}, differ)
	} else {
		c.gocahce.SetDefault(k, GoCacheValue[T]{Value: v})
	}
}

// Delete 删除 cache
func (c *LocalCache[T]) Delete(k string) {
	c.gocahce.Delete(k)
}

// Add 加入缓存
func (c *LocalCache[T]) Add(k string, v T, d time.Duration) {
	err := c.gocahce.Add(k, GoCacheValue[T]{Value: v}, d)
	if err != nil {
		return
	}
}

// IncrementInt 对已存在的key 值自增n
func (c *LocalCache[T]) IncrementInt(k string, n int) (num int, err error) {
	return c.gocahce.IncrementInt(k, n)
}

// ItemCount 获取已存在key的数量
func (c *LocalCache[T]) ItemCount() (int64, error) {
	return int64(c.gocahce.ItemCount()), nil
}

// Flush 删除当前已存在的所有key
func (c *LocalCache[T]) Flush() {
	c.gocahce.Flush()
}

func (c *LocalCache[T]) encode(value T) ([]byte, error) {
	return json.Marshal(value)
}

func (c *LocalCache[T]) decode(valueByte []byte, value T) error {
	return json.Unmarshal(valueByte, value)
}
