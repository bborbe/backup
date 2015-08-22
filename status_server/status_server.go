package status_server

import (
	backup_service "github.com/bborbe/backup/service"
	backup_status_checker "github.com/bborbe/backup/status_checker"
	backup_status_handler "github.com/bborbe/backup/status_handler"
	"github.com/bborbe/server"
)

func NewServer(port int, rootdir string) server.Server {
	backupService := backup_service.NewBackupService(rootdir)
	statusChecker := backup_status_checker.NewStatusChecker(backupService)
	handler := backup_status_handler.NewStatusHandler(statusChecker)
	return server.NewServerPort(port, handler)
}
