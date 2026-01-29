package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/D3rise/dchat/internal/server"
)

// EchoHandler is an http.Handler that copies its request body
// back to the response.
type EchoHandler struct{}

var _ server.Route = &EchoHandler{}

// NewEchoHandler builds a new EchoHandler.
func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (*EchoHandler) Pattern() string {
	return "/echo"
}

// ServeHTTP handles an HTTP request to the /echo endpoint.
func (*EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := io.Copy(w, r.Body); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to handle request:", err)
	}
}
