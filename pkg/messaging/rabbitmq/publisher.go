package rabbitmq

import (
	"context"
	"encoding/json"
	"github.com/jfbramlett/go-aop/pkg/messaging"
	"github.com/streadway/amqp"
	"strings"
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

	content := &strings.Builder{}
	enc := json.NewEncoder(content)
	enc.SetIndent("", "    ")
	if err := enc.Encode(envelope); err != nil {
		return err
	}

	err := r.channel.Publish(
		"",     // exchange
		r.queue.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(content.String()),
		})

	return err
}

func (r *rabbitMQSender) Close() {
	_ = r.channel.Close()
	_ = r.connection.Close()
}