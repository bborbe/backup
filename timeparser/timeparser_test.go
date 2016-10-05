package timeparser

import (
	"testing"

	"os"

	. "github.com/bborbe/assert"
	"github.com/golang/glog"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}

func TestImplementsTimeParse(t *testing.T) {
	object := New()
	var expected *TimeParser
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
