package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nebisin/api_structure/internal/store"
	"github.com/nebisin/api_structure/pkg/response"
	"net/http"
	"strconv"
	"time"
)

func (h *handler) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "create a post\n")
}

func (h *handler) handleShowPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		h.Logger.Println(err)
		response.NotFoundResponse(w)
		return
	}

	post := store.Post{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Developing Rest API with Golang",
		Body:      "Lorem ipsum dolor sit amet.",
		Tags:      []string{"golang", "rest", "api"},
		Version:   1,
	}

	if err := response.JSONResponse(w, http.StatusOK, post); err != nil {
		h.Logger.Println(err)
		response.ServerErrorResponse(w)
	}
}
