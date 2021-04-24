package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "I'm a web server.")
	})

	timeOut := time.Second * 45

	srv := &http.Server{
		Addr:           ":3001",
		Handler:        mux,
		ReadTimeout:    timeOut,
		WriteTimeout:   timeOut,
		IdleTimeout:    timeOut * 2,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf(" listen and serve http server fail:\n %v ", err)
		}
	}()

	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt)
	<-exit
	ctx, cacel := context.WithTimeout(context.Background(), timeOut)
	defer cacel()
	err := srv.Shutdown(ctx)
	log.Println("shutting down now. ", err)
	os.Exit(0)
}
