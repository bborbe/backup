package main

import (
	"fmt"
	"testing"

	. "github.com/bborbe/assert"
	backup_config "github.com/bborbe/backup/config"
	backup_service "github.com/bborbe/backup/service"
	io_mock "github.com/bborbe/io/mock"
)

func TestResumeFail(t *testing.T) {
	writer := io_mock.NewWriter()
	backupService := backup_service.NewBackupServiceMock()
	backupService.SetResume(fmt.Errorf("error"))

	err := do(writer, backupService, backup_config.DEFAULT_HOST)
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(writer.Content(), NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(string(writer.Content()), Is("resume backup for host all failed\n"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestResumeSuccess(t *testing.T) {
	writer := io_mock.NewWriter()
	backupService := backup_service.NewBackupServiceMock()
	backupService.SetResume(nil)

	err := do(writer, backupService, backup_config.DEFAULT_HOST)
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(writer.Content(), NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(string(writer.Content()), Is("resume backup for host all success\n"))
	if err != nil {
		t.Fatal(err)
	}
}
