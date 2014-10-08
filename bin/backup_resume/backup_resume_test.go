package main

import (
	"testing"
	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/config"
	"github.com/bborbe/backup/dto"
	backup_mock "github.com/bborbe/backup/service"
	server_mock "github.com/bborbe/server/mock"
)

func TestDoEmpty(t *testing.T) {
	writer := server_mock.NewWriter()
	backupService := backup_mock.NewBackupServiceMock()
	backupService.SetListHosts(make([]dto.Host, 0), nil)
	err := do(writer, backupService, config.DEFAULT_ROOT_DIR, config.DEFAULT_HOST)
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
