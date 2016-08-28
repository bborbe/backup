package backup

import (
	"sort"
	"strconv"
	"time"

	backup_timeparser "github.com/bborbe/backup/timeparser"
	"github.com/golang/glog"
)

func latestBackup(backups []Backup) Backup {
	if backups != nil && len(backups) > 0 {
		sort.Sort(BackupByName(backups))
		return backups[len(backups)-1]
	}
	return nil
}

func getKeepToday(backups []Backup, now time.Time, timeParser backup_timeparser.TimeParser) ([]Backup, error) {
	var result []Backup
	for _, b := range backups {
		t, err := timeParser.TimeByName(b.Name())
		if err != nil {
			return nil, err
		}
		if isToday(t, now) {
			result = append(result, b)
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

func getKeepDay(backups []Backup, now time.Time, timeParser backup_timeparser.TimeParser) ([]Backup, error) {
	sort.Sort(BackupByName(backups))
	var result []Backup

	var lastYear int = -1
	var lastMonth time.Month = -1
	var lastDay int = -1

	for _, b := range backups {
		t, err := timeParser.TimeByName(b.Name())
		if err != nil {
			return nil, err
		}
		if ageLessThanDays(t, now, 7) && (t.Year() != lastYear || t.Month() != lastMonth || t.Day() != lastDay) {
			result = append(result, b)
			lastYear = t.Year()
			lastMonth = t.Month()
			lastDay = t.Day()
		}
	}
	return result, nil
}

func getKeepBackups(backups []Backup, timeParser backup_timeparser.TimeParser) ([]Backup, error) {
	keep := make(map[string]Backup, 0)
	now := time.Now()
	// keep all backups from today
	{
		bs, err := getKeepToday(backups, now, timeParser)
		if err != nil {
			return nil, err
		}
		for _, b := range bs {
			keep[b.Name()] = b
		}
	}
	// keep first backup per day if age <= 7 days
	{
		bs, err := getKeepDay(backups, now, timeParser)
		if err != nil {
			return nil, err
		}
		for _, b := range bs {
			keep[b.Name()] = b
		}
	}
	// keep first backup per week if age <= 28 days
	{
		bs, err := getKeepWeek(backups, now, timeParser)
		if err != nil {
			return nil, err
		}
		for _, b := range bs {
			keep[b.Name()] = b
		}
	}
	// keep first backup per month
	{
		bs, err := getKeepMonth(backups)
		if err != nil {
			return nil, err
		}
		for _, b := range bs {
			keep[b.Name()] = b
		}
	}
	// keep latest backup / current
	{
		b := latestBackup(backups)
		if b != nil {
			keep[b.Name()] = b
		}
	}

	var result []Backup
	for _, backup := range keep {
		result = append(result, backup)
	}
	return result, nil
}

func getKeepMonth(backups []Backup) ([]Backup, error) {
	sort.Sort(BackupByName(backups))
	var lastYear int64 = -1
	var lastMonth int64 = -1
	var result []Backup
	for _, b := range backups {
		name := b.Name()
		year, err := strconv.ParseInt(name[0:4], 10, 64)
		if err != nil {
			return nil, err
		}
		month, err := strconv.ParseInt(name[5:7], 10, 64)
		if err != nil {
			return nil, err
		}
		glog.V(4).Infof("year %d month %d", year, month)

		if year != lastYear || month != lastMonth {
			result = append(result, b)
		}
		lastYear = year
		lastMonth = month
	}
	return result, nil
}

func getKeepWeek(backups []Backup, now time.Time, timeParser backup_timeparser.TimeParser) ([]Backup, error) {
	sort.Sort(BackupByName(backups))
	var result []Backup
	var lastWeek int = -1
	for _, b := range backups {
		t, err := timeParser.TimeByName(b.Name())
		if err != nil {
			return nil, err
		}
		_, week := t.ISOWeek()
		if ageLessThanDays(t, now, 40) && week != lastWeek {
			result = append(result, b)
			lastWeek = week
		}
	}
	return result, nil
}
