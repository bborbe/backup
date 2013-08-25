package main

import (
	"flag"
	"fmt"
	"github.com/bborbe/backup/service"
	"github.com/bborbe/log"
	"io"
	"os"
)

var logger = log.DefaultLogger

const (
	DEFAULT_LOG_LEVEL = log.OFF
	DEFAULT_ROOT_DIR  = "/rsync"
)

func main() {
	logLevelPtr := flag.Int("loglevel", DEFAULT_LOG_LEVEL, "int")
	rootdirPtr := flag.String("loglevel", DEFAULT_ROOT_DIR, "string")
	flag.Parse()
	logger.SetLevelThreshold(*logLevelPtr)
	logger.Debugf("set log level to %s", log.LogLevelToString(*logLevelPtr))

	writer := os.Stdout
	logger.Debugf("use backup dir %s", *rootdirPtr)
	backupService := service.NewBackupService(*rootdirPtr)
	err := do(writer, backupService)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
}

func do(writer io.Writer, backupService service.BackupService) error {
	logger.Debug("start")
	hosts, err := backupService.ListHosts()
	if err != nil {
		return err
	}
	for _, host := range hosts {
		backups, err := backupService.ListBackups(host)
		if err != nil {
			return err
		}
		for _, backup := range backups {
			fmt.Fprintf(writer, "%s => %s\n", host.GetName(), backup.GetName())
		}
	}
	logger.Debug("done")
	return nil
}
