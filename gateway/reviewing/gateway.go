package reviewing

import (
	"errors"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pejovski/catalog/model"
)

type Gateway interface {
	Rating(productId string) (*model.Rating, error)
}

type gateway struct {
	client *retryablehttp.Client
	host   string
}

func NewGateway(c *retryablehttp.Client, host string) Gateway {
	return gateway{client: c, host: host}
}

func (g gateway) Rating(productId string) (*model.Rating, error) {

	// ToDo
	return nil, errors.New("method not implemented yet")
}
