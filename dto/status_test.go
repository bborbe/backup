package dto

import (
	. "github.com/bborbe/assert"
	"testing"
)

func TestBackupDateTime(t *testing.T) {
	date := BackupDate("2016-10-25T00:01:51")
	time, err := date.Time()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(time.String(), Is("2016-10-25 00:01:51 +0000 UTC")); err != nil {
		t.Fatal(err)
	}
}
