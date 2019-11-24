package messaging

import "context"


type Callback func(ctx context.Context, msg interface{}) error

type MessageReceiver interface {
	OnMessage(ctx context.Context, callback Callback) error
	Close()
}