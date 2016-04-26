package status_server

import (
	"testing"

	"net/http"

	. "github.com/bborbe/assert"
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
