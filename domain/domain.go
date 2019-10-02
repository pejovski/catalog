package domain

import (
	"errors"
	"github.com/streadway/amqp"
	"net/http"
)

var (
	ErrNotFound = errors.New("not found")
)

type Product struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Brand    string  `json:"brand"`
	Price    float32 `json:"price"`
	Category string  `json:"category"`
	Image    string  `json:"image"`
	Rating
}

type Rating struct {
	// out of 5 (e.g. 3.9)
	Stars float32 `json:"rating"`
	// number of customers who reviewed the product
	Customers int `json:"customers"`
}

type Receiver interface {
	Receive()
}

type AmqpHandler interface {
	RatingUpdated(d *amqp.Delivery)
}

type HttpHandler interface {
	Products() http.HandlerFunc
	Product() http.HandlerFunc
	CreateProduct() http.HandlerFunc
	UpdateProduct() http.HandlerFunc
	UpdateProductPrice() http.HandlerFunc
	DeleteProduct() http.HandlerFunc
}

type CatalogController interface {
	GetProduct(id string) (*Product, error)
	GetProducts(category string) ([]*Product, error)
	CreateProduct(p *Product) (id string, err error)
	UpdateProduct(p *Product) error
	DeleteProduct(id string) error
	UpdateProductPrice(id string, price float32) error
	UpdateRating(id string) error
}

type CatalogRepository interface {
	Get(id string) (*Product, error)
	Create(p *Product) (id string, err error)
	Update(p *Product) error
	Delete(id string) error
	GetByCategory(category string) ([]*Product, error)
	UpdatePrice(id string, price float32) error
	UpdateRating(id string, r *Rating) error
}

type AmqpEmitter interface {
	ProductPriceUpdated(id string, price float32)
	ProductUpdated(id string)
	ProductDeleted(id string)
}

type ReviewingGateway interface {
	Rating(productId string) (*Rating, error)
}
