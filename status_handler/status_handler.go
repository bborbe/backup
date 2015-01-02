package status_handler

import (
	"net/http"

	"github.com/bborbe/backup/dto"
	"github.com/bborbe/backup/status_checker"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
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
		e := error_handler.NewErrorMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(responseWriter, request)
		return
	}
	status = filter(status, request.FormValue("status"))
	handler := json.NewJsonHandler(status)
	handler.ServeHTTP(responseWriter, request)
}

func filter(list []dto.Status, status string) []dto.Status {
	if list == nil {
		return list
	}
	result := make([]dto.Status, 0)
	for _, s := range list {
		if "true" == status {
			if s.GetStatus() {
				result = append(result, s)
			}
		} else if "false" == status {
			if !s.GetStatus() {
				result = append(result, s)
			}
		} else {
			result = append(result, s)
		}
	}
	return result
}
