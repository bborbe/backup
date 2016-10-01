package main

import (
	"encoding/json"
	"github.com/golang/glog"
	"os"
)

type configPath string

func (c configPath) String() string {
	return string(c)
}

func (c configPath) ParseHosts() ([]host, error) {
	file, err := os.Open(c.String())
	if err != nil {
		glog.V(1).Infof("open file %v failed: %v", c, err)
		return nil, err
	}
	var hosts []host
	if err := json.NewDecoder(file).Decode(&hosts); err != nil {
		glog.V(1).Infof("parse json failed: %v", err)
		return nil, err
	}
	glog.V(2).Infof("found %d hosts", len(hosts))
	return hosts, nil
}

func (c configPath) IsValue() bool {
	return len(c) > 0
}
