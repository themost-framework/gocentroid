package events

type EventSubscription[T interface{}] struct {
	Emitter *SyncEventEmitter[T]
	Handler EventCallable[T]
}

func (subscription *EventSubscription[T]) unsubscribe() {
	subscription.Emitter.unsubscribe(subscription.Handler)
}
