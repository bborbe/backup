package service

import (
	"fmt"
	. "github.com/bborbe/assert"
	"github.com/bborbe/log"
	"os"
	"testing"
)

var logger = log.DefaultLogger

const BACKUP_ROOT_DIR = "/tmp/backuproot"

func clearBackupRootDir(root string) error {
	return os.RemoveAll(root)
}

func createBackupRootDir(root string) error {
	var fileMode os.FileMode
	fileMode = 0777
	return os.Mkdir(root, fileMode)
}

func createHostDir(root string, host string) error {
	var fileMode os.FileMode
	fileMode = 0777
	dir := fmt.Sprintf("%s%c%s", root, os.PathSeparator, host)
	logger.Debugf("create hostdir %s", dir)
	return os.Mkdir(dir, fileMode)
}

func TestImplementsBackupService(t *testing.T) {
	service := NewBackupService(BACKUP_ROOT_DIR)
	var expected *BackupService
	err := AssertThat(service, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestListHosts(t *testing.T) {
	clearBackupRootDir(BACKUP_ROOT_DIR)
	service := NewBackupService(BACKUP_ROOT_DIR)
	{
		_, err := service.ListHosts()
		err = AssertThat(err, NotNilValue().Message("expect error, backup dir not existing"))
		if err != nil {
			t.Fatal(err)
		}
	}
	createBackupRootDir(BACKUP_ROOT_DIR)
	{
		hosts, err := service.ListHosts()
		err = AssertThat(err, NilValue().Message("expect no error, backup dir existing"))
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(hosts), Is(0))
		if err != nil {
			t.Fatal(err)
		}
	}
	hostName := "firewall.example.com"
	createHostDir(BACKUP_ROOT_DIR, hostName)
	{
		hosts, err := service.ListHosts()
		err = AssertThat(err, NilValue().Message("expect no error, backup dir existing"))
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(hosts), Is(1))
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(hosts[0].GetName(), Is(hostName))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestListBackups(t *testing.T) {

}
