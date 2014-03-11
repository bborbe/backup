package status_handler

import (
	"net/http"

	"github.com/bborbe/backup/status_checker"
	"github.com/bborbe/log"
	"github.com/bborbe/server/handler/error"
	"github.com/bborbe/server/handler/json"
)

var logger = log.DefaultLogger

type statusHandler struct {
	statusChecker status_checker.StatusChecker
}

func NewStatusHandler(statusChecker status_checker.StatusChecker) http.Handler {
	s := new(statusHandler)
	s.statusChecker = statusChecker
	return s
}

func (s *statusHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	status, err := s.statusChecker.Check()
	if err != nil {
		logger.Debugf("check status failed: %v", err)
		e := error.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(responseWriter, request)
		return
	}
	handler := json.NewJsonHandler(status)
	handler.ServeHTTP(responseWriter, request)
}
