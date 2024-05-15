package event

type Event interface {
	GetName() string
	GetStep() EventStep
	GetPayloadFormat() string
	GetPayload() []byte
}
