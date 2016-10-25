package dto

import (
	. "github.com/bborbe/assert"
	"testing"
	"time"
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

func TestFormatDurationDays(t *testing.T) {
	d, err := time.ParseDuration("128h")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(FormatDuration(d), Is("5d")); err != nil {
		t.Fatal(err)
	}
}

func TestFormatDurationHours(t *testing.T) {
	d, err := time.ParseDuration("8h59m")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(FormatDuration(d), Is("8h")); err != nil {
		t.Fatal(err)
	}
}

func TestFormatDurationMinutes(t *testing.T) {
	d, err := time.ParseDuration("59m30s")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(FormatDuration(d), Is("59m")); err != nil {
		t.Fatal(err)
	}
}
