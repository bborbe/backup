package backup

import (
	"errors"
	"fmt"
	"os"

	"regexp"

	"sort"
	"time"

	"github.com/bborbe/backup/fileutil"
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
	Delete() error
	IsDir() (bool, error)
	ValidName() bool
}

const INCOMPLETE = "incomplete"
const CURRENT = "current"

var logger = log.DefaultLogger

func ByTime(h host.Host, t time.Time) Backup {
	return ByName(h, t.Format("2006-01-02T15:04:05"))
}

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
	if !fileutil.Exists(b.Path()) {
		return false, nil
	}
	return fileutil.IsDir(b.Path())
}

func KeepBackups(h host.Host) ([]Backup, error) {
	backups, err := All(h)
	if err != nil {
		return nil, err
	}
	return getKeepBackups(backups)
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

func Resume(h host.Host) error {
	existsIncomplete, err := existsIncomplete(h)
	if err != nil {
		return err
	}
	if existsIncomplete {
		logger.Debug("skip resume => incomplete dir exists")
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
		logger.Debug("skip resume => no backup for today exists")
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

func removeCurrentSymlink(h host.Host) error {
	return os.Remove(current(h).Path())
}

func symlinkBeforeLastToCurrent(h host.Host, backups []Backup) error {
	beforeLastBackup := backups[len(backups)-2]
	return os.Symlink(beforeLastBackup.Path(), current(h).Path())
}

func renameLastBackupToIncomplete(h host.Host, backups []Backup) error {
	lastBackup := backups[len(backups)-1]
	return os.Rename(lastBackup.Path(), incomplete(h).Path())
}

func existsIncomplete(h host.Host) (bool, error) {
	logger.Debugf("existsIncomplete host: %s", h.Name())
	return incomplete(h).IsDir()
}

func incomplete(h host.Host) Backup {
	return ByName(h, INCOMPLETE)
}

func current(h host.Host) Backup {
	return ByName(h, CURRENT)
}

func existsBackupForToday(backups []Backup) (bool, error) {
	now := time.Now()
	backupsToday, err := getKeepToday(backups, now)
	if err != nil {
		return false, err
	}
	return len(backupsToday) > 0, nil
}
