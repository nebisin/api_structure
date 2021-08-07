package response

import (
	"encoding/json"
	"net/http"
)

type Envelope map[string]interface{}

func JSONResponse(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)

	return err
}