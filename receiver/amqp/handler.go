package amqp

import (
	"encoding/json"
	"github.com/pejovski/catalog/domain"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

const rejectSleepTime = 5 * time.Second

type Handler struct {
	controller domain.CatalogController
}

func NewHandler(c domain.CatalogController) *Handler {
	return &Handler{
		controller: c,
	}
}

func (h Handler) RatingUpdated(d *amqp.Delivery) {

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

func (h Handler) reject(d *amqp.Delivery) {
	time.Sleep(rejectSleepTime)
	if err := d.Reject(true); err != nil {
		logrus.Errorln("Failed to reject msg", err)
	}
}

func (h Handler) ack(d *amqp.Delivery) {
	if err := d.Ack(false); err != nil {
		logrus.Errorln("Failed to ack msg", err)
	}
}
