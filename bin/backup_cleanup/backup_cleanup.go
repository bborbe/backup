package main

import (
	"flag"
	"github.com/bborbe/backup/config"
	"github.com/bborbe/backup/service"
	"github.com/bborbe/log"
	"os"
)

var logger = log.DefaultLogger

func main() {
	logLevelPtr := flag.Int("loglevel", config.DEFAULT_LOG_LEVEL, "int")
	rootdirPtr := flag.String("rootdir", config.DEFAULT_ROOT_DIR, "string")
	flag.Parse()
	logger.SetLevelThreshold(*logLevelPtr)
	logger.Debugf("set log level to %s", log.LogLevelToString(*logLevelPtr))

	logger.Debugf("use backup dir %s", *rootdirPtr)
	backupService := service.NewBackupService(*rootdirPtr)
	logger.Debug("start")
	err := backupService.Cleanup()
	logger.Debug("done")
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
}
