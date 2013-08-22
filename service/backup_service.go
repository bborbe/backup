package service

import "github.com/bborbe/backup/object"

type BackupService interface {
	ListHosts() ([]object.Host, error)
}

type backupService struct {
	rootdir string
}

func NewBackupService(rootdir string) *backupService {
	s := new(backupService)
	s.rootdir = rootdir
	return s
}

func (s *backupService) ListHosts() ([]object.Host, error) {
	return nil, nil
}
