package checker

import (
	"testing"

	"os"

	. "github.com/bborbe/assert"
	backup_dto "github.com/bborbe/backup/dto"
	backup_service "github.com/bborbe/backup/service"
	backup_timeparser "github.com/bborbe/backup/timeparser"
	"github.com/golang/glog"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	glog.Flush()
	os.Exit(exit)
}

func TestImplementsStatusChecker(t *testing.T) {
	var backupService backup_service.BackupService
	object := New(backupService)
	var expected *StatusChecker
	err := AssertThat(object, Implements(expected))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateStatusDtoTrue(t *testing.T) {
	var err error
	hostname := "test"
	status := true
	statusDto := createStatusDto(backup_dto.CreateHost(hostname), backup_dto.CreateBackup("2014-01-01T12:23:45"), status)
	err = AssertThat(statusDto, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(statusDto.Status, Is(status))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(statusDto.Host, Is(hostname))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateStatusDtoFalse(t *testing.T) {
	var err error
	hostname := "test"
	status := false
	statusDto := createStatusDto(backup_dto.CreateHost(hostname), nil, status)
	err = AssertThat(statusDto, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(statusDto.Status, Is(status))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(statusDto.Host, Is(hostname))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateStatusDtoForHostsEmptyHosts(t *testing.T) {
	var (
		err           error
		backupService backup_service.BackupService
		hostDtos      []backup_dto.Host
		statusDtos    []backup_dto.Status
	)
	hostDtos = []backup_dto.Host{}
	statusDtos, err = createStatusDtoForHosts(backupService, hostDtos)
	err = AssertThat(statusDtos, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(statusDtos), Is(len(hostDtos)))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateStatusDtoForHostsOneHost(t *testing.T) {
	var (
		err        error
		hostDtos   []backup_dto.Host
		statusDtos []backup_dto.Status
	)
	hostName := "fire.example.com"
	backupName := "2014-01-10T23:15:35"
	hostDtos = []backup_dto.Host{
		createHostDto(hostName),
	}
	backupService := backup_service.NewBackupServiceMock()
	backupService.SetLatestBackup(createBackupDto(backupName), nil)
	statusDtos, err = createStatusDtoForHosts(backupService, hostDtos)
	err = AssertThat(statusDtos, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(len(statusDtos), Is(len(hostDtos)))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(statusDtos[0].Host, Is(hostName))
	if err != nil {
		t.Fatal(err)
	}
}

func createHostDto(name string) backup_dto.Host {
	host := backup_dto.NewHost()
	host.SetName(name)
	return host
}

func createBackupDto(name string) backup_dto.Backup {
	backup := backup_dto.NewBackup()
	backup.SetName(name)
	return backup
}

func TestBackupIsInLastReturnTrueIfBackupDateIsNow(t *testing.T) {
	timeParser := backup_timeparser.New()
	now, err := timeParser.TimeByName("2014-01-01T12:45:59")
	if err != nil {
		t.Fatal(err)
	}
	result, err := backupIsInLastDays(createBackupDto("2014-01-01T12:45:59"), timeParser, now)
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(result, Is(true))
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackupIsInLastDaysReturnTrueIfDivIsLessThanSevenDays(t *testing.T) {
	timeParser := backup_timeparser.New()
	now, err := timeParser.TimeByName("2014-01-02T12:45:59")
	if err != nil {
		t.Fatal(err)
	}
	result, err := backupIsInLastDays(createBackupDto("2014-01-01T12:45:59"), timeParser, now)
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(result, Is(true))
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackupIsInLastDaysReturnFalseIfBackupDateIfMoreThanSevenDaysBefore(t *testing.T) {
	timeParser := backup_timeparser.New()
	now, err := timeParser.TimeByName("2014-01-09T12:45:59")
	if err != nil {
		t.Fatal(err)
	}
	result, err := backupIsInLastDays(createBackupDto("2014-01-01T12:45:59"), timeParser, now)
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(result, Is(false))
	if err != nil {
		t.Fatal(err)
	}
}
