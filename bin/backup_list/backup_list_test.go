package main

import (
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
	backupService.SetListBackups(make([]backup_dto.Backup, 0), nil)
	err := do(writer, backupService, backup_config.DEFAULT_HOST)
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
	}
	backupService.SetListHosts(hosts, nil)
	backups := []backup_dto.Backup{
		backup_service.CreateBackup("2013-12-24T20:15:59"),
		backup_service.CreateBackup("2013-12-25T20:15:59"),
	}
	backupService.SetListBackups(backups, nil)
	err := do(writer, backupService, backup_config.DEFAULT_HOST)
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
	err = AssertThat(string(writer.Content()), Is("hostA/2013-12-24T20:15:59\nhostA/2013-12-25T20:15:59\n"))
	if err != nil {
		t.Fatal(err)
	}
}
