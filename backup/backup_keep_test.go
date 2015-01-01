package backup

import (
	"testing"
	"time"

	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/host"
	"github.com/bborbe/backup/rootdir"
	"github.com/bborbe/backup/timeparser"
)

func TestAgeLessThan7Days(t *testing.T) {
	var err error
	{
		ti, _ := timeparser.New().TimeByName("2013-12-24T20:15:59")
		now, _ := timeparser.New().TimeByName("2013-12-24T20:15:59")
		err = AssertThat(ageLessThanDays(ti, now, 7), Is(true))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		ti, _ := timeparser.New().TimeByName("2013-12-17T20:15:59")
		now, _ := timeparser.New().TimeByName("2013-12-24T20:15:59")
		err = AssertThat(ageLessThanDays(ti, now, 7), Is(true))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		ti, _ := timeparser.New().TimeByName("2013-12-17T20:15:58")
		now, _ := timeparser.New().TimeByName("2013-12-24T20:15:59")
		err = AssertThat(ageLessThanDays(ti, now, 7), Is(false))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetKeepMonth(t *testing.T) {
	var err error
	var result []Backup

	h := host.ByName(rootdir.ByName("/rootdir"), "hostname")

	{
		backups := []Backup{}
		result, err = getKeepMonth(backups)
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(result, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(result), Is(len(backups)))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		backups := []Backup{
			ByName(h, "2013-12-12T24:15:59"),
		}
		result, err = getKeepMonth(backups)
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(result, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(result), Is(len(backups)))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		backups := []Backup{
			ByName(h, "2013-12-12T24:15:59"),
			ByName(h, "2013-12-01T24:15:59"),
		}
		result, err = getKeepMonth(backups)
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(result, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(result), Is(1))
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat("2013-12-01T24:15:59", Is(result[0].Name()))
	}
	{
		backups := []Backup{
			ByName(h, "2012-11-12T24:15:59"),
			ByName(h, "2013-11-01T24:15:59"),
			ByName(h, "2013-05-28T24:15:59"),
			ByName(h, "2013-05-29T24:15:59"),
		}
		result, err = getKeepMonth(backups)
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(result, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(result), Is(3))
		if err != nil {
			t.Fatal(err)
		}
	}
}
func TestGetKeepToday(t *testing.T) {
	var result []Backup
	now, err := timeparser.New().TimeByName("2013-12-24T20:15:59")
	if err != nil {
		t.Fatal(err)
	}

	h := host.ByName(rootdir.ByName("/rootdir"), "hostname")

	{
		backups := []Backup{}
		result, err = getKeepToday(backups, now, timeparser.New())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(result, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(result), Is(len(backups)))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		backups := []Backup{
			ByName(h, "2013-12-23T10:15:59"),
			ByName(h, "2013-12-24T15:15:59"),
			ByName(h, "2013-12-25T20:15:59"),
		}
		result, err = getKeepToday(backups, now, timeparser.New())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(result, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(result), Is(1))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestgetTimeByName(t *testing.T) {
	{
		_, err := timeparser.New().TimeByName("")
		err = AssertThat(err, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		calcTime, err := timeparser.New().TimeByName("2013-07-01T00:24:52")
		err = AssertThat(err, NilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(calcTime.Year(), Is(2013))
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(calcTime.Month(), Is(time.July))
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(calcTime.Day(), Is(1))
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(calcTime.Hour(), Is(0))
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(calcTime.Minute(), Is(24))
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(calcTime.Second(), Is(52))
		if err != nil {
			t.Fatal(err)
		}
	}

}

func TestLatestBackup(t *testing.T) {
	var (
		err     error
		backups []Backup
		b       Backup
	)
	{
		backups = []Backup{}
		b = latestBackup(backups)
		err = AssertThat(b, NilValue())
		if err != nil {
			t.Fatal(err)
		}
	}

	h := host.ByName(rootdir.ByName("/rootdir"), "hostname")

	{
		backups := []Backup{
			ByName(h, "2013-12-06T20:15:59"),
		}
		b = latestBackup(backups)
		err = AssertThat(b, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(b.Name(), Is("2013-12-06T20:15:59"))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		backups := []Backup{
			ByName(h, "2013-12-06T20:15:55"),
			ByName(h, "2013-12-06T20:15:54"),
			ByName(h, "2013-12-06T20:15:53"),
			ByName(h, "2013-12-06T20:15:56"),
			ByName(h, "2013-12-06T20:15:52"),
		}
		b = latestBackup(backups)
		err = AssertThat(b, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(b.Name(), Is("2013-12-06T20:15:56"))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetKeepWeek(t *testing.T) {
	var result []Backup
	now, err := timeparser.New().TimeByName("2013-12-24T20:15:59")
	if err != nil {
		t.Fatal(err)
	}

	h := host.ByName(rootdir.ByName("/rootdir"), "hostname")

	{
		backups := []Backup{}
		result, err = getKeepWeek(backups, now, timeparser.New())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(result, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(result), Is(len(backups)))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		backups := []Backup{
			ByName(h, "2013-12-06T20:15:59"),
			ByName(h, "2013-12-07T20:15:59"),
			ByName(h, "2013-12-08T20:15:59"),
			ByName(h, "2013-12-09T20:15:59"),
			ByName(h, "2013-12-10T20:15:59"),
			ByName(h, "2013-12-11T20:15:59"),
			ByName(h, "2013-12-12T20:15:59"),
			ByName(h, "2013-12-13T20:15:59"),
			ByName(h, "2013-12-14T20:15:59"),
			ByName(h, "2013-12-15T20:15:59"),
			ByName(h, "2013-12-16T20:15:59"),
			ByName(h, "2013-12-17T20:15:59"),
			ByName(h, "2013-12-18T20:15:59"),
			ByName(h, "2013-12-19T20:15:59"),
			ByName(h, "2013-12-20T20:15:59"),
			ByName(h, "2013-12-21T20:15:59"),
			ByName(h, "2013-12-22T20:15:59"),
			ByName(h, "2013-12-23T20:15:59"),
			ByName(h, "2013-12-24T20:15:59"),
		}
		result, err = getKeepWeek(backups, now, timeparser.New())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(result, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(result), Is(4))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetKeepDay(t *testing.T) {
	var result []Backup
	now, err := timeparser.New().TimeByName("2013-12-24T20:15:59")
	if err != nil {
		t.Fatal(err)
	}

	h := host.ByName(rootdir.ByName("/rootdir"), "hostname")

	{
		backups := []Backup{}
		result, err = getKeepDay(backups, now, timeparser.New())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(result, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(result), Is(len(backups)))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		backups := []Backup{
			ByName(h, "2013-12-16T20:15:59"),
			ByName(h, "2013-12-17T20:15:58"),
			ByName(h, "2013-12-17T20:15:59"),
			ByName(h, "2013-12-18T20:15:59"),
			ByName(h, "2013-12-19T20:15:59"),
			ByName(h, "2013-12-20T20:15:59"),
			ByName(h, "2013-12-21T20:15:59"),
			ByName(h, "2013-12-22T20:15:59"),
			ByName(h, "2013-12-23T20:15:59"),
			ByName(h, "2013-12-24T20:15:59"),
		}
		result, err = getKeepDay(backups, now, timeparser.New())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(result, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(len(result), Is(8))
		if err != nil {
			t.Fatal(err)
		}
	}
}
