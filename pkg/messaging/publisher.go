package messaging

import "context"

type MessageSender interface {
	Publish(ctx context.Context, msg interface{}) error
	Close()
}
