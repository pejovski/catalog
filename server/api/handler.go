package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/pejovski/catalog/controller"
	myerr "github.com/pejovski/catalog/error"
	"github.com/pejovski/catalog/model"
)

type Handler interface {
	Products() http.HandlerFunc
	Product() http.HandlerFunc
	CreateProduct() http.HandlerFunc
	UpdateProduct() http.HandlerFunc
	UpdateProductPrice() http.HandlerFunc
	DeleteProduct() http.HandlerFunc
}

type handler struct {
	controller controller.Controller
	mapper     Mapper
}

func newHandler(c controller.Controller) Handler {
	return handler{
		controller: c,
		mapper:     newMapper(),
	}
}

func (h handler) Products() http.HandlerFunc {
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

		h.respond(w, r, h.mapper.mapDomainProductsToProducts(dps), http.StatusOK)
	}
}

func (h handler) Product() http.HandlerFunc {
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
			if err == myerr.ErrNotFound {
				http.Error(w, "Product not found", http.StatusNotFound)
				return
			}
			logrus.Errorf("Failed to get product with id %s. Error: %s", id, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		h.respond(w, r, h.mapper.mapDomainProductToProduct(p), http.StatusOK)
	}
}

func (h handler) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var p *model.Product
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

func (h handler) UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var p *model.Product
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

func (h handler) UpdateProductPrice() http.HandlerFunc {

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

func (h handler) DeleteProduct() http.HandlerFunc {
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

func (h handler) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
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

func (h handler) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
