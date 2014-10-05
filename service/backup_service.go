package service

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/bborbe/backup/dto"
	"github.com/bborbe/backup/host"
	"github.com/bborbe/backup/rootdir"
	"github.com/bborbe/backup/util"
	"github.com/bborbe/log"
)

type BackupService interface {
	GetRootdir(rootdir string) (rootdir.Rootdir, error)
	GetHost(host string) (dto.Host, error)
	ListHosts() ([]dto.Host, error)
	ListBackups(host dto.Host) ([]dto.Backup, error)
	ListOldBackups(host dto.Host) ([]dto.Backup, error)
	ListKeepBackups(host dto.Host) ([]dto.Backup, error)
	GetLatestBackup(host dto.Host) (dto.Backup, error)
	Cleanup(host dto.Host) error
	Resume(host dto.Host) error
}

type backupService struct {
	rootdir rootdir.Rootdir
}

var logger = log.DefaultLogger

func NewBackupService(rootdirectory string) *backupService {
	s := new(backupService)
	s.rootdir = rootdir.New(rootdirectory)
	return s
}

func (s *backupService) GetRootdir(dir string) (rootdir.Rootdir, error) {
	isDir, err := isDir(dir)
	if err != nil {
		return nil, err
	}
	if !isDir {
		return nil, fmt.Errorf("dir %s is not a directory", dir)
	}
	return rootdir.New(dir), nil
}

func (s *backupService) Resume(host dto.Host) error {
	return nil
}

func (s *backupService) ListHosts() ([]dto.Host, error) {
	hosts, err := host.All(s.rootdir)
	if err != nil {
		return nil, err
	}
	return s.createHosts(hosts)
}

func (s *backupService) createHosts(hosts []host.Host) ([]dto.Host, error) {
	result := []dto.Host{}
	for _, h := range hosts {
		dir := h.Path()
		isDir, err := isDir(dir)
		if err != nil {
			logger.Debugf("is dir failed: %v", err)
			return nil, err
		}
		if isDir {
			result = append(result, createHost(h.Name()))
		} else {
			logger.Debugf("createHost for %s failed, is not a directory", h)
		}
	}
	return result, nil
}

func isDir(dir string) (bool, error) {
	file, err := os.Open(dir)
	if err != nil {
		logger.Debugf("open dir %s failed: %v", dir, err)
		return false, nil
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

func (s *backupService) ListBackups(h dto.Host) ([]dto.Backup, error) {
	if h == nil {
		return nil, errors.New("parameter host missing")
	}
	dir := fmt.Sprintf("%s%c%s", s.rootdir.Path(), os.PathSeparator, h.GetName())
	file, err := os.Open(dir)
	if err != nil {
		logger.Debugf("open dir failed: %v", err)
		return nil, err
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		logger.Debugf("file stat failed: %v", err)
		return nil, err
	}
	if !fileinfo.IsDir() {
		return nil, fmt.Errorf("dir %s is not a directory", dir)
	}
	dirnames, err := file.Readdirnames(0)
	if err != nil {
		logger.Debugf("read dir names failed: %v", err)
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
	dir := fmt.Sprintf("%s%c%s", s.rootdir.Path(), os.PathSeparator, host)
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
		logger.Debugf("list backups failed: %v", err)
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
	return backups[names[len(names)-1]], nil
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
		logger.Tracef("year %d month %d", year, month)

		if year != lastYear || month != lastMonth {
			result = append(result, backup)
		}
		lastYear = year
		lastMonth = month
	}
	return result, nil
}

func (s *backupService) Cleanup(hostDto dto.Host) error {
	if hostDto == nil {
		return errors.New("parameter host missing")
	}
	backups, err := s.ListOldBackups(hostDto)
	if err != nil {
		return err
	}
	logger.Debugf("found %d backup to delete for host %s", len(backups), hostDto.GetName())
	for _, backup := range backups {
		dir := fmt.Sprintf("%s%c%s%c%s", s.rootdir.Path(), os.PathSeparator, hostDto.GetName(), os.PathSeparator, backup.GetName())
		logger.Infof("delete %s started", dir)
		os.RemoveAll(dir)
		logger.Infof("delete %s finished", dir)
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

func latestBackup(backups []dto.Backup) dto.Backup {
	if backups != nil && len(backups) > 0 {
		sort.Sort(util.BackupByDate(backups))
		return backups[len(backups)-1]
	}
	return nil
}
