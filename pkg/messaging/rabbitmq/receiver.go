package rabbitmq

import (
	"context"
	"github.com/jfbramlett/go-aop/pkg/common"
	"github.com/jfbramlett/go-aop/pkg/messaging"
	"github.com/streadway/amqp"
)


func NewRabbitMQReceiver(config Config) (messaging.MessageReceiver, error) {
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


	return &rabbitMQReceiver{connection: conn, channel: ch, queue: q}, nil
}



type rabbitMQReceiver struct {
	connection		*amqp.Connection
	channel 		*amqp.Channel
	queue 			amqp.Queue
	callback		messaging.Callback

}


func (rr *rabbitMQReceiver) OnMessage(ctx context.Context, callback messaging.Callback) error {
	msgs, err := rr.channel.Consume(
		rr.queue.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			envelope := messaging.Envelope{}
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

	<-forever
}

func (rr *rabbitMQReceiver) Close() {
	_ = rr.channel.Close()
	_ = rr.connection.Close()
}
