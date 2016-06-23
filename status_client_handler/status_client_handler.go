package status_client_handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"sort"

	backup_dto "github.com/bborbe/backup/dto"
	"github.com/bborbe/log"
)

type Download func(url string) (resp *http.Response, err error)

var logger = log.DefaultLogger

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
	logger.Debug("handle request")
	err := s.serveHTTP(responseWriter, request)
	if err != nil {
		logger.Debug(err)
		responseWriter.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(responseWriter, "%v", err)
		return
	}

}
func (s *statusHandler) serveHTTP(responseWriter http.ResponseWriter, request *http.Request) error {
	statusList, err := getStatusList(s.download, s.address)
	if err != nil {
		return err
	}
	sort.Sort(backup_dto.StatusByName(statusList))
	responseWriter.Header().Set("Content-Type", "text/html")
	fmt.Fprint(responseWriter, "<html><body>")
	fmt.Fprint(responseWriter, "<h1>Backup-Status</h1>")
	fmt.Fprint(responseWriter, "<ul>")
	for _, status := range statusList {
		if status.Status {
			fmt.Fprint(responseWriter, "<li style=\"color:green\">")
			fmt.Fprint(responseWriter, status.Host)
			fmt.Fprint(responseWriter, "</li>")
		} else {
			fmt.Fprint(responseWriter, "<li style=\"color:red\">")
			fmt.Fprint(responseWriter, status.Host)
			fmt.Fprint(responseWriter, " ")
			fmt.Fprint(responseWriter, status.LatestBackup)
			fmt.Fprint(responseWriter, "</li>")
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
	logger.Tracef(string(content))
	var statusList []*backup_dto.Status
	err = json.Unmarshal(content, &statusList)
	if err != nil {
		return nil, err
	}
	return statusList, nil
}
