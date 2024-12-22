package pkg

import (
	"context"
	v1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
)

type BackupFinder interface {
	List(ctx context.Context, host v1.BackupHost) (Backups, error)
}

func NewBackupFinder(backupRootDir Path) BackupFinder {
	return &backupFinder{backupRootDir: backupRootDir}
}

type backupFinder struct {
	backupRootDir Path
}

func (b *backupFinder) List(ctx context.Context, host v1.BackupHost) (Backups, error) {
	//hostDir := b.backupRootDir.Join(host.String())

	return nil, nil
}
