package main

import (
	"flag"
	"github.com/bborbe/backup/config"
	"github.com/bborbe/backup/status_server"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

const DEFAULT_PORT int = 8002

func main() {
	defer logger.Close()
	logLevelPtr := flag.String("loglevel", log.LogLevelToString(config.DEFAULT_LOG_LEVEL), "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	rootdirPtr := flag.String("rootdir", config.DEFAULT_ROOT_DIR, "root directory for backups")
	portnumberPtr := flag.Int("port", DEFAULT_PORT, "server port")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Tracef("set log level to %s", *logLevelPtr)
	logger.Tracef("rootdir %s", *rootdirPtr)
	logger.Tracef("portnumberPtr %d", *portnumberPtr)
	logger.Debugf("backup status server started at port %d", *portnumberPtr)
	srv := status_server.NewServer(*portnumberPtr, *rootdirPtr)
	{
		err := srv.Start()
		if err != nil {
			logger.Errorf("start server failed, %v", err)
			return
		}
	}
	srv.Wait()
	{
		err := srv.Stop()
		if err != nil {
			logger.Errorf("stop server failed, %v", err)
			return
		}
	}
	logger.Debug("backup status server finished")
}
