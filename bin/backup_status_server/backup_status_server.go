package main

import (
	"flag"

	backup_config "github.com/bborbe/backup/config"
	backup_status_server "github.com/bborbe/backup/status_server"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

const DEFAULT_PORT int = 8002

func main() {
	logLevelPtr := flag.String("loglevel", log.LogLevelToString(backup_config.DEFAULT_LOG_LEVEL), "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	rootdirPtr := flag.String("rootdir", backup_config.DEFAULT_ROOT_DIR, "root directory for backups")
	portnumberPtr := flag.Int("port", DEFAULT_PORT, "server port")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Tracef("set log level to %s", *logLevelPtr)
	logger.Tracef("rootdir %s", *rootdirPtr)
	logger.Tracef("portnumberPtr %d", *portnumberPtr)
	logger.Debugf("backup status server started at port %d", *portnumberPtr)
	srv := backup_status_server.NewServer(*portnumberPtr, *rootdirPtr)
	srv.Run()
}
