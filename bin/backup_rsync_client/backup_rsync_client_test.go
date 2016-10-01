package main

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestDo(t *testing.T) {
	err := do()
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
}
