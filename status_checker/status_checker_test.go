package status_checker

import (
	. "github.com/bborbe/assert"
	"testing"
)

func TestImplementsStatusChecker(t *testing.T) {
	var rootdir string
	object := NewStatusChecker(rootdir)
	var expected *StatusChecker
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
