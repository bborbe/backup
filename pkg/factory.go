package pkg

import (
	"context"

	"github.com/bborbe/cron"
	"github.com/bborbe/errors"
	"github.com/bborbe/k8s"
	"github.com/bborbe/run"
	libsentry "github.com/bborbe/sentry"
	libtime "github.com/bborbe/time"
)

func CreateBackupCron(
	sentryClient libsentry.Client,
	currentTimeGetter libtime.CurrentTimeGetter,
	kubeConfig string,
	backupRootDirectory BackupRootDirectory,
	sshKeyPath SSHPrivateKey,
	namespace k8s.Namespace,
	cronExpression string,
) run.Func {
	return func(ctx context.Context) error {
		return cron.NewExpressionCron(
			cronExpression,
			libsentry.NewSkipErrorAndReport(
				sentryClient,
				CreateBackupAction(
					currentTimeGetter,
					kubeConfig,
					backupRootDirectory,
					sshKeyPath,
					namespace,
				),
			),
		).Run(ctx)
	}
}

func CreateBackupAction(
	currentTimeGetter libtime.CurrentTimeGetter,
	kubeConfig string,
	backupRootDirectory BackupRootDirectory,
	sshKeyPath SSHPrivateKey,
	namespace k8s.Namespace,
) run.Runnable {
	return NewBackupAction(
		NewK8sConnector(
			kubeConfig,
			namespace,
		),
		CreateBackupExectuor(
			currentTimeGetter,
			backupRootDirectory,
			sshKeyPath,
		),
	)
}

func CreateSetupResourceDefinition(
	kubeConfig string,
	namespace k8s.Namespace,
	trigger run.Fire,
) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		k8sConnector := NewK8sConnector(
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
	backupRootDirectory BackupRootDirectory,
	sshPrivateKey SSHPrivateKey,
) BackupExectuor {
	return NewBackupExectuor(
		currentTimeGetter,
		NewRsyncExectuor(),
		backupRootDirectory,
		sshPrivateKey,
	)
}
