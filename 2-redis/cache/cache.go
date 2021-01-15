package cache

type StringCache interface {
	Set(key string, value string)
	Get(key string) string
}
