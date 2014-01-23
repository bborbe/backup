package status_handler

import (
	"github.com/bborbe/backup/status_checker"
	"github.com/bborbe/server/handler/error"
	"github.com/bborbe/server/handler/json"
	"net/http"
)

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
		e := error.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(responseWriter, request)
		return
	}
	handler := json.NewJsonHandler(status)
	handler.ServeHTTP(responseWriter, request)
}
