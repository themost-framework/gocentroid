package events

type EventArgs struct{}

type EventCallable[T interface{}] func(event *T)

type SyncEventHandler[T interface{}] struct {
	Handler EventCallable[T]
}

func (eventHandler *SyncEventHandler[T]) execute(event *T) {
	eventHandler.Handler(event)
}
