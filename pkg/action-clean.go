package pkg

import (
	"context"

	"github.com/bborbe/errors"
	"github.com/bborbe/run"
	libsentry "github.com/bborbe/sentry"
	"github.com/getsentry/sentry-go"
	"github.com/golang/glog"
)

func NewCleanAction(
	sentryClient libsentry.Client,
	k8sConnector K8sConnector,
	cleanExectuor BackupCleaner,
) run.Runnable {
	return run.Func(func(ctx context.Context) error {
		glog.V(2).Infof("clean cron started")
		targets, err := k8sConnector.Targets(ctx)
		if err != nil {
			return errors.Wrapf(ctx, err, "get target failed")
		}
		glog.V(2).Infof("found %d targets to clean", len(targets))
		for _, target := range targets {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				glog.V(2).Infof("clean %s started", target.Name)
				if err := cleanExectuor.Clean(ctx, target.Spec.Host); err != nil {
					sentryClient.CaptureException(
						err,
						&sentry.EventHint{
							Context: ctx,
							Data: map[string]interface{}{
								"name":     target.Name,
								"host":     target.Spec.Host,
								"port":     target.Spec.Port,
								"user":     target.Spec.User,
								"dirs":     target.Spec.Dirs,
								"excludes": target.Spec.Excludes,
							},
						},
						nil,
					)
					glog.Warningf("clean %s failed: %v", target.Name, err)
					continue
				}
				glog.V(2).Infof("clean %s completed", target.Name)
			}
		}
		glog.V(2).Infof("clean cron completed")
		return nil
	})
}
