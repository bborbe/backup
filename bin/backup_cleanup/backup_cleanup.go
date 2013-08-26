package main

import (
	"flag"
	"fmt"
	"github.com/bborbe/backup/config"
	"github.com/bborbe/backup/service"
	"github.com/bborbe/log"
	"io"
	"os"
)

var logger = log.DefaultLogger

func main() {
	logLevelPtr := flag.Int("loglevel", config.DEFAULT_LOG_LEVEL, "int")
	rootdirPtr := flag.String("rootdir", config.DEFAULT_ROOT_DIR, "string")
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
		err := backupService.Cleanup(host)
		if err != nil {
			return err
		}
		fmt.Fprintf(writer, "%s cleaned\n", host.GetName())
	}
	logger.Debug("done")
	return nil
}
