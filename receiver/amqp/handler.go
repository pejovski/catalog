package amqp

import (
	"encoding/json"
	"github.com/pejovski/catalog/controller"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

const rejectSleepTime = 5 * time.Second

type Handler interface {
	RatingUpdated(d *amqp.Delivery)
}

type handler struct {
	controller controller.Controller
}

func NewHandler(c controller.Controller) Handler {
	return handler{
		controller: c,
	}
}

func (h handler) RatingUpdated(d *amqp.Delivery) {

	msq := struct {
		ProductId string `json:"product_id"`
	}{}

	err := json.Unmarshal(d.Body, &msq)
	if err != nil {
		logrus.Errorln("Failed to read body", err)
		h.reject(d)
		return
	}

	err = h.controller.UpdateRating(msq.ProductId)
	if err != nil {
		logrus.Errorln("Failed to update product", err)
		h.reject(d)
		return
	}

	h.ack(d)
}

func (h handler) reject(d *amqp.Delivery) {
	time.Sleep(rejectSleepTime)
	if err := d.Reject(true); err != nil {
		logrus.Errorln("Failed to reject msg", err)
	}
}

func (h handler) ack(d *amqp.Delivery) {
	if err := d.Ack(false); err != nil {
		logrus.Errorln("Failed to ack msg", err)
	}
}
