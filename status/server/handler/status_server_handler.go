package handler

import (
	"net/http"

	"sort"

	backup_dto "github.com/bborbe/backup/dto"
	backup_status_checker "github.com/bborbe/backup/status/server/checker"
	error_handler "github.com/bborbe/http_handler/error"
	json_handler "github.com/bborbe/http_handler/json"
	"github.com/golang/glog"
)

type handler struct {
	statusChecker backup_status_checker.StatusChecker
}

func New(statusChecker backup_status_checker.StatusChecker) http.Handler {
	h := new(handler)
	h.statusChecker = statusChecker
	return h
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	status, err := h.statusChecker.Check()
	if err != nil {
		glog.V(2).Infof("check status failed: %v", err)
		e := error_handler.NewMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(responseWriter, request)
		return
	}
	status = filter(status, request.FormValue("status"))
	handler := json_handler.New(status)
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
	sort.Sort(backup_dto.StatusByBackupDate(result))
	return result
}
