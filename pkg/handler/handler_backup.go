package handler

import (
	"context"
	"net/http"

	"github.com/bborbe/errors"
	libhttp "github.com/bborbe/http"
	"github.com/gorilla/mux"

	"github.com/bborbe/backup/pkg"
)

func NewBackupHandler(k8sConnector pkg.K8sConnector, backupExectuor pkg.BackupExectuor) libhttp.WithError {
	return libhttp.WithErrorFunc(func(ctx context.Context, resp http.ResponseWriter, req *http.Request) error {
		vars := mux.Vars(req)

		target, err := k8sConnector.Target(ctx, vars["name"])
		if err != nil {
			return errors.Wrapf(ctx, err, "get target failed")
		}

		if err := backupExectuor.Backup(ctx, target.Spec); err != nil {
			return errors.Wrapf(ctx, err, "backup %s failed", target.Name)
		}
		libhttp.WriteAndGlog(resp, "backup %s completed", target.Name)
		return nil
	})
}
