package internal

import (
	"crypto/sha256"
	"github.com/yourbasic/radix"
	"hash"
	"sync"
	"unsafe"
)

var mu sync.Mutex
var blocks = map[string][]interface{}{}

func PutAndGetStates(b []interface{}) []interface{} {
	h := hashBlockStates(b)
	mu.Lock()
	defer mu.Unlock()

	if states, ok := blocks[h]; ok {
		return states
	}
	blocks[h] = b
	return b
}

func hashBlockStates(b []interface{}) string {
	h := sha256.New()
	var k []string
	for _, block := range b {
		m, _ := block.(map[string]interface{})
		blockData, _ := m["block"]
		blockMap, _ := blockData.(map[string]interface{})
		nameData, _ := blockMap["name"]
		name, _ := nameData.(string)
		propertyData, _ := blockMap["states"]
		properties, _ := propertyData.(map[string]interface{})

		h.Write(*(*[]byte)(unsafe.Pointer(&name)))
		k = hashProperties(properties, h, k)
	}
	return string(h.Sum(nil))
}

// hashProperties produces a hash for the block properties map passed.
// Passing the same map into hashProperties will always result in the same hash.
func hashProperties(properties map[string]interface{}, h hash.Hash, keys []string) []string {
	keys = keys[:0]
	for k := range properties {
		keys = append(keys, k)
	}
	radix.Sort(keys)

	for _, k := range keys {
		switch v := properties[k].(type) {
		case bool:
			if v {
				h.Write([]byte{1})
			} else {
				h.Write([]byte{0})
			}
		case uint8:
			h.Write([]byte{v})
		case int32:
			a := uint32(v)
			h.Write([]byte{byte(a), byte(a >> 8), byte(a >> 16), byte(a >> 24)})
		case string:
			h.Write(*(*[]byte)(unsafe.Pointer(&v)))
		}
	}
	return keys
}
