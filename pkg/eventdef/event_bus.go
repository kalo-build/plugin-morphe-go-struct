package eventdef

import (
	"sync"

	"github.com/google/uuid"
)

type EventBus struct {
	mutex sync.RWMutex

	subscribers map[string]EventHandlerRegistry
}

func (b *EventBus) Subscribe(eventName string, handler EventHandler) (string, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	handlerUUID := uuid.NewString()
	eventSubscribers, subscribersExist := b.subscribers[eventName]
	if subscribersExist {
		eventSubscribers.allHandlers[handlerUUID] = handler
		eventSubscribers.handlerOrder = append(eventSubscribers.handlerOrder, handlerUUID)
		b.subscribers[eventName] = eventSubscribers
		return handlerUUID, nil
	}

	b.subscribers[eventName] = EventHandlerRegistry{
		allHandlers: map[string]EventHandler{
			handlerUUID: handler,
		},
		handlerOrder: []string{handlerUUID},
	}
	return handlerUUID, nil
}

func (b *EventBus) Publish(event Event) error {
	eventName := event.GetName()
	b.mutex.Lock()
	eventSubscribers, subscribersExist := b.subscribers[eventName]
	b.mutex.Unlock()

	if !subscribersExist {
		return nil
	}

	for subscriberIdx := 0; subscriberIdx < len(eventSubscribers.handlerOrder); subscriberIdx++ {
		handlerUUID := eventSubscribers.handlerOrder[subscriberIdx]
		subscriberHandler, subscriberExists := eventSubscribers.allHandlers[handlerUUID]
		if !subscriberExists {
			return ErrSubscriberNotFound(handlerUUID)
		}

		subscriberErr := subscriberHandler(event)
		if subscriberErr != nil {
			return ErrWrapSubscriberError(eventName, handlerUUID, subscriberErr)
		}
	}

	return nil
}
