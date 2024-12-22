package pkg

import (
	"context"

	libtime "github.com/bborbe/time"

	v1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
)

type BackupCleaner interface {
	Clean(ctx context.Context, backupSpec v1.BackupSpec) error
}

func NewBackupCleaner(
	currentTimeGetter libtime.CurrentTimeGetter,
	backupRootDirectory Path,
) BackupCleaner {
	return &backupCleaner{
		currentTimeGetter:   currentTimeGetter,
		backupRootDirectory: backupRootDirectory,
	}
}

type backupCleaner struct {
	currentTimeGetter   libtime.CurrentTimeGetter
	backupRootDirectory Path
}

func (b *backupCleaner) Clean(ctx context.Context, backupSpec v1.BackupSpec) error {
	return nil
}
