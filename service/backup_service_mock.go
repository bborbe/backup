package service

import (
	backup_dto "github.com/bborbe/backup/dto"
)

type backupServiceMock struct {
	listBackupsDtos      []backup_dto.Backup
	listBackupsError     error
	listOldBackupsDtos   []backup_dto.Backup
	listOldBackupsError  error
	listKeepBackupsDtos  []backup_dto.Backup
	listKeepBackupsError error
	listHostsDtos        []backup_dto.Host
	listHostsError       error
	latestBackup         backup_dto.Backup
	latestBackupError    error
	cleanupErr           error
	resumeErr            error
	getHostErr           error
}

func NewBackupServiceMock() *backupServiceMock {
	return new(backupServiceMock)
}

func (s *backupServiceMock) Resume(host backup_dto.Host) error {
	return s.resumeErr
}

func (s *backupServiceMock) SetResume(err error) {
	s.resumeErr = err
}

func (s *backupServiceMock) GetHost(host string) (backup_dto.Host, error) {
	return nil, nil
}

func (s *backupServiceMock) SetGetHost(err error) {
	s.getHostErr = err
}

func (s *backupServiceMock) ListHosts() ([]backup_dto.Host, error) {
	return s.listHostsDtos, s.listHostsError
}

func (s *backupServiceMock) SetListHosts(dtos []backup_dto.Host, err error) {
	s.listHostsDtos = dtos
	s.listHostsError = err
}

func (s *backupServiceMock) ListBackups(host backup_dto.Host) ([]backup_dto.Backup, error) {
	return s.listBackupsDtos, s.listBackupsError
}

func (s *backupServiceMock) SetListBackups(backups []backup_dto.Backup, err error) {
	s.listBackupsDtos = backups
	s.listBackupsError = err
}

func (s *backupServiceMock) GetLatestBackup(host backup_dto.Host) (backup_dto.Backup, error) {
	return s.latestBackup, s.latestBackupError
}

func (s *backupServiceMock) SetLatestBackup(backup backup_dto.Backup, err error) {
	s.latestBackup = backup
	s.latestBackupError = err
}

func CreateBackup(name string) backup_dto.Backup {
	b := backup_dto.NewBackup()
	b.SetName(name)
	return b
}

func CreateHost(name string) backup_dto.Host {
	b := backup_dto.NewHost()
	b.SetName(name)
	return b
}

func (s *backupServiceMock) ListOldBackups(host backup_dto.Host) ([]backup_dto.Backup, error) {
	return s.listOldBackupsDtos, s.listOldBackupsError
}

func (s *backupServiceMock) SetListOldBackups(backups []backup_dto.Backup, err error) {
	s.listOldBackupsDtos = backups
	s.listOldBackupsError = err
}

func (s *backupServiceMock) Cleanup(host backup_dto.Host) error {
	return s.cleanupErr
}

func (s *backupServiceMock) SetCleanup(err error) {
	s.cleanupErr = err
}

func (s *backupServiceMock) ListKeepBackups(host backup_dto.Host) ([]backup_dto.Backup, error) {
	return s.listKeepBackupsDtos, s.listKeepBackupsError
}

func (s *backupServiceMock) SetListKeepBackups(backups []backup_dto.Backup, err error) {
	s.listKeepBackupsDtos = backups
	s.listKeepBackupsError = err
}
