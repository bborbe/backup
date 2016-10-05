package rootdir

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

func TestImplementsRootdir(t *testing.T) {
	object := ByName("/rootdir")
	var expected *Rootdir
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestPath(t *testing.T) {
	for _, name := range []string{"/rootdirA", "/rootdirB"} {
		object := ByName(name)
		err := AssertThat(object.Path(), Is(name))
		if err != nil {
			t.Fatal(err)
		}
	}
}
