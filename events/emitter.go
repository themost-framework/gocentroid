package events

import (
	"reflect"

	"github.com/Goldziher/go-utils/sliceutils"
)

type SyncEventEmitter[T interface{}] struct {
	Handlers []SyncEventHandler[T]
}

func (emitter *SyncEventEmitter[T]) subscribe(handler EventCallable[T]) *EventSubscription[T] {
	if emitter.Handlers == nil {
		emitter.Handlers = []SyncEventHandler[T]{}
	}
	eventHandler := &SyncEventHandler[T]{
		Handler: handler,
	}
	emitter.Handlers = append(emitter.Handlers, *eventHandler)
	subscription := &EventSubscription[T]{
		Emitter: emitter,
		Handler: handler,
	}
	return subscription
}

func (emitter *SyncEventEmitter[T]) unsubscribe(handler EventCallable[T]) {
	if emitter.Handlers == nil {
		emitter.Handlers = []SyncEventHandler[T]{}
	}
	var comparer = func(value SyncEventHandler[T], index int, slice []SyncEventHandler[T]) bool {
		return reflect.ValueOf(value.Handler) == reflect.ValueOf(handler)
	}
	var index = sliceutils.FindIndex[SyncEventHandler[T]](emitter.Handlers, comparer)
	if index >= 0 {
		var handlers = sliceutils.Remove[SyncEventHandler[T]](emitter.Handlers, index)
		emitter.Handlers = handlers
	}
}

func (emitter *SyncEventEmitter[T]) emit(event *T) {
	for _, v := range emitter.Handlers {
		v.execute(event)
	}
}
