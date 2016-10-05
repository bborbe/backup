package service

import (
	"errors"
	"fmt"
	"sort"

	backup_backup "github.com/bborbe/backup/backup"
	backup_dto "github.com/bborbe/backup/dto"
	backup_host "github.com/bborbe/backup/host"
	backup_rootdir "github.com/bborbe/backup/rootdir"
	io_util "github.com/bborbe/io/util"
	"github.com/golang/glog"
)

type BackupService interface {
	GetHost(host string) (backup_dto.Host, error)
	ListHosts() ([]backup_dto.Host, error)
	ListBackups(host backup_dto.Host) ([]backup_dto.Backup, error)
	ListOldBackups(host backup_dto.Host) ([]backup_dto.Backup, error)
	ListKeepBackups(host backup_dto.Host) ([]backup_dto.Backup, error)
	GetLatestBackup(host backup_dto.Host) (backup_dto.Backup, error)
	Cleanup(host backup_dto.Host) error
	Resume(host backup_dto.Host) error
}

type backupService struct {
	rootdir backup_rootdir.Rootdir
}

func NewBackupService(rootdirectory string) *backupService {
	s := new(backupService)
	s.rootdir = backup_rootdir.ByName(rootdirectory)
	return s
}

func (s *backupService) Resume(hostDto backup_dto.Host) error {
	h := backup_host.ByName(s.rootdir, hostDto.GetName())
	return backup_backup.Resume(h)
}

func (s *backupService) ListHosts() ([]backup_dto.Host, error) {
	hosts, err := backup_host.All(s.rootdir)
	if err != nil {
		return nil, err
	}
	hostDtos := make([]backup_dto.Host, len(hosts))
	for i := 0; i < len(hosts); i++ {
		hostDtos[i] = backup_dto.CreateHost(hosts[i].Name())
	}
	return hostDtos, nil
}

func (s *backupService) ListBackups(hostDto backup_dto.Host) ([]backup_dto.Backup, error) {
	if hostDto == nil {
		return nil, errors.New("parameter host missing")
	}
	h := backup_host.ByName(s.rootdir, hostDto.GetName())
	backups, err := backup_backup.All(h)
	if err != nil {
		return nil, err
	}
	return convertBackupsToBackupDtos(backups), nil
}

func convertBackupToBackupDto(b backup_backup.Backup) backup_dto.Backup {
	return backup_dto.CreateBackup(b.Name())
}

func convertBackupsToBackupDtos(backups []backup_backup.Backup) []backup_dto.Backup {
	backupDtos := make([]backup_dto.Backup, len(backups))
	for i := 0; i < len(backups); i++ {
		backupDtos[i] = convertBackupToBackupDto(backups[i])
	}
	return backupDtos
}

func (s *backupService) GetHost(hostname string) (backup_dto.Host, error) {
	h := backup_host.ByName(s.rootdir, hostname)
	isDir, err := io_util.IsDirectory(h.Path())
	if err != nil {
		return nil, err
	}
	if !isDir {
		return nil, fmt.Errorf("dir %s is not a directory", h.Path())
	}
	hostDto := backup_dto.NewHost()
	hostDto.SetName(hostname)
	return hostDto, nil
}

func (s *backupService) GetLatestBackup(hostDto backup_dto.Host) (backup_dto.Backup, error) {
	if hostDto == nil {
		return nil, errors.New("parameter host missing")
	}
	h := backup_host.ByName(s.rootdir, hostDto.GetName())
	list, err := backup_backup.All(h)
	if err != nil {
		glog.V(2).Infof("list backups failed: %v", err)
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	var names []string
	backups := make(map[string]backup_backup.Backup, 0)
	for _, b := range list {
		backups[b.Name()] = b
		names = append(names, b.Name())
	}
	sort.Strings(names)
	return convertBackupToBackupDto(backups[names[len(names)-1]]), nil
}

func (s *backupService) ListOldBackups(hostDto backup_dto.Host) ([]backup_dto.Backup, error) {
	if hostDto == nil {
		return nil, errors.New("parameter host missing")
	}
	h := backup_host.ByName(s.rootdir, hostDto.GetName())
	result, err := backup_backup.OldBackups(h)
	if err != nil {
		return nil, err
	}
	return convertBackupsToBackupDtos(result), nil
}

func (s *backupService) Cleanup(hostDto backup_dto.Host) error {
	if hostDto == nil {
		return errors.New("parameter host missing")
	}
	backups, err := s.ListOldBackups(hostDto)
	if err != nil {
		return err
	}
	h := backup_host.ByName(s.rootdir, hostDto.GetName())
	glog.V(2).Infof("found %d backup to delete for host %s", len(backups), hostDto.GetName())
	for _, backupDto := range backups {
		b := backup_backup.ByName(h, backupDto.GetName())
		if err := b.Delete(); err != nil {
			return err
		}
		glog.V(1).Infof("backup %s deleted", b.Path())
	}
	return nil
}

func (s *backupService) ListKeepBackups(hostDto backup_dto.Host) ([]backup_dto.Backup, error) {
	if hostDto == nil {
		return nil, errors.New("parameter host missing")
	}
	h := backup_host.ByName(s.rootdir, hostDto.GetName())
	keepBackups, err := backup_backup.KeepBackups(h)
	if err != nil {
		return nil, err
	}
	return convertBackupsToBackupDtos(keepBackups), nil
}
