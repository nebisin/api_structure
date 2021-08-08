package main

import (
	"fmt"
	"github.com/nebisin/api_structure/internal/app"
	"net/http"
	"time"
)

func main() {
	h := app.NewServer()

	h.Run()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", h.Config.Port),
		Handler:      h.Router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	h.Logger.WithField("port", srv.Addr).Info("starting the server")
	h.Logger.Fatal(srv.ListenAndServe())
}
