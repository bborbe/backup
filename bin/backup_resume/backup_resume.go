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

const NO_HOST = "-"

func main() {
	defer logger.Close()
	logLevelPtr := flag.String("loglevel", log.LogLevelToString(config.DEFAULT_LOG_LEVEL), "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	rootdirPtr := flag.String("rootdir", config.DEFAULT_ROOT_DIR, "string")
	hostPtr := flag.String("host", NO_HOST, "string")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	writer := os.Stdout
	logger.Debugf("use backup dir %s", *rootdirPtr)
	backupService := service.NewBackupService(*rootdirPtr)
	err := do(writer, backupService, *hostPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, backupService service.BackupService, hostname string) error {
	logger.Debug("start")
	if hostname == NO_HOST {
		return fmt.Errorf("parameter host missing")
	}
	host, err := backupService.GetHost(hostname)
	if err != nil {
		fmt.Fprintf(writer, "host %s not found", hostname)
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
