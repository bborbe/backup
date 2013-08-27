package service

import (
	"errors"
	"fmt"
	"github.com/bborbe/backup/dto"
	"github.com/bborbe/backup/util"
	"github.com/bborbe/log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

type BackupService interface {
	GetHost(host string) (dto.Host, error)
	ListHosts() ([]dto.Host, error)
	ListBackups(host dto.Host) ([]dto.Backup, error)
	ListOldBackups(host dto.Host) ([]dto.Backup, error)
	ListKeepBackups(host dto.Host) ([]dto.Backup, error)
	GetLatestBackup(host dto.Host) (dto.Backup, error)
	Cleanup(host dto.Host) error
}

type backupService struct {
	rootdir string
}

var logger = log.DefaultLogger

func NewBackupService(rootdir string) *backupService {
	s := new(backupService)
	s.rootdir = rootdir
	return s
}

func (s *backupService) ListHosts() ([]dto.Host, error) {
	file, err := os.Open(s.rootdir)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if !fileinfo.IsDir() {
		return nil, fmt.Errorf("rootdir %s is not a directory", s.rootdir)
	}
	names, err := file.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	return s.createHosts(names)
}

func (s *backupService) createHosts(hosts []string) ([]dto.Host, error) {
	var result []dto.Host
	for _, host := range hosts {
		dir := fmt.Sprintf("%s%c%s", s.rootdir, os.PathSeparator, host)
		isDir, err := isDir(dir)
		if err != nil {
			return nil, err
		}
		if isDir {
			result = append(result, createHost(host))
		} else {
			logger.Warnf("createHost for %s failed, is not a directory", host)
		}
	}
	return result, nil
}

func isDir(dir string) (bool, error) {
	file, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		return false, err
	}
	return fileinfo.IsDir(), nil
}

func createHost(host string) dto.Host {
	h := dto.NewHost()
	h.SetName(host)
	return h
}

func (s *backupService) ListBackups(host dto.Host) ([]dto.Backup, error) {
	if host == nil {
		return nil, errors.New("parameter host missing")
	}
	dir := fmt.Sprintf("%s%c%s", s.rootdir, os.PathSeparator, host.GetName())
	file, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if !fileinfo.IsDir() {
		return nil, fmt.Errorf("dir %s is not a directory", dir)
	}
	dirnames, err := file.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, name := range dirnames {
		if validBackupName(name) {
			names = append(names, name)
		}
	}
	backups := createBackups(names)
	return backups, nil
}

func validBackupName(name string) bool {
	re := regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}")
	return re.MatchString(name)
}

func createBackups(backups []string) []dto.Backup {
	result := make([]dto.Backup, len(backups))
	for i, backup := range backups {
		result[i] = createBackup(backup)
	}
	return result
}

func createBackup(backup string) dto.Backup {
	h := dto.NewBackup()
	h.SetName(backup)
	return h
}

func (s *backupService) GetHost(host string) (dto.Host, error) {
	dir := fmt.Sprintf("%s%c%s", s.rootdir, os.PathSeparator, host)
	file, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if !fileinfo.IsDir() {
		return nil, fmt.Errorf("dir %s is not a directory", dir)
	}
	h := dto.NewHost()
	h.SetName(host)
	return h, nil
}

func (s *backupService) GetLatestBackup(host dto.Host) (dto.Backup, error) {
	list, err := s.ListBackups(host)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	var names []string
	backups := make(map[string]dto.Backup, 0)
	for _, backup := range list {
		backups[backup.GetName()] = backup
		names = append(names, backup.GetName())
	}
	sort.Strings(names)
	return backups[names[len(names) - 1]], nil
}

func (s *backupService) ListOldBackups(host dto.Host) ([]dto.Backup, error) {
	backups, err := s.ListBackups(host)
	if err != nil {
		return nil, err
	}
	keep, err := getKeepBackups(backups)
	if err != nil {
		return nil, err
	}
	keepMap := make(map[string]bool)
	for _, b := range keep {
		keepMap[b.GetName()] = true
	}
	var result []dto.Backup
	for _, backup := range backups {
		if !keepMap[backup.GetName()] {
			result = append(result, backup)
		}
	}
	return result, nil
}

func getKeepBackups(backups []dto.Backup) ([]dto.Backup, error) {
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

	var result []dto.Backup
	for _, backup := range keep {
		result = append(result, backup)
	}
	return result, nil
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

func getKeepDay(backups []dto.Backup, now time.Time) ([]dto.Backup, error) {
	sort.Sort(util.BackupByDate(backups))
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

func ageLessThanDays(t time.Time, now time.Time, days int64) bool {
	diff := now.Unix() - t.Unix()
	return diff <= 24*60*60*days
}

func getKeepMonth(backups []dto.Backup) ([]dto.Backup, error) {
	sort.Sort(util.BackupByDate(backups))
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
		logger.Debugf("year %d month %d", year, month)

		if year != lastYear || month != lastMonth {
			result = append(result, backup)
		}
		lastYear = year
		lastMonth = month
	}
	return result, nil
}

func (s *backupService) Cleanup(host dto.Host) error {
	if host == nil {
		return errors.New("parameter host missing")
	}
	backups, err := s.ListOldBackups(host)
	if err != nil {
		return err
	}
	logger.Debugf("found %d backup to delete for host %s", len(backups), host.GetName())
	for _, backup := range backups {
		dir := fmt.Sprintf("%s%c%s%c%s", s.rootdir, os.PathSeparator, host.GetName(), os.PathSeparator, backup.GetName())
		logger.Debugf("delete %s started", dir)
		//os.RemoveAll(dir)
		logger.Debugf("delete %s finished", dir)
	}
	return nil
}

func (s *backupService) ListKeepBackups(host dto.Host) ([]dto.Backup, error) {
	backups, err := s.ListBackups(host)
	if err != nil {
		return nil, err
	}
	return getKeepBackups(backups)
}

func getTimeByName(backupName string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05", backupName)
}

func getKeepWeek(backups []dto.Backup, now time.Time) ([]dto.Backup, error) {
	sort.Sort(util.BackupByDate(backups))
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
