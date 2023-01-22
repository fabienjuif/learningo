package utils

import "sync"

// TODO: it should exist a lib doing this better

type SafeMap[K comparable, V interface{}] struct {
	mu sync.Mutex
	v  map[K]*V
}

func NewSafeMap[K comparable, V interface{}](m map[K]*V) SafeMap[K, V] {
	return SafeMap[K, V]{v: m}
}

func (safeMap *SafeMap[K, V]) Add(key K, value *V) (*V, error) {
	safeMap.mu.Lock()
	defer safeMap.mu.Unlock()
	_, exists := safeMap.v[key]
	if exists {
		return nil, &SafeMapAlreadyExistsError{}
	}
	safeMap.v[key] = value
	return value, nil
}

func (safeMap *SafeMap[K, V]) Get(key K) *V {
	safeMap.mu.Lock()
	defer safeMap.mu.Unlock()
	return safeMap.v[key]
}

func (safeMap *SafeMap[K, V]) Del(key K) {
	safeMap.mu.Lock()
	defer safeMap.mu.Unlock()
	delete(safeMap.v, key)
}

type SafeMapAlreadyExistsError struct{}

func (m *SafeMapAlreadyExistsError) Error() string {
	return "Key already exists"
}
