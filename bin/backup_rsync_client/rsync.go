package main

import (
	"github.com/golang/glog"
	"os/exec"
)

func runRsync(args ...string) error {
	cmd := exec.Command("rsync", args...)
	if err := cmd.Start(); err != nil {
		return err
	}
	glog.V(2).Infof("rsync started")
	if err := cmd.Wait(); err != nil {
		return err
	}
	glog.V(2).Infof("rsync finished")
	return nil
}
