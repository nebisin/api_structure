package response

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

// ErrorResponse is the default response function for errors. It takes http.ResponseWriter,
// status code as an int and message interface.
func ErrorResponse(w http.ResponseWriter, status int, message interface{}) {
	err := JSONResponse(w, status, Envelope{"error": message})
	if err != nil {
		w.WriteHeader(500)
	}
}

// ServerErrorResponse function is for sending the 500 internal server error to the client.
// ServerErrorResponse logs the request method, request url with the error message.
func ServerErrorResponse(w http.ResponseWriter, r *http.Request, log *logrus.Logger, err error) {
	log.WithFields(map[string]interface{}{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	}).WithError(err).Error("server error response")

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

func RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	ErrorResponse(w, http.StatusTooManyRequests, message)
}

func InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	ErrorResponse(w, http.StatusUnauthorized, message)
}

func InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	ErrorResponse(w, http.StatusUnauthorized, message)
}

func AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	ErrorResponse(w, http.StatusUnauthorized, message)
}

func InactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	ErrorResponse(w, http.StatusForbidden, message)
}

func NotPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your account does not have the necessary permissions to access this resource"
	ErrorResponse(w, http.StatusForbidden, message)
}
