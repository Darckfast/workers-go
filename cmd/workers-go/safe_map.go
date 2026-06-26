//go:build !js && !wasm

package main

import (
	_ "embed"
	"go/ast"
	"sync"
)

type SafeMap[K comparable, V any] struct {
	m  map[K]V
	mu sync.RWMutex
}

func (sm *SafeMap[K, V]) Set(key K, val V) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.m == nil {
		sm.m = make(map[K]V)
	}

	sm.m[key] = val
}

func (sm *SafeMap[K, V]) Get(key K) (V, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	val, ok := sm.m[key]

	return val, ok
}

var AllStructs = SafeMap[string, *ast.StructType]{}

type DOM = SafeMap[string, []DurableObjectFunc]
type WM = SafeMap[string, []string]
type A2K = map[handlerName]struct{}
