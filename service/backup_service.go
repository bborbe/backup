package service

import (
	"fmt"
	"github.com/bborbe/backup/dto"
	"os"
)

type BackupService interface {
	ListHosts() ([]dto.Host, error)
	ListBackups(host dto.Host) ([]dto.Backup, error)
	GetLatestBackup(host dto.Host) (dto.Backup, error)
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
	return nil, nil
}

func (s *backupService) GetLatestBackup(host dto.Host) (dto.Backup, error) {
	return nil, nil
}
