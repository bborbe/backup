package backup

import (
	"testing"
	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/host"
	"github.com/bborbe/backup/rootdir"
)

func TestImplementsBackup(t *testing.T) {
	backup := ByName(host.ByName(rootdir.ByName("/rootdir"), "hostname"), "backupname")
	var expected *Backup
	err := AssertThat(backup, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestName(t *testing.T) {
	backup := ByName(host.ByName(rootdir.ByName("/rootdir"), "hostname"), "backupname")
	err := AssertThat(backup.Name(), Is("backupname"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestPath(t *testing.T) {
	backup := ByName(host.ByName(rootdir.ByName("/rootdir"), "hostname"), "backupname")
	err := AssertThat(backup.Path(), Is("/rootdir/hostname/backupname"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestValidBackupName(t *testing.T) {
	var err error
	err = AssertThat(validName("foo"), Is(false))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(validName("2013-12-12T24:15:59"), Is(true))
	if err != nil {
		t.Fatal(err)
	}
}

func TestResume(t *testing.T) {
	err := Resume(host.ByName(rootdir.ByName("/rootdir"), "hostname"))
	err = AssertThat(err, NilValue())
	if err != nil {
		t.Fatal(err)
	}
}
