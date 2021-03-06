package backup

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/bborbe/backup/constants"
	backup_host "github.com/bborbe/backup/host"
	backup_timeparser "github.com/bborbe/backup/timeparser"
	io_util "github.com/bborbe/io/util"
	"github.com/golang/glog"
)

type backup struct {
	name string
	host backup_host.Host
}

type Backup interface {
	Path() string
	Name() string
	Delete() error
	IsDir() (bool, error)
	ValidName() bool
}

const INCOMPLETE = "incomplete"
const CURRENT = "current"

func ByTime(h backup_host.Host, t time.Time) Backup {
	return ByName(h, t.Format(constants.DATEFORMAT))
}

func ByName(h backup_host.Host, name string) Backup {
	b := new(backup)
	b.host = h
	b.name = name
	return b
}

func All(h backup_host.Host) ([]Backup, error) {
	file, err := os.Open(h.Path())
	defer file.Close()
	if err != nil {
		glog.V(2).Infof("open host %s failed: %v", h.Path(), err)
		return nil, err
	}
	fileinfo, err := file.Stat()
	if err != nil {
		glog.V(2).Infof("file stat failed: %v", err)
		return nil, err
	}
	if !fileinfo.IsDir() {
		msg := fmt.Sprintf("host %s is not a directory", h.Path())
		glog.V(2).Info(msg)
		return nil, errors.New(msg)
	}
	names, err := file.Readdirnames(0)
	if err != nil {
		glog.V(2).Infof("read dir names failed: %v", err)
		return nil, err
	}
	backups := make([]Backup, 0)
	for _, name := range names {
		backup := ByName(h, name)
		isDir, err := backup.IsDir()
		if err != nil {
			return nil, err
		}
		if isDir && backup.ValidName() {
			backups = append(backups, backup)
		}
	}
	return backups, nil
}

func (b *backup) IsDir() (bool, error) {
	if !io_util.Exists(b.Path()) {
		return false, nil
	}
	return io_util.IsDirectory(b.Path())
}

func (b *backup) ValidName() bool {
	return validName(b.Name())
}

func validName(name string) bool {
	re := regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}")
	return re.MatchString(name)
}

func (b *backup) Path() string {
	return fmt.Sprintf("%s%c%s", b.host.Path(), os.PathSeparator, b.name)
}

func (b *backup) Name() string {
	return b.name
}

func (b *backup) Delete() error {
	return os.RemoveAll(b.Path())
}

func Resume(h backup_host.Host) error {
	existsIncomplete, err := existsIncomplete(h)
	if err != nil {
		return err
	}
	if existsIncomplete {
		glog.V(2).Info("skip resume => incomplete dir exists")
		return nil
	}
	backups, err := All(h)
	if err != nil {
		return err
	}
	if len(backups) < 2 {
		return fmt.Errorf("can't resume with less than two existing backups")
	}
	existsBackupForToday, err := existsBackupForToday(backups)
	if err != nil {
		return err
	}
	if !existsBackupForToday {
		glog.V(2).Info("skip resume => no backup for today exists")
		return nil
	}
	sort.Sort(BackupByName(backups))
	err = renameLastBackupToIncomplete(h, backups)
	if err != nil {
		return err
	}
	err = removeCurrentSymlink(h)
	if err != nil {
		return err
	}
	err = symlinkBeforeLastToCurrent(h, backups)
	if err != nil {
		return err
	}
	return nil
}

func removeCurrentSymlink(h backup_host.Host) error {
	return os.Remove(current(h).Path())
}

func symlinkBeforeLastToCurrent(h backup_host.Host, backups []Backup) error {
	beforeLastBackup := backups[len(backups)-2]
	return os.Symlink(beforeLastBackup.Path(), current(h).Path())
}

func renameLastBackupToIncomplete(h backup_host.Host, backups []Backup) error {
	lastBackup := backups[len(backups)-1]
	return os.Rename(lastBackup.Path(), incomplete(h).Path())
}

func existsIncomplete(h backup_host.Host) (bool, error) {
	glog.V(2).Infof("existsIncomplete host: %s", h.Name())
	return incomplete(h).IsDir()
}

func incomplete(h backup_host.Host) Backup {
	return ByName(h, INCOMPLETE)
}

func current(h backup_host.Host) Backup {
	return ByName(h, CURRENT)
}

func existsBackupForToday(backups []Backup) (bool, error) {
	now := time.Now()
	backupsToday, err := getKeepToday(backups, now, backup_timeparser.New())
	if err != nil {
		return false, err
	}
	return len(backupsToday) > 0, nil
}

func KeepBackups(h backup_host.Host) ([]Backup, error) {
	backups, err := All(h)
	if err != nil {
		return nil, err
	}
	return getKeepBackups(backups, backup_timeparser.New())
}

func OldBackups(h backup_host.Host) ([]Backup, error) {
	allBackups, err := All(h)
	if err != nil {
		return nil, err
	}
	keepBackups, err := getKeepBackups(allBackups, backup_timeparser.New())
	if err != nil {
		return nil, err
	}
	keepMap := make(map[string]bool)
	for _, b := range keepBackups {
		keepMap[b.Name()] = true
	}
	var result []Backup
	for _, b := range allBackups {
		if !keepMap[b.Name()] {
			result = append(result, b)
		}
	}
	return result, nil
}
