package service

import (
	"fmt"
	. "github.com/bborbe/assert"
	"github.com/bborbe/log"
	"os"
	"testing"
)

var logger = log.DefaultLogger

const BACKUP_ROOT_DIR = "/tmp/backuproot"

func clearBackupRootDir(root string) error {
	return os.RemoveAll(root)
}

func createBackupRootDir(root string) error {
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
	createBackupRootDir(BACKUP_ROOT_DIR)
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
	{
		_, err := service.GetHost(hostName)
		err = AssertThat(err, NotNilValue().Message("expect error, backup dir not existing"))
		if err != nil {
			t.Fatal(err)
		}
	}
	createBackupRootDir(BACKUP_ROOT_DIR)
	createHostDir(BACKUP_ROOT_DIR, hostName)
	{
		host, err := service.GetHost(hostName)
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
	createBackupRootDir(BACKUP_ROOT_DIR)
	createHostDir(BACKUP_ROOT_DIR, hostName)
	{
		host, err := service.GetHost(hostName)
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
		host, err := service.GetHost(hostName)
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
		host, err := service.GetHost(hostName)
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
	createBackupRootDir(BACKUP_ROOT_DIR)
	{
		_, err := service.GetLatestBackup(nil)
		err = AssertThat(err, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
	}
	hostName := "firewall.example.com"
	createHostDir(BACKUP_ROOT_DIR, hostName)
	{
		host, err := service.GetHost(hostName)
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
		host, err := service.GetHost(hostName)
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
		host, err := service.GetHost(hostName)
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
	createBackupRootDir(BACKUP_ROOT_DIR)
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
	createBackupRootDir(BACKUP_ROOT_DIR)
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
