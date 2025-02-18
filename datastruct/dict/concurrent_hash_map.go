package dict

import (
	"crypto/sha256"
	"sync"
)

type ConcurrentHashMap struct {
	mps   []map[string]any
	seg   int
	locks []sync.RWMutex
	seed  int
}

func NewConcurrentHashMap(seg, cap int) *ConcurrentHashMap {
	mps := make([]map[string]any, seg)
	locks := make([]sync.RWMutex, seg)
	for i := 0; i < seg; i++ {
		mps[i] = make(map[string]any, cap/seg)
	}
	return &ConcurrentHashMap{
		mps:   mps,
		seg:   seg,
		locks: locks,
		seed:  0,
	}
}

func (cm *ConcurrentHashMap) getSegIndex(key string) int {
	sha := sha256.New()
	hash, _ := sha.Write([]byte(key))
	return hash % cm.seg
}

func (cm *ConcurrentHashMap) Set(key string, value any) {
	segIndex := cm.getSegIndex(key)
	cm.locks[segIndex].Lock()
	defer cm.locks[segIndex].Unlock()
	cm.mps[segIndex][key] = value
}

func (cm *ConcurrentHashMap) Get(key string) (any, bool) {
	segIndex := cm.getSegIndex(key)
	cm.locks[segIndex].RLock()
	defer cm.locks[segIndex].RUnlock()
	value, ok := cm.mps[segIndex][key]
	return value, ok
}

func (cm *ConcurrentHashMap) Remove(key string) {
	segIndex := cm.getSegIndex(key)
	cm.locks[segIndex].Lock()
	defer cm.locks[segIndex].Unlock()
	delete(cm.mps[segIndex], key)
}

func (cm *ConcurrentHashMap) CreateIterator() MapIterator {
	keys := make([][]string, 0, len(cm.mps))
	for i, mp := range cm.mps {
		for key, _ := range mp {
			keys[i] = append(keys[i], key)
		}
	}
	return &ConcurrentMapIterator{
		cm:       cm,
		keys:     keys,
		rowIndex: 0,
		colIndex: 0,
	}
}

type MapEntry struct {
	Key   string
	Value any
}

type MapIterator interface {
	Next() *MapEntry
}
type ConcurrentMapIterator struct {
	cm       *ConcurrentHashMap
	keys     [][]string
	rowIndex int
	colIndex int
}

func (it *ConcurrentMapIterator) Next() *MapEntry {
	if it.rowIndex >= len(it.keys) {
		return nil
	}
	row := it.keys[it.rowIndex] // 当前行
	if len(row) == 0 {          //当前行没有元素，跳过
		it.rowIndex++
		return it.Next()
	}

	key := row[it.colIndex]
	value, _ := it.cm.Get(key)
	if it.colIndex < len(row)-1 {
		it.colIndex++
	} else {
		it.colIndex = 0
		it.rowIndex++
	}
	return &MapEntry{
		Key:   key,
		Value: value,
	}
}
