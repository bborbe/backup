package pkg

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/bborbe/errors"
	libtime "github.com/bborbe/time"
	"github.com/golang/glog"

	v1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
)

type BackupExectuor interface {
	Backup(ctx context.Context, target v1.BackupSpec) error
}

func NewBackupExectuor(
	currentTimeGetter libtime.CurrentTimeGetter,
	rsyncExectuor RsyncExectuor,
	backupRootDirectory BackupRootDirectory,
	sshPrivateKey SSHPrivateKey,
) BackupExectuor {
	return &backupExectuor{
		sshPrivateKey:       sshPrivateKey,
		currentTimeGetter:   currentTimeGetter,
		backupRootDirectory: backupRootDirectory,
		rsyncExectuor:       rsyncExectuor,
	}
}

type backupExectuor struct {
	currentTimeGetter   libtime.CurrentTimeGetter
	rsyncExectuor       RsyncExectuor
	backupRootDirectory BackupRootDirectory
	sshPrivateKey       SSHPrivateKey
}

func (b *backupExectuor) Backup(ctx context.Context, backupSpec v1.BackupSpec) error {
	if err := backupSpec.Validate(ctx); err != nil {
		return errors.Wrapf(ctx, err, "valid backup faild")
	}

	backupPath := b.backupPath(backupSpec)
	exists, err := exists(backupPath)
	if err != nil {
		return errors.Wrapf(ctx, err, "exists failed")
	}
	if exists {
		glog.V(2).Infof("backup %s already exists", backupPath)
		return nil
	}

	if err := b.createIncompleteIfNotExists(ctx, backupSpec); err != nil {
		return errors.Wrapf(ctx, err, "create incomplete if not exists failed")
	}

	if err := b.createCurrentIfNotExists(ctx, backupSpec); err != nil {
		return errors.Wrapf(ctx, err, "create current if not exists failed")
	}

	if err := b.runRsync(ctx, backupSpec); err != nil {
		return errors.Wrapf(ctx, err, "run rsync failed")
	}

	if err := b.renameIncomplete(ctx, backupSpec); err != nil {
		return errors.Wrapf(ctx, err, "rename incomplete if not exists failed")
	}

	if err := b.updateCurrentSymlink(ctx, backupSpec); err != nil {
		return errors.Wrapf(ctx, err, "update current symlink if not exists failed")
	}

	if err := b.removeEmpty(ctx, backupSpec); err != nil {
		return errors.Wrapf(ctx, err, "remove empty failed")
	}
	return nil
}

func (b *backupExectuor) createIncompleteIfNotExists(ctx context.Context, backupSpec v1.BackupSpec) error {
	incompletePath := b.incompletePath(backupSpec)
	if err := os.MkdirAll(incompletePath, os.ModePerm); err != nil {
		return errors.Wrapf(ctx, err, "create incomplete directory failed")
	}
	glog.V(3).Infof("create incomplete directory completed")
	return nil
}

func (b *backupExectuor) createCurrentIfNotExists(ctx context.Context, backupSpec v1.BackupSpec) error {
	currentPath := b.currentPath(backupSpec)
	emptyPath := b.emptyPath(backupSpec)
	currentExists, err := exists(currentPath)
	if err != nil {
		return errors.Wrapf(ctx, err, "check current exsits failed")
	}
	if currentExists {
		glog.V(3).Infof("current directory already exists")
		return nil
	}
	if err := os.MkdirAll(emptyPath, os.ModePerm); err != nil {
		return errors.Wrapf(ctx, err, "create incomplete directory failed")
	}
	glog.V(3).Infof("create empty directory completed")
	if err := os.Symlink("empty", currentPath); err != nil {
		return errors.Wrapf(ctx, err, "create symlink from empty to current failed")
	}
	glog.V(3).Infof("create current directory completed")
	return nil
}

