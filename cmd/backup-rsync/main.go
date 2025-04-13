// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"os"

	"github.com/bborbe/errors"
	libsentry "github.com/bborbe/sentry"
	"github.com/bborbe/service"
	libtime "github.com/bborbe/time"
	"github.com/golang/glog"

	v1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
	"github.com/bborbe/backup/pkg"
	"github.com/bborbe/backup/pkg/factory"
)

func main() {
	app := &application{}
	os.Exit(service.Main(context.Background(), app, &app.SentryDSN, &app.SentryProxy))
}

type application struct {
	SentryDSN      string `required:"true" arg:"sentry-dsn" env:"SENTRY_DSN" usage:"SentryDSN" display:"length"`
	SentryProxy    string `required:"false" arg:"sentry-proxy" env:"SENTRY_PROXY" usage:"Sentry Proxy"`
	BackupRootDir  string `required:"true" arg:"backup-root-dir" env:"BACKUP_ROOT_DIR" usage:"Directory all backups are stored"`
	BackupHost     string `required:"true" arg:"backup-host" env:"BACKUP_HOST" usage:"host used to connect to client to backup"`
	BackupPort     int    `required:"true" arg:"backup-port" env:"BACKUP_PORT" usage:"port used to connect to client to backup" default:"22"`
	BackupUser     string `required:"true" arg:"backup-user" env:"BACKUP_USER" usage:"user used to connect to client to backup" default:"root"`
	BackupDirs     string `required:"true" arg:"backup-dirs" env:"BACKUP_DIRS" usage:"comma seperated list of directoies to backup"`
	BackupExcludes string `required:"false" arg:"backup-excludes" env:"BACKUP_EXCLUDES" usage:"comma seperated excludes"`
	SSHPrivateKey  string `required:"true" arg:"ssh-key" env:"SSH_KEY" usage:"path to ssh private key"`
}

func (a *application) Run(ctx context.Context, sentryClient libsentry.Client) error {
	currentTime := libtime.NewCurrentTime()
	backupExectuor := factory.CreateBackupExectuor(
		currentTime,
		pkg.Path(a.BackupRootDir),
		pkg.SSHPrivateKey(a.SSHPrivateKey),
	)
	err := backupExectuor.Backup(ctx, v1.BackupSpec{
		Host:     v1.BackupHost(a.BackupHost),
		Port:     v1.BackupPort(a.BackupPort),
		User:     v1.BackupUser(a.BackupUser),
		Dirs:     v1.ParseBackupDirsFromString(a.BackupDirs),
		Excludes: v1.ParseBackupExcludesFromString(a.BackupExcludes),
	})
	if err != nil {
		return errors.Wrapf(ctx, err, "backup failed")
	}
	glog.V(2).Infof("backup completed")
	return nil
}
