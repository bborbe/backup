package main

import (
	"flag"
	"io"
	"os"

	"fmt"

	"github.com/bborbe/backup/config"
	"github.com/bborbe/backup/service"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

func main() {
	defer logger.Close()
	logLevelPtr := flag.Int("loglevel", config.DEFAULT_LOG_LEVEL, "int")
	rootdirPtr := flag.String("rootdir", config.DEFAULT_ROOT_DIR, "string")
	hostPtr := flag.String("host", config.DEFAULT_HOST, "string")
	flag.Parse()
	logger.SetLevelThreshold(*logLevelPtr)
	logger.Debugf("set log level to %s", log.LogLevelToString(*logLevelPtr))

	writer := os.Stdout
	logger.Debugf("use backup dir %s", *rootdirPtr)
	backupService := service.NewBackupService(*rootdirPtr)
	err := do(writer, backupService, *hostPtr)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
}

func do(writer io.Writer, backupService service.BackupService, hostname string) error {
	logger.Debug("start")
	host, err := backupService.GetHost(hostname)
	if err != nil {
		return err
	}
	err = backupService.Resume(host)
	if err != nil {
		fmt.Fprintf(writer, "resume backup for host %s failed\n", hostname)
		logger.Warn(err)
	} else {
		fmt.Fprintf(writer, "resume backup for host %s success\n", hostname)
	}
	logger.Debug("done")
	return nil
}
