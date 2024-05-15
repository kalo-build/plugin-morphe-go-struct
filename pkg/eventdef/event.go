package eventdef

type Event interface {
	GetName() string
	GetStep() EventStep
	GetPayloadFormat() string
	GetPayload() []byte
}
