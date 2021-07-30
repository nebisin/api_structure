package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func ReadJSON(r *http.Request, dst interface{}) error {
	// decode the request body
	err := json.NewDecoder(r.Body).Decode(dst)
	// if there is an error during decoding:
	if err != nil {
		// There is a syntax problem with the JSON
		var syntaxError *json.SyntaxError
		// JSON value is not appropriate for the destination Go type
		var unmarshalTypeError *json.UnmarshalTypeError
		// The decode destination is not valid
		// usually because it is not a pointer.
		// This is a problem with our application.
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// check if error has the type json.SyntaxError
		// if it does, return a plain error message
		// which includes the location of problem.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// this occurs when the wrong type for the target destination
		case errors.As(err, &unmarshalTypeError):
			// If the error relates to a specific field include that in the message
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// If the request body is empty
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// this will happen if we pass a non-nil pointer to Decode()
		// we catch this and panic instead of returning the error
		// bacause this error should not happen during normal operation
		// and probably the result of a developer mistake
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		// For anything else return the error message as it is
		default:
			return err
		}

	}

	return nil
}
