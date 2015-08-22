package main

import (
	"os"
	"testing"

	. "github.com/bborbe/assert"
	backup_config "github.com/bborbe/backup/config"
	backup_dto "github.com/bborbe/backup/dto"
	backup_service "github.com/bborbe/backup/service"
	io_mock "github.com/bborbe/io/mock"
)

func TestDoEmpty(t *testing.T) {
	writer := io_mock.NewWriter()
	backupService := backup_service.NewBackupServiceMock()
	backupService.SetListHosts(make([]backup_dto.Host, 0), nil)
	err := do(writer, backupService, backup_config.DEFAULT_ROOT_DIR, backup_config.DEFAULT_HOST, os.TempDir()+"/bla.lock")
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(writer.Content(), NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(writer.Content()), Is(0))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDoNotEmpty(t *testing.T) {
	writer := io_mock.NewWriter()
	backupService := backup_service.NewBackupServiceMock()
	hosts := []backup_dto.Host{
		backup_service.CreateHost("hostA"),
		backup_service.CreateHost("hostB"),
	}
	backupService.SetListHosts(hosts, nil)
	backups := []backup_dto.Backup{
		backup_service.CreateBackup("backupA"),
		backup_service.CreateBackup("backupB"),
	}
	backupService.SetListOldBackups(backups, nil)
	backupService.SetCleanup(nil)
	err := do(writer, backupService, backup_config.DEFAULT_ROOT_DIR, backup_config.DEFAULT_HOST, os.TempDir()+"/bla.lock")
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(writer.Content(), NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(writer.Content()), Gt(0))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(string(writer.Content()), Is("hostA cleaned\nhostB cleaned\n"))
	if err != nil {
		t.Fatal(err)
	}
}
