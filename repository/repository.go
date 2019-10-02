package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/pejovski/catalog/domain"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"net/http"
)

const index = "products"

type ESProductRepository struct {
	client *elasticsearch.Client
}

func NewESProductRepository(es *elasticsearch.Client) *ESProductRepository {
	return &ESProductRepository{client: es}
}

func (r ESProductRepository) Get(id string) (*domain.Product, error) {
	var h *Hit

	res, err := r.client.Get(index, id)
	if err != nil {
		logrus.Errorf("Failed to get product %s", id)
		return nil, err
	}

	if res.IsError() {
		if res.StatusCode == http.StatusNotFound {
			return nil, domain.ErrNotFound
		}
		logrus.Errorf("Error in the response for product with id: %s. Status code: %d. Response: %s", id, res.StatusCode, res.String())
		return nil, errors.New("response error")
	}

	if err := json.NewDecoder(res.Body).Decode(&h); err != nil {
		logrus.Errorf("Failed to decode body for product %s", id)
		return nil, err
	}
	defer res.Body.Close()

	return mapHitToProduct(h), nil
}

func (r ESProductRepository) Create(p *domain.Product) (id string, err error) {
	d := mapProductToDocument(p)

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(d); err != nil {
		logrus.Errorf("Failed to decode body for product %s", id)
		return "", err
	}

	id = ksuid.New().String()

	res, err := r.client.Create(index, id, &buf)
	if err != nil {
		logrus.Errorf("Failed to create product %s", id)
		return "", err
	}

	if res.IsError() {
		logrus.Errorf("Error in the response for product with id: %s. Status code: %d. Response: %s", id, res.StatusCode, res.String())
		return "", errors.New("response error")
	}

	return id, nil
}

// Update currently updates only name, brand, price, category and image
func (r ESProductRepository) Update(p *domain.Product) error {
	d := mapProductToDocument(p)
	u := Update{Doc: d}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(u); err != nil {
		logrus.Errorf("Failed to decode body for product %s", p.Id)
		return err
	}

	res, err := r.client.Update(index, p.Id, &buf)
	if err != nil {
		logrus.Errorf("Failed to update product %s", p.Id)
		return err
	}

	if res.IsError() {
		logrus.Errorf("Error in the response for product with id: %s. Status code: %d. Response: %s", p.Id, res.StatusCode, res.String())
		return errors.New("response error")
	}

	return nil
}

func (r ESProductRepository) UpdatePrice(id string, price float32) error {

	up := map[string]map[string]float32{
		"doc": {"price": price},
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(up); err != nil {
		logrus.Errorf("Failed to encode update for product %s", id)
		return err
	}

	res, err := r.client.Update(index, id, &buf)
	if err != nil {
		logrus.Errorf("Failed to update product %s", id)
		return err
	}

	if res.IsError() {
		logrus.Errorf("Error in the response for product with id: %s. Status code: %d. Response: %s", id, res.StatusCode, res.String())
		return errors.New("response error")
	}

	return nil
}

func (r ESProductRepository) Delete(id string) error {
	res, err := r.client.Delete(index, id)
	if err != nil {
		logrus.Errorf("Failed to delete product %s", id)
		return err
	}

	if res.IsError() {
		logrus.Errorf("Error in the response for product with id: %s. Status code: %d. Response: %s", id, res.StatusCode, res.String())
		return errors.New("response error")
	}

	return nil
}

func (r ESProductRepository) GetByCategory(category string) ([]*domain.Product, error) {

	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"category": category,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		logrus.Errorf("Failed to encode query for category %s", category)
		return nil, err
	}

	// Perform the search request.
	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex(index),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
		r.client.Search.WithPretty(),
	)
	if err != nil {
		logrus.Errorf("Failed to get response for category %s", category)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		logrus.Errorf("Error in the response for category %s. Status code: %d. Response: %s", category, res.StatusCode, res.String())
		return nil, errors.New("response error")
	}

	var result *Result

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		logrus.Errorf("Failed to decode result for category %s", category)
		return nil, err
	}

	products := []*domain.Product{}

	for _, hit := range result.Hits.Hits {
		products = append(products, mapHitToProduct(&hit))
	}

	return products, nil

}

func (r ESProductRepository) UpdateRating(id string, rating *domain.Rating) error {
	// ToDo
	return errors.New("not implemented yet")
}
