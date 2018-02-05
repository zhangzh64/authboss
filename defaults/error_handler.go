package defaults

import (
	"fmt"
	"net/http"

	"github.com/volatiletech/authboss"
)

// ErrorHandler wraps http handlers with errors with itself
// to provide error handling.
//
// The pieces provided to this struct must be thread-safe
// since they will be handed to many pointers to themselves.
type ErrorHandler struct {
	LogWriter authboss.Logger
}

// Wrap an http handler with an error
func (e ErrorHandler) Wrap(handler func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return errorHandler{
		Handler:   handler,
		LogWriter: e.LogWriter,
	}
}

type errorHandler struct {
	Handler   func(w http.ResponseWriter, r *http.Request) error
	LogWriter authboss.Logger
}

// ServeHTTP handles errors
func (e errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := e.Handler(w, r)
	if err == nil {
		return
	}

	e.LogWriter.Error(fmt.Sprintf("error at %s: %+v", r.URL.String(), err))
}