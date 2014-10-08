package keep

import (
	"sort"
	"strconv"
	"time"

	"github.com/bborbe/backup/dto"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

func getTimeByName(backupName string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05", backupName)
}

func latestBackup(backups []dto.Backup) dto.Backup {
	if backups != nil && len(backups) > 0 {
		sort.Sort(dto.BackupByName(backups))
		return backups[len(backups)-1]
	}
	return nil
}

func getKeepToday(backups []dto.Backup, now time.Time) ([]dto.Backup, error) {
	var result []dto.Backup
	for _, backup := range backups {
		t, err := getTimeByName(backup.GetName())
		if err != nil {
			return nil, err
		}
		if isToday(t, now) {
			result = append(result, backup)
		}
	}
	return result, nil
}

func isToday(t time.Time, now time.Time) bool {
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

func ageLessThanDays(t time.Time, now time.Time, days int64) bool {
	diff := now.Unix() - t.Unix()
	return diff <= 24*60*60*days
}

func getKeepDay(backups []dto.Backup, now time.Time) ([]dto.Backup, error) {
	sort.Sort(dto.BackupByName(backups))
	var result []dto.Backup

	var lastYear int = -1
	var lastMonth time.Month = -1
	var lastDay int = -1

	for _, backup := range backups {
		t, err := getTimeByName(backup.GetName())
		if err != nil {
			return nil, err
		}
		if ageLessThanDays(t, now, 7) && (t.Year() != lastYear || t.Month() != lastMonth || t.Day() != lastDay) {
			result = append(result, backup)
			lastYear = t.Year()
			lastMonth = t.Month()
			lastDay = t.Day()
		}
	}
	return result, nil
}

func GetKeepBackups(backups []dto.Backup) ([]dto.Backup, error) {
	keep := make(map[string]dto.Backup, 0)
	now := time.Now()
	// keep all backups from today
	{
		b, err := getKeepToday(backups, now)
		if err != nil {
			return nil, err
		}
		for _, backup := range b {
			keep[backup.GetName()] = backup
		}
	}
	// keep first backup per day if age <= 7 days
	{
		b, err := getKeepDay(backups, now)
		if err != nil {
			return nil, err
		}
		for _, backup := range b {
			keep[backup.GetName()] = backup
		}
	}
	// keep first backup per week if age <= 28 days
	{
		b, err := getKeepWeek(backups, now)
		if err != nil {
			return nil, err
		}
		for _, backup := range b {
			keep[backup.GetName()] = backup
		}
	}
	// keep first backup per month
	{
		b, err := getKeepMonth(backups)
		if err != nil {
			return nil, err
		}
		for _, backup := range b {
			keep[backup.GetName()] = backup
		}
	}
	// keep latest backup / current
	{
		backup := latestBackup(backups)
		if backup != nil {
			keep[backup.GetName()] = backup
		}
	}

	var result []dto.Backup
	for _, backup := range keep {
		result = append(result, backup)
	}
	return result, nil
}

func getKeepMonth(backups []dto.Backup) ([]dto.Backup, error) {
	sort.Sort(dto.BackupByName(backups))
	var lastYear int64 = -1
	var lastMonth int64 = -1
	var result []dto.Backup
	for _, backup := range backups {
		name := backup.GetName()
		year, err := strconv.ParseInt(name[0:4], 10, 64)
		if err != nil {
			return nil, err
		}
		month, err := strconv.ParseInt(name[5:7], 10, 64)
		if err != nil {
			return nil, err
		}
		logger.Tracef("year %d month %d", year, month)

		if year != lastYear || month != lastMonth {
			result = append(result, backup)
		}
		lastYear = year
		lastMonth = month
	}
	return result, nil
}

func getKeepWeek(backups []dto.Backup, now time.Time) ([]dto.Backup, error) {
	sort.Sort(dto.BackupByName(backups))
	var result []dto.Backup
	var lastWeek int = -1
	for _, backup := range backups {
		t, err := getTimeByName(backup.GetName())
		if err != nil {
			return nil, err
		}
		_, week := t.ISOWeek()
		if ageLessThanDays(t, now, 40) && week != lastWeek {
			result = append(result, backup)
			lastWeek = week
		}
	}
	return result, nil
}
