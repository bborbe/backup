package main

import (
	"fmt"
	backup_config "github.com/bborbe/backup/constants"
	backup_dto "github.com/bborbe/backup/dto"
	backup_service "github.com/bborbe/backup/service"
	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/lock"
	"github.com/golang/glog"
	"runtime"
	"sort"
	"time"
)

const (
	defaultWait      = time.Minute * 5
	defaultLockName  = "/var/run/backup_cleanup.lock"
	parameterLock    = "lock"
	parameterHost    = "host"
	parameterRootdir = "target"
	parameterWait    = "wait"
	parameterOneTime = "one-time"
)

var (
	rootdirPtr = flag.String(parameterRootdir, backup_config.DEFAULT_ROOT_DIR, "backup root directory")
	hostPtr    = flag.String(parameterHost, backup_config.DEFAULT_HOST, "host to cleanup")
	lockPtr    = flag.String(parameterLock, defaultLockName, "lock file")
	waitPtr    = flag.Duration(parameterWait, defaultWait, "wait")
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
	glog.V(2).Infof("do started")
	l := lock.NewLock(*lockPtr)
	if err := l.Lock(); err != nil {
		return err
	}
	defer l.Unlock()

	for {
		glog.V(1).Infof("backup started")
		if err := cleanup(); err != nil {
			return err
		}
		glog.V(1).Infof("backup finished")

		if *oneTimePtr {
			glog.V(2).Infof("one-time => exit")
			return nil
		}

		glog.V(2).Infof("wait %v", *waitPtr)
		time.Sleep(*waitPtr)
		glog.V(2).Infof("sleep done")
	}
	glog.V(2).Infof("do finished")
	return nil
}

func cleanup() error {
	glog.V(2).Info("backup cleanup started")
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
		err := backupService.Cleanup(host)
		if err != nil {
			return err
		}
		glog.V(1).Infof("clean backups of host %s completed", host.GetName())
	}
	glog.V(2).Info("backup cleanup finished")
	return nil
}
