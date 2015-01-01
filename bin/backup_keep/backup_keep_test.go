package main

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/config"
	"github.com/bborbe/backup/dto"
	backup_mock "github.com/bborbe/backup/service"
	"github.com/bborbe/io"
)

func TestDoEmpty(t *testing.T) {
	writer := io.NewWriter()
	backupService := backup_mock.NewBackupServiceMock()
	backupService.SetListHosts(make([]dto.Host, 0), nil)
	backupService.SetListBackups(make([]dto.Backup, 0), nil)
	err := do(writer, backupService, config.DEFAULT_HOST)
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
	writer := io.NewWriter()
	backupService := backup_mock.NewBackupServiceMock()
	hosts := []dto.Host{
		backup_mock.CreateHost("hostA"),
	}
	backupService.SetListHosts(hosts, nil)
	backups := []dto.Backup{
		backup_mock.CreateBackup("2013-12-25T20:15:59"),
		backup_mock.CreateBackup("2013-12-24T20:15:59"),
	}
	backupService.SetListKeepBackups(backups, nil)
	err := do(writer, backupService, config.DEFAULT_HOST)
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
