package util

import (
	"sort"
	"testing"
	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/dto"
)

func TestHostSortEmpty(t *testing.T) {
	backups := make([]dto.Host, 0)
	sort.Sort(HostByDate(backups))
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
	backups := []dto.Host{createHost("test")}
	sort.Sort(HostByDate(backups))
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
	backups := []dto.Host{createHost("c"), createHost("a"), createHost("b")}
	sort.Sort(HostByDate(backups))
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
	backups := []dto.Host{createHost("2013-08-25T16:33:26"), createHost("2013-07-29T10:20:15"), createHost("2013-08-23T07:45:48")}
	sort.Sort(HostByDate(backups))
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
	backups := []dto.Host{createHost("aaa"), createHost("a"), createHost("aa")}
	sort.Sort(HostByDate(backups))
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

func createHost(name string) dto.Host {
	host := dto.NewHost()
	host.SetName(name)
	return host
}
