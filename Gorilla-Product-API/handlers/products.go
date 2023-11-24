package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tejasneema/GoWebApp/Gorilla-Product-API/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l: l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()

	rw.Header().Add("content-type", "application/json")
	err := lp.ToJSON(rw)
	if nil != err {
		http.Error(rw, "Unable to marshal the json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Post Products")

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	p.l.Printf("Product: %#v", &prod)

	data.AddProduct(&prod)
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {

	vers := mux.Vars(r)
	id, err := strconv.Atoi(vers["id"])
	if nil != err {
		http.Error(rw, "unable to conver int to string", http.StatusInternalServerError)
	}

	p.l.Println("Got ID", id)
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Unable to update product", http.StatusNotFound)
		return
	}

	if nil != err {
		http.Error(rw, "Some internal server error", http.StatusInternalServerError)
	}
	return
}

type KeyProduct struct{}

func (p *Products) MiddlewareProductJSONConversion(next http.Handler) http.Handler {
	p.l.Println("Inside MiddlewareProductJSONConversion")
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if nil != err {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		// validate the product
		error := prod.Validate()
		if nil != error {
			p.l.Println("[ERROR] validating product", err)
			http.Error(rw, fmt.Sprintf("validation failed: %s", error), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		p.l.Println("conversion successful")
		next.ServeHTTP(rw, r)
	})
}
