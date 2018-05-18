package main

import (
	"fmt"
	"os"

	"github.com/golang/glog"
)

type targetDirectory string

func (t targetDirectory) String() string {
	return string(t)
}

func (t targetDirectory) IsValid() error {
	if _, err := os.Stat(t.String()); err != nil {
		glog.V(2).Infof("target %s invalid: %v", t, err)
		return fmt.Errorf("target %s invalid: %v", t, err)
	}
	glog.V(2).Infof("target %s valid", t)
	return nil
}
