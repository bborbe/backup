package status_server

import (
	"testing"

	. "github.com/bborbe/assert"
	"net/http"
)

func TestImplementsServer(t *testing.T) {
	var expected *http.Server
	var port int
	var rootdir string
	expected = NewServer(port, rootdir)
	if err := AssertThat(expected, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
