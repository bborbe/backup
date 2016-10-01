package main

import (
	"flag"

	"runtime"

	"fmt"
	"github.com/bborbe/lock"
	"github.com/golang/glog"
	"time"
)

const (
	defaultWait          = time.Minute * 5
	defaultLockName      = "/var/run/backup_rsync_client.lock"
	parameterConfigPath  = "config"
	parameterTarget      = "target"
	parameterUser        = "user"
	parameterHost        = "host"
	parameterPort        = "port"
	parameterDirectory   = "dir"
	parameterExcludeFrom = "exclude_from"
	parameterWait        = "wait"
	parameterOneTime     = "one-time"
	parameterLock        = "lock"
)

var (
	configPathPtr  = flag.String(parameterConfigPath, "", "path to json config")
	targetPtr      = flag.String(parameterTarget, "", "target")
	userPtr        = flag.String(parameterUser, "", "user")
	hostPtr        = flag.String(parameterHost, "", "host")
	portPtr        = flag.Int(parameterPort, 22, "port")
	dirPtr         = flag.String(parameterDirectory, "", "dir")
	excludeFromPtr = flag.String(parameterExcludeFrom, "", "exclude_from")
	waitPtr        = flag.Duration(parameterWait, defaultWait, "wait")
	oneTimePtr     = flag.Bool(parameterOneTime, false, "exit after first fetch")
	lockPtr        = flag.String(parameterLock, defaultLockName, "lock")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := do()
	if err != nil {
		glog.Exit(err)
	}
}

func do() error {
	l := lock.NewLock(*lockPtr)
	if err := l.Lock(); err != nil {
		return err
	}
	defer l.Unlock()

	for {
		glog.V(1).Infof("backup started")
		if err := backup(); err != nil {
			return err
		}
		glog.V(1).Infof("backup finished")

		if *oneTimePtr {
			return nil
		}

		glog.V(2).Infof("wait %v", *waitPtr)
		time.Sleep(*waitPtr)
		glog.V(2).Infof("sleep done")
	}

	return nil
}

func backup() error {
	hosts, err := getHosts()
	if err != nil {
		return err
	}
	if err := validateHosts(hosts); err != nil {
		return err
	}
	targetDirectory, err := getTargetDirectory()
	if err != nil {
		return err
	}
	if err := backupHosts(hosts, *targetDirectory); err != nil {
		return err
	}
	return nil
}

func getTargetDirectory() (*targetDirectory, error) {
	target := targetDirectory(*targetPtr)
	if err := target.IsValid(); err != nil {
		return nil, err
	}
	return &target, nil
}

func backupHosts(hosts []host, targetDirectory targetDirectory) error {
	for _, host := range hosts {
		glog.V(1).Infof("backup %s started", host.Host)
		if err := host.Backup(targetDirectory); err != nil {
			return err
		}
		glog.V(1).Infof("backup %s finished", host.Host)
	}
	return nil
}

func validateHosts(hosts []host) error {
	if len(hosts) == 0 {
		return fmt.Errorf("no host found in config")
	}
	for _, host := range hosts {
		if err := host.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func getHosts() ([]host, error) {
	if configPathPtr != nil {
		configPath := configPath(*configPathPtr)
		if configPath.IsValue() {
			return configPath.ParseHosts()
		}
	}
	return []host{{
		Active:      true,
		User:        *userPtr,
		Host:        *hostPtr,
		Port:        *portPtr,
		Directory:   *dirPtr,
		ExcludeFrom: *excludeFromPtr,
	}}, nil
}
