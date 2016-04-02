package status_client_handler

import (
	"net/http"

	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type statusHandler struct {
	address string
}

func NewStatusHandler(address string) http.Handler {
	s := new(statusHandler)
	s.address = address
	return s
}

func (s *statusHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	logger.Debug("handle request")
}

