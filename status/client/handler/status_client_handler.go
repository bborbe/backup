package status_client_handler

import (
	"fmt"
	"net/http"

	"sort"

	"time"

	backup_dto "github.com/bborbe/backup/dto"
	"github.com/golang/glog"
)

type statusList func() ([]backup_dto.Status, error)

type handler struct {
	statusList statusList
}

func New(statusList statusList) http.Handler {
	h := new(handler)
	h.statusList = statusList
	return h
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	glog.V(2).Info("handle status request")
	err := h.serveHTTP(responseWriter, request)
	if err != nil {
		glog.V(1).Info(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(responseWriter, "%v", err)
		return
	}
	glog.V(2).Info("handle status request completed")
}

func (h *handler) serveHTTP(responseWriter http.ResponseWriter, request *http.Request) error {
	statusList, err := h.statusList()
	if err != nil {
		glog.V(1).Infof("get status list failed: %v", err)
		return err
	}
	sort.Sort(backup_dto.StatusByDate(statusList))
	responseWriter.Header().Set("Content-Type", "text/html")
	fmt.Fprint(responseWriter, "<html><body>")
	fmt.Fprint(responseWriter, "<h1>Backup-Status</h1>")
	fmt.Fprint(responseWriter, "<ul>")
	now := time.Now()
	for _, status := range statusList {
		if status.Status {
			fmt.Fprint(responseWriter, "<li>")
			fmt.Fprint(responseWriter, "<span style=\"color:green\">")
			fmt.Fprint(responseWriter, status.Host)
			fmt.Fprint(responseWriter, "</span> (")
			fmt.Fprint(responseWriter, status.LatestBackup.Age(now))
			fmt.Fprint(responseWriter, ")</li>")
		} else {
			fmt.Fprint(responseWriter, "<li>")
			fmt.Fprint(responseWriter, "<span style=\"color:red\">")
			fmt.Fprint(responseWriter, status.Host)
			fmt.Fprint(responseWriter, "</span> (")
			fmt.Fprint(responseWriter, status.LatestBackup)
			fmt.Fprint(responseWriter, ")</li>")
		}
	}
	fmt.Fprint(responseWriter, "</ul>")
	fmt.Fprint(responseWriter, "</body></html>")
	return nil
}
