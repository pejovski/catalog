package factory

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
	"net/http"
)

const maxIdleConnsPerHost = 10

func CreateESClient(url string) *elasticsearch.Client {
	esConfig := elasticsearch.Config{
		Addresses: []string{
			url,
		},
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConnsPerHost,
		},
	}

	es, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		logrus.Fatalf("Error creating the client: %s", err)
	}

	res, err := es.Ping()
	if err != nil {
		logrus.Fatalf("Error pinging elastic server: %s", err)
	}
	if res.IsError() {
		logrus.Fatalf("Error pinging elastic server")
	}

	return es
}
