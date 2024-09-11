package commands

import "sync"

/*
* @Author: zouyx
* @Email: Set 数据结构
* @Date:   2024/9/10 下午3:33
* @Package:
 */
type set struct {
	data map[string]string
	sync.RWMutex
}

func newSet() *set {
	return &set{
		data:    make(map[string]string),
		RWMutex: sync.RWMutex{},
	}
}
func (s *set) set(key, value string) {
	s.Lock()
	defer s.Unlock()
	s.data[key] = value
}

func (s *set) get(key string) (string, bool) {
	s.RLock()
	defer s.RUnlock()
	if value, ok := s.data[key]; ok {
		return value, true
	}
	return "", false
}

func (s *set) del(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.data, key)
}
