// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"context"
	"net/http"

	"github.com/bborbe/errors"
	libhttp "github.com/bborbe/http"
	"github.com/gorilla/mux"

	"github.com/bborbe/backup/pkg"
)

func NewBackupHandler(targetFinder pkg.TargetFinder, backupExectuor pkg.BackupExectuor) libhttp.WithError {
	return libhttp.WithErrorFunc(func(ctx context.Context, resp http.ResponseWriter, req *http.Request) error {
		vars := mux.Vars(req)

		target, err := targetFinder.Target(ctx, vars["name"])
		if err != nil {
			return err
		}

		if err := backupExectuor.Backup(ctx, target.Spec); err != nil {
			return errors.Wrapf(ctx, err, "backup %s failed", target.Name)
		}
		_, _ = libhttp.WriteAndGlog(resp, "backup %s completed", target.Name)
		return nil
	})
}
