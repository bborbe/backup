package status_client

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/server"
)

func TestImplementsServer(t *testing.T) {
	var port int
	var rootdir string
	object := NewServer(nil, port, rootdir)
	var expected *server.Server
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
