package status_client

import (
	"testing"

	"net/http"

	. "github.com/bborbe/assert"
	"github.com/golang/glog"
	"os"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}

func TestImplementsServer(t *testing.T) {
	var expected *http.Server
	var port int
	var rootdir string
	expected = NewServer(nil, port, rootdir)
	if err := AssertThat(expected, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
