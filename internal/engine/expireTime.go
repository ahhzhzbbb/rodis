package engine

import (
	"sync"
	"time"
)

type ExpireTime struct {
	Et map[string]time.Time
	Mu sync.RWMutex
}

func NewExpireTime() *ExpireTime {
	return &ExpireTime{
		Et: make(map[string]time.Time),
	}
}
