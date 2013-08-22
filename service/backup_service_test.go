package service

import (
	. "github.com/bborbe/assert"
	"testing"
)

func TestImplementsBackupService(t *testing.T) {
	service := NewBackupService("/rsync")
	var expected *BackupService
	err := AssertThat(service, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
