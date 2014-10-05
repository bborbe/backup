package rootdir

import (
	"testing"
	. "github.com/bborbe/assert"
)

func TestImplementsRootdir(t *testing.T) {
	object := New("/test")
	var expected *Rootdir
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
