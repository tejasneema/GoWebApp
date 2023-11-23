package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/tejasneema/GoWebApp/Product-API/handlers"
)

func main() {
	l := log.New(os.Stdout, "Go-APIs: ", log.LstdFlags)
	ph := handlers.NewProducts(l)

	sm := http.NewServeMux()
	sm.Handle("/", ph)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      sm,
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		l.Println("Starting the server on port 8080")

		err := s.ListenAndServe()
		if nil != err {
			l.Fatal(err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate signal, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
