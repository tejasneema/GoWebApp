package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/tejasneema/GoWebApp/Product-API/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l: l}
}

// json.marshal will utilize the application storage to keep the data and then write to the responsewriter
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
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		p.updateProduct(rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)

}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()

	rw.Header().Add("content-type", "application/json")
	err := lp.ToJSON(rw)
	if nil != err {
		http.Error(rw, "Unable to marshal the json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Post Products")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)

	if nil != err {
		http.Error(rw, "Unable to unMarshal the JSON", http.StatusBadRequest)
		return
	}

	p.l.Printf("Product: %#v", prod)

	data.AddProduct(prod)
}

func (p *Products) updateProduct(rw http.ResponseWriter, r *http.Request) {
	rex := regexp.MustCompile(`/([0-9]+)`)
	g := rex.FindAllStringSubmatch(r.URL.Path, -1)

	if len(g) != 1 {
		http.Error(rw, "Invalid URI 1", http.StatusBadRequest)
		return
	}

	idString := g[0][1]

	if len(g[0]) != 2 {
		http.Error(rw, "Invalid URI 2", http.StatusBadRequest)
		return

	}

	id, _ := strconv.Atoi(idString)

	p.l.Println("Got ID", id)

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)

	if nil != err {
		http.Error(rw, "Unable to unMarshal the JSON", http.StatusBadRequest)
		return
	}
	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Unable to update product", http.StatusNotFound)
		return
	}

	if nil != err {
		http.Error(rw, "Some internal server error", http.StatusInternalServerError)
	}
	return
}
