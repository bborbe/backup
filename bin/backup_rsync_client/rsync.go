package main

import (
	"github.com/golang/glog"
	"os"
	"os/exec"
	"strings"
)

func runRsync(args ...string) error {
	glog.V(1).Infof("run: rsync %s", strings.Join(args, " "))
	cmd := exec.Command("rsync", args...)
	if glog.V(2) {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}
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
