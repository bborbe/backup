package main

import (
	"fmt"
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/config"
	backup_mock "github.com/bborbe/backup/service"
	io "github.com/bborbe/io/mock"
)

func TestResumeFail(t *testing.T) {
	writer := io.NewWriter()
	backupService := backup_mock.NewBackupServiceMock()
	backupService.SetResume(fmt.Errorf("error"))

	err := do(writer, backupService, config.DEFAULT_HOST)
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
	writer := io.NewWriter()
	backupService := backup_mock.NewBackupServiceMock()
	backupService.SetResume(nil)

	err := do(writer, backupService, config.DEFAULT_HOST)
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
