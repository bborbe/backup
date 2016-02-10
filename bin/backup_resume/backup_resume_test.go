package main

import (
	"fmt"
	"testing"

	"bytes"

	. "github.com/bborbe/assert"
	backup_config "github.com/bborbe/backup/config"
	backup_service "github.com/bborbe/backup/service"
)

func TestResumeFail(t *testing.T) {
	writer := bytes.NewBufferString("")
	backupService := backup_service.NewBackupServiceMock()
	backupService.SetResume(fmt.Errorf("error"))

	err := do(writer, backupService, backup_config.DEFAULT_HOST)
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(writer.String(), NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(writer.String(), Is("resume backup for host all failed\n"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestResumeSuccess(t *testing.T) {
	writer := bytes.NewBufferString("")
	backupService := backup_service.NewBackupServiceMock()
	backupService.SetResume(nil)

	err := do(writer, backupService, backup_config.DEFAULT_HOST)
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(writer.String(), NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(writer.String(), Is("resume backup for host all success\n"))
	if err != nil {
		t.Fatal(err)
	}
}
