package main

import (
	"fmt"
	"runtime"
	"time"

	"context"

	"github.com/bborbe/cron"
	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/lock"
	"github.com/golang/glog"
)

const (
	defaultWait          = time.Minute * 5
	defaultLockName      = "/var/run/backup_rsync_client.lock"
	parameterConfigPath  = "config"
	parameterTarget      = "target"
	parameterUser        = "user"
	parameterHost        = "host"
	parameterIp          = "ip"
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
	ipPtr          = flag.String(parameterIp, "", "ip")
	portPtr        = flag.Int(parameterPort, 22, "port")
	dirPtr         = flag.String(parameterDirectory, "", "dir")
	excludeFromPtr = flag.String(parameterExcludeFrom, "", "exclude from")
	waitPtr        = flag.Duration(parameterWait, defaultWait, "wait")
	oneTimePtr     = flag.Bool(parameterOneTime, false, "exit after first fetch")
	lockPtr        = flag.String(parameterLock, defaultLockName, "lock")
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
	defer func() {
		if err := l.Unlock(); err != nil {
			glog.Warningf("unlock failed: %v", err)
		}
	}()

	var c cron.Cron
	if *oneTimePtr {
		c = cron.NewOneTimeCron(backup)
	} else {
		c = cron.NewWaitCron(
			*waitPtr,
			backup,
		)
	}
	return c.Run(context.Background())
}

func backup(ctx context.Context) error {
	hosts, err := getHosts()
	if err != nil {
		return err
	}
	hosts = parseActive(hosts)
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

func parseActive(hosts []host) []host {
	result := []host{}
	for _, host := range hosts {
		if host.Active {
			result = append(result, host)
		}
	}
	glog.V(2).Infof("%d of %d hosts are active", len(result), len(hosts))
	return result
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
		glog.V(2).Infof("backup %s started", host.Host)
		if err := host.Backup(targetDirectory); err != nil {
			glog.Warningf("backup %s failed: %v", host.Host, err)
		} else {
			glog.V(2).Infof("backup %s finished", host.Host)
		}
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
	glog.V(2).Infof("get hosts")
	if configPathPtr != nil {
		configPath := configPath(*configPathPtr)
		if configPath.IsValue() {
			glog.V(2).Infof("read config %s", configPath)
			return configPath.ParseHosts()
		}
	}
	glog.V(2).Infof("create hosts from args")
	return []host{{
		Active:      true,
		User:        *userPtr,
		Host:        *hostPtr,
		Ip:          *ipPtr,
		Port:        *portPtr,
		Directory:   *dirPtr,
		ExcludeFrom: *excludeFromPtr,
	}}, nil
}
