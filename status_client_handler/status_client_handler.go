package status_client_handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
	resp, err := s.download(s.address)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("request failed: %s", (content))
	}

	logger.Debugf(string(content))

	var statusList []*backup_dto.Status
	err = json.Unmarshal(content, &statusList)
	if err != nil {
		return err
	}
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
