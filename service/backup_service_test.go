package service

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/dto"
	"github.com/bborbe/backup/testutil"
)

func TestImplementsBackupService(t *testing.T) {
	service := NewBackupService(testutil.BACKUP_ROOT_DIR)
	var expected *BackupService
	err := AssertThat(service, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestListHosts(t *testing.T) {
	testutil.ClearRootDir(testutil.BACKUP_ROOT_DIR)
	service := NewBackupService(testutil.BACKUP_ROOT_DIR)
	{
		_, err := service.ListHosts()
		err = AssertThat(err, NotNilValue().Message("expect error, backup dir not existing"))
		if err != nil {
			t.Fatal(err)
		}
	}
	testutil.CreateRootDir(testutil.BACKUP_ROOT_DIR)
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
	testutil.CreateHostDir(testutil.BACKUP_ROOT_DIR, hostName)
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
	testutil.ClearRootDir(testutil.BACKUP_ROOT_DIR)
	service := NewBackupService(testutil.BACKUP_ROOT_DIR)
	hostName := "firewall.example.com"
	testutil.CreateRootDir(testutil.BACKUP_ROOT_DIR)
	{
		_, err := service.GetHost(hostName)
		err = AssertThat(err, NotNilValue().Message("expect error, backup dir not existing"))
		if err != nil {
			t.Fatal(err)
		}
	}
	testutil.CreateHostDir(testutil.BACKUP_ROOT_DIR, hostName)
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
	var err error
	err = testutil.ClearRootDir(testutil.BACKUP_ROOT_DIR)
	if err != nil {
		t.Fatal(err)
	}
	service := NewBackupService(testutil.BACKUP_ROOT_DIR)
	{
		_, err := service.ListBackups(nil)
		err = AssertThat(err, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
	}
	hostName := "firewall.example.com"
	err = testutil.CreateRootDir(testutil.BACKUP_ROOT_DIR)
	if err != nil {
		t.Fatal(err)
	}
	err = testutil.CreateHostDir(testutil.BACKUP_ROOT_DIR, hostName)
	if err != nil {
		t.Fatal(err)
	}
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
	err = testutil.CreateBackupDir(testutil.BACKUP_ROOT_DIR, hostName, backupName)
	if err != nil {
		t.Fatal(err)
	}
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
	err = testutil.CreateBackupDir(testutil.BACKUP_ROOT_DIR, hostName, "incomplete")
	if err != nil {
		t.Fatal(err)
	}
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
	testutil.ClearRootDir(testutil.BACKUP_ROOT_DIR)
	service := NewBackupService(testutil.BACKUP_ROOT_DIR)
	{
		_, err := service.GetLatestBackup(nil)
		err = AssertThat(err, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
	}
	testutil.CreateRootDir(testutil.BACKUP_ROOT_DIR)
	{
		_, err := service.GetLatestBackup(nil)
		err = AssertThat(err, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
	}
	hostName := "firewall.example.com"
	testutil.CreateHostDir(testutil.BACKUP_ROOT_DIR, hostName)
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
	testutil.CreateBackupDir(testutil.BACKUP_ROOT_DIR, hostName, backupName)
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
	testutil.CreateBackupDir(testutil.BACKUP_ROOT_DIR, hostName, backupName)
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
	testutil.ClearRootDir(testutil.BACKUP_ROOT_DIR)
	service := NewBackupService(testutil.BACKUP_ROOT_DIR)
	{
		_, err := service.ListOldBackups(nil)
		err = AssertThat(err, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
	}
	testutil.CreateRootDir(testutil.BACKUP_ROOT_DIR)
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
	testutil.ClearRootDir(testutil.BACKUP_ROOT_DIR)
	service := NewBackupService(testutil.BACKUP_ROOT_DIR)
	{
		err := service.Cleanup(nil)
		err = AssertThat(err, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
	}
	testutil.CreateRootDir(testutil.BACKUP_ROOT_DIR)
	{
		err := service.Cleanup(nil)
		err = AssertThat(err, NotNilValue())
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
		err      error
	)
	hostName = "firewall.example.com"
	testutil.ClearRootDir(testutil.BACKUP_ROOT_DIR)
	testutil.CreateRootDir(testutil.BACKUP_ROOT_DIR)
	testutil.CreateHostDir(testutil.BACKUP_ROOT_DIR, hostName)
	service = NewBackupService(testutil.BACKUP_ROOT_DIR)

	host, err = service.GetHost(hostName)
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
