package pkg

import (
	"context"

	"github.com/bborbe/errors"
	"github.com/bborbe/k8s"
	"github.com/bborbe/run"
	libtime "github.com/bborbe/time"
)

func CreateBackupCron(
	currentTimeGetter libtime.CurrentTimeGetter,
	kubeConfig string,
	backupRootDirectory BackupRootDirectory,
	sshKeyPath SSHPrivateKey,
	namespace k8s.Namespace,
) BackupCron {
	return NewBackupCron(
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
