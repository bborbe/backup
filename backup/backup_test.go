package backup

import (
	"testing"
	"time"

	. "github.com/bborbe/assert"
	backup_host "github.com/bborbe/backup/host"
	backup_rootdir "github.com/bborbe/backup/rootdir"
	backup_testutil "github.com/bborbe/backup/testutil"
)

func TestByTime(t *testing.T) {
	rootdirName := backup_testutil.BACKUP_ROOT_DIR
	hostName := "hostname"
	h := backup_host.ByName(backup_rootdir.ByName(rootdirName), hostName)
	ti, err := time.Parse("2006-01-02T15:04:05", "2010-12-24T10:11:12")
	if err != nil {
		t.Fatal(err)
	}
	backup := ByTime(h, ti)
	err = AssertThat(backup.Name(), Is("2010-12-24T10:11:12"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestImplementsBackup(t *testing.T) {
	backup := ByName(backup_host.ByName(backup_rootdir.ByName("/rootdir"), "hostname"), "backupname")
	var expected *Backup
	err := AssertThat(backup, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestName(t *testing.T) {
	backup := ByName(backup_host.ByName(backup_rootdir.ByName("/rootdir"), "hostname"), "backupname")
	err := AssertThat(backup.Name(), Is("backupname"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestPath(t *testing.T) {
	backup := ByName(backup_host.ByName(backup_rootdir.ByName("/rootdir"), "hostname"), "backupname")
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
	rootdirName := backup_testutil.BACKUP_ROOT_DIR
	hostName := "hostname"
	h := backup_host.ByName(backup_rootdir.ByName(rootdirName), hostName)
	backupNameA := "2014-12-24T13:14:15"
	backupNameB := ByTime(h, time.Now()).Name()
	var err error
	err = backup_testutil.ClearRootDir(rootdirName)
	if err != nil {
		t.Fatal(err)
	}
	err = backup_testutil.CreateRootDir(rootdirName)
	if err != nil {
		t.Fatal(err)
	}
	err = backup_testutil.CreateHostDir(rootdirName, hostName)
	if err != nil {
		t.Fatal(err)
	}
	err = backup_testutil.CreateBackupDir(rootdirName, hostName, backupNameA)
	if err != nil {
		t.Fatal(err)
	}
	err = backup_testutil.CreateBackupDir(rootdirName, hostName, backupNameB)
	if err != nil {
		t.Fatal(err)
	}
	err = backup_testutil.CreateBackupCurrentSymlink(rootdirName, hostName, backupNameB)
	if err != nil {
		t.Fatal(err)
	}
	exists, err := existsIncomplete(h)
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(exists, Is(false))
	if err != nil {
		t.Fatal(err)
	}
	err = Resume(h)
	err = AssertThat(err, NilValue())
	if err != nil {
		t.Fatal(err)
	}
	exists, err = existsIncomplete(h)
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(exists, Is(true))
	if err != nil {
		t.Fatal(err)
	}
}
