package main

import (
	"fmt"
	"runtime"
	"sort"
	"time"

	backup_config "github.com/bborbe/backup/constants"
	backup_dto "github.com/bborbe/backup/dto"
	backup_service "github.com/bborbe/backup/service"
	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/lock"
	"github.com/golang/glog"
	"github.com/bborbe/cron"
	"context"
)

const (
	defaultWait = time.Minute * 5
	defaultLockName = "/var/run/backup_cleanup.lock"
	parameterLock = "lock"
	parameterHost = "host"
	parameterRootdir = "target"
	parameterWait = "wait"
	parameterOneTime = "one-time"
)

var (
	rootdirPtr = flag.String(parameterRootdir, backup_config.DEFAULT_ROOT_DIR, "backup root directory")
	hostPtr = flag.String(parameterHost, backup_config.DEFAULT_HOST, "host to cleanup")
	lockPtr = flag.String(parameterLock, defaultLockName, "lock file")
	waitPtr = flag.Duration(parameterWait, defaultWait, "wait")
	oneTimePtr = flag.Bool(parameterOneTime, false, "exit after first fetch")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := do(); err != nil {
		glog.Exit(err)
	}
}

func do() error {
	l := lock.NewLock(*lockPtr)
	if err := l.Lock(); err != nil {
		return err
	}
	defer func() {
		if err := l.Unlock(); err != nil {
			glog.Warningf("unlock failed: %v", err)
		}
	}()

	cron := cron.New(
		*oneTimePtr,
		*waitPtr,
		cleanup,
	)
	return cron.Run(context.Background())
}

func cleanup(ctx context.Context) error {
	glog.V(1).Info("backup cleanup started")
	defer glog.V(1).Info("backup cleanup finished")
	if len(*rootdirPtr) == 0 {
		return fmt.Errorf("parameter %s missing", parameterRootdir)
	}
	glog.V(2).Infof("use backup dir %s", *rootdirPtr)
	backupService := backup_service.NewBackupService(*rootdirPtr)

	var hosts []backup_dto.Host
	var err error
	if *hostPtr == backup_config.DEFAULT_HOST {
		hosts, err = backupService.ListHosts()
		if err != nil {
			return err
		}
	} else {
		host, err := backupService.GetHost(*hostPtr)
		if err != nil {
			return err
		}
		hosts = []backup_dto.Host{host}
	}
	sort.Sort(backup_dto.HostByName(hosts))
	for _, host := range hosts {
		glog.V(1).Infof("clean backups of host %s stared", host.GetName())
		err := backupService.Cleanup(host)
		if err != nil {
			return err
		}
		glog.V(1).Infof("clean backups of host %s finished", host.GetName())
	}
	return nil
}
