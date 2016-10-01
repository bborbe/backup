package main

import (
	"flag"

	"runtime"

	"fmt"
	"github.com/golang/glog"
)

const (
	parameterConfigPath  = "config"
	parameterTarget      = "target"
	parameterUser        = "user"
	parameterHost        = "host"
	parameterPort        = "port"
	parameterDirectory   = "dir"
	parameterExcludeFrom = "exclude_from"
)

var (
	configPathPtr  = flag.String(parameterConfigPath, "", "path to json config")
	targetPtr      = flag.String(parameterTarget, "", "target")
	userPtr        = flag.String(parameterUser, "", "user")
	hostPtr        = flag.String(parameterHost, "", "host")
	portPtr        = flag.Int(parameterPort, 22, "port")
	dirPtr         = flag.String(parameterDirectory, "", "dir")
	excludeFromPtr = flag.String(parameterExcludeFrom, "", "exclude_from")
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
	glog.Infof("backup started")
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
	glog.Infof("backup finished")
	return nil
}

func getTargetDirectory() (*targetDirectory, error) {
	targetDirectory := targetDirectory(*targetPtr)
	if err := targetDirectory.IsValid(); err != nil {
		return nil, err
	}
	return &targetDirectory, nil
}

func backupHosts(hosts []host, targetDirectory targetDirectory) error {
	for _, host := range hosts {
		glog.V(2).Infof("backup host started")
		if err := host.Backup(targetDirectory); err != nil {
			return err
		}
		glog.V(2).Infof("backup host finished")
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
