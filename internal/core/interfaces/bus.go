package interfaces

type EventBus interface {
	Send(v Event)
	Close()
	Receive() <-chan Event
}
