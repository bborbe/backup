package status_client_handler

import (
	"net/http"
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsStatusHandler(t *testing.T) {
	object := NewStatusHandler(nil, "")
	var expected *http.Handler
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
