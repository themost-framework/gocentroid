package events

type EventArgs struct{}

type EventCallable[T interface{}] func(event *T)

// A class which wraps an event emitter handler
type SyncEventHandler[T interface{}] struct {
	Handler EventCallable[T]
}

// Executes event handler
func (eventHandler *SyncEventHandler[T]) execute(event *T) {
	eventHandler.Handler(event)
}
