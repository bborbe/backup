package main

import (
	"fmt"
	"github.com/bborbe/backup/service"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

func main() {
	logger.Debug("start")

	backupRootDir := "/rsync"
	backupService := service.NewBackupService(backupRootDir)
	hosts, err := backupService.ListHosts()
	if err != nil {
		logger.Fatal(err)
		return
	}

	for _, host := range hosts {
		backup, err := backupService.GetLatestBackup(host)
		if err != nil {
			logger.Fatal(err)
			return
		}
		fmt.Printf("%s => %s", host.GetName(), backup.GetName())
	}

	logger.Debug("done")
}
