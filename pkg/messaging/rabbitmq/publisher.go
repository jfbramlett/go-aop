package rabbitmq

import (
	"context"
	"github.com/jfbramlett/go-aop/pkg/common"
	"github.com/jfbramlett/go-aop/pkg/messaging"
	"github.com/streadway/amqp"
)

func NewRabbitMQSender(config Config) (messaging.MessageSender, error) {
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


	return &rabbitMQSender{connection: conn, channel: ch, queue: q}, nil
}

type rabbitMQSender struct {
	connection		*amqp.Connection
	channel 		*amqp.Channel
	queue 			amqp.Queue
}

func (r *rabbitMQSender) Publish(ctx context.Context, msg interface{}) error {
	envelope := messaging.Envelope{Content: msg}

	content, err := common.ToJSON(envelope)
	if err != nil {
		return err
	}
	err = r.channel.Publish(
		"",     // exchange
		r.queue.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(content),
		})

	return err
}

func (r *rabbitMQSender) Close() {
	_ = r.channel.Close()
	_ = r.connection.Close()
}