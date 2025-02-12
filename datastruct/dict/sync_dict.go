package dict

import "sync"

type SyncDict struct {
	m sync.Map
}

func NewSyncDict() *SyncDict {
	return &SyncDict{}
}

func (s *SyncDict) Get(key string) (interface{}, bool) {
	value, ok := s.m.Load(key)
	return value, ok
}

func (s *SyncDict) Len() int {
	length := 0
	s.m.Range(func(key, value interface{}) bool {
		length++
		return true
	})
	return length
}

func (s *SyncDict) Put(key string, value interface{}) int {
	_, existed := s.m.Load(key)
	if existed {
		return 0
	}
	s.m.Store(key, value)
	return 1
}

func (s *SyncDict) PutIfAbsent(key string, value interface{}) int {
	_, existed := s.m.LoadOrStore(key, value)
	if existed {
		return 0
	}
	return 1
}
func (s *SyncDict) PutIfExists(key string, value interface{}) int {
	_, existed := s.m.Load(key)
	if !existed {
		return 0
	}
	s.m.Store(key, value)
	return 1
}
func (s *SyncDict) Remove(key string) int {
	_, existed := s.m.Load(key)
	if !existed {
		return 0
	}
	s.m.Delete(key)
	return 1
}

func (s *SyncDict) ForEach(consumer Consumer) {
	s.m.Range(func(key, value interface{}) bool {
		consumer(key.(string), value)
		return true
	})
}

func (s *SyncDict) Keys() []string {
	result := make([]string, s.Len())
	s.m.Range(func(key, value interface{}) bool {
		result = append(result, key.(string))
		return true
	})
	return result
}

func (s *SyncDict) RandomKeys(size int) []string {
	result := make([]string, size)
	for i := 0; i < size; i++ {
		result[i] = s.Keys()[i]
	}
	return result
}

func (s *SyncDict) RandomDistinctKeys(size int) []string {
	result := make([]string, size)
	s.m.Range(func(key, value interface{}) bool {
		result = append(result, key.(string))
		if len(result) == size {
			return false
		}
		return true
	})
	return result
}

func (s *SyncDict) Clear() {
	//s.m.Range(func(key, value interface{}) bool {
	//	s.m.Delete(key)
	//	return true
	//})
	s.m = sync.Map{} // 清空所有元素 上一个m自动gc
}
