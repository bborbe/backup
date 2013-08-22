package main

import (
	"fmt"
	"os"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

func main() {
	logger.Debug("start")

	backups, err := getBackups("/rsync")
	if err != nil {
		logger.Errorf("read backups in /rsync failed, %v", err)
		return
	}
	for _, backup := range backups {
		fmt.Println(backup)
	}
	logger.Debug("done")
}

func getBackups(backupRootDir string) ([]string, error) {
	file, err := os.Open(backupRootDir)
	if err != nil {
		return nil, err
	}
	return file.Readdirnames(0)
}

