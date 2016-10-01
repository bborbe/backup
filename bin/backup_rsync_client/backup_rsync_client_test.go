package main

import (
	"testing"

	"flag"
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

func TestDo(t *testing.T) {
	file, err := ioutil.TempFile("", "config")
	defer os.Remove(file.Name())
	if err != nil {
		t.Fatal("create temp file faileD")
	}
	file.WriteString(`[{}]`)
	file.Close()
	flag.Set(parameterConfigPath, file.Name())
	if err := AssertThat(do(), NilValue()); err != nil {
		t.Fatal(err)
	}
}
