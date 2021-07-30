package response

import (
	"fmt"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, status int, message interface{}) {
	envelope := map[string]interface{}{
		"error": message,
	}
	err := JSONResponse(w, status, envelope)
	if err != nil {
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(w http.ResponseWriter) {
	message := "something went wrong"
	ErrorResponse(w, http.StatusInternalServerError, message)
}

func NotFoundResponse(w http.ResponseWriter) {
	message := "the requested resource could not be found"
	ErrorResponse(w, http.StatusNotFound, message)
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this endpoint", r.Method)
	ErrorResponse(w, http.StatusMethodNotAllowed, message)
}

func BadRequestResponse(w http.ResponseWriter, err error) {
	ErrorResponse(w, http.StatusBadRequest, err.Error())
}