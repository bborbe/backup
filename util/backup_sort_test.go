package util

import (
	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/dto"
	"github.com/bborbe/backup/mock"
	"sort"
	"testing"
)

func TestBackupSortEmpty(t *testing.T) {
	backups := make([]dto.Backup, 0)
	sort.Sort(BackupByDate(backups))
	err := AssertThat(backups, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(backups), Is(0))
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackupSortOne(t *testing.T) {
	backups := []dto.Backup{mock.CreateBackup("test")}
	sort.Sort(BackupByDate(backups))
	err := AssertThat(backups, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(backups), Is(1))
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackupSort(t *testing.T) {
	backups := []dto.Backup{mock.CreateBackup("c"), mock.CreateBackup("a"), mock.CreateBackup("b")}
	sort.Sort(BackupByDate(backups))
	err := AssertThat(backups, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(backups), Is(3))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(backups[0].GetName(), Is("a"))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(backups[1].GetName(), Is("b"))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(backups[2].GetName(), Is("c"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackupSortReal(t *testing.T) {
	backups := []dto.Backup{mock.CreateBackup("2013-08-25T16:33:26"), mock.CreateBackup("2013-07-29T10:20:15"), mock.CreateBackup("2013-08-23T07:45:48")}
	sort.Sort(BackupByDate(backups))
	err := AssertThat(backups, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(backups), Is(3))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(backups[0].GetName(), Is("2013-07-29T10:20:15"))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(backups[1].GetName(), Is("2013-08-23T07:45:48"))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(backups[2].GetName(), Is("2013-08-25T16:33:26"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackupSortSameLetter(t *testing.T) {
	backups := []dto.Backup{mock.CreateBackup("aaa"), mock.CreateBackup("a"), mock.CreateBackup("aa")}
	sort.Sort(BackupByDate(backups))
	err := AssertThat(backups, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(backups), Is(3))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(backups[0].GetName(), Is("a"))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(backups[1].GetName(), Is("aa"))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(backups[2].GetName(), Is("aaa"))
	if err != nil {
		t.Fatal(err)
	}
}