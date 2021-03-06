package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	"runtime"

	backup_config "github.com/bborbe/backup/constants"
	backup_dto "github.com/bborbe/backup/dto"
	backup_service "github.com/bborbe/backup/service"
	"github.com/golang/glog"
)

const (
	parameterTarget = "target"
	parameterHost   = "host"
)

var (
	rootdirPtr = flag.String(parameterTarget, backup_config.DEFAULT_ROOT_DIR, "backup directory")
	hostPtr    = flag.String(parameterHost, backup_config.DEFAULT_HOST, "host")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	writer := os.Stdout
	glog.V(2).Infof("use backup dir %s", *rootdirPtr)
	backupService := backup_service.NewBackupService(*rootdirPtr)
	err := do(writer, backupService, *hostPtr)
	if err != nil {
		glog.Exit(err)
	}
}

func do(writer io.Writer, backupService backup_service.BackupService, hostname string) error {
	glog.V(2).Info("start")
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
	glog.V(2).Info("done")
	return nil
}
