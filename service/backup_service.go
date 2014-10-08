package service

import (
	"errors"
	"fmt"
	"sort"

	"github.com/bborbe/backup/backup"
	"github.com/bborbe/backup/dto"
	"github.com/bborbe/backup/fileutil"
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
	hostDtos := make([]dto.Host, len(hosts))
	for i := 0; i < len(hosts); i++ {
		hostDtos[i] = dto.CreateHost(hosts[i].Name())
	}
	return hostDtos, nil
}

func (s *backupService) ListBackups(hostDto dto.Host) ([]dto.Backup, error) {
	if hostDto == nil {
		return nil, errors.New("parameter host missing")
	}
	h := host.ByName(s.rootdir, hostDto.GetName())
	backups, err := backup.All(h)
	if err != nil {
		return nil, err
	}
	backupDtos := make([]dto.Backup, len(backups))
	for i := 0; i < len(backups); i++ {
		backupDtos[i] = dto.CreateBackup(backups[i].Name())
	}
	return backupDtos, nil
}

func (s *backupService) GetHost(hostname string) (dto.Host, error) {
	h := host.ByName(s.rootdir, hostname)
	isDir, err := fileutil.IsDir(h.Path())
	if err != nil {
		return nil, err
	}
	if !isDir {
		return nil, fmt.Errorf("dir %s is not a directory", h.Path())
	}
	hostDto := dto.NewHost()
	hostDto.SetName(hostname)
	return hostDto, nil
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
	for _, b := range list {
		backups[b.GetName()] = b
		names = append(names, b.GetName())
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
	for _, b := range backups {
		if !keepMap[b.GetName()] {
			result = append(result, b)
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

	h := host.ByName(s.rootdir, hostDto.GetName())

	logger.Debugf("found %d backup to delete for host %s", len(backups), hostDto.GetName())
	for _, backupDto := range backups {
		b := backup.ByName(h, backupDto.GetName())
		logger.Infof("delete %s started", b.Path())
		b.Delete()
		if err != nil {
			return err
		}
		logger.Infof("delete %s finished", b.Path())
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
