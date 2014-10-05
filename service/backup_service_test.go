package service

import (
	"fmt"
	"os"
	"testing"
	"time"
	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/dto"
)

const BACKUP_ROOT_DIR = "/tmp/backuproot"

func clearBackupRootDir(root string) error {
	return os.RemoveAll(root)
}

func createRootDir(root string) error {
	var fileMode os.FileMode
	fileMode = 0777
	return os.Mkdir(root, fileMode)
}

func createHostDir(root string, host string) error {
	var fileMode os.FileMode
	fileMode = 0777
	dir := fmt.Sprintf("%s%c%s", root, os.PathSeparator, host)
	logger.Debugf("create hostdir %s", dir)
	return os.Mkdir(dir, fileMode)
}

func createBackupDir(root string, host string, backup string) error {
	var fileMode os.FileMode
	fileMode = 0777
	dir := fmt.Sprintf("%s%c%s%c%s", root, os.PathSeparator, host, os.PathSeparator, backup)
	logger.Debugf("create backupdir %s", dir)
	return os.Mkdir(dir, fileMode)
}

func TestImplementsBackupService(t *testing.T) {
	service := NewBackupService(BACKUP_ROOT_DIR)
	var expected *BackupService
	err := AssertThat(service, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestListHosts(t *testing.T) {
	clearBackupRootDir(BACKUP_ROOT_DIR)
	service := NewBackupService(BACKUP_ROOT_DIR)
	{
		_, err := service.ListHosts()
		err = AssertThat(err, NotNilValue().Message("expect error, backup dir not existing"))
		if err != nil {
			t.Fatal(err)
		}
	}
	createRootDir(BACKUP_ROOT_DIR)
	{
		hosts, err := service.ListHosts()
		err = AssertThat(err, NilValue().Message("expect no error, backup dir existing"))
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(hosts), Is(0))
		if err != nil {
			t.Fatal(err)
		}
	}
	hostName := "firewall.example.com"
	createHostDir(BACKUP_ROOT_DIR, hostName)
	{
		hosts, err := service.ListHosts()
		err = AssertThat(err, NilValue().Message("expect no error, backup dir existing"))
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(hosts), Is(1))
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(hosts[0].GetName(), Is(hostName))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetHost(t *testing.T) {
	clearBackupRootDir(BACKUP_ROOT_DIR)
	service := NewBackupService(BACKUP_ROOT_DIR)
	hostName := "firewall.example.com"
	createRootDir(BACKUP_ROOT_DIR)
	rootdir, err := service.GetRootdir(BACKUP_ROOT_DIR)
	if err != nil {
		t.Fatal(err)
	}
	{
		_, err := service.GetHost(rootdir, hostName)
		err = AssertThat(err, NotNilValue().Message("expect error, backup dir not existing"))
		if err != nil {
			t.Fatal(err)
		}
	}
	createHostDir(BACKUP_ROOT_DIR, hostName)
	{
		host, err := service.GetHost(rootdir, hostName)
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(host, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestListBackups(t *testing.T) {
	clearBackupRootDir(BACKUP_ROOT_DIR)
	service := NewBackupService(BACKUP_ROOT_DIR)
	{
		_, err := service.ListBackups(nil)
		err = AssertThat(err, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
	}
	hostName := "firewall.example.com"
	createRootDir(BACKUP_ROOT_DIR)
	createHostDir(BACKUP_ROOT_DIR, hostName)
	rootDir, err := service.GetRootdir(BACKUP_ROOT_DIR)
	if err != nil {
		t.Fatal(err)
	}
	{
		host, err := service.GetHost(rootDir, hostName)
		if err != nil {
			t.Fatal(err)
		}
		backups, err := service.ListBackups(host)
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(backups, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(backups), Is(0))
		if err != nil {
			t.Fatal(err)
		}
	}
	backupName := "2013-07-28T00:24:52"
	createBackupDir(BACKUP_ROOT_DIR, hostName, backupName)
	{
		host, err := service.GetHost(rootDir, hostName)
		if err != nil {
			t.Fatal(err)
		}
		backups, err := service.ListBackups(host)
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(backups, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(backups), Is(1))
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(backups[0].GetName(), Is(backupName))
		if err != nil {
			t.Fatal(err)
		}
	}
	createBackupDir(BACKUP_ROOT_DIR, hostName, "incomplete")
	{
		host, err := service.GetHost(rootDir, hostName)
		if err != nil {
			t.Fatal(err)
		}
		backups, err := service.ListBackups(host)
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(backups, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(backups), Is(1))
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(backups[0].GetName(), Is(backupName))
		if err != nil {
			t.Fatal(err)
		}
	}
}

//GetLatestBackup(host dto.Host) (dto.Backup, error)
func TestGetLatestBackup(t *testing.T) {
	var backupName string
	clearBackupRootDir(BACKUP_ROOT_DIR)
	service := NewBackupService(BACKUP_ROOT_DIR)
	{
		_, err := service.GetLatestBackup(nil)
		err = AssertThat(err, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
	}
	createRootDir(BACKUP_ROOT_DIR)
	{
		_, err := service.GetLatestBackup(nil)
		err = AssertThat(err, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
	}
	rootDir, err := service.GetRootdir(BACKUP_ROOT_DIR)
	if err != nil {
		t.Fatal(err)
	}
	hostName := "firewall.example.com"
	createHostDir(BACKUP_ROOT_DIR, hostName)
	{
		host, err := service.GetHost(rootDir, hostName)
		if err != nil {
			t.Fatal(err)
		}
		backup, err := service.GetLatestBackup(host)
		err = AssertThat(err, NilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(backup, NilValue())
		if err != nil {
			t.Fatal(err)
		}
	}
	backupName = "2013-07-01T00:24:52"
	createBackupDir(BACKUP_ROOT_DIR, hostName, backupName)
	{
		host, err := service.GetHost(rootDir, hostName)
		if err != nil {
			t.Fatal(err)
		}
		backup, err := service.GetLatestBackup(host)
		err = AssertThat(err, NilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(backup, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(backup.GetName(), Is(backupName))
		if err != nil {
			t.Fatal(err)
		}
	}
	backupName = "2013-07-02T00:24:52"
	createBackupDir(BACKUP_ROOT_DIR, hostName, backupName)
	{
		host, err := service.GetHost(rootDir, hostName)
		if err != nil {
			t.Fatal(err)
		}
		backup, err := service.GetLatestBackup(host)
		err = AssertThat(err, NilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(backup, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(backup.GetName(), Is(backupName))
		if err != nil {
			t.Fatal(err)
		}
	}
}

//ListOldBackups(host dto.Host) ([]dto.Backup, error)
func TestListOldBackups(t *testing.T) {
	clearBackupRootDir(BACKUP_ROOT_DIR)
	service := NewBackupService(BACKUP_ROOT_DIR)
	{
		_, err := service.ListOldBackups(nil)
		err = AssertThat(err, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
	}
	createRootDir(BACKUP_ROOT_DIR)
	{
		_, err := service.ListOldBackups(nil)
		err = AssertThat(err, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
	}
}

//Cleanup() error
func TestCleanup(t *testing.T) {
	clearBackupRootDir(BACKUP_ROOT_DIR)
	service := NewBackupService(BACKUP_ROOT_DIR)
	{
		err := service.Cleanup(nil)
		err = AssertThat(err, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
	}
	createRootDir(BACKUP_ROOT_DIR)
	{
		err := service.Cleanup(nil)
		err = AssertThat(err, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestValidBackupName(t *testing.T) {
	var err error
	err = AssertThat(validBackupName("foo"), Is(false))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(validBackupName("2013-12-12T24:15:59"), Is(true))
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetKeepMonth(t *testing.T) {
	var err error
	var result []dto.Backup
	{
		backups := []dto.Backup{}
		result, err = getKeepMonth(backups)
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(result, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(result), Is(len(backups)))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		backups := []dto.Backup{
			createBackup("2013-12-12T24:15:59"),
		}
		result, err = getKeepMonth(backups)
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(result, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(result), Is(len(backups)))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		backups := []dto.Backup{
			createBackup("2013-12-12T24:15:59"),
			createBackup("2013-12-01T24:15:59"),
		}
		result, err = getKeepMonth(backups)
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(result, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(result), Is(1))
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat("2013-12-01T24:15:59", Is(result[0].GetName()))
	}
	{
		backups := []dto.Backup{
			createBackup("2012-11-12T24:15:59"),
			createBackup("2013-11-01T24:15:59"),
			createBackup("2013-05-28T24:15:59"),
			createBackup("2013-05-29T24:15:59"),
		}
		result, err = getKeepMonth(backups)
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(result, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(result), Is(3))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetTimeByName(t *testing.T) {
	{
		_, err := getTimeByName("")
		err = AssertThat(err, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		calcTime, err := getTimeByName("2013-07-01T00:24:52")
		err = AssertThat(err, NilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(calcTime.Year(), Is(2013))
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(calcTime.Month(), Is(time.July))
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(calcTime.Day(), Is(1))
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(calcTime.Hour(), Is(0))
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(calcTime.Minute(), Is(24))
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(calcTime.Second(), Is(52))
		if err != nil {
			t.Fatal(err)
		}
	}

}

func TestGetKeepToday(t *testing.T) {
	var result []dto.Backup
	now, err := getTimeByName("2013-12-24T20:15:59")
	if err != nil {
		t.Fatal(err)
	}
	{
		backups := []dto.Backup{}
		result, err = getKeepToday(backups, now)
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(result, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(result), Is(len(backups)))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		backups := []dto.Backup{
			createBackup("2013-12-23T10:15:59"),
			createBackup("2013-12-24T15:15:59"),
			createBackup("2013-12-25T20:15:59"),
		}
		result, err = getKeepToday(backups, now)
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(result, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(result), Is(1))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetKeepDay(t *testing.T) {
	var result []dto.Backup
	now, err := getTimeByName("2013-12-24T20:15:59")
	if err != nil {
		t.Fatal(err)
	}
	{
		backups := []dto.Backup{}
		result, err = getKeepDay(backups, now)
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(result, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(result), Is(len(backups)))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		backups := []dto.Backup{
			createBackup("2013-12-16T20:15:59"),
			createBackup("2013-12-17T20:15:58"),
			createBackup("2013-12-17T20:15:59"),
			createBackup("2013-12-18T20:15:59"),
			createBackup("2013-12-19T20:15:59"),
			createBackup("2013-12-20T20:15:59"),
			createBackup("2013-12-21T20:15:59"),
			createBackup("2013-12-22T20:15:59"),
			createBackup("2013-12-23T20:15:59"),
			createBackup("2013-12-24T20:15:59"),
		}
		result, err = getKeepDay(backups, now)
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(result, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(result), Is(8))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestAgeLessThan7Days(t *testing.T) {
	var err error
	{
		ti, _ := getTimeByName("2013-12-24T20:15:59")
		now, _ := getTimeByName("2013-12-24T20:15:59")
		err = AssertThat(ageLessThanDays(ti, now, 7), Is(true))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		ti, _ := getTimeByName("2013-12-17T20:15:59")
		now, _ := getTimeByName("2013-12-24T20:15:59")
		err = AssertThat(ageLessThanDays(ti, now, 7), Is(true))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		ti, _ := getTimeByName("2013-12-17T20:15:58")
		now, _ := getTimeByName("2013-12-24T20:15:59")
		err = AssertThat(ageLessThanDays(ti, now, 7), Is(false))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetKeepWeek(t *testing.T) {
	var result []dto.Backup
	now, err := getTimeByName("2013-12-24T20:15:59")
	if err != nil {
		t.Fatal(err)
	}
	{
		backups := []dto.Backup{}
		result, err = getKeepWeek(backups, now)
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(result, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(result), Is(len(backups)))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		backups := []dto.Backup{
			createBackup("2013-12-06T20:15:59"),
			createBackup("2013-12-07T20:15:59"),
			createBackup("2013-12-08T20:15:59"),
			createBackup("2013-12-09T20:15:59"),
			createBackup("2013-12-10T20:15:59"),
			createBackup("2013-12-11T20:15:59"),
			createBackup("2013-12-12T20:15:59"),
			createBackup("2013-12-13T20:15:59"),
			createBackup("2013-12-14T20:15:59"),
			createBackup("2013-12-15T20:15:59"),
			createBackup("2013-12-16T20:15:59"),
			createBackup("2013-12-17T20:15:59"),
			createBackup("2013-12-18T20:15:59"),
			createBackup("2013-12-19T20:15:59"),
			createBackup("2013-12-20T20:15:59"),
			createBackup("2013-12-21T20:15:59"),
			createBackup("2013-12-22T20:15:59"),
			createBackup("2013-12-23T20:15:59"),
			createBackup("2013-12-24T20:15:59"),
		}
		result, err = getKeepWeek(backups, now)
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(result, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(result), Is(4))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestListKeepBackups(t *testing.T) {
	var (
		service  BackupService
		backups  []dto.Backup
		host     dto.Host
		hostName string
	)
	hostName = "firewall.example.com"
	clearBackupRootDir(BACKUP_ROOT_DIR)
	createRootDir(BACKUP_ROOT_DIR)
	createHostDir(BACKUP_ROOT_DIR, hostName)
	service = NewBackupService(BACKUP_ROOT_DIR)

	rootDir, err := service.GetRootdir(BACKUP_ROOT_DIR)
	if err != nil {
		t.Fatal(err)
	}
	host, err = service.GetHost(rootDir, hostName)
	if err != nil {
		t.Fatal(err)
	}
	{
		backups, err = service.ListKeepBackups(host)
		err = AssertThat(err, NilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(backups, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(backups), Is(0))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestLatestBackup(t *testing.T) {
	var (
		err     error
		backups []dto.Backup
		backup  dto.Backup
	)
	{
		backups = []dto.Backup{}
		backup = latestBackup(backups)
		err = AssertThat(backup, NilValue())
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		backups := []dto.Backup{
			createBackup("2013-12-06T20:15:59"),
		}
		backup = latestBackup(backups)
		err = AssertThat(backup, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(backup.GetName(), Is("2013-12-06T20:15:59"))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		backups := []dto.Backup{
			createBackup("2013-12-06T20:15:55"),
			createBackup("2013-12-06T20:15:54"),
			createBackup("2013-12-06T20:15:53"),
			createBackup("2013-12-06T20:15:56"),
			createBackup("2013-12-06T20:15:52"),
		}
		backup = latestBackup(backups)
		err = AssertThat(backup, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(backup.GetName(), Is("2013-12-06T20:15:56"))
		if err != nil {
			t.Fatal(err)
		}
	}
}
