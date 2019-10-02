package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"

	"github.com/pejovski/catalog/controller"
	amqpEmitter "github.com/pejovski/catalog/emitter/amqp"
	"github.com/pejovski/catalog/factory"
	"github.com/pejovski/catalog/gateway/reviewing"
	amqpReceiver "github.com/pejovski/catalog/receiver/amqp"
	"github.com/pejovski/catalog/repository"
	httpServer "github.com/pejovski/catalog/server/http"
)

const (
	serverShutdownTimeout = 3 * time.Second
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
	catalogRepository := repository.NewESProductRepository(esClient)
	reviewingGateway := reviewing.NewGateway(retryablehttp.NewClient(), os.Getenv("REVIEWING_API_HOST"))

	catalogController := controller.NewCatalog(catalogRepository, emitter, reviewingGateway)

	amqpHandler := amqpReceiver.NewHandler(catalogController)
	receiver := amqpReceiver.NewReceiver(amqpCh, amqpHandler)
	// receive messages in goroutines
	receiver.Receive()

	serverHandler := httpServer.NewHandler(catalogController)
	serverRouter := httpServer.NewRouter(serverHandler)

	server := factory.CreateHttpServer(serverRouter, fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf(err.Error())
		}
	}()
	logrus.Infof("Server started at port: %s", os.Getenv("APP_PORT"))

	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	// Create channel for shutdown signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Receive shutdown signals.
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorln("Server shutdown failed", err)
	}
	logrus.Println("Server exited properly")
}
