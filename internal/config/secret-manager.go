package config

import "sync"

type Keys struct {
	sync.RWMutex
	key   string // []byte
	appID int
}

func (k *Keys) GetKey() string {
	k.RLock()
	defer k.RUnlock()
	return k.key
}

func (k *Keys) SetKey(key string) {
	k.Lock()
	defer k.Unlock()
	k.key = key
}
