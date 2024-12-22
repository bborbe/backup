package factory

import (
	"context"

	"github.com/bborbe/cron"
	libcron "github.com/bborbe/cron"
	"github.com/bborbe/errors"
	libhttp "github.com/bborbe/http"
	"github.com/bborbe/k8s"
	"github.com/bborbe/run"
	libsentry "github.com/bborbe/sentry"
	libtime "github.com/bborbe/time"

	"github.com/bborbe/backup/pkg"
	"github.com/bborbe/backup/pkg/handler"
)

func CreateCleanupCron(
	sentryClient libsentry.Client,
	backupCleaner pkg.BackupCleaner,
	kubeConfig string,
	namespace k8s.Namespace,
	cronExpression libcron.Expression,
) run.Func {
	return func(ctx context.Context) error {
		backupAction := CreateCleanAction(
			sentryClient,
			backupCleaner,
			kubeConfig,
			namespace,
		)
		parallelSkipper := run.NewParallelSkipper()
		return cron.NewExpressionCron(
			cronExpression,
			libsentry.NewSkipErrorAndReport(
				sentryClient,
				parallelSkipper.SkipParallel(backupAction.Run),
			),
		).Run(ctx)
	}
}

func CreateBackupCron(
	sentryClient libsentry.Client,
	backupExectuor pkg.BackupExectuor,
	kubeConfig string,
	namespace k8s.Namespace,
	cronExpression libcron.Expression,
) run.Func {
	return func(ctx context.Context) error {
		backupAction := CreateBackupAction(
			sentryClient,
			backupExectuor,
			kubeConfig,
			namespace,
		)
		parallelSkipper := run.NewParallelSkipper()
		return cron.NewExpressionCron(
			cronExpression,
			libsentry.NewSkipErrorAndReport(
				sentryClient,
				parallelSkipper.SkipParallel(backupAction.Run),
			),
		).Run(ctx)
	}
}

func CreateCleanAction(
	sentryClient libsentry.Client,
	backupCleaner pkg.BackupCleaner,
	kubeConfig string,
	namespace k8s.Namespace,
) run.Runnable {
	return pkg.NewCleanAction(
		sentryClient,
		pkg.NewK8sConnector(
			kubeConfig,
			namespace,
		),
		backupCleaner,
	)
}

func CreateBackupAction(
	sentryClient libsentry.Client,
	backupExectuor pkg.BackupExectuor,
	kubeConfig string,
	namespace k8s.Namespace,
) run.Runnable {
	return pkg.NewBackupAction(
		sentryClient,
		pkg.NewK8sConnector(
			kubeConfig,
			namespace,
		),
		backupExectuor,
	)
}

func CreateSetupResourceDefinition(
	kubeConfig string,
	namespace k8s.Namespace,
	trigger run.Fire,
) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		k8sConnector := pkg.NewK8sConnector(
			kubeConfig,
			namespace,
		)
		if err := k8sConnector.SetupCustomResourceDefinition(ctx); err != nil {
			return errors.Wrap(ctx, err, "setup resource definition failed")
		}
		trigger.Fire()
		<-ctx.Done()
		return nil
	}
}

func CreateBackupExectuor(
	currentTimeGetter libtime.CurrentTimeGetter,
	backupRootDirectory pkg.Path,
	sshPrivateKey pkg.SSHPrivateKey,
) pkg.BackupExectuor {
	return pkg.NewBackupExectuorOnlyOnce(
		pkg.NewBackupExectuor(
			currentTimeGetter,
			pkg.NewRsyncExectuor(),
			backupRootDirectory,
			sshPrivateKey,
		),
	)
}

func CreateStatusHandler(
	kubeconfig string,
	namespace k8s.Namespace,
	backupRootDir pkg.Path,
) libhttp.WithError {
	return handler.NewStatusHandler(
		pkg.NewK8sConnector(
			kubeconfig,
			namespace,
		),
		pkg.NewBackupFinder(
			backupRootDir,
		),
	)
}

func CreateListHandler(kubeconfig string, namespace k8s.Namespace) libhttp.WithError {
	k8sConnector := pkg.NewK8sConnector(
		kubeconfig,
		namespace,
	)
	return handler.NewListHandler(k8sConnector)
}

func CreateBackupHandler(kubeconfig string, namespace k8s.Namespace, backupExectuor pkg.BackupExectuor) libhttp.WithError {
	return handler.NewBackupHandler(pkg.NewK8sConnector(
		kubeconfig,
		namespace,
	), backupExectuor)
}

func CreateBackupCleaner(
	currentTimeGetter libtime.CurrentTimeGetter,
	backupRootDirectory pkg.Path,
) pkg.BackupCleaner {
	return pkg.NewBackupCleanerOnlyOnce(
		pkg.NewBackupCleaner(
			currentTimeGetter,
			pkg.NewBackupFinder(
				backupRootDirectory,
			),
		),
	)
}
