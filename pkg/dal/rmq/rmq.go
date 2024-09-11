package rmq

import (
	"fmt"

	"github.com/asciifaceman/emri/pkg/global"
	"github.com/rabbitmq/amqp091-go"
)

func NewWithConfig(c *QueueConfig) (*RMQ, error) {
	r := &RMQ{
		conf: c,
	}

	return r, r.Connect()
}

func New(queue string, workers int, handler func(amqp091.Delivery)) (*RMQ, error) {
	c := &QueueConfig{
		Workers:   workers,
		msgFunc:   handler,
		Name:      queue,
		MaxLength: 50000,
		TTL:       300000,
		Overflow:  "reject-publish",
	}
	return NewWithConfig(c)
}

type RMQ struct {
	c    *amqp091.Connection
	ch   *amqp091.Channel
	conf *QueueConfig
}

func (r *RMQ) Connect() error {
	conn, err := amqp091.Dial(global.C().AMQPUrl())
	if err != nil {
		return err
	}
	r.c = conn

	ch, err := r.c.Channel()
	if err != nil {
		return err
	}
	r.ch = ch

	return nil
}

type QueueConfig struct {
	Workers int

	msgFunc func(amqp091.Delivery)

	Name      string
	MaxLength int
	TTL       int
	Overflow  string
}

func (c *QueueConfig) BuildArgs() amqp091.Table {
	queueArgs := amqp091.Table{}
	queueArgs["x-dead-letter-exchange"] = fmt.Sprintf("%s-dlx", c.Name)
	queueArgs["x-max-length"] = c.MaxLength
	queueArgs["x-message-ttl"] = c.TTL
	queueArgs["x-overflow"] = c.Overflow
	return queueArgs
}