func (b *backupExectuor) runRsync(ctx context.Context, backupSpec v1.BackupSpec) error {
	glog.V(3).Infof("rsync started")

	excludePath := b.excludePath(backupSpec)
	if err := os.WriteFile(excludePath, backupSpec.Excludes.Bytes(), 0644); err != nil {
		return errors.Wrapf(ctx, err, "write exclude failed")
	}

	args := []string{
		"-a",
		"-m",
		"--progress",
		"--compress",
		"--numeric-ids",
		"--delete",
		"--delete-excluded",
		"-e",
		fmt.Sprintf("ssh -T -x -o StrictHostKeyChecking=no -p %d -i %s", backupSpec.Port, b.sshPrivateKey),
		fmt.Sprintf("--exclude-from=%s", excludePath),
		fmt.Sprintf("--port=%d", backupSpec.Port),
		fmt.Sprintf("--link-dest=%s", b.currentPath(backupSpec)),
	}
	for _, dir := range backupSpec.Dirs {
		args = append(args, fmt.Sprintf("%s@%s:%s", backupSpec.User, backupSpec.Host, dir))
	}
	args = append(args, b.incompletePath(backupSpec))

	if err := b.rsyncExectuor.Rsync(ctx, args...); err != nil {
		return errors.Wrapf(ctx, err, "rsync failed")
	}

	glog.V(3).Infof("rsync completed")
	return nil
}

func (b *backupExectuor) renameIncomplete(ctx context.Context, backupSpec v1.BackupSpec) error {
	incompletePath := b.incompletePath(backupSpec)
	backupPath := b.backupPath(backupSpec)
	if err := os.Rename(incompletePath, backupPath); err != nil {
		return errors.Wrapf(ctx, err, "rename incomplete to backup failed")
	}
	return nil
}

func (b *backupExectuor) updateCurrentSymlink(ctx context.Context, backupSpec v1.BackupSpec) error {
	currentPath := b.currentPath(backupSpec)
	if err := os.Remove(currentPath); err != nil {
		return errors.Wrapf(ctx, err, "remove current path failed")
	}
	if err := os.Symlink(b.currentTimeGetter.Now().Format(time.DateOnly), currentPath); err != nil {
		return errors.Wrapf(ctx, err, "create symlink from empty to current failed")
	}
	glog.V(3).Infof("create current directory completed")
	return nil
}

func (b *backupExectuor) removeEmpty(ctx context.Context, backupSpec v1.BackupSpec) error {
	emptyPath := b.emptyPath(backupSpec)
	exists, err := exists(emptyPath)
	if err != nil {
		return errors.Wrapf(ctx, err, "check empty exists failed")
	}
	if exists == false {
		return nil
	}
	if err := os.Remove(emptyPath); err != nil {
		return errors.Wrapf(ctx, err, "remove empty failed")
	}
	glog.V(3).Infof("remove empty completed")
	return nil
}

func exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (b *backupExectuor) path(spec v1.BackupSpec, folderName string) string {
	return path.Join(
		b.backupRootDirectory.String(),
		spec.Host.String(),
		folderName,
	)
}

func (b *backupExectuor) emptyPath(backupSpec v1.BackupSpec) string {
	return b.path(
		backupSpec,
		"empty",
	)
}

func (b *backupExectuor) incompletePath(backupSpec v1.BackupSpec) string {
	return b.path(
		backupSpec,
		"incomplete",
	)
}

func (b *backupExectuor) currentPath(backupSpec v1.BackupSpec) string {
	return b.path(
		backupSpec,
		"current",
	)
}

func (b *backupExectuor) backupPath(backupSpec v1.BackupSpec) string {
	return b.path(
		backupSpec,
		b.currentTimeGetter.Now().Format(time.DateOnly),
	)
}

func (b *backupExectuor) excludePath(spec v1.BackupSpec) string {
	return fmt.Sprintf("/tmp/%s.excludes", spec.Host)
}
