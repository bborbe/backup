package host

import (
	"testing"
	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/rootdir"
)

func TestImplementsHost(t *testing.T) {
	object := ByName(rootdir.New("/rootdir"), "hostname")
	var expected *Host
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestName(t *testing.T) {
	host := ByName(rootdir.New("/rootdir"), "hostname")
	err := AssertThat(host.Name(), Is("hostname"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestPath(t *testing.T) {
	host := ByName(rootdir.New("/rootdir"), "hostname")
	err := AssertThat(host.Path(), Is("/rootdir/hostname"))
	if err != nil {
		t.Fatal(err)
	}
}
