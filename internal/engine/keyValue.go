package engine

import (
	"time"

	"github.com/tidwall/shardmap"
)

type KeyValue struct {
	kv shardmap.Map
	et shardmap.Map
	// Mu sync.RWMutex
}

func NewKeyValue() *KeyValue {
	return &KeyValue{
		kv: *shardmap.New(1024),
		et: *shardmap.New(1024),
	}
}

func (k *KeyValue) Get(key string) (string, bool) {
	val, ok := k.et.Get(key)
	t, _ := val.(time.Time)

	if ok && t.Before(time.Now()) {
		k.kv.Delete(key)
		k.et.Delete(key)
		return "", false
	}
	temp, exists := k.kv.Get(key)
	res, _ := temp.(string)
	return res, exists
}

func (k *KeyValue) Set(key, value string) {
	if _, exists := k.kv.Get(key); exists {
		k.et.Delete(key)
	}
	k.kv.Set(key, value)
}

func (k *KeyValue) Del(key string) bool {
	var rs bool
	rs = false

	if _, ok := k.kv.Get(key); ok {
		rs = true
		k.kv.Delete(key)
		k.et.Delete(key)
	}

	return rs
}

func (k *KeyValue) DelValue(key string) {
	k.kv.Delete(key)
}

func (k *KeyValue) SetExpireTime(key string, t time.Time) bool {
	if _, exists := k.kv.Get(key); !exists {
		return false
	}

	k.et.Set(key, t)

	return true
}
