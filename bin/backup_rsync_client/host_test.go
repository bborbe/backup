package main

import (
	"testing"

	. "github.com/bborbe/assert"
)

func getValidHost() host {
	return host{
		Active:      true,
		User:        "backupuser",
		Host:        "example.com",
		Port:        1337,
		Directory:   "/data/",
		ExcludeFrom: "exclude_from",
	}
}

func TestFrom(t *testing.T) {
	h := getValidHost()
	if err := AssertThat(h.from(), Is("backupuser@example.com:/data/")); err != nil {
		t.Fatal(err)
	}
}

func TestFromWithIp(t *testing.T) {
	h := getValidHost()
	h.Ip = "192.168.2.1"
	if err := AssertThat(h.from(), Is("backupuser@192.168.2.1:/data/")); err != nil {
		t.Fatal(err)
	}
}

func TestTo(t *testing.T) {
	h := getValidHost()
	if err := AssertThat(h.to(targetDirectory("/backup")), Is("/backup/example.com/incomplete/data/")); err != nil {
		t.Fatal(err)
	}
}

func TestLinkDest(t *testing.T) {
	h := getValidHost()
	if err := AssertThat(h.linkDest(targetDirectory("/backup")), Is("/backup/example.com/current/data/")); err != nil {
		t.Fatal(err)
	}
}

func TestValidateSuccess(t *testing.T) {
	h := getValidHost()
	if err := AssertThat(h.Validate(), NilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateUserInvalid(t *testing.T) {
	h := getValidHost()
	h.User = ""
	if err := AssertThat(h.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateHostInvalid(t *testing.T) {
	h := getValidHost()
	h.Host = ""
	if err := AssertThat(h.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidatePortInvalid(t *testing.T) {
	h := getValidHost()
	h.Port = 0
	if err := AssertThat(h.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateDirectoryInvalid(t *testing.T) {
	h := getValidHost()
	h.Directory = ""
	if err := AssertThat(h.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestValidateDirectorySlashMissing(t *testing.T) {
	h := getValidHost()
	h.Directory = "/data"
	if err := AssertThat(h.Validate(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
