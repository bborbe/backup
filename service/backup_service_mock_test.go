package service

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsBackupServiceMock(t *testing.T) {
	s := NewBackupServiceMock()
	var expected *BackupService
	err := AssertThat(s, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
