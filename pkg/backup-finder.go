package pkg

import (
	"context"
	"os"

	"github.com/bborbe/errors"
	libtime "github.com/bborbe/time"
	"github.com/golang/glog"

	v1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
)

type BackupFinder interface {
	List(ctx context.Context, host v1.BackupHost) ([]libtime.Date, error)
}

func NewBackupFinder(backupRootDir Path) BackupFinder {
	return &backupFinder{backupRootDir: backupRootDir}
}

type backupFinder struct {
	backupRootDir Path
}

func (b *backupFinder) List(ctx context.Context, backupHost v1.BackupHost) ([]libtime.Date, error) {
	host := backupHost
	backupDir := b.backupRootDir.Join(host.String())
	glog.V(4).Infof("search for backups in %s", backupDir)
	entries, err := os.ReadDir(backupDir.String())
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "list failed")
	}
	glog.V(4).Infof("found %d entries in %s", len(entries), backupDir)

	var result []libtime.Date
	for _, entry := range entries {
		date, err := libtime.ParseDate(ctx, entry.Name())
		if err != nil {
			glog.V(4).Infof("name(%s) is not valid  => skip", entry.Name())
			continue
		}
		result = append(result, *date)
	}
	return result, nil
}
