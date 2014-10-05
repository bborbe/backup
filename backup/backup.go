package backup

import (
	"errors"
	"fmt"
	"os"

	"github.com/bborbe/backup/host"
	"github.com/bborbe/log"
)

type backup struct {
	name string
	host host.Host
}

type Backup interface {
	Path() string
	Name() string
}

var logger = log.DefaultLogger

func ByName(h host.Host, name string) Backup {
	b := new(backup)
	b.host = h
	b.name = name
	return b
}

func All(h host.Host) ([]Backup, error) {
	file, err := os.Open(h.Path())
	if err != nil {
		logger.Debugf("open host %s failed: %v", h.Path(), err)
		return nil, err
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		logger.Debugf("file stat failed: %v", err)
		return nil, err
	}
	if !fileinfo.IsDir() {
		msg := fmt.Sprintf("host %s is not a directory", h.Path())
		logger.Debug(msg)
		return nil, errors.New(msg)
	}
	names, err := file.Readdirnames(0)
	if err != nil {
		logger.Debugf("read dir names failed: %v", err)
		return nil, err
	}
	backups := make([]Backup, 0)
	for _, name := range names {
		backups = append(backups, ByName(h, name))
	}
	return backups, nil
}

func (h *backup) Path() string {
	return fmt.Sprintf("%s%c%s", h.host.Path(), os.PathSeparator, h.name)
}

func (h *backup) Name() string {
	return h.name
}
