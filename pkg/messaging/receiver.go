package messaging

import "context"

type Callback func(ctx context.Context, msg interface{}) error
type MsgContentTypeCreator func() interface{}

type MessageReceiver interface {
	Run() error
	Close()
}
