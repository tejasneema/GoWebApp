package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tejasneema/GoWebApp/Product-API/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l: l}
}

//json.marshal will utilize the application storage to keep the data and then write to the responsewriter
func (p *Products) ServeHTTPMarshaller(rw http.ResponseWriter, r *http.Request) {

	lp := data.GetProducts()

	d, err := json.Marshal(lp)
	if nil != err {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}

	rw.Header().Add("content-type", "application/json")
	rw.Write(d)
}

// use encoder that writes directly to the stream rather then keeping something in the buffer
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
	} else if r.Method == http.MethodPost {
		rw.Write([]byte("Yet to implement post calls"))
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()

	rw.Header().Add("content-type", "application/json")
	err := lp.ToJSON(rw)
	if nil != err {
		http.Error(rw, "Unable to marshal the json", http.StatusInternalServerError)
	}
}
