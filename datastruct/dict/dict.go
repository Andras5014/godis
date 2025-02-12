package dict

type Consumer func(key string, value interface{}) bool
type Dict interface {
	Get(key string) (interface{}, bool)
	Len() int
	Put(key string, value interface{}) int
	PutIfAbsent(key string, value interface{}) int
	PutIfExists(key string, value interface{}) int
	Remove(key string) int
	ForEach(consumer Consumer)
	Keys() []string
	RandomKeys(size int) []string
	RandomDistinctKeys(size int) []string
	clear()
}
