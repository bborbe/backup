package cache

import (
	"testing"

	"os"

	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/model"
	"github.com/bborbe/backup/status/client/fetcher"
	"github.com/golang/glog"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}

func TestImplementsStatusHandler(t *testing.T) {
	object := New(nil, model.CacheTTL(1234))
	var expected *fetcher.BackupStatusFetcher
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
