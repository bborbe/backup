package status_server_handler

import (
	"net/http"

	backup_dto "github.com/bborbe/backup/dto"
	backup_status_checker "github.com/bborbe/backup/status_checker"
	"github.com/bborbe/log"
	error_handler "github.com/bborbe/server/handler/error"
	json_handler "github.com/bborbe/server/handler/json"
)

var logger = log.DefaultLogger

type statusHandler struct {
	statusChecker backup_status_checker.StatusChecker
}

func NewStatusHandler(statusChecker backup_status_checker.StatusChecker) http.Handler {
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
	handler := json_handler.NewJsonHandler(status)
	handler.ServeHTTP(responseWriter, request)
}

func filter(list []*backup_dto.Status, status string) []*backup_dto.Status {
	if list == nil {
		return list
	}
	result := make([]*backup_dto.Status, 0)
	for _, s := range list {
		if "true" == status {
			if s.Status {
				result = append(result, s)
			}
		} else if "false" == status {
			if !s.Status {
				result = append(result, s)
			}
		} else {
			result = append(result, s)
		}
	}
	return result
}
