package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/bborbe/backup/config"
	"github.com/bborbe/backup/dto"
	"github.com/bborbe/backup/service"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

func main() {
	defer logger.Close()
	logLevelPtr := flag.String("loglevel", log.LogLevelToString(config.DEFAULT_LOG_LEVEL), "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	rootdirPtr := flag.String("rootdir", config.DEFAULT_ROOT_DIR, "string")
	hostPtr := flag.String("host", config.DEFAULT_HOST, "string")
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
	var hosts []dto.Host
	var err error

	if hostname == config.DEFAULT_HOST {
		hosts, err = backupService.ListHosts()
		if err != nil {
			return err
		}
	} else {
		host, err := backupService.GetHost(hostname)
		if err != nil {
			return err
		}
		hosts = []dto.Host{host}
	}
	sort.Sort(dto.HostByName(hosts))
	for _, host := range hosts {
		backups, err := backupService.ListBackups(host)
		if err != nil {
			return err
		}
		sort.Sort(dto.BackupByName(backups))
		for _, backup := range backups {
			fmt.Fprintf(writer, "%s/%s\n", host.GetName(), backup.GetName())
		}
	}
	logger.Debug("done")
	return nil
}
