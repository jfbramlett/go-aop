package lock

import "sync"

// NewConditionalMutex creates a new locker/mutex where the lock/unlock is conditional
// based on whether we are flagged as locking or not
func NewConditionalMutex(enabled bool) sync.Locker {
    return &conditionalMutex{enabled: enabled}
}

type conditionalMutex struct {
    enabled     bool
    mux         sync.Mutex
}

func (o *conditionalMutex) Lock() {
    if o.enabled {
        o.mux.Lock()
    }
}

func (o *conditionalMutex) Unlock() {
    if o.enabled {
        o.mux.Unlock()
    }
}
