package api

import (
	"github.com/gorilla/mux"
	_ "github.com/pejovski/catalog/app/statik"
	"github.com/pejovski/catalog/controller"
	"github.com/rakyll/statik/fs"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Router interface {
	routes()
	swagger()
	health()

	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type router struct {
	router  *mux.Router
	handler Handler
}

func newRouter(c controller.Controller) Router {
	s := &router{
		router:  mux.NewRouter(),
		handler: newHandler(c),
	}

	s.health()
	s.swagger()
	s.routes()

	return s
}

func (rtr *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rtr.router.ServeHTTP(w, r)
}

func (rtr *router) routes() {
	rtr.router.Path("/products").Queries("category", "{category}").Methods("GET").HandlerFunc(rtr.handler.Products()).Name("products")
	rtr.router.HandleFunc("/products", rtr.handler.CreateProduct()).Methods("POST")
	rtr.router.HandleFunc("/products/{id}", rtr.handler.Product()).Methods("GET")
	rtr.router.HandleFunc("/products/{id}", rtr.handler.UpdateProduct()).Methods("PUT")
	rtr.router.HandleFunc("/products/{id}", rtr.handler.UpdateProductPrice()).Methods("PATCH")
	rtr.router.HandleFunc("/products/{id}", rtr.handler.DeleteProduct()).Methods("DELETE")
}

func (rtr *router) swagger() {
	// swagger handler
	statikFS, err := fs.New()
	if err != nil {
		logrus.Fatalf("%s: %s", "Failed to find statik", err)
	}
	sh := http.FileServer(statikFS)

	rtr.router.Handle("/", sh).Methods("GET")
	rtr.router.PathPrefix("/swagger").Handler(sh)
}

func (rtr *router) health() {
	rtr.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Up"))
	}).Methods("GET")
}
