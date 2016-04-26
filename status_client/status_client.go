package status_client

import (
	"net/http"

	"fmt"

	backup_status_handler "github.com/bborbe/backup/status_client_handler"
)

func NewServer(download func(url string) (resp *http.Response, err error), port int, address string) *http.Server {
	handler := backup_status_handler.NewStatusHandler(download, address)
	return &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}
}
