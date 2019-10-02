package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pejovski/catalog/domain"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	controller domain.CatalogController
}

func NewHandler(c domain.CatalogController) *Handler {
	return &Handler{
		controller: c,
	}
}

func (h Handler) Products() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category := r.FormValue("category")
		if category == "" {
			logrus.Warnln("Category not found")
			http.Error(w, "Category not found", http.StatusBadRequest)
			return
		}

		dps, err := h.controller.GetProducts(category)
		if err != nil {
			logrus.Errorf("Failed to get products for category %s. Error: %s", category, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		h.respond(w, r, mapDomainProductsToProducts(dps), http.StatusOK)
	}
}

func (h Handler) Product() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		id := params["id"]
		if id == "" {
			logrus.Warnln("Product id not found")
			http.Error(w, "Product id not found", http.StatusBadRequest)
			return
		}

		p, err := h.controller.GetProduct(id)
		if err != nil {
			if err == domain.ErrNotFound {
				http.Error(w, "Product not found", http.StatusNotFound)
				return
			}
			logrus.Errorf("Failed to get product with id %s. Error: %s", id, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		h.respond(w, r, mapDomainProductToProduct(p), http.StatusOK)
	}
}

func (h Handler) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var p *domain.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			logrus.Warnln("Failed to decode request body")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		id, err := h.controller.CreateProduct(p)
		if err != nil {
			logrus.Errorf("Failed to create product for id %s. Error: %s", p.Id, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", fmt.Sprintf("/products/%s", id))
		h.respond(w, r, nil, http.StatusCreated)
	}
}

func (h Handler) UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var p *domain.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			logrus.Warnln("Failed to decode request body")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		params := mux.Vars(r)
		id := params["id"]
		if id == "" {
			logrus.Warnln("Product id not found")
			http.Error(w, "Product id not found", http.StatusBadRequest)
			return
		}

		p.Id = id

		if err := h.controller.UpdateProduct(p); err != nil {
			logrus.Errorf("Failed to update product for id %s. Error: %s", p.Id, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		h.respond(w, r, nil, http.StatusNoContent)
	}
}

func (h Handler) UpdateProductPrice() http.HandlerFunc {

	var request struct {
		Price float32 `json:"price"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			logrus.Warnln("Failed to decode request body")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		params := mux.Vars(r)
		id := params["id"]
		if id == "" {
			logrus.Warnln("Product id not found")
			http.Error(w, "Product id not found", http.StatusBadRequest)
			return
		}

		if err := h.controller.UpdateProductPrice(id, request.Price); err != nil {
			logrus.Errorf("Failed to update product price for product %s. Error: %s", id, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		h.respond(w, r, nil, http.StatusNoContent)
	}
}

func (h Handler) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id := params["id"]
		if id == "" {
			logrus.Warnln("Product id not found")
			http.Error(w, "Product id not found", http.StatusBadRequest)
			return
		}

		if err := h.controller.DeleteProduct(id); err != nil {
			logrus.Errorf("Failed to delete product %s. Error: %s", id, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		h.respond(w, r, nil, http.StatusNoContent)
	}
}

func (h Handler) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logrus.Errorf("Failed to encode data. Error: %s", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func (h Handler) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
