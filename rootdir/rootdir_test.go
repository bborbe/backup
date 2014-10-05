package rootdir

import (
	"testing"
	. "github.com/bborbe/assert"
)

func TestImplementsRootdir(t *testing.T) {
	object := New("/rootdir")
	var expected *Rootdir
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestPath(t *testing.T) {
	for _, name := range []string{"/rootdirA", "/rootdirB"} {
		object := New(name)
		err := AssertThat(object.Path(), Is(name))
		if err != nil {
			t.Fatal(err)
		}
	}
}
