// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package factory

import (
	"context"
	"net/http"

	"github.com/bborbe/cron"
	libcron "github.com/bborbe/cron"
	"github.com/bborbe/errors"
	libhttp "github.com/bborbe/http"
	"github.com/bborbe/k8s"
	libk8s "github.com/bborbe/k8s"
	"github.com/bborbe/run"
	libsentry "github.com/bborbe/sentry"
	libtime "github.com/bborbe/time"

	"github.com/bborbe/backup/pkg"
	"github.com/bborbe/backup/pkg/handler"
)

func CreateCleanupCron(
	sentryClient libsentry.Client,
	backupCleaner pkg.BackupCleaner,
	backupClientset pkg.BackupClientset,
	apiextensionsInterface libk8s.ApiextensionsInterface,
	namespace k8s.Namespace,
	cronExpression libcron.Expression,
) run.Func {
	return func(ctx context.Context) error {
		backupAction := CreateCleanAction(
			sentryClient,
			backupCleaner,
			backupClientset,
			apiextensionsInterface,
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
	backupClientset pkg.BackupClientset,
	apiextensionsInterface libk8s.ApiextensionsInterface,
	namespace k8s.Namespace,
	cronExpression libcron.Expression,
) run.Func {
	return func(ctx context.Context) error {
		backupAction := CreateBackupAction(
			sentryClient,
			backupExectuor,
			backupClientset,
			apiextensionsInterface,
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
	backupClientset pkg.BackupClientset,
	apiextensionsInterface libk8s.ApiextensionsInterface,
	namespace k8s.Namespace,
) run.Runnable {
	return pkg.NewCleanAction(
		sentryClient,
		pkg.NewK8sConnector(
			backupClientset,
			apiextensionsInterface,
			namespace,
		),
		backupCleaner,
	)
}

func CreateBackupAction(
	sentryClient libsentry.Client,
	backupExectuor pkg.BackupExectuor,
	backupClientset pkg.BackupClientset,
	apiextensionsInterface libk8s.ApiextensionsInterface,
	namespace k8s.Namespace,
) run.Runnable {
	return pkg.NewBackupAction(
		sentryClient,
		pkg.NewK8sConnector(
			backupClientset,
			apiextensionsInterface,
			namespace,
		),
		backupExectuor,
	)
}

func CreateSetupResourceDefinition(
	backupClientset pkg.BackupClientset,
	apiextensionsInterface libk8s.ApiextensionsInterface,
	namespace k8s.Namespace,
	trigger run.Fire,
) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		k8sConnector := pkg.NewK8sConnector(
			backupClientset,
			apiextensionsInterface,
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
	backupClientset pkg.BackupClientset,
	apiextensionsInterface libk8s.ApiextensionsInterface,
	namespace k8s.Namespace,
	backupRootDir pkg.Path,
) http.Handler {
	return libhttp.NewErrorHandler(
		handler.NewStatusHandler(
			pkg.NewK8sConnector(
				backupClientset,
				apiextensionsInterface,
				namespace,
			),
			pkg.NewBackupFinder(
				backupRootDir,
			),
		),
	)
}

func CreateListHandler(
	backupClientset pkg.BackupClientset,
	apiextensionsInterface libk8s.ApiextensionsInterface,
	namespace k8s.Namespace) http.Handler {
	return libhttp.NewErrorHandler(
		handler.NewListHandler(
			pkg.NewK8sConnector(
				backupClientset,
				apiextensionsInterface,
				namespace,
			),
		),
	)
}

func CreateBackupHandler(backupClientset pkg.BackupClientset, apiextensionsInterface libk8s.ApiextensionsInterface, namespace k8s.Namespace, backupExectuor pkg.BackupExectuor) http.Handler {
	return libhttp.NewErrorHandler(
		handler.NewBackupHandler(
			pkg.NewK8sConnector(
				backupClientset,
				apiextensionsInterface,
				namespace,
			),
			backupExectuor,
		),
	)
}

func CreateCleanupHandler(backupClientset pkg.BackupClientset, apiextensionsInterface libk8s.ApiextensionsInterface, namespace k8s.Namespace, backupCleaner pkg.BackupCleaner) http.Handler {
	return libhttp.NewErrorHandler(
		handler.NewCleanupHandler(
			pkg.NewK8sConnector(
				backupClientset,
				apiextensionsInterface,
				namespace,
			),
			backupCleaner,
		),
	)
}

func CreateBackupCleaner(
	currentTimeGetter libtime.CurrentTimeGetter,
	backupRootDirectory pkg.Path,
	backupKeepAmount int,
	backupCleanEnabled bool,
) pkg.BackupCleaner {
	return pkg.NewBackupCleanerOnlyOnce(
		pkg.NewBackupCleaner(
			currentTimeGetter,
			pkg.NewBackupFinder(
				backupRootDirectory,
			),
			backupRootDirectory,
			backupKeepAmount,
			backupCleanEnabled,
		),
	)
}

func CreateBackupActionHandler(ctx context.Context, sentryClient libsentry.Client, backupExectuor pkg.BackupExectuor, backupClientset pkg.BackupClientset, apiextensionClientset libk8s.ApiextensionsInterface, namespace k8s.Namespace) http.Handler {
	return libhttp.NewBackgroundRunHandler(ctx,
		CreateBackupAction(
			sentryClient,
			backupExectuor,
			backupClientset,
			apiextensionClientset,
			namespace,
		).Run,
	)
}

func CreateCleanActionHandler(ctx context.Context, sentryClient libsentry.Client, backupCleaner pkg.BackupCleaner, backupClientset pkg.BackupClientset, apiextensionClientset libk8s.ApiextensionsInterface, namespace libk8s.Namespace) http.Handler {
	return libhttp.NewBackgroundRunHandler(ctx,
		CreateCleanAction(
			sentryClient,
			backupCleaner,
			backupClientset,
			apiextensionClientset,
			namespace,
		).Run,
	)
}
