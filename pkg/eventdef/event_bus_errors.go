package eventdef

import "fmt"

func ErrSubscriberNotFound(handlerUUID string) error {
	return fmt.Errorf("subscriber handler with uuid '%s' not found", handlerUUID)
}

func ErrWrapSubscriberError(eventName string, handlerUUID string, subscriberErr error) error {
	return fmt.Errorf("event '%s' subscriber with uuid '%s' encountered an error: %e", eventName, handlerUUID, subscriberErr)
}
