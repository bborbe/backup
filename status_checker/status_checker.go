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
	status := make([]dto.Status, 0)
	for _, host := range hosts {
		backup, err := backupService.GetLatestBackup(host)
		if err != nil {
			logger.Debugf("get latest backup failed: %v", err)
			return nil, err
		}
		if backup != nil {
			logger.Debugf("host: %s backup: %s", host.GetName(), backup.GetName())
			status = append(status, createStatusDto(host.GetName(), true))
		} else {
			logger.Debugf("no backup for host %s found", host.GetName())
			status = append(status, createStatusDto(host.GetName(), false))
		}
	}
	return status, nil
}

func createStatusDto(hostname string, status bool) dto.Status {
	statusDto := dto.NewStatus()
	statusDto.SetStatus(status)
	statusDto.SetHost(hostname)
	return statusDto
}
