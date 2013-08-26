package mock

import "github.com/bborbe/backup/dto"

type backupServiceMock struct {
	listBackupsDtos      []dto.Backup
	listBackupsError     error
	listOldBackupsDtos   []dto.Backup
	listOldBackupsError  error
	listKeepBackupsDtos  []dto.Backup
	listKeepBackupsError error
	listHostsDtos        []dto.Host
	listHostsError       error
	latestBackup         dto.Backup
	latestBackupError    error
	cleanupErr           error
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
	return s.latestBackup, s.latestBackupError
}

func (s *backupServiceMock) SetLatestBackup(backup dto.Backup, err error) {
	s.latestBackup = backup
	s.latestBackupError = err
}

func CreateBackup(name string) dto.Backup {
	b := dto.NewBackup()
	b.SetName(name)
	return b
}

func CreateHost(name string) dto.Host {
	b := dto.NewHost()
	b.SetName(name)
	return b
}

func (s *backupServiceMock) ListOldBackups(host dto.Host) ([]dto.Backup, error) {
	return s.listOldBackupsDtos, s.listOldBackupsError
}

func (s *backupServiceMock) SetListOldBackups(backups []dto.Backup, err error) {
	s.listOldBackupsDtos = backups
	s.listOldBackupsError = err
}

func (s *backupServiceMock) Cleanup(host dto.Host) error {
	return s.cleanupErr
}

func (s *backupServiceMock) SetCleanup(err error) {
	s.cleanupErr = err
}

func (s *backupServiceMock) ListKeepBackups(host dto.Host) ([]dto.Backup, error) {
	return s.listKeepBackupsDtos, s.listKeepBackupsError
}

func (s *backupServiceMock) SetListKeepBackups(backups []dto.Backup, err error) {
	s.listKeepBackupsDtos = backups
	s.listKeepBackupsError = err
}
