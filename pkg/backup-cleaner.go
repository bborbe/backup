package pkg

import (
	"context"
	"sort"

	"github.com/bborbe/errors"
	libtime "github.com/bborbe/time"
	"github.com/golang/glog"

	v1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
)

type BackupCleaner interface {
	Clean(ctx context.Context, backupSpec v1.BackupHost) error
}

func NewBackupCleaner(
	currentTimeGetter libtime.CurrentTimeGetter,
	backupFinder BackupFinder,
) BackupCleaner {
	return &backupCleaner{
		backupFinder:      backupFinder,
		currentTimeGetter: currentTimeGetter,
		keepAmount:        2,
	}
}

type backupCleaner struct {
	currentTimeGetter libtime.CurrentTimeGetter
	backupFinder      BackupFinder
	keepAmount        int
}

func (b *backupCleaner) Clean(ctx context.Context, backupHost v1.BackupHost) error {
	glog.V(2).Infof("backup host %s started", backupHost)
	dates, err := b.backupFinder.List(ctx, backupHost)
	if err != nil {
		return errors.Wrapf(ctx, err, "list backups failed")
	}
	if len(dates) <= b.keepAmount {
		glog.V(2).Infof("host %s has nothing to clean", backupHost)
		return nil
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Time().After(dates[j].Time())
	})
	for i, date := range dates {
		if i < b.keepAmount {
			glog.V(2).Infof("keep backup %s/%s", backupHost, date)
			continue
		}
		glog.V(2).Infof("delete backup %s/%s", backupHost, date)
	}
	glog.V(2).Infof("backup host %s completed", backupHost)
	return nil
}
