package rabbitmq

type Config struct {
	Host       	string `config:"rabbitHost"`
	Channel   	string `config:"rabbitChannel"`
}
