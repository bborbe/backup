package status_client

import (
	backup_status_handler "github.com/bborbe/backup/status_client_handler"
	"github.com/bborbe/server"
)

func NewServer(port int, address string) server.Server {
	handler := backup_status_handler.NewStatusHandler(address)
	return server.NewServerPort(port, handler)
}
