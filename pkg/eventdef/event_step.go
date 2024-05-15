package eventdef

type EventStep string

const (
	EventStepStart   EventStep = "Start"
	EventStepSuccess EventStep = "Success"
	EventStepError   EventStep = "Error"
)
