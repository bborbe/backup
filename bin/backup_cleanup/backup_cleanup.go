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
	"github.com/bborbe/lock"
	"github.com/golang/glog"
)

const (
	LOCK_NAME = "/var/run/backup_clean.lock"
)

var (
	rootdirPtr = flag.String("rootdir", backup_config.DEFAULT_ROOT_DIR, "string")
	hostPtr    = flag.String("host", backup_config.DEFAULT_HOST, "string")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	writer := os.Stdout
	glog.V(2).Infof("use backup dir %s", *rootdirPtr)
	backupService := backup_service.NewBackupService(*rootdirPtr)
	err := do(writer, backupService, *rootdirPtr, *hostPtr, LOCK_NAME)
	if err != nil {
		glog.Exit(err)
	}
}

func do(writer io.Writer, backupService backup_service.BackupService, rootdirName string, hostName string, lockName string) error {
	var err error
	var hosts []backup_dto.Host

	l := lock.NewLock(lockName)
	err = l.Lock()
	if err != nil {
		return err
	}
	defer l.Unlock()
	glog.V(2).Info("start")

	if hostName == backup_config.DEFAULT_HOST {
		hosts, err = backupService.ListHosts()
		if err != nil {
			return err
		}
	} else {
		host, err := backupService.GetHost(hostName)
		if err != nil {
			return err
		}
		hosts = []backup_dto.Host{host}
	}
	sort.Sort(backup_dto.HostByName(hosts))
	for _, host := range hosts {
		err := backupService.Cleanup(host)
		if err != nil {
			return err
		}
		fmt.Fprintf(writer, "%s cleaned\n", host.GetName())
	}
	glog.V(2).Info("done")
	return nil
}
