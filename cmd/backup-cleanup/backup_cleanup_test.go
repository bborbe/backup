package main

import (
	"os"
	"testing"

	. "github.com/bborbe/assert"
	"github.com/golang/glog"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}

func TestDoEmpty(t *testing.T) {
	err := do()
	if err = AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
