package main

import (
	"github.com/nebisin/api_structure/internal/app"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	h := app.NewHandler()

	h.Run()

	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)
	h.Logger = logger

	srv := &http.Server{
		Addr:              ":3000",
		Handler:           h.Router,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       time.Minute,
	}

	h.Logger.Println("Server is running on port :3000")
	h.Logger.Fatal(srv.ListenAndServe())
}
