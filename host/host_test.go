package host

import (
	"testing"
	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/rootdir"
)

func TestImplementsHost(t *testing.T) {
	object := ByName(rootdir.New("/test"), "test")
	var expected *Host
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
