package amqp

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	exRatingUpdated = "rating_updated"

	queueName = "catalog"

	exKind        = "fanout"
	prefetchCount = 5
)

type Receiver struct {
	ch      *amqp.Channel
	handler *Handler
}

func NewReceiver(ch *amqp.Channel, h *Handler) *Receiver {
	return &Receiver{
		ch:      ch,
		handler: h,
	}
}

func (r *Receiver) Receive() {
	if err := r.ch.Qos(
		prefetchCount,
		0,
		false,
	); err != nil {
		logrus.Fatalln("Failed to set Qos", err)
	}

	exchanges := []string{exRatingUpdated}

	for _, ex := range exchanges {

		dCh := r.deliveryCh(ex)

		switch ex {
		case exRatingUpdated:
			go func() {
				for d := range dCh {
					r.handler.RatingUpdated(&d)
				}
			}()
		default:
			return
		}
	}

}

func (r *Receiver) deliveryCh(ex string) <-chan amqp.Delivery {
	queue := fmt.Sprintf("%s:%s", ex, queueName)

	_, err := r.ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.Fatalln("Failed to declare a queue", err)
	}

	err = r.ch.ExchangeDeclare(
		ex,
		exKind,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.Fatalf("%s %s: %s", "Failed to declare an exchange", ex, err)
	}

	err = r.ch.QueueBind(
		queue,
		"",
		ex,
		false,
		nil,
	)
	if err != nil {
		logrus.Fatalln("Failed to bind a queue", err)
	}

	logrus.Infof("RabbitMQ queue %s declared\n", queue)

	msgs, err := r.ch.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.Fatalln("Failed to register a consumer", err)
	}

	return msgs
}
