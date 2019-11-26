package rabbitmq

import (
	"context"

	"github.com/jfbramlett/go-aop/pkg/common"
	"github.com/jfbramlett/go-aop/pkg/messaging"
	"github.com/streadway/amqp"
)

func NewRabbitMQReceiver(config Config, callback messaging.Callback, contentType messaging.MsgContentTypeCreator) (messaging.MessageReceiver, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, err
	}

	return &rabbitMQReceiver{connection: conn, channel: ch, queue: q, callback: callback, contentType: contentType}, nil
}

type rabbitMQReceiver struct {
	connection  *amqp.Connection
	channel     *amqp.Channel
	queue       amqp.Queue
	callback    messaging.Callback
	contentType messaging.MsgContentTypeCreator
}

func (rr *rabbitMQReceiver) Run() error {
	msgs, err := rr.channel.Consume(
		rr.queue.Name, // queue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			content := rr.contentType()
			envelope := messaging.Envelope{Content: &content}
			err := common.FromJSON(string(d.Body), &envelope)
			if err != nil {
				continue
			}
			err = rr.callback(context.Background(), envelope.Content)
			if err != nil {
				continue
			}
		}
	}()

	return nil
}

func (rr *rabbitMQReceiver) Close() {
	_ = rr.channel.Close()
	_ = rr.connection.Close()
}
