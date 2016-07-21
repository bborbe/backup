package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	"runtime"

	backup_config "github.com/bborbe/backup/config"
	backup_dto "github.com/bborbe/backup/dto"
	backup_service "github.com/bborbe/backup/service"
	"github.com/bborbe/log"
)

const (
	PARAMETER_LOGLEVEL = "loglevel"
)

var (
	logger      = log.DefaultLogger
	logLevelPtr = flag.String(PARAMETER_LOGLEVEL, log.LogLevelToString(backup_config.DEFAULT_LOG_LEVEL), log.FLAG_USAGE)
	rootdirPtr  = flag.String("rootdir", backup_config.DEFAULT_ROOT_DIR, "string")
	hostPtr     = flag.String("host", backup_config.DEFAULT_HOST, "string")
)

func main() {
	defer logger.Close()
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	writer := os.Stdout
	logger.Debugf("use backup dir %s", *rootdirPtr)
	backupService := backup_service.NewBackupService(*rootdirPtr)
	err := do(writer, backupService, *hostPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, backupService backup_service.BackupService, hostname string) error {
	logger.Debug("start")
	var hosts []backup_dto.Host
	var err error

	if hostname == backup_config.DEFAULT_HOST {
		hosts, err = backupService.ListHosts()
		if err != nil {
			return err
		}
	} else {
		host, err := backupService.GetHost(hostname)
		if err != nil {
			return err
		}
		hosts = []backup_dto.Host{host}
	}
	sort.Sort(backup_dto.HostByName(hosts))
	for _, host := range hosts {
		backup, err := backupService.GetLatestBackup(host)
		if err != nil {
			return err
		}
		if backup != nil {
			fmt.Fprintf(writer, "%s/%s\n", host.GetName(), backup.GetName())
		}
	}
	logger.Debug("done")
	return nil
}
