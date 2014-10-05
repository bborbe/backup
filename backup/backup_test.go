package backup

import (
	"testing"
	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/host"
	"github.com/bborbe/backup/rootdir"
)

func TestImplementsBackup(t *testing.T) {
	backup := ByName(host.ByName(rootdir.New("/rootdir"), "hostname"), "backupname")
	var expected *Backup
	err := AssertThat(backup, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestName(t *testing.T) {
	backup := ByName(host.ByName(rootdir.New("/rootdir"), "hostname"), "backupname")
	err := AssertThat(backup.Name(), Is("backupname"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestPath(t *testing.T) {
	backup := ByName(host.ByName(rootdir.New("/rootdir"), "hostname"), "backupname")
	err := AssertThat(backup.Path(), Is("/rootdir/hostname/backupname"))
	if err != nil {
		t.Fatal(err)
	}
}
