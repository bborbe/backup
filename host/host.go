package host

import (
	"errors"
	"fmt"
	"os"

	backup_rootdir "github.com/bborbe/backup/rootdir"
	io_util "github.com/bborbe/io/util"
	"github.com/golang/glog"
)

type host struct {
	name    string
	rootdir backup_rootdir.Rootdir
}

type Host interface {
	Path() string
	Name() string
}

func ByName(rootdir backup_rootdir.Rootdir, name string) Host {
	h := new(host)
	h.rootdir = rootdir
	h.name = name
	return h
}

func All(root backup_rootdir.Rootdir) ([]Host, error) {
	file, err := os.Open(root.Path())
	defer file.Close()
	if err != nil {
		glog.V(2).Infof("open rootdir %s failed: %v", root.Path(), err)
		return nil, err
	}
	fileinfo, err := file.Stat()
	if err != nil {
		glog.V(2).Infof("file stat failed: %v", err)
		return nil, err
	}
	if !fileinfo.IsDir() {
		msg := fmt.Sprintf("rootdir %s is not a directory", root.Path())
		glog.V(2).Info(msg)
		return nil, errors.New(msg)
	}
	names, err := file.Readdirnames(0)
	if err != nil {
		glog.V(2).Infof("read dir names failed: %v", err)
		return nil, err
	}
	hosts := make([]Host, 0)
	for _, name := range names {
		host := ByName(root, name)
		isDir, err := io_util.IsDirectory(host.Path())
		if err != nil {
			return nil, err
		}
		if isDir {
			hosts = append(hosts, host)
		}
	}
	return hosts, nil
}

func (h *host) Path() string {
	return fmt.Sprintf("%s%c%s", h.rootdir.Path(), os.PathSeparator, h.name)
}

func (h *host) Name() string {
	return h.name
}
