package main

import (
	"flag"

	backup_config "github.com/bborbe/backup/config"
	backup_status_client "github.com/bborbe/backup/status_client"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

const (
	DEFAULT_PORT int = 8080
	DEFAULT_SERVER = "http://backup.pn.benjamin-borbe.de:7777"
	PARAMETER_LOGLEVEL = "loglevel"
	PARAMETER_PORT = "port"
	PARAMETER_SERVER = "server"
)

func main() {
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.LogLevelToString(backup_config.DEFAULT_LOG_LEVEL), log.FLAG_USAGE)
	serverPtr := flag.String(PARAMETER_SERVER, DEFAULT_SERVER, "backup status server address")
	portnumberPtr := flag.Int(PARAMETER_PORT, DEFAULT_PORT, "server port")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Tracef("set log level to %s", *logLevelPtr)
	logger.Tracef("server %s", *serverPtr)
	logger.Tracef("portnumberPtr %d", *portnumberPtr)
	logger.Debugf("backup status server started at port %d", *portnumberPtr)
	srv := backup_status_client.NewServer(*portnumberPtr, *serverPtr)
	srv.Run()
}
