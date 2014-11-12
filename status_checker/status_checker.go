package status_checker

import (
	"github.com/bborbe/backup/dto"
	"github.com/bborbe/backup/service"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type StatusChecker interface {
	Check() ([]dto.Status, error)
}

type statusChecker struct {
	backupService service.BackupService
}

func NewStatusChecker(backupService service.BackupService) StatusChecker {
	s := new(statusChecker)
	s.backupService = backupService
	return s
}

func (s *statusChecker) Check() ([]dto.Status, error) {
	hosts, err := s.backupService.ListHosts()
	if err != nil {
		logger.Debugf("list hosts failed: %v", err)
		return nil, err
	}
	return createStatusDtoForHosts(s.backupService, hosts)
}

func createStatusDtoForHosts(backupService service.BackupService, hosts []dto.Host) ([]dto.Status, error) {
	result := make([]dto.Status, 0)
	for _, host := range hosts {
		status, err := createStatusDtoForHost(backupService, host)
		if err != nil {
			return nil, err
		}
		result = append(result, status)
	}
	return result, nil
}

func createStatusDtoForHost(backupService service.BackupService, host dto.Host) (dto.Status, error) {
	backup, err := backupService.GetLatestBackup(host)
	if err != nil {
		logger.Debugf("get latest backup failed: %v", err)
		return nil, err
	}
	if backup == nil {
		logger.Debugf("no backup for host %s found", host.GetName())
		return createStatusDto(host, nil, false), nil
	}

	logger.Debugf("host: %s backup: %s", host.GetName(), backup.GetName())
	return createStatusDto(host, backup, true), nil
}

func createStatusDto(host dto.Host, backup dto.Backup, status bool) dto.Status {
	statusDto := dto.NewStatus()
	statusDto.SetStatus(status)
	statusDto.SetHost(host.GetName())
	if backup != nil {
		statusDto.SetLatestBackup(backup.GetName())
	}
	return statusDto
}
