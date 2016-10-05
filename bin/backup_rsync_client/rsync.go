package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/golang/glog"
)

func runRsync(args ...string) error {
	glog.V(2).Infof("run: rsync %s", strings.Join(args, " "))
	cmd := exec.Command("rsync", args...)
	if glog.V(4) {
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
