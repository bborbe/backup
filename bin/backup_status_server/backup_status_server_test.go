package main

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/golang/glog"
	"os"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}

func TestResumeFail(t *testing.T) {
	if err := AssertThat(DEFAULT_PORT, Is(DEFAULT_PORT)); err != nil {
		t.Fatal(err)
	}
}
