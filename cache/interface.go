package cache

// Cache 缓存， 用于存储验证码缓存
type Cache interface {
	Set(key string, value interface{}, ttl int) error
	Get(key string) string
	Del(key string) error
}
