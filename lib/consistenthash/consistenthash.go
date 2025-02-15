package consistenthash

import (
	"hash/crc32"
	"sort"
)

type HashFunc func(data []byte) uint32
type NodeMap struct {
	hashFunc    HashFunc
	nodeHashs   []int
	nodeHashMap map[int]string
}

func NewNodeMap(hashFunc HashFunc) *NodeMap {

	nm := &NodeMap{
		hashFunc:    hashFunc,
		nodeHashs:   make([]int, 0),
		nodeHashMap: make(map[int]string),
	}
	if nm.hashFunc == nil {
		nm.hashFunc = crc32.ChecksumIEEE
	}
	return nm
}

func (nm *NodeMap) IsEmpty() bool {
	return len(nm.nodeHashs) == 0
}

func (nm *NodeMap) AddNode(nodes ...string) {
	for _, node := range nodes {
		if node == "" {
			continue
		}
		hash := int(nm.hashFunc([]byte(node)))
		nm.nodeHashs = append(nm.nodeHashs, hash)
		nm.nodeHashMap[hash] = node
	}
	sort.Ints(nm.nodeHashs)
}

func (nm *NodeMap) PickNode(key string) string {
	if nm.IsEmpty() {
		return ""
	}
	hash := int(nm.hashFunc([]byte(key)))
	idx := sort.Search(len(nm.nodeHashs), func(i int) bool {
		return nm.nodeHashs[i] >= hash
	})
	if idx == len(nm.nodeHashs) {
		return nm.nodeHashMap[nm.nodeHashs[0]]
	}
	return nm.nodeHashMap[nm.nodeHashs[idx]]
}
func (nm *NodeMap) RemoveNode(nodes ...string) {
	for _, node := range nodes {
		if node == "" {
			continue
		}
		hash := int(nm.hashFunc([]byte(node)))
		for i, v := range nm.nodeHashs {
			if v == hash {
				nm.nodeHashs = append(nm.nodeHashs[:i], nm.nodeHashs[i+1:]...)
				delete(nm.nodeHashMap, hash)
				break
			}
		}
	}
}
