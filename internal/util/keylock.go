package util

import "sync"

func NewKeysLocker() *KeysLocker {
	return &KeysLocker{keys: make(map[string]*KeyLocker)}
}

type KeysLocker struct {
	keys map[string]*KeyLocker
}

func (k *KeysLocker) Lock(key string) {
	kl, ok := k.keys[key]
	if ok {
		kl.Lock()

		return
	}

	k.keys[key] = &KeyLocker{}

	k.keys[key].Lock()
}

func (k *KeysLocker) Unlock(key string) {
	kl, ok := k.keys[key]
	if !ok {
		return
	}

	kl.Unlock()
}

type KeyLocker struct {
	sync.Mutex
}
