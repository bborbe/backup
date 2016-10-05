package dto

import (
	"sort"
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

func TestHostSortEmpty(t *testing.T) {
	backups := make([]Host, 0)
	sort.Sort(HostByName(backups))
	err := AssertThat(backups, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(backups), Is(0))
	if err != nil {
		t.Fatal(err)
	}
}

func TestHostSortOne(t *testing.T) {
	backups := []Host{createHost("test")}
	sort.Sort(HostByName(backups))
	err := AssertThat(backups, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(backups), Is(1))
	if err != nil {
		t.Fatal(err)
	}
}

func TestHostSort(t *testing.T) {
	backups := []Host{createHost("c"), createHost("a"), createHost("b")}
	sort.Sort(HostByName(backups))
	err := AssertThat(backups, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(backups), Is(3))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(backups[0].GetName(), Is("a"))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(backups[1].GetName(), Is("b"))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(backups[2].GetName(), Is("c"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestHostSortReal(t *testing.T) {
	backups := []Host{createHost("2013-08-25T16:33:26"), createHost("2013-07-29T10:20:15"), createHost("2013-08-23T07:45:48")}
	sort.Sort(HostByName(backups))
	err := AssertThat(backups, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(backups), Is(3))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(backups[0].GetName(), Is("2013-07-29T10:20:15"))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(backups[1].GetName(), Is("2013-08-23T07:45:48"))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(backups[2].GetName(), Is("2013-08-25T16:33:26"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestHostSortSameLetter(t *testing.T) {
	backups := []Host{createHost("aaa"), createHost("a"), createHost("aa")}
	sort.Sort(HostByName(backups))
	err := AssertThat(backups, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(backups), Is(3))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(backups[0].GetName(), Is("a"))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(backups[1].GetName(), Is("aa"))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(backups[2].GetName(), Is("aaa"))
	if err != nil {
		t.Fatal(err)
	}
}

func createHost(name string) Host {
	host := NewHost()
	host.SetName(name)
	return host
}
