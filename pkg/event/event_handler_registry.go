package event

type EventHandlerRegistry struct {
	allHandlers  map[string]EventHandler
	handlerOrder []string
}
