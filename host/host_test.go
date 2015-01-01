package host

import (
	"fmt"
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/rootdir"
	"github.com/bborbe/backup/testutil"
)

func TestImplementsHost(t *testing.T) {
	h := ByName(rootdir.ByName(testutil.BACKUP_ROOT_DIR), "hostname")
	var expected *Host
	err := AssertThat(h, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestName(t *testing.T) {
	h := ByName(rootdir.ByName(testutil.BACKUP_ROOT_DIR), "hostname")
	err := AssertThat(h.Name(), Is("hostname"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestPath(t *testing.T) {
	h := ByName(rootdir.ByName(testutil.BACKUP_ROOT_DIR), "hostname")
	err := AssertThat(h.Path(), Is("/tmp/backuproot/hostname"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestAllRootdirDoesNotExists(t *testing.T) {
	var err error
	err = testutil.ClearRootDir(testutil.BACKUP_ROOT_DIR)
	if err != nil {
		t.Fatal(err)
	}
	_, err = All(rootdir.ByName(testutil.BACKUP_ROOT_DIR))
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}

func TestAllEmpty(t *testing.T) {
	var err error
	err = testutil.ClearRootDir(testutil.BACKUP_ROOT_DIR)
	if err != nil {
		t.Fatal(err)
	}
	err = testutil.CreateRootDir(testutil.BACKUP_ROOT_DIR)
	if err != nil {
		t.Fatal(err)
	}
	hosts, err := All(rootdir.ByName(testutil.BACKUP_ROOT_DIR))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(hosts), Is(0))
	if err != nil {
		t.Fatal(err)
	}
}

func TestAllBackups(t *testing.T) {
	var err error
	err = testutil.ClearRootDir(testutil.BACKUP_ROOT_DIR)
	if err != nil {
		t.Fatal(err)
	}
	err = testutil.CreateRootDir(testutil.BACKUP_ROOT_DIR)
	if err != nil {
		t.Fatal(err)
	}
	err = testutil.CreateHostDir(testutil.BACKUP_ROOT_DIR, "hostA")
	if err != nil {
		t.Fatal(err)
	}
	err = testutil.CreateHostDir(testutil.BACKUP_ROOT_DIR, "hostB")
	if err != nil {
		t.Fatal(err)
	}
	hosts, err := All(rootdir.ByName(testutil.BACKUP_ROOT_DIR))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(hosts), Is(2))
	if err != nil {
		t.Fatal(err)
	}
}

func TestAllFile(t *testing.T) {
	var err error
	err = testutil.ClearRootDir(testutil.BACKUP_ROOT_DIR)
	if err != nil {
		t.Fatal(err)
	}
	err = testutil.CreateRootDir(testutil.BACKUP_ROOT_DIR)
	if err != nil {
		t.Fatal(err)
	}
	err = testutil.CreateFile(fmt.Sprintf("%s/file", testutil.BACKUP_ROOT_DIR))
	if err != nil {
		t.Fatal(err)
	}
	hosts, err := All(rootdir.ByName(testutil.BACKUP_ROOT_DIR))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(hosts), Is(0))
	if err != nil {
		t.Fatal(err)
	}
}
