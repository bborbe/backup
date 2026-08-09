// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"context"
	"errors"
	"net/http"

	bberrors "github.com/bborbe/errors"
	libhttp "github.com/bborbe/http"
	"github.com/gorilla/mux"

	"github.com/bborbe/backup/pkg"
)

// ErrorCodeBackupAlreadyRunning is the stable, machine-readable error code returned when a
// backup trigger is refused because a backup for the same host is already in flight.
const ErrorCodeBackupAlreadyRunning = "BACKUP_ALREADY_RUNNING"

func NewBackupHandler(
	targetFinder pkg.TargetFinder,
	backupExectuor pkg.BackupExectuor,
) libhttp.WithError {
	return libhttp.WithErrorFunc(
		func(ctx context.Context, resp http.ResponseWriter, req *http.Request) error {
			vars := mux.Vars(req)

			target, err := targetFinder.Target(ctx, vars["name"])
			if err != nil {
				return err
			}

			if err := backupExectuor.Backup(ctx, target.Spec); err != nil {
				if errors.Is(err, pkg.BackupAlreadyRunningError) {
					return libhttp.WrapWithCode(
						bberrors.Errorf(ctx, "backup for %s is already running", target.Spec.Host),
						ErrorCodeBackupAlreadyRunning,
						http.StatusConflict,
					)
				}
				return bberrors.Wrapf(ctx, err, "backup %s failed", target.Name)
			}
			_, _ = libhttp.WriteAndGlog(resp, "backup %s completed", target.Name)
			return nil
		},
	)
}
