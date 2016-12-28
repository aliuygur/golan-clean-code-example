package handlers

import (
	"app"
	"net/http"

	"strings"

	"app/usecases"

	"github.com/alioygur/gores"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

type catalogService interface {
	OneActiveProduct(id interface{}) (*app.Product, error)
	FindActiveProducts(*app.DBFilter) ([]app.Product, error)
	FindActiveProductsByCategory([]interface{}, *app.DBFilter) ([]app.Product, error)
	FindActiveCategories(*app.DBFilter) ([]app.Category, error)
	CreateProduct(*usecases.ProductForm) (*app.Product, error)
	DeleteProduct(id interface{}) error
	UpdateProduct(*usecases.ProductForm) (*app.Product, error)
}

func NewCatalog(srv catalogService, eh app.ErrorHandler) *Catalog {
	return &Catalog{srv, eh}
}

type Catalog struct {
	srv catalogService
	eh  app.ErrorHandler
}

func (ch *Catalog) SetRoutes(r *mux.Router, mid ...alice.Constructor) {
	h := alice.New(mid...)
	r.Handle("/v1/products/{id}", h.ThenFunc(ch.getProduct)).Methods("GET")
	r.Handle("/v1/products", h.ThenFunc(ch.getProducts)).Methods("GET")
	r.Handle("/v1/categories", h.ThenFunc(ch.getCategories)).Methods("GET")

	r.Handle("/v1/admin/products", h.ThenFunc(ch.createProduct)).Methods("POST")
	r.Handle("/v1/admin/products/{id}", h.ThenFunc(ch.updateProduct)).Methods("PATCH", "PUT")
	r.Handle("/v1/admin/products/{id}", h.ThenFunc(ch.deleteProduct)).Methods("DELETE")
}

func (ch *Catalog) createProduct(w http.ResponseWriter, r *http.Request) {
	f := new(usecases.ProductForm)
	if err := decodeR(r, f); err != nil {
		ch.eh.Handle(w, err)
		return
	}

	p, err := ch.srv.CreateProduct(f)
	if err != nil {
		ch.eh.Handle(w, err)
		return
	}

	gores.JSON(w, 201, response{p})
}

func (ch *Catalog) getProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	p, err := ch.srv.OneActiveProduct(params["id"])
	if err != nil {
		ch.eh.Handle(w, err)
		return
	}

	gores.JSON(w, 200, response{p})
}

func (ch *Catalog) getProducts(w http.ResponseWriter, r *http.Request) {
	cids := qCategoryParam(r)
	ps, err := func() ([]app.Product, error) {
		if len(cids) > 0 {
			return ch.srv.FindActiveProductsByCategory(cids, nil)
		}
		return ch.srv.FindActiveProducts(nil)
	}()
	if err != nil {
		ch.eh.Handle(w, err)
		return
	}
	gores.JSON(w, 200, response{ps})
}

func (ch *Catalog) updateProduct(w http.ResponseWriter, r *http.Request) {
	f := new(usecases.ProductForm)
	if err := decodeR(r, f); err != nil {
		ch.eh.Handle(w, err)
		return
	}

	f.ID = 54

	p, err := ch.srv.UpdateProduct(f)
	if err != nil {
		ch.eh.Handle(w, err)
		return
	}

	gores.JSON(w, 200, response{p})
}

func (ch *Catalog) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := ch.srv.DeleteProduct(id); err != nil {
		ch.eh.Handle(w, err)
		return
	}

	gores.NoContent(w)
}

func (ch *Catalog) getCategories(w http.ResponseWriter, r *http.Request) {
	cs, err := ch.srv.FindActiveCategories(nil)
	if err != nil {
		ch.eh.Handle(w, err)
		return
	}

	gores.JSON(w, 200, response{cs})
}

// qCategoryParam gets category param like 1,2,3 as []interface{1, 2, 3}
func qCategoryParam(r *http.Request) []interface{} {
	var cs []interface{}
	c := qParam("category", r)
	if c != "" {
		for _, id := range strings.Split(c, ",") {
			if id == "" {
				continue
			}
			cs = append(cs, strings.Trim(id, " "))
		}
	}
	return cs
}
