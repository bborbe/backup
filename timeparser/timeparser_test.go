package timeparser

import (
	"testing"
	. "github.com/bborbe/assert"
)

func TestImplementsTimeParse(t *testing.T) {
	object := New()
	var expected *TimeParser
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
