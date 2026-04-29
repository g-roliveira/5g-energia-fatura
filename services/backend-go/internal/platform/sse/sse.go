package sse

import (
	"fmt"
	"io"
	"sync"
)

// Broker is a pub/sub dispatcher for SSE events.
// Clients subscribe by key and receive byte payloads on channels.
type Broker struct {
	mu          sync.Mutex
	subscribers map[string]map[chan []byte]struct{}
}

// NewBroker creates a new Broker.
func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string]map[chan []byte]struct{}),
	}
}

// Subscribe returns a channel that receives all payloads published for key.
func (b *Broker) Subscribe(key string) chan []byte {
	ch := make(chan []byte, 64)
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.subscribers[key] == nil {
		b.subscribers[key] = make(map[chan []byte]struct{})
	}
	b.subscribers[key][ch] = struct{}{}
	return ch
}

// Unsubscribe removes a channel from the broker for the given key and closes it.
func (b *Broker) Unsubscribe(key string, ch chan []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if subs, ok := b.subscribers[key]; ok {
		if _, exists := subs[ch]; exists {
			delete(subs, ch)
			close(ch)
		}
		if len(subs) == 0 {
			delete(b.subscribers, key)
		}
	}
}

// Publish sends data to all subscribers of key. Drops if a subscriber's buffer is full.
func (b *Broker) Publish(key string, data []byte) {
	b.mu.Lock()
	subs := b.subscribers[key]
	b.mu.Unlock()
	for ch := range subs {
		select {
		case ch <- data:
		default:
			// Drop slow consumers
		}
	}
}

// Close removes and closes all subscriber channels.
func (b *Broker) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for key, subs := range b.subscribers {
		for ch := range subs {
			close(ch)
		}
		delete(b.subscribers, key)
	}
}

// WriteSSE writes a single SSE event (event type + data) to w.
func WriteSSE(w io.Writer, eventType, data string) {
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, data)
}
