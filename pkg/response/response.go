package response

import (
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = w.Write(js)

	return err
}