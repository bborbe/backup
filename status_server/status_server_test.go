package status_server

import (
	. "github.com/bborbe/assert"
	"github.com/bborbe/server"
	"testing"
)

func TestImplementsServer(t *testing.T) {
	var port int
	var rootdir string
	object := NewServer(port, rootdir)
	var expected *server.Server
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
