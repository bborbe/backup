package main

import (
	"fmt"
	"github.com/bborbe/backup/constants"
	"github.com/golang/glog"
	"os"
	"time"
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

	if h.todayHasAllreadyBackup(targetDirectory) {
		glog.V(1).Infof("skip backup => allready exists")
		return nil
	}

	if err := h.createCurrentDirectory(targetDirectory); err != nil {
		glog.V(1).Infof("create current failed: %v", err)
		return err
	}

	if err := h.createToDirectory(targetDirectory); err != nil {
		return err
	}

	if err := runRsync(
		"-azP",
		"-e",
		fmt.Sprintf("ssh -p %d", h.Port),
		"--delete",
		"--delete-excluded",
		fmt.Sprintf("--port=%d", h.Port),
		fmt.Sprintf("--exclude-from=%s", h.ExcludeFrom),
		fmt.Sprintf("--link-dest=%s", h.linkDest(targetDirectory)),
		h.from(),
		h.to(targetDirectory),
	); err != nil {
		return err
	}

	backupDate := time.Now()

	if err := h.renameIncompleteToDate(targetDirectory, backupDate); err != nil {
		return err
	}

	if err := h.deleteCurrentSymlink(targetDirectory); err != nil {
		return err
	}

	if err := h.createCurrentSymlink(targetDirectory, backupDate); err != nil {
		return err
	}

	if err := h.deleteEmpty(targetDirectory); err != nil {
		return err
	}

	glog.V(1).Infof("backup completed")
	return nil
}

func (h *host) renameIncompleteToDate(targetDirectory targetDirectory, date time.Time) error {
	return os.Rename(h.incomplete(targetDirectory), h.date(targetDirectory, date))
}

func (h *host) deleteCurrentSymlink(targetDirectory targetDirectory) error {
	return os.Remove(h.current(targetDirectory))
}

func (h *host) createCurrentSymlink(targetDirectory targetDirectory, date time.Time) error {
	if err := os.Symlink(h.date(targetDirectory, date), h.current(targetDirectory)); err != nil {
		glog.V(2).Infof("create symlink to current failed: %v", err)
		return err
	}
	return nil
}

func (h *host) deleteEmpty(targetDirectory targetDirectory) error {
	os.RemoveAll(h.empty(targetDirectory))
	return nil
}

func (h *host) createCurrentDirectory(targetDirectory targetDirectory) error {
	_, err := os.Stat(h.current(targetDirectory))
	if os.IsNotExist(err) {
		glog.V(2).Infof("current not existing")
		if err := os.MkdirAll(h.empty(targetDirectory)+h.Directory, 0700); err != nil {
			glog.V(2).Infof("create empty directory failed: %v", err)
			return err
		}
		if err := os.Symlink(h.empty(targetDirectory), h.current(targetDirectory)); err != nil {
			glog.V(2).Infof("create symlink from empty to current failed: %v", err)
			return err
		}
	}
	glog.V(2).Infof("create current completed")
	return nil
}
func (h *host) todayHasAllreadyBackup(targetDirectory targetDirectory) bool {
	return false
}

func (h *host) createToDirectory(targetDirectory targetDirectory) error {
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
