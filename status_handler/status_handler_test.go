package status_handler

import (
	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/status_checker"
	"net/http"
	"testing"
)

func TestImplementsStatusHandler(t *testing.T) {
	var statusChecker status_checker.StatusChecker
	object := NewStatusHandler(statusChecker)
	var expected *http.Handler
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}
