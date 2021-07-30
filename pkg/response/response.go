package response

import (
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, status int, data interface{}) error {
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	return err
}