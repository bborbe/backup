package util

import (
	"sort"
	"testing"
	. "github.com/bborbe/assert"
)

func TestStringSort(t *testing.T) {
	var err error
	names := []string{"aa", "a", "aaa"}
	sort.Strings(names)
	err = AssertThat(names[0], Is("a"))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(names[1], Is("aa"))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(names[2], Is("aaa"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestStringLess(t *testing.T) {
	var err error
	err = AssertThat(StringLess("-", "-"), Is(false))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(StringLess("7", "8"), Is(true))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(StringLess("0", "8"), Is(true))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(StringLess("a", "b"), Is(true))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(StringLess("b", "a"), Is(false))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(StringLess("a", "a"), Is(false))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(StringLess("a", "aa"), Is(true))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(StringLess("aa", "a"), Is(false))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(StringLess("2013-07-29T10:20:15", "2013-08-23T07:45:48"), Is(true))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(StringLess("2013-08-23T07:45:48", "2013-07-29T10:20:15"), Is(false))
	if err != nil {
		t.Fatal(err)
	}
}
