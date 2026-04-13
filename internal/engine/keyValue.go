package engine

import (
	"fmt"
	"sync"
	"time"
)

type KeyValue struct {
	kv map[string]string
	et map[string]time.Time
	Mu sync.RWMutex
}

func NewKeyValue() *KeyValue {
	return &KeyValue{
		kv: make(map[string]string),
		et: make(map[string]time.Time),
	}
}

func (k *KeyValue) Get(key string) (string, bool) {
	k.Mu.Lock()
	defer k.Mu.Unlock()

	t, ok := k.et[key]

	if ok && t.Before(time.Now()) {
		fmt.Println("hello")
		delete(k.kv, key)
		delete(k.et, key)
		return "", false
	}

	res, exists := k.kv[key]

	return res, exists
}

func (k *KeyValue) Set(key, value string) {
	k.Mu.Lock()
	defer k.Mu.Unlock()

	if _, exists := k.kv[key]; exists {
		delete(k.et, key)
	}
	k.kv[key] = value
}

func (k *KeyValue) Del(key string) bool {
	k.Mu.Lock()
	defer k.Mu.Unlock()

	var rs bool
	rs = false

	if _, ok := k.kv[key]; ok {
		rs = true
		delete(k.kv, key)
		delete(k.et, key)
	}

	return rs
}

func (k *KeyValue) DelValue(key string) {
	k.Mu.Lock()
	defer k.Mu.Unlock()

	delete(k.kv, key)
}

func (k *KeyValue) SetExpireTime(key string, t time.Time) bool {
	k.Mu.Lock()
	defer k.Mu.Unlock()

	if _, exists := k.kv[key]; !exists {
		return false
	}

	k.et[key] = t

	return true
}
