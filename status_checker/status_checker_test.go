package status_checker

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/backup/dto"
	"github.com/bborbe/backup/service"
	"github.com/bborbe/backup/timeparser"
)

func TestImplementsStatusChecker(t *testing.T) {
	var backupService service.BackupService
	object := NewStatusChecker(backupService)
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
	statusDto := createStatusDto(dto.CreateHost(hostname), dto.CreateBackup("2014-01-01T12:23:45"), status)
	err = AssertThat(statusDto, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(statusDto.GetStatus(), Is(status))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(statusDto.GetHost(), Is(hostname))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateStatusDtoFalse(t *testing.T) {
	var err error
	hostname := "test"
	status := false
	statusDto := createStatusDto(dto.CreateHost(hostname), nil, status)
	err = AssertThat(statusDto, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(statusDto.GetStatus(), Is(status))
	if err != nil {
		t.Fatal(err)
	}
	err = AssertThat(statusDto.GetHost(), Is(hostname))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateStatusDtoForHostsEmptyHosts(t *testing.T) {
	var (
		err           error
		backupService service.BackupService
		hostDtos      []dto.Host
		statusDtos    []dto.Status
	)
	hostDtos = []dto.Host{}
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
		hostDtos   []dto.Host
		statusDtos []dto.Status
	)
	hostName := "fire.example.com"
	backupName := "2014-01-10T23:15:35"
	hostDtos = []dto.Host{
		createHostDto(hostName),
	}
	backupService := service.NewBackupServiceMock()
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
	err = AssertThat(statusDtos[0].GetHost(), Is(hostName))
	if err != nil {
		t.Fatal(err)
	}
}

func createHostDto(name string) dto.Host {
	host := dto.NewHost()
	host.SetName(name)
	return host
}

func createBackupDto(name string) dto.Backup {
	backup := dto.NewBackup()
	backup.SetName(name)
	return backup
}

func TestBackupIsInLastReturnTrueIfBackupDateIsNow(t *testing.T) {
	timeParser := timeparser.New()
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
	timeParser := timeparser.New()
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
	timeParser := timeparser.New()
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
