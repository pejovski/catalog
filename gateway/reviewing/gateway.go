package reviewing

import (
	"errors"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pejovski/catalog/domain"
)

type Gateway struct {
	client *retryablehttp.Client
	host   string
}

func NewGateway(c *retryablehttp.Client, host string) Gateway {
	return Gateway{client: c, host: host}
}

func (g Gateway) Rating(productId string) (*domain.Rating, error) {

	// ToDo
	return nil, errors.New("method not implemented yet")
}
