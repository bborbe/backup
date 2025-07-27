// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import (
	"context"
	"os"
	"sort"

	"github.com/bborbe/errors"
	libtime "github.com/bborbe/time"
	"github.com/golang/glog"

	v1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
)

//counterfeiter:generate -o ../mocks/backup-cleaner.go --fake-name BackupCleaner . BackupCleaner
type BackupCleaner interface {
	Clean(ctx context.Context, backupSpec v1.BackupHost) error
}

func NewBackupCleaner(
	currentTimeGetter libtime.CurrentTimeGetter,
	backupFinder BackupFinder,
	backupRootDir Path,
	backupKeepAmount int,
	backupCleanEnabled bool,
) BackupCleaner {
	return &backupCleaner{
		backupFinder:       backupFinder,
		currentTimeGetter:  currentTimeGetter,
		backupKeepAmount:   backupKeepAmount,
		backupRootDir:      backupRootDir,
		backupCleanEnabled: backupCleanEnabled,
	}
}

type backupCleaner struct {
	currentTimeGetter  libtime.CurrentTimeGetter
	backupFinder       BackupFinder
	backupKeepAmount   int
	backupCleanEnabled bool
	backupRootDir      Path
}

func (b *backupCleaner) Clean(ctx context.Context, backupHost v1.BackupHost) error {
	glog.V(2).Infof("backup host %s started", backupHost)
	dates, err := b.backupFinder.List(ctx, backupHost)
	if err != nil {
		return errors.Wrapf(ctx, err, "list backups failed")
	}
	if len(dates) <= b.backupKeepAmount {
		glog.V(2).Infof("host %s has nothing to clean", backupHost)
		return nil
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Time().After(dates[j].Time())
	})
	for i, date := range dates {
		if i < b.backupKeepAmount {
			glog.V(2).Infof("keep backup %s/%s", backupHost, date)
			continue
		}
		if b.backupCleanEnabled == false {
			glog.V(2).Infof("would delete backup %s/%s", backupHost, date)
			continue
		}
		glog.V(2).Infof("delete backup %s/%s", backupHost, date)
		dir := b.backupRootDir.Join(backupHost.String(), date.String())
		glog.V(2).Infof("os.RemoveAll(%s) started", dir)
		if err := os.RemoveAll(dir.String()); err != nil {
			return errors.Wrapf(ctx, err, "os.RemoveAll(%s) failed", dir)
		}
		glog.V(2).Infof("os.RemoveAll(%s) completed", dir)
	}
	glog.V(2).Infof("backup host %s completed", backupHost)
	return nil
}
