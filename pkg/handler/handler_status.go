// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/bborbe/errors"
	libhttp "github.com/bborbe/http"
	libtime "github.com/bborbe/time"

	v1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
	"github.com/bborbe/backup/pkg"
)

func NewStatusHandler(
	k8sConnector pkg.K8sConnector,
	backupFinder pkg.BackupFinder,
) libhttp.WithError {
	return libhttp.NewJsonHandler(
		libhttp.JsonHandlerFunc(func(ctx context.Context, req *http.Request) (interface{}, error) {

			targets, err := k8sConnector.Targets(ctx)
			if err != nil {
				return nil, errors.Wrapf(ctx, err, "list targets failed")
			}
			result := map[v1.BackupHost]string{}
			for _, target := range targets {
				dates, err := backupFinder.List(ctx, target.Spec.Host)
				if err != nil {
					return nil, errors.Wrapf(ctx, err, "get dates failed")
				}
				var latestBackup *libtime.Date
				for _, backupTime := range dates {
					if latestBackup == nil || backupTime.Time().After(latestBackup.Time()) {
						latestBackup = backupTime.Ptr()
						result[target.Spec.Host] = backupTime.Format(time.DateOnly)
					}
				}
			}
			return result, nil
		}),
	)
}
