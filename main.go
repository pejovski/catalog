package main

import (
	"fmt"
	"github.com/pejovski/catalog/pkg/signals"
	"github.com/pejovski/catalog/repository/es"
	"github.com/pejovski/catalog/server/api"
	"os"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"

	"github.com/pejovski/catalog/controller"
	amqpEmitter "github.com/pejovski/catalog/emitter/amqp"
	"github.com/pejovski/catalog/factory"
	"github.com/pejovski/catalog/gateway/reviewing"
	amqpReceiver "github.com/pejovski/catalog/receiver/amqp"
)

const (
	shutdownDuration = 3 * time.Second
)

func main() {
	esClient := factory.CreateESClient(fmt.Sprintf(
		"http://%s:%s",
		os.Getenv("ES_HOST"),
		os.Getenv("ES_PORT"),
	))
	amqpCh := factory.CreateAmqpChannel(fmt.Sprintf(
		"amqp://%s:%s@%s:%s/%s",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
		os.Getenv("RABBITMQ_VHOST"),
	))

	emitter := amqpEmitter.NewEmitter(amqpCh)
	catalogRepository := es.NewRepository(esClient)
	reviewingGateway := reviewing.NewGateway(retryablehttp.NewClient(), os.Getenv("REVIEWING_API_HOST"))

	catalogController := controller.New(catalogRepository, emitter, reviewingGateway)

	amqpHandler := amqpReceiver.NewHandler(catalogController)
	receiver := amqpReceiver.NewReceiver(amqpCh, amqpHandler)
	// receive messages in goroutines
	receiver.Receive()

	ctx := signals.Context()

	serverAPI := api.NewServer(catalogController)
	serverAPI.Run(ctx)

	logrus.Infof("allowing %s for graceful shutdown to complete", shutdownDuration)
	<-time.After(shutdownDuration)
}
