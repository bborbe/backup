package main

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestFrom(t *testing.T) {
	h := host{
		Active:      true,
		User:        "backupuser",
		Host:        "example.com",
		Port:        1337,
		Directory:   "/data",
		ExcludeFrom: "exclude_from",
	}
	if err := AssertThat(h.from(), Is("backupuser@example.com:/data")); err != nil {
		t.Fatal(err)
	}
}

func TestTo(t *testing.T) {
	h := host{
		Active:      true,
		User:        "backupuser",
		Host:        "example.com",
		Port:        1337,
		Directory:   "/data",
		ExcludeFrom: "exclude_from",
	}
	if err := AssertThat(h.to(targetDirectory("/backup")), Is("/backup/example.com/incomplete/data")); err != nil {
		t.Fatal(err)
	}
}
