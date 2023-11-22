package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/tejasneema/GoWebApp/HttpServerCore/handlers"
)

func main() {
	l := log.New(os.Stdout, "Go-APIs: ", log.LstdFlags)
	hh := handlers.NewHandler(l)
	gh := handlers.NewGoodbye(l)

	l.Println("Starting the server")

	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/bye", gh)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if nil != err {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate signal, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
