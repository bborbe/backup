package main

import (
	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/config"
	"github.com/bborbe/backup/dto"
	backup_mock "github.com/bborbe/backup/mock"
	server_mock "github.com/bborbe/server/mock"
	"os"
	"testing"
	"time"
)

func TestDoEmpty(t *testing.T) {
	writer := server_mock.NewWriter()
	backupService := backup_mock.NewBackupServiceMock()
	backupService.SetListHosts(make([]dto.Host, 0), nil)
	err := do(writer, backupService, config.DEFAULT_HOST, os.TempDir()+"/bla.lock")
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
		backup_mock.CreateHost("hostA"),
		backup_mock.CreateHost("hostB"),
	}
	backupService.SetListHosts(hosts, nil)
	backups := []dto.Backup{
		backup_mock.CreateBackup("backupA"),
		backup_mock.CreateBackup("backupB"),
	}
	backupService.SetListOldBackups(backups, nil)
	backupService.SetCleanup(nil)
	err := do(writer, backupService, config.DEFAULT_HOST, os.TempDir()+"/bla.lock")
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

func TestLocking(t *testing.T) {
	var err error
	name := os.TempDir() + "/bla.lock"
	result := true
	err = lock(name)
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		erro := lock(name)
		if erro != nil {
			t.Fatal(erro)
		}
		result = false
		erro = unlock(name)
		if erro != nil {
			t.Fatal(erro)
		}
	}()
	err = AssertThat(result, Is(true))
	if err != nil {
		t.Fatal(err)
	}
	err = unlock(name)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	err = AssertThat(result, Is(false))
	if err != nil {
		t.Fatal(err)
	}
}
