package service

import (
	"errors"
	"fmt"
	"github.com/bborbe/backup/dto"
	"os"
	"regexp"
	"sort"
)

type BackupService interface {
	GetHost(host string) (dto.Host, error)
	ListHosts() ([]dto.Host, error)
	ListBackups(host dto.Host) ([]dto.Backup, error)
	ListOldBackups(host dto.Host) ([]dto.Backup, error)
	GetLatestBackup(host dto.Host) (dto.Backup, error)
	Cleanup(host dto.Host) error
}

type backupService struct {
	rootdir string
}

func NewBackupService(rootdir string) *backupService {
	s := new(backupService)
	s.rootdir = rootdir
	return s
}

func (s *backupService) ListHosts() ([]dto.Host, error) {
	file, err := os.Open(s.rootdir)
	if err != nil {
		return nil, err
	}
	fileinfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if !fileinfo.IsDir() {
		return nil, fmt.Errorf("rootdir %s is not a directory", s.rootdir)
	}
	names, err := file.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	return createHosts(names), nil
}

func createHosts(hosts []string) []dto.Host {
	result := make([]dto.Host, len(hosts))
	for i, host := range hosts {
		result[i] = createHost(host)
	}
	return result
}

func createHost(host string) dto.Host {
	h := dto.NewHost()
	h.SetName(host)
	return h
}

func (s *backupService) ListBackups(host dto.Host) ([]dto.Backup, error) {
	if host == nil {
		return nil, errors.New("parameter host missing")
	}
	dir := fmt.Sprintf("%s%c%s", s.rootdir, os.PathSeparator, host.GetName())
	file, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	fileinfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if !fileinfo.IsDir() {
		return nil, fmt.Errorf("dir %s is not a directory", dir)
	}
	dirnames, err := file.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, name := range dirnames {
		if validBackupName(name) {
			names = append(names, name)
		}
	}

	return createBackups(names), nil
}

func validBackupName(name string) bool {
	re := regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}")
	return re.MatchString(name)
}

func createBackups(backups []string) []dto.Backup {
	result := make([]dto.Backup, len(backups))
	for i, backup := range backups {
		result[i] = createBackup(backup)
	}
	return result
}

func createBackup(backup string) dto.Backup {
	h := dto.NewBackup()
	h.SetName(backup)
	return h
}

func (s *backupService) GetHost(host string) (dto.Host, error) {
	dir := fmt.Sprintf("%s%c%s", s.rootdir, os.PathSeparator, host)
	file, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	fileinfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if !fileinfo.IsDir() {
		return nil, fmt.Errorf("dir %s is not a directory", dir)
	}
	h := dto.NewHost()
	h.SetName(host)
	return h, nil
}

func (s *backupService) GetLatestBackup(host dto.Host) (dto.Backup, error) {
	list, err := s.ListBackups(host)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	var names []string
	backups := make(map[string]dto.Backup, 0)
	for _, backup := range list {
		backups[backup.GetName()] = backup
		names = append(names, backup.GetName())
	}
	sort.Strings(names)
	return backups[names[len(names)-1]], nil
}

func (s *backupService) ListOldBackups(host dto.Host) ([]dto.Backup, error) {
	if host == nil {
		return nil, errors.New("parameter host missing")
	}
	return nil, nil
}

func (s *backupService) Cleanup(host dto.Host) error {
	if host == nil {
		return errors.New("parameter host missing")
	}
	return nil
}
