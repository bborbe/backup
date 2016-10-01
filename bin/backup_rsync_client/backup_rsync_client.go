package main

import (
	"flag"

	"runtime"

	"github.com/golang/glog"
)

const ()

var ()

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
	glog.Infof("run")
	return nil
}
