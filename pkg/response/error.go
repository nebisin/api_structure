package response

import (
	"fmt"
	"log"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, status int, message interface{}) {
	err := JSONResponse(w, status, Envelope{"error": message})
	if err != nil {
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(w http.ResponseWriter, log *log.Logger, err error) {
	log.Println(err.Error())
	message := "something went wrong"
	ErrorResponse(w, http.StatusInternalServerError, message)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
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

func FailedValidationResponse(w http.ResponseWriter, errs map[string]string) {
	err := JSONResponse(w, http.StatusUnprocessableEntity, Envelope{"errors": errs})
	if err != nil {
		w.WriteHeader(500)
	}
}

func EditConflictResponse(w http.ResponseWriter) {
	message := "unable to update the record due to an edit conflict, please try again"
	ErrorResponse(w, http.StatusConflict, message)
}