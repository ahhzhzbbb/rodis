package engine

import (
	"sync"
)

type KeyValue struct {
	Kv map[string]string
	Mu sync.RWMutex
}

func NewKeyValue() *KeyValue {
	return &KeyValue{
		Kv: make(map[string]string),
	}
}
