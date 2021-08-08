package main

import (
	"github.com/nebisin/api_structure/internal/app"
)

func main() {
	h := app.NewServer()

	h.Run()
}
