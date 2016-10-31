package fetcher

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

func TestImplementsStatusHandler(t *testing.T) {
	object := New(nil, "http://www.example.com/status")
	var expected *BackupStatusFetcher
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
