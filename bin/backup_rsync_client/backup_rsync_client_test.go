package main

import (
	"testing"

	"fmt"
	. "github.com/bborbe/assert"
	"github.com/golang/glog"
	"io/ioutil"
	"os"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}

func TestGetHostsByConfig(t *testing.T) {
	host := "example.com"
	excludeFrom := "exclude"
	user := "backupuser"
	directory := "/data/"
	port := 1337

	file, err := ioutil.TempFile("", "config")
	filename := file.Name()
	defer os.Remove(filename)
	if err != nil {
		t.Fatal("create temp file faileD")
	}
	file.WriteString(fmt.Sprintf(`[{"host":"%s","port":%d,"dir":"%s","user":"%s","exclude_from":"%s","active":true}]`, host, port, directory, user, excludeFrom))
	file.Close()
	configPathPtr = &filename

	hosts, err := getHosts()

	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(len(hosts), Is(1)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(hosts[0].Active, Is(true)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(hosts[0].User, Is(user)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(hosts[0].Host, Is(host)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(hosts[0].Port, Is(port)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(hosts[0].Directory, Is(directory)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(hosts[0].ExcludeFrom, Is(excludeFrom)); err != nil {
		t.Fatal(err)
	}
}

func TestGetHostsByArgs(t *testing.T) {
	configPathPtr = nil
	host := "example.com"
	excludeFrom := "exclude"
	user := "backupuser"
	directory := "/data/"
	port := 1337
	hostPtr = &host
	excludeFromPtr = &excludeFrom
	dirPtr = &directory
	portPtr = &port
	userPtr = &user
	hosts, err := getHosts()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(len(hosts), Is(1)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(hosts[0].Active, Is(true)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(hosts[0].User, Is(user)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(hosts[0].Host, Is(host)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(hosts[0].Port, Is(port)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(hosts[0].Directory, Is(directory)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(hosts[0].ExcludeFrom, Is(excludeFrom)); err != nil {
		t.Fatal(err)
	}
}

func TestGetTargetDirectoryFailed(t *testing.T) {
	target := "/backup_not_existing"
	targetPtr = &target
	dir, err := getTargetDirectory()
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(dir, NilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestGetTargetDirectorySuccess(t *testing.T) {
	target := os.TempDir()
	targetPtr = &target
	dir, err := getTargetDirectory()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(dir, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(dir.IsValid(), NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(dir.String(), Is(target)); err != nil {
		t.Fatal(err)
	}
}
