package backup

import (
	"errors"
	"fmt"
	"os"

	"regexp"

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
	// falls bereits ein fertiges backup fuer heute existiert
	// symlink von current auf das vorletzte backup umbiegen
	// backup von heute in incomplete umbennen
	return nil
}

func ExistsIncomplete(h host.Host) (bool, error) {
	//	i := ByName(h, INCOMPLETE)
	return false, nil
}
