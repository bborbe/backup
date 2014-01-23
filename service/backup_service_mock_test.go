package service

import (
	. "github.com/bborbe/assert"
	"testing"
)

func TestImplementsBackupServiceMock(t *testing.T) {
	s := NewBackupServiceMock()
	var expected *BackupService
	err := AssertThat(s, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
