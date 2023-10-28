package events

type EventSubscription[T interface{}] struct {
	Emitter *SyncEventEmitter[T]
	Handler EventCallable[T]
}

// Unsubscribes an event handler
func (subscription *EventSubscription[T]) unsubscribe() {
	subscription.Emitter.unsubscribe(subscription.Handler)
}
