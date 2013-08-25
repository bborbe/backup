package mock

import "github.com/bborbe/backup/dto"

type backupServiceMock struct {
	listBackupsDtos  []dto.Backup
	listBackupsError error
	listHostsDtos    []dto.Host
	listHostsError   error
}

func NewBackupServiceMock() *backupServiceMock {
	return new(backupServiceMock)
}

func (s *backupServiceMock) GetHost(host string) (dto.Host, error) {
	return nil, nil
}

func (s *backupServiceMock) ListHosts() ([]dto.Host, error) {
	return s.listHostsDtos, s.listHostsError
}

func (s *backupServiceMock) SetListHosts(dtos []dto.Host, err error) {
	s.listHostsDtos = dtos
	s.listHostsError = err
}

func (s *backupServiceMock) ListBackups(host dto.Host) ([]dto.Backup, error) {
	return s.listBackupsDtos, s.listBackupsError
}

func (s *backupServiceMock) SetListBackups(backups []dto.Backup, err error) {
	s.listBackupsDtos = backups
	s.listBackupsError = err
}

func (s *backupServiceMock) GetLatestBackup(host dto.Host) (dto.Backup, error) {
	return nil, nil
}
