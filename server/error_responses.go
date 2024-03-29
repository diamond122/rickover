// Helpers for building various types of error responses.

package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kevinburke/rest"
)

func new405(r *http.Request) *rest.Error {
	return &rest.Error{
		Title:    "Method not allowed",
		ID:       "method_not_allowed",
		Instance: r.URL.Path,
		Status:   405,
	}
}

func new404(r *http.Request) *rest.Error {
	return &rest.Error{
		Title:    "Resource not found",
		ID:       "not_found",
		Instance: r.URL.Path,
		Status:   404,
	}
}

func insecure403(r *http.Request) *rest.Error {
	return &rest.Error{
		Title:    "Server not available over HTTP",
		ID:       "insecure_request",
		Detail:   "For your security, please use an encrypted connection",
		Instance: r.URL.Path,
		Status:   403,
	}
}

func new401(r *http.Request) *rest.Error {
	return &rest.Error{
		Title:    "Unauthorized. Please include your API credentials",
		ID:       "unauthorized",
		Instance: r.URL.Path,
		Status:   401,
	}
}

// createEmptyErr returns a rest.Error indicating the request omits a required
// field.
func createEmptyErr(field string, path string) *rest.Error {
	return &rest.Error{
		Title:    fmt.Sprintf("Missing required field: %s", field),
		Detail:   fmt.Sprintf("Please include a %s in the request body", field),
		ID:       "missing_parameter",
		Instance: path,
	}
}

func createPositiveIntErr(field string, path string) *rest.Error {
	return &rest.Error{
		Title:    fmt.Sprintf("%s must be set to a number greater than zero", field),
		ID:       "invalid_parameter",
		Instance: path,
	}
}

func notFound(w http.ResponseWriter, err *rest.Error) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(err)
}

func authenticate(w http.ResponseWriter, err *rest.Error) {
	w.Header().Set("WWW-Authenticate", "Basic realm=\"rickover\"")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(err)
}

func forbidden(w http.ResponseWriter, err *rest.Error) {
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(err)
}

var serverError = rest.Error{
	Status: http.StatusInternalServerError,
	ID:     "server_error",
	Title:  "Unexpected server error. Please try again",
}

func tooLarge(w http.ResponseWriter) {
	resp := &rest.Error{
		ID:    "entity_too_large",
		Title: "Data parameter is too large (100KB max)",
	}
	data, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusRequestEntityTooLarge)
	w.Write(data)
}
