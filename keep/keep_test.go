package keep

import (
	"testing"
	"time"
	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/dto"
)

func TestAgeLessThan7Days(t *testing.T) {
	var err error
	{
		ti, _ := getTimeByName("2013-12-24T20:15:59")
		now, _ := getTimeByName("2013-12-24T20:15:59")
		err = AssertThat(ageLessThanDays(ti, now, 7), Is(true))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		ti, _ := getTimeByName("2013-12-17T20:15:59")
		now, _ := getTimeByName("2013-12-24T20:15:59")
		err = AssertThat(ageLessThanDays(ti, now, 7), Is(true))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		ti, _ := getTimeByName("2013-12-17T20:15:58")
		now, _ := getTimeByName("2013-12-24T20:15:59")
		err = AssertThat(ageLessThanDays(ti, now, 7), Is(false))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetKeepMonth(t *testing.T) {
	var err error
	var result []dto.Backup
	{
		backups := []dto.Backup{}
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
		backups := []dto.Backup{
			dto.CreateBackup("2013-12-12T24:15:59"),
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
		backups := []dto.Backup{
			dto.CreateBackup("2013-12-12T24:15:59"),
			dto.CreateBackup("2013-12-01T24:15:59"),
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
		err = AssertThat("2013-12-01T24:15:59", Is(result[0].GetName()))
	}
	{
		backups := []dto.Backup{
			dto.CreateBackup("2012-11-12T24:15:59"),
			dto.CreateBackup("2013-11-01T24:15:59"),
			dto.CreateBackup("2013-05-28T24:15:59"),
			dto.CreateBackup("2013-05-29T24:15:59"),
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
	var result []dto.Backup
	now, err := getTimeByName("2013-12-24T20:15:59")
	if err != nil {
		t.Fatal(err)
	}
	{
		backups := []dto.Backup{}
		result, err = getKeepToday(backups, now)
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
		backups := []dto.Backup{
			dto.CreateBackup("2013-12-23T10:15:59"),
			dto.CreateBackup("2013-12-24T15:15:59"),
			dto.CreateBackup("2013-12-25T20:15:59"),
		}
		result, err = getKeepToday(backups, now)
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
		_, err := getTimeByName("")
		err = AssertThat(err, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		calcTime, err := getTimeByName("2013-07-01T00:24:52")
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
		backups []dto.Backup
		backup  dto.Backup
	)
	{
		backups = []dto.Backup{}
		backup = latestBackup(backups)
		err = AssertThat(backup, NilValue())
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		backups := []dto.Backup{
			dto.CreateBackup("2013-12-06T20:15:59"),
		}
		backup = latestBackup(backups)
		err = AssertThat(backup, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(backup.GetName(), Is("2013-12-06T20:15:59"))
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		backups := []dto.Backup{
			dto.CreateBackup("2013-12-06T20:15:55"),
			dto.CreateBackup("2013-12-06T20:15:54"),
			dto.CreateBackup("2013-12-06T20:15:53"),
			dto.CreateBackup("2013-12-06T20:15:56"),
			dto.CreateBackup("2013-12-06T20:15:52"),
		}
		backup = latestBackup(backups)
		err = AssertThat(backup, NotNilValue())
		if err != nil {
			t.Fatal(err)
		}
		err = AssertThat(backup.GetName(), Is("2013-12-06T20:15:56"))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetKeepWeek(t *testing.T) {
	var result []dto.Backup
	now, err := getTimeByName("2013-12-24T20:15:59")
	if err != nil {
		t.Fatal(err)
	}
	{
		backups := []dto.Backup{}
		result, err = getKeepWeek(backups, now)
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
		backups := []dto.Backup{
			dto.CreateBackup("2013-12-06T20:15:59"),
			dto.CreateBackup("2013-12-07T20:15:59"),
			dto.CreateBackup("2013-12-08T20:15:59"),
			dto.CreateBackup("2013-12-09T20:15:59"),
			dto.CreateBackup("2013-12-10T20:15:59"),
			dto.CreateBackup("2013-12-11T20:15:59"),
			dto.CreateBackup("2013-12-12T20:15:59"),
			dto.CreateBackup("2013-12-13T20:15:59"),
			dto.CreateBackup("2013-12-14T20:15:59"),
			dto.CreateBackup("2013-12-15T20:15:59"),
			dto.CreateBackup("2013-12-16T20:15:59"),
			dto.CreateBackup("2013-12-17T20:15:59"),
			dto.CreateBackup("2013-12-18T20:15:59"),
			dto.CreateBackup("2013-12-19T20:15:59"),
			dto.CreateBackup("2013-12-20T20:15:59"),
			dto.CreateBackup("2013-12-21T20:15:59"),
			dto.CreateBackup("2013-12-22T20:15:59"),
			dto.CreateBackup("2013-12-23T20:15:59"),
			dto.CreateBackup("2013-12-24T20:15:59"),
		}
		result, err = getKeepWeek(backups, now)
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
	var result []dto.Backup
	now, err := getTimeByName("2013-12-24T20:15:59")
	if err != nil {
		t.Fatal(err)
	}
	{
		backups := []dto.Backup{}
		result, err = getKeepDay(backups, now)
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
		backups := []dto.Backup{
			dto.CreateBackup("2013-12-16T20:15:59"),
			dto.CreateBackup("2013-12-17T20:15:58"),
			dto.CreateBackup("2013-12-17T20:15:59"),
			dto.CreateBackup("2013-12-18T20:15:59"),
			dto.CreateBackup("2013-12-19T20:15:59"),
			dto.CreateBackup("2013-12-20T20:15:59"),
			dto.CreateBackup("2013-12-21T20:15:59"),
			dto.CreateBackup("2013-12-22T20:15:59"),
			dto.CreateBackup("2013-12-23T20:15:59"),
			dto.CreateBackup("2013-12-24T20:15:59"),
		}
		result, err = getKeepDay(backups, now)
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
