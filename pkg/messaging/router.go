package messaging

import (
	"context"
)

const (
	serviceName = "ingestion"
)

type OutboxHandlerFunc func(ctx context.Context, msg *Message) error

// outboxManager is a local interface used to abstract the outbox.Manager in a manner that allows
// for mocking/testing
type outboxManager interface {
	Advance(id int64)
	Messages() chan *Message
}

// OutboxRouter is the interface for a service that routes outbox messages to a handler
type OutboxRouter interface {
	HandlerFunc(event string, f OutboxHandlerFunc)
	WithMiddleware(f MiddlewareFunc) OutboxRouter``
	ListenAndServeAsync()
	ListenAndServe()
	Shutdown()
}

// NewOutboxRouter constructs a new outbox router off the given table and lock
func NewOutboxRouter(dbWriter *sqlx.DB, table string, lock string) OutboxRouter {
	manager := outbox.NewManager(dbWriter, table, lock)
	return &outboxRouter{manager: manager, handlers: make(map[string]OutboxHandlerFunc, 0)}
}

// OutboxRouter mimics the behavior of mux.Router but for outbox messages. Instead of requiring
// a switch statement for processing of msgs this works via registration of a handler for a
// regex pattern for the event.
type outboxRouter struct {
	manager    outboxManager
	handlers   map[string]OutboxHandlerFunc
	middleware []MiddlewareFunc
}

func (o *outboxRouter) WithMiddleware(f MiddlewareFunc) OutboxRouter {
	o.middleware = append(o.middleware, f)
}

// HandlerFunc adds a new handler for a given event - event can be a regex
func (o *outboxRouter) HandlerFunc(event string, f OutboxHandlerFunc) {
	o.handlers[event] = f
}

// ListenAndServeAsync starts listening for outbox messages asynchronously (i.e. as a go func)
func (o *outboxRouter) ListenAndServeAsync() {
	go o.ListenAndServe()
}

// ListenAndServe listens for messages from the outbox and processes them - this is a blocking call
// so if needing to run in a goroutine use ListenAndServeAsync
func (o *outboxRouter) ListenAndServe() {
	for message := range o.manager.Messages() {
		if err := o.handleMessage(message); err == nil {
			o.manager.Advance(message.ID)
		}
	}
}

// Shutdown stops listening for messages
func (o *outboxRouter) Shutdown() {
	close(o.manager.Messages())
}

// handleMessage takes an incoming outbox message and determines if it is smoething we have a registered
// handler for and if so invokes the handler
func (o *outboxRouter) handleMessage(msg *Message) error {
	if f, found := o.handlers[msg.Event]; found {

		for i := len(o.middleware) - 1; i >= 0; i-- {
			f = o.middleware[i].Middleware(f)
		}

		err := f(context.Background(), msg)
		return err
	}

	return nil
}
