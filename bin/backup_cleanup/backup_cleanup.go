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
	"github.com/bborbe/backup/util"
	"github.com/bborbe/lock"
	"github.com/bborbe/log"
)

const LOCK_NAME = "/var/run/backup_clean.lock"

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
	err := do(writer, backupService, *rootdirPtr, *hostPtr, LOCK_NAME)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
}

func do(writer io.Writer, backupService service.BackupService, rootdirName string, hostName string, lockName string) error {
	var err error
	var hosts []dto.Host

	l := lock.NewLock(lockName)
	err = l.Lock()
	if err != nil {
		return err
	}
	defer l.Unlock()
	logger.Debug("start")

	rootdir, err := backupService.GetRootdir(rootdirName)
	if err != nil {
		return err
	}

	if hostName == config.DEFAULT_HOST {
		hosts, err = backupService.ListHosts()
		if err != nil {
			return err
		}
	} else {
		host, err := backupService.GetHost(rootdir, hostName)
		if err != nil {
			return err
		}
		hosts = []dto.Host{host}
	}
	sort.Sort(util.HostByDate(hosts))
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
