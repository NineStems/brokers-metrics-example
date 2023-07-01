package rabbit

import (
	"context"
	"encoding/json"
	"mb-and-metrics/common"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"mb-and-metrics/internal/config"
	"mb-and-metrics/pkg/errors"
)

const (
	rabbitMark            = "RabbitMQ"
	rabbitOpenChannel     = "open channel"
	rabbitCloseChannel    = "close channel"
	rabbitExchangeDeclare = "exchange declare"
	rabbitQueueDeclare    = "queue declare"
	rabbitQueueBind       = "queue bind"
	rabbitBindConsumer    = "consumer bind"
	rabbitStartConsume    = "consume started"
	rabbitEndConsume      = "consume ended"

	rabbitPublish = "Publish message"
)

type Rabbit struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	log     common.Logger
	cfg     config.Rabbit
}

func New(conn *amqp.Connection, log common.Logger, cfg config.Rabbit) *Rabbit {
	return &Rabbit{
		conn: conn,
		log:  log,
		cfg:  cfg,
	}
}

// Start инициализирует работу с rabbitMQ.
func (r *Rabbit) Start() error {
	var err error
	r.channel, err = r.conn.Channel()
	if err != nil {
		return errors.Wrap(err, "a.rabbit.Channel")
	}

	r.log.Infow(rabbitMark, rabbitOpenChannel, "success")

	err = r.channel.ExchangeDeclare(
		r.cfg.Exchange, // name
		"topic",        // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return errors.Wrap(err, "channel.ExchangeDeclare")
	}

	r.log.Infow(rabbitMark, rabbitExchangeDeclare, r.cfg.Exchange)

	queue, err := r.channel.QueueDeclare(
		r.cfg.Queue, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return errors.Wrap(err, "channel.QueueDeclare")
	}

	r.log.Infow(rabbitMark, rabbitQueueDeclare, r.cfg.Queue)

	err = r.channel.QueueBind(
		queue.Name,     // queue name
		r.cfg.Key,      // routing key
		r.cfg.Exchange, // exchange
		false,
		nil)
	if err != nil {
		return errors.Wrap(err, "a.rabbit.Channel")
	}

	r.log.Infow(rabbitMark, rabbitQueueBind, "success")

	return nil
}

func (r *Rabbit) Close() {
	if r.channel != nil {
		r.channel.Close()
		r.log.Infow(rabbitMark, rabbitCloseChannel, "success")
	}
}

// Publish публикует сообщение в exchange заданный в конфиге.
func (r *Rabbit) Publish(ctx context.Context, message interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	body, err := json.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	r.log.Infow(rabbitMark, rabbitPublish, string(body))

	err = r.channel.PublishWithContext(ctx,
		r.cfg.Exchange, // exchange
		r.cfg.Key,      // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return errors.Wrap(err, "r.channel.PublishWithContext")
	}

	return nil
}

// Read читает сообщение из заданной в конфиге очереди.
func (r *Rabbit) Read(_ context.Context, message chan []byte) error {
	msgs, err := r.channel.Consume(
		r.cfg.Queue, // queue
		"sender",    // consumer
		true,        // auto ack
		false,       // exclusive
		false,       // no local
		false,       // no wait
		nil,         // args
	)
	if err != nil {
		return errors.Wrap(err, "r.channel.Consume")
	}
	r.log.Infow(rabbitMark, rabbitBindConsumer, "success")

	go func() {
		r.log.Infow(rabbitMark, rabbitStartConsume, "success")
		for body := range msgs {
			message <- body.Body
		}
		r.log.Infow(rabbitMark, rabbitEndConsume, "success")
	}()

	return nil
}
