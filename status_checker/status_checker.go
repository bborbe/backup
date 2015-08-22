package status_checker

import (
	"time"

	backup_dto "github.com/bborbe/backup/dto"
	backup_service "github.com/bborbe/backup/service"
	backup_timeparser "github.com/bborbe/backup/timeparser"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type StatusChecker interface {
	Check() ([]backup_dto.Status, error)
}

type statusChecker struct {
	backupService backup_service.BackupService
}

func NewStatusChecker(backupService backup_service.BackupService) StatusChecker {
	s := new(statusChecker)
	s.backupService = backupService
	return s
}

func (s *statusChecker) Check() ([]backup_dto.Status, error) {
	hosts, err := s.backupService.ListHosts()
	if err != nil {
		logger.Debugf("list hosts failed: %v", err)
		return nil, err
	}
	return createStatusDtoForHosts(s.backupService, hosts)
}

func createStatusDtoForHosts(backupService backup_service.BackupService, hosts []backup_dto.Host) ([]backup_dto.Status, error) {
	result := make([]backup_dto.Status, 0)
	for _, host := range hosts {
		status, err := createStatusDtoForHost(backupService, host)
		if err != nil {
			return nil, err
		}
		result = append(result, status)
	}
	return result, nil
}

func createStatusDtoForHost(backupService backup_service.BackupService, host backup_dto.Host) (backup_dto.Status, error) {
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
	result, err := backupIsInLastDays(backup, backup_timeparser.New(), time.Now())
	if err != nil {
		logger.Debug("parse backup failed")
		return nil, err
	}
	return createStatusDto(host, backup, result), nil
}

func backupIsInLastDays(backup backup_dto.Backup, timeParser backup_timeparser.TimeParser, now time.Time) (bool, error) {
	t, err := timeParser.TimeByName(backup.GetName())
	if err != nil {
		return false, err
	}
	return !now.After(t.Add(time.Duration(time.Hour * 24 * 7))), nil
}

func createStatusDto(host backup_dto.Host, backup backup_dto.Backup, status bool) backup_dto.Status {
	statusDto := backup_dto.NewStatus()
	statusDto.SetStatus(status)
	statusDto.SetHost(host.GetName())
	if backup != nil {
		statusDto.SetLatestBackup(backup.GetName())
	}
	return statusDto
}
