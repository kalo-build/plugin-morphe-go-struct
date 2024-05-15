package eventdef

import "sync"

// InitEventBus creates a new instance of EventBus.
func InitEventBus() *EventBus {
	return &EventBus{
		mutex:       sync.RWMutex{},
		subscribers: make(map[string]EventHandlerRegistry),
	}
}
