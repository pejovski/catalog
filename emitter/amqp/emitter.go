package amqp

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"sync"
)

const (
	exProductUpdated      = "product_updated"
	exProductDeleted      = "product_deleted"
	exProductPriceUpdated = "product_price_updated"

	exKind = "fanout"
)

type Emitter interface {
	ProductPriceUpdated(id string, price float32)
	ProductUpdated(id string)
	ProductDeleted(id string)
}

type emitter struct {
	ch   *amqp.Channel
	once sync.Once
}

func NewEmitter(ch *amqp.Channel) Emitter {
	return emitter{ch: ch}
}

func (e emitter) ProductUpdated(id string) {
	e.once.Do(e.declareExchange(exProductUpdated))

	msg := struct {
		Id string `json:"id"`
	}{Id: id}

	b, err := json.Marshal(&msg)
	if err != nil {
		logrus.Errorf("Failed to json marshal product %s; Error: %s", id, err)
		return
	}

	if err = e.ch.Publish(
		exProductUpdated,
		"",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         b,
		}); err != nil {
		logrus.Errorf("Failed publish event %s; Error: %s", string(b), err)
		return
	}

	logrus.Infof("ProductUpdated event for product %s sent. Body: %s", id, string(b))
}

func (e emitter) ProductDeleted(id string) {
	e.once.Do(e.declareExchange(exProductDeleted))

	msg := struct {
		Id string `json:"id"`
	}{Id: id}

	b, err := json.Marshal(&msg)
	if err != nil {
		logrus.Errorf("Failed to json marshal product %s; Error: %s", id, err)
		return
	}

	if err = e.ch.Publish(
		exProductDeleted,
		"",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         b,
		}); err != nil {
		logrus.Errorf("Failed publish event %s; Error: %s", string(b), err)
		return
	}

	logrus.Infof("ProductDeleted event for product %s sent. Body: %s", id, string(b))
}

func (e emitter) ProductPriceUpdated(id string, price float32) {
	e.once.Do(e.declareExchange(exProductPriceUpdated))

	msg := struct {
		Id    string  `json:"id"`
		Price float32 `json:"price"`
	}{Id: id, Price: price}

	b, err := json.Marshal(&msg)
	if err != nil {
		logrus.Errorf("Failed to json marshal product %s; Error: %s", id, err)
		return
	}

	if err = e.ch.Publish(
		exProductPriceUpdated,
		"",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         b,
		}); err != nil {
		logrus.Errorf("Failed publish event %s; Error: %s", string(b), err)
		return
	}

	logrus.Infof("ProductPriceUpdated event for product %s sent. Body: %s", id, string(b))
}

func (e emitter) declareExchange(ex string) func() {
	return func() {
		err := e.ch.ExchangeDeclare(
			ex,
			exKind,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			logrus.Errorf("%s %s: %s", "Failed to declare an exchange", ex, err)
			return
		}
		logrus.Infof("RabbitMq exchange %s declared", ex)
	}
}
