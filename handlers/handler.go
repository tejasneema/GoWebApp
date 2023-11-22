package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Handler struct {
	l *log.Logger
}

func NewHandler(l *log.Logger) *Handler {
	return &Handler{l}
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello world from handleRootCall")

	d, err := ioutil.ReadAll(r.Body)

	if nil != err {
		http.Error(rw, "Ooops something went wrong", http.StatusBadGateway)
		return
	}

	h.l.Printf("Data is %s\n", d)
	fmt.Fprintf(rw, "Hello, received your response successfully, which is: %s\n", d)

}
