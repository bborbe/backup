package service

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"

	"github.com/bborbe/backup/dto"
	"github.com/bborbe/backup/host"
	"github.com/bborbe/backup/keep"
	"github.com/bborbe/backup/rootdir"
	"github.com/bborbe/log"
)

type BackupService interface {
	GetHost(host string) (dto.Host, error)
	ListHosts() ([]dto.Host, error)
	ListBackups(host dto.Host) ([]dto.Backup, error)
	ListOldBackups(host dto.Host) ([]dto.Backup, error)
	ListKeepBackups(host dto.Host) ([]dto.Backup, error)
	GetLatestBackup(host dto.Host) (dto.Backup, error)
	Cleanup(host dto.Host) error
	Resume(host dto.Host) error
}

type backupService struct {
	rootdir rootdir.Rootdir
}

var logger = log.DefaultLogger

func NewBackupService(rootdirectory string) *backupService {
	s := new(backupService)
	s.rootdir = rootdir.New(rootdirectory)
	return s
}

func (s *backupService) Resume(host dto.Host) error {
	return nil
}

func (s *backupService) ListHosts() ([]dto.Host, error) {
	hosts, err := host.All(s.rootdir)
	if err != nil {
		return nil, err
	}
	return s.createHosts(hosts)
}

func (s *backupService) createHosts(hosts []host.Host) ([]dto.Host, error) {
	result := []dto.Host{}
	for _, h := range hosts {
		dir := h.Path()
		isDir, err := isDir(dir)
		if err != nil {
			logger.Debugf("is dir failed: %v", err)
			return nil, err
		}
		if isDir {
			result = append(result, dto.CreateHost(h.Name()))
		} else {
			logger.Debugf("createHost for %s failed, is not a directory", h)
		}
	}
	return result, nil
}

func isDir(dir string) (bool, error) {
	file, err := os.Open(dir)
	if err != nil {
		logger.Debugf("open dir %s failed: %v", dir, err)
		return false, nil
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		return false, err
	}
	return fileinfo.IsDir(), nil
}

func (s *backupService) ListBackups(h dto.Host) ([]dto.Backup, error) {
	if h == nil {
		return nil, errors.New("parameter host missing")
	}
	dir := fmt.Sprintf("%s%c%s", s.rootdir.Path(), os.PathSeparator, h.GetName())
	file, err := os.Open(dir)
	if err != nil {
		logger.Debugf("open dir failed: %v", err)
		return nil, err
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		logger.Debugf("file stat failed: %v", err)
		return nil, err
	}
	if !fileinfo.IsDir() {
		return nil, fmt.Errorf("dir %s is not a directory", dir)
	}
	dirnames, err := file.Readdirnames(0)
	if err != nil {
		logger.Debugf("read dir names failed: %v", err)
		return nil, err
	}
	var names []string
	for _, name := range dirnames {
		if validBackupName(name) {
			names = append(names, name)
		}
	}
	backups := dto.CreateBackups(names)
	return backups, nil
}

func validBackupName(name string) bool {
	re := regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}")
	return re.MatchString(name)
}

func (s *backupService) GetHost(host string) (dto.Host, error) {
	dir := fmt.Sprintf("%s%c%s", s.rootdir.Path(), os.PathSeparator, host)
	file, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer file.Close()
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
		logger.Debugf("list backups failed: %v", err)
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
	backups, err := s.ListBackups(host)
	if err != nil {
		return nil, err
	}
	keepBackups, err := keep.GetKeepBackups(backups)
	if err != nil {
		return nil, err
	}
	keepMap := make(map[string]bool)
	for _, b := range keepBackups {
		keepMap[b.GetName()] = true
	}
	var result []dto.Backup
	for _, backup := range backups {
		if !keepMap[backup.GetName()] {
			result = append(result, backup)
		}
	}
	return result, nil
}

func (s *backupService) Cleanup(hostDto dto.Host) error {
	if hostDto == nil {
		return errors.New("parameter host missing")
	}
	backups, err := s.ListOldBackups(hostDto)
	if err != nil {
		return err
	}
	logger.Debugf("found %d backup to delete for host %s", len(backups), hostDto.GetName())
	for _, backup := range backups {
		dir := fmt.Sprintf("%s%c%s%c%s", s.rootdir.Path(), os.PathSeparator, hostDto.GetName(), os.PathSeparator, backup.GetName())
		logger.Infof("delete %s started", dir)
		os.RemoveAll(dir)
		logger.Infof("delete %s finished", dir)
	}
	return nil
}

func (s *backupService) ListKeepBackups(host dto.Host) ([]dto.Backup, error) {
	backups, err := s.ListBackups(host)
	if err != nil {
		return nil, err
	}
	return keep.GetKeepBackups(backups)
}
