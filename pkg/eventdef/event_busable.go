package eventdef

type EventBusable interface {
	Subscribe(eventType string, handler EventHandler) error
	Unsubscribe(handlerUUID string) error
	Publish(event Event) error
}
