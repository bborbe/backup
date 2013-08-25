package mock

import (
	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/service"
	"testing"
)

func TestImplementsBackupService(t *testing.T) {
	s := NewBackupServiceMock()
	var expected *service.BackupService
	err := AssertThat(s, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
