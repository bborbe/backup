package main

import (
	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/dto"
	backup_mock "github.com/bborbe/backup/mock"
	server_mock "github.com/bborbe/server/mock"
	"testing"
)

func TestDoEmpty(t *testing.T) {
	writer := server_mock.NewWriter()
	backupService := backup_mock.NewBackupServiceMock()
	backupService.SetListHosts(make([]dto.Host, 0), nil)
	backupService.SetListBackups(make([]dto.Backup, 0), nil)
	err := do(writer, backupService)
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
	writer := server_mock.NewWriter()
	backupService := backup_mock.NewBackupServiceMock()
	hosts := []dto.Host{
		createHost("hostA"),
	}
	backupService.SetListHosts(hosts, nil)
	backups := []dto.Backup{
		createBackup("backupA"),
		createBackup("backupB"),
	}
	backupService.SetListBackups(backups, nil)
	err := do(writer, backupService)
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
	err = AssertThat(string(writer.Content()), Is("hostA => backupA\nhostA => backupB\n"))
	if err != nil {
		t.Fatal(err)
	}
}

func createBackup(name string) dto.Backup {
	b := dto.NewBackup()
	b.SetName(name)
	return b
}

func createHost(name string) dto.Host {
	b := dto.NewHost()
	b.SetName(name)
	return b
}
