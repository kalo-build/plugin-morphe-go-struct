package eventdef

type EventHandlerRegistry struct {
	allHandlers  map[string]EventHandler
	handlerOrder []string
}
