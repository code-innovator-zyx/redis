package commands

import (
	"sync"
)

/*
* @Author: zouyx
* @Email:
* @Date:   2024/9/10 下午3:56
* @Package:
 */

type hset struct {
	data map[string]map[string]string
	sync.RWMutex
}

func newHset() *hset {
	return &hset{
		data:    make(map[string]map[string]string),
		RWMutex: sync.RWMutex{},
	}
}

func (h *hset) hset(hash, key, value string) {
	h.Lock()
	defer h.Unlock()
	if _, ok := h.data[hash]; !ok {
		// 不存在，初始化map
		h.data[hash] = make(map[string]string)
	}
	h.data[hash][key] = value
}
func (h *hset) hget(hash, key string) (string, bool) {
	h.RLock()
	defer h.RUnlock()
	if v, ok := h.data[hash]; ok {
		if value, ok := v[key]; ok {
			return value, true
		}
	}
	return "", false
}

func (h *hset) hdel(hash, key string) {
	h.Lock()
	defer h.Unlock()
	if _, ok := h.data[hash]; ok {
		// 存在删除
		delete(h.data[hash], key)
	}
}
