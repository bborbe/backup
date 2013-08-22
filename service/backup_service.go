package service

import "github.com/bborbe/backup/dto"

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
	return nil, nil
}

func (s *backupService) ListBackups(host dto.Host) ([]dto.Backup, error) {
	return nil, nil
}

func (s *backupService) GetLatestBackup(host dto.Host) (dto.Backup, error) {
	return nil, nil
}
