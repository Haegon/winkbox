package simplegamed

import (
	"sync"
)

// SafeMap은 map에 데이터를 넣고 읽고 지울때 락을 걸어 동시성을
// 제어하기 위한 구조체 이다.
type SafeMap struct {
	m  map[interface{}]interface{}
	mu sync.RWMutex
}

// NewSafeMap는 새 SafeMap을 생성한다.
func NewSafeMap() *SafeMap {
	return &SafeMap{m: make(map[interface{}]interface{})}
}

// get은 map에서 값을 쓸 때 사용한다.
func (s *SafeMap) get(key interface{}) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.m[key]
	return value, ok
}

// set은 map에 값을 넣을떄 사용한다.
func (s *SafeMap) set(key, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[key] = value
}

// del는 map에서 값을 삭제 할 때 사용한다.
func (s *SafeMap) del(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.m, key)
}
