package main

import (
	"flag"

	"runtime"

	"fmt"
	"github.com/golang/glog"
)

const (
	parameterConfigPath = "config"
)

var (
	configPathPtr = flag.String(parameterConfigPath, "", "path to json config")
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

	configPath := configPath(*configPathPtr)
	hosts, err := configPath.ParseHosts()
	if err != nil {
		return fmt.Errorf("parse config %v failed: %v", configPath, err)
	}
	if len(hosts) == 0 {
		return fmt.Errorf("no host found in config")
	}
	for _, host := range hosts {
		glog.V(2).Infof("backup host started")
		if err := host.Backup(); err != nil {
			return err
		}
		glog.V(2).Infof("backup host finished")
	}

	glog.Infof("backup finished")
	return nil
}
