package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/bborbe/backup/constants"
	"github.com/bborbe/backup/date"
	"github.com/golang/glog"
)

type host struct {
	Active      bool   `json:"active"`
	User        string `json:"user"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Directory   string `json:"dir"`
	ExcludeFrom string `json:"exclude_from"`
}

func (h *host) Backup(targetDirectory targetDirectory) error {
	glog.V(2).Infof("create backup of %s on target: %s", h.Host, targetDirectory)

	if err := h.createCurrentDirectory(targetDirectory); err != nil {
		glog.V(2).Infof("create current failed: %v", err)
		return err
	}

	found, err := h.todayHasAlreadyBackup(targetDirectory)
	if err != nil {
		glog.V(2).Infof("search for existing backups failed: %v", err)
		return err
	}
	if found {
		glog.V(2).Infof("skip backup => already exists")
		return nil
	}

	glog.V(1).Infof("backup %s started", h.Host)

	if err := h.createToDirectory(targetDirectory); err != nil {
		glog.V(2).Infof("create target directory failed: %v", err)
		return err
	}

	if err := h.rsync(targetDirectory); err != nil {
		glog.V(2).Infof("rsync failed: %v", err)
		return err
	}

	backupDate := time.Now()

	if err := h.renameIncompleteToDate(targetDirectory, backupDate); err != nil {
		glog.V(2).Infof("rename incomplete to date failed: %v", err)
		return err
	}

	if err := h.deleteCurrentSymlink(targetDirectory); err != nil {
		glog.V(2).Infof("delete current symlink failed: %v", err)
		return err
	}

	if err := h.createCurrentSymlink(targetDirectory, backupDate); err != nil {
		glog.V(2).Infof("create new current symlink failed: %v", err)
		return err
	}

	if err := h.deleteEmpty(targetDirectory); err != nil {
		glog.V(2).Infof("delete empty dir failed: %v", err)
		return err
	}

	glog.V(1).Infof("backup %s finished", h.Host)
	return nil
}

func (h *host) rsync(targetDirectory targetDirectory) error {
	return runRsync(
		"-azP",
		"--no-p",
		"--numeric-ids",
		"-e",
		fmt.Sprintf("ssh -T -x -o StrictHostKeyChecking=no -p %d", h.Port),
		"--delete",
		"--delete-excluded",
		fmt.Sprintf("--port=%d", h.Port),
		fmt.Sprintf("--exclude-from=%s", h.ExcludeFrom),
		fmt.Sprintf("--link-dest=%s", h.linkDest(targetDirectory)),
		h.from(),
		h.to(targetDirectory),
	)
}

func (h *host) renameIncompleteToDate(targetDirectory targetDirectory, date time.Time) error {
	glog.V(2).Infof("rename incomplete on target: %s", targetDirectory)
	return os.Rename(h.incomplete(targetDirectory), h.date(targetDirectory, date))
}

func (h *host) deleteCurrentSymlink(targetDirectory targetDirectory) error {
	glog.V(2).Infof("delete current symlink on target: %s", targetDirectory)
	return os.Remove(h.current(targetDirectory))
}

func (h *host) createCurrentSymlink(targetDirectory targetDirectory, date time.Time) error {
	glog.V(2).Infof("create current symlink on target: %s", targetDirectory)
	if err := os.Symlink(date.Format(constants.DATEFORMAT), h.current(targetDirectory)); err != nil {
		glog.V(2).Infof("create symlink to current failed: %v", err)
		return err
	}
	return nil
}

func (h *host) deleteEmpty(targetDirectory targetDirectory) error {
	glog.V(2).Infof("delete empty directory on target: %s", targetDirectory)
	os.RemoveAll(h.empty(targetDirectory))
	return nil
}

func (h *host) createCurrentDirectory(targetDirectory targetDirectory) error {
	glog.V(2).Infof("create current directory on target: %s", targetDirectory)
	_, err := os.Stat(h.current(targetDirectory))
	if os.IsNotExist(err) {
		glog.V(2).Infof("current not existing")
		if err := os.MkdirAll(h.empty(targetDirectory)+h.Directory, 0700); err != nil {
			glog.V(2).Infof("create empty directory failed: %v", err)
			return err
		}
		if err := os.Symlink("empty", h.current(targetDirectory)); err != nil {
			glog.V(2).Infof("create symlink from empty to current failed: %v", err)
			return err
		}
	}
	glog.V(2).Infof("create current completed")
	return nil
}

func (h *host) todayHasAlreadyBackup(targetDirectory targetDirectory) (bool, error) {
	glog.V(2).Infof("search if backup already exists on target: %s", targetDirectory)
	hostDirectory := fmt.Sprintf("%s/%s", targetDirectory, h.Host)
	dirs, err := ioutil.ReadDir(hostDirectory)
	if err != nil {
		glog.V(2).Infof("read directory %s failed: %v", targetDirectory, err)
		return false, err
	}
	for _, dir := range dirs {
		timeOfDir, err := time.Parse(constants.DATEFORMAT, dir.Name())
		if err != nil {
			glog.V(2).Infof("parse date of dir failed")
			continue
		}
		if date.DayEqual(timeOfDir, time.Now()) {
			glog.V(2).Infof("found backup for today")
			return true, nil
		}
	}
	glog.V(2).Infof("found no backup for today")
	return false, nil
}

func (h *host) createToDirectory(targetDirectory targetDirectory) error {
	glog.V(2).Infof("create to directory on target: %s", targetDirectory)
	return os.MkdirAll(h.to(targetDirectory), 0700)
}

func (h *host) from() string {
	return fmt.Sprintf("%s@%s:%s", h.User, h.Host, h.Directory)
}

func (h *host) to(targetDirectory targetDirectory) string {
	return fmt.Sprintf("%s/%s/incomplete%s", targetDirectory, h.Host, h.Directory)
}

func (h *host) incomplete(targetDirectory targetDirectory) string {
	return fmt.Sprintf("%s/%s/incomplete", targetDirectory, h.Host)
}

func (h *host) current(targetDirectory targetDirectory) string {
	return fmt.Sprintf("%s/%s/current", targetDirectory, h.Host)
}

func (h *host) empty(targetDirectory targetDirectory) string {
	return fmt.Sprintf("%s/%s/empty", targetDirectory, h.Host)
}

func (h *host) date(targetDirectory targetDirectory, date time.Time) string {
	return fmt.Sprintf("%s/%s/%s", targetDirectory, h.Host, date.Format(constants.DATEFORMAT))
}

func (h *host) linkDest(targetDirectory targetDirectory) string {
	return fmt.Sprintf("%s/%s/current%s", targetDirectory, h.Host, h.Directory)
}

func (h *host) Validate() error {
	glog.V(2).Infof("validate host: %s", h.Host)
	if len(h.User) == 0 {
		return fmt.Errorf("user invalid")
	}
	if len(h.Host) == 0 {
		return fmt.Errorf("host invalid")
	}
	if h.Port <= 0 {
		return fmt.Errorf("port invalid")
	}
	if len(h.Directory) == 0 {
		return fmt.Errorf("directory invalid")
	}
	if h.Directory[len(h.Directory)-1:] != "/" {
		return fmt.Errorf("directory not ending with '/'")
	}
	return nil
}
