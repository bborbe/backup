package v1

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/bborbe/collection"
	"github.com/bborbe/errors"
	"github.com/bborbe/validation"
)

type BackupSpecs []BackupSpec

type BackupHost string

func (f BackupHost) String() string {
	return string(f)
}

type BackupPort int

func (f BackupPort) Int() int {
	return int(f)
}

func (f BackupPort) String() string {
	return strconv.Itoa(f.Int())
}

type BackupUser string

func (f BackupUser) String() string {
	return string(f)
}

func ParseBackupDirectoriesFromString(value string) BackupDirectories {
	return ParseBackupDirectories(strings.FieldsFunc(value, func(r rune) bool {
		return r == ','
	}))
}

func ParseBackupDirectories(values []string) BackupDirectories {
	result := make(BackupDirectories, len(values))
	for i, value := range values {
		result[i] = BackupDirectory(value)
	}
	return result
}

type BackupDirectories []BackupDirectory

func (a BackupDirectories) Len() int { return len(a) }
func (a BackupDirectories) Less(i, j int) bool {
	return strings.Compare(strings.ToLower(a[i].String()), strings.ToLower(a[j].String())) < 0
}
func (a BackupDirectories) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a BackupDirectories) Strings() []string {
	result := make([]string, len(a))
	for i, aa := range a {
		result[i] = aa.String()
	}
	return result
}

type BackupDirectory string

func (f BackupDirectory) String() string {
	return string(f)
}

// BackupSpec is the spec for a Foo resource
type BackupSpec struct {
	Host BackupHost        `json:"host" yaml:"host"`
	Port BackupPort        `json:"port" yaml:"port"`
	User BackupUser        `json:"user" yaml:"user"`
	Dirs BackupDirectories `json:"dirs" yaml:"dirs"`
}

func (a BackupSpec) Equal(backup BackupSpec) bool {
	if a.Host != backup.Host {
		return false
	}
	if a.Port != backup.Port {
		return false
	}
	if a.User != backup.User {
		return false
	}
	if collection.Equal(sortStrings(a.Dirs), sortStrings(backup.Dirs)) == false {
		return false
	}
	return true
}

func (a BackupSpec) Validate(ctx context.Context) error {
	if a.Host == "" {
		return errors.Wrap(ctx, validation.Error, "Host is empty")
	}
	if a.Port <= 0 {
		return errors.Wrap(ctx, validation.Error, "Port is invalid")
	}
	if a.User == "" {
		return errors.Wrap(ctx, validation.Error, "User is empty")
	}
	return nil
}

func sortStrings(backupDirectories BackupDirectories) BackupDirectories {
	sort.Sort(backupDirectories)
	return backupDirectories
}
