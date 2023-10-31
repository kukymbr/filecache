package filecache

import "sync"

func newKeysLocker() *keysLocker {
	return &keysLocker{keys: make(map[string]*keyLocker)}
}

type keysLocker struct {
	keys map[string]*keyLocker
}

func (k *keysLocker) lock(key string) {
	kl, ok := k.keys[key]
	if ok {
		kl.Lock()

		return
	}

	k.keys[key] = &keyLocker{}

	k.keys[key].Lock()
}

func (k *keysLocker) unlock(key string) {
	kl, ok := k.keys[key]
	if !ok {
		return
	}

	kl.Unlock()
}

type keyLocker struct {
	sync.Mutex
}
