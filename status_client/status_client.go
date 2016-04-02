package status_client

import (
	"net/http"

	backup_status_handler "github.com/bborbe/backup/status_client_handler"
	"github.com/bborbe/server"
)

func NewServer(download func(url string) (resp *http.Response, err error), port int, address string) server.Server {
	handler := backup_status_handler.NewStatusHandler(download, address)
	return server.NewServerPort(port, handler)
}
