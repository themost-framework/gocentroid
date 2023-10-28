package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEventEmitter(t *testing.T) {
	var emitter = &SyncEventEmitter[EventArgs]{}
	var counter = 0
	handler := func(event *EventArgs) {
		counter++
	}
	subscription := emitter.subscribe(handler)
	assert.NotEmpty(t, emitter.Handlers)
	subscription.unsubscribe()
	assert.Empty(t, emitter.Handlers)
}

func TestEmitEvent(t *testing.T) {
	var emitter = &SyncEventEmitter[EventArgs]{}
	var counter = 0
	handler := func(event *EventArgs) {
		counter++
	}
	var subscription = emitter.subscribe(handler)
	var event = &EventArgs{}
	emitter.emit(event)
	assert.Equal(t, 1, counter)
	subscription.unsubscribe()
	emitter.emit(event)
	assert.Equal(t, 1, counter)

}

type CustomEventArgs struct {
	EventArgs,
	count int
}

func TestEmitCustomEvent(t *testing.T) {
	var emitter = &SyncEventEmitter[CustomEventArgs]{}
	handler := func(event *CustomEventArgs) {
		event.count++
	}
	var subscription = emitter.subscribe(handler)
	event := &CustomEventArgs{
		count: 0,
	}
	emitter.emit(event)
	assert.Equal(t, 1, event.count)
	subscription.unsubscribe()
	emitter.emit(event)
	assert.Equal(t, 1, event.count)

}
