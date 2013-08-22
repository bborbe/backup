package main

import (
	"github.com/bborbe/log"
	"fmt"
	"github.com/bborbe/backup/service"
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
		fmt.Printf("%s => %s", host.GetName(), "2013-06-01 12:23:45")
	}


	logger.Debug("done")
}
