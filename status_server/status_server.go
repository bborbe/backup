package status_server

import (
	"fmt"
	"net/http"

	backup_service "github.com/bborbe/backup/service"
	backup_status_checker "github.com/bborbe/backup/status_checker"
	backup_status_handler "github.com/bborbe/backup/status_server_handler"
)

func NewServer(port int, rootdir string) *http.Server {
	backupService := backup_service.NewBackupService(rootdir)
	statusChecker := backup_status_checker.NewStatusChecker(backupService)
	handler := backup_status_handler.NewStatusHandler(statusChecker)
	return &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}
}
