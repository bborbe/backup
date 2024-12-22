package pkg

import (
	"context"
	stderrors "errors"
	"sync"

	v1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
)

var CleanupAlreadyRunningError = stderrors.New("cleanup already running")

func NewBackupCleanerOnlyOnce(
	backupCleaner BackupCleaner,
) BackupCleaner {
	return &backupCleanerOnlyOnce{
		backupCleaner: backupCleaner,
	}
}

type backupCleanerOnlyOnce struct {
	mux           sync.Mutex
	running       bool
	backupCleaner BackupCleaner
}

func (b *backupCleanerOnlyOnce) Clean(ctx context.Context, backupHost v1.BackupHost) error {
	b.mux.Lock()
	if b.running {
		b.mux.Unlock()
		return CleanupAlreadyRunningError
	}
	b.running = true
	b.mux.Unlock()
	err := b.backupCleaner.Clean(ctx, backupHost)
	b.mux.Lock()
	b.running = false
	b.mux.Unlock()
	return err
}
