package pkg

import (
	"context"

	"github.com/bborbe/errors"
	"github.com/bborbe/run"
	"github.com/golang/glog"
)

func NewBackupAction(
	k8sConnector K8sConnector,
	backupExectuor BackupExectuor,
) run.Runnable {
	return run.Func(func(ctx context.Context) error {
		glog.V(2).Infof("backup cron started")
		targets, err := k8sConnector.Targets(ctx)
		if err != nil {
			return errors.Wrapf(ctx, err, "get target failed")
		}
		glog.V(2).Infof("found %d targets to backup", len(targets))
		for _, target := range targets {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				glog.V(2).Infof("backup %s started", target.Name)
				if err := backupExectuor.Backup(ctx, target.Spec); err != nil {
					return errors.Wrapf(ctx, err, "backup %s failed", target.Name)
				}
				glog.V(2).Infof("backup %s completed", target.Name)
			}
		}
		glog.V(2).Infof("backup cron completed")
		return nil
	})
}
