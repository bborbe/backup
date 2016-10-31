package status_client_handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"sort"

	"time"

	backup_dto "github.com/bborbe/backup/dto"
	"github.com/golang/glog"
)

type Download func(url string) (resp *http.Response, err error)

type statusHandler struct {
	address  string
	download Download
}

func NewStatusHandler(download Download, address string) http.Handler {
	s := new(statusHandler)
	s.address = address
	s.download = download
	return s
}

func (s *statusHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	glog.V(2).Info("handle status request")
	err := s.serveHTTP(responseWriter, request)
	if err != nil {
		glog.V(1).Info(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(responseWriter, "%v", err)
		return
	}
	glog.V(2).Info("handle status request completed")
}

func (s *statusHandler) serveHTTP(responseWriter http.ResponseWriter, request *http.Request) error {
	statusList, err := getStatusList(s.download, s.address)
	if err != nil {
		glog.V(1).Infof("get status list from %v failed: %v", s.address, err)
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

func getStatusList(download Download, address string) ([]*backup_dto.Status, error) {
	resp, err := download(address)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("request failed: %s", (content))
	}
	glog.V(4).Infof(string(content))
	var statusList []*backup_dto.Status
	err = json.Unmarshal(content, &statusList)
	if err != nil {
		glog.V(1).Infof("unmarshal jsoni failed: %v", err)
		return nil, err
	}
	return statusList, nil
}
