// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"os"
	"time"

	libcron "github.com/bborbe/cron"
	libhttp "github.com/bborbe/http"
	"github.com/bborbe/k8s"
	"github.com/bborbe/log"
	"github.com/bborbe/run"
	libsentry "github.com/bborbe/sentry"
	"github.com/bborbe/service"
	libtime "github.com/bborbe/time"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/bborbe/backup/pkg"
	"github.com/bborbe/backup/pkg/factory"
)

func main() {
	app := &application{}
	os.Exit(service.Main(context.Background(), app, &app.SentryDSN, &app.SentryProxy))
}

type application struct {
	SentryDSN            string `required:"true" arg:"sentry-dsn" env:"SENTRY_DSN" usage:"SentryDSN" display:"length"`
	SentryProxy          string `required:"false" arg:"sentry-proxy" env:"SENTRY_PROXY" usage:"Sentry Proxy"`
	Listen               string `required:"true" arg:"listen" env:"LISTEN" usage:"address to listen to"`
	Kubeconfig           string `required:"false" arg:"kubeconfig" env:"KUBECONFIG" usage:"Path to k8s config"`
	CronExpression       string `required:"true" arg:"cron-expression" env:"CRON_EXPRESSION" usage:"Cron expression to determine when service is run" default:"@every 1h"`
	SSHPrivateKey        string `required:"true" arg:"ssh-key" env:"SSH_KEY" usage:"path to ssh private key"`
	Namespace            string `required:"true" arg:"namespace" env:"NAMESPACE" usage:"kubernetes namespace"`
	BackupRootDir        string `required:"true" arg:"backup-root-dir" env:"BACKUP_ROOT_DIR" usage:"Directory all backups are stored"`
	BackupCleanupEnabled bool   `required:"true" arg:"backup-cleanup-enabled" env:"BACKUP_CLEANUP_ENABLED" usage:"allow enable backup cleanup"`
	BackupKeepAmount     int    `required:"true" arg:"backup-keep-amount" env:"BACKUP_KEEP_AMOUNT" usage:"how many backups to keep" default:"3"`
}

func (a *application) Run(ctx context.Context, sentryClient libsentry.Client) error {
	currentTime := libtime.NewCurrentTime()

	backupExectuor := factory.CreateBackupExectuor(
		currentTime,
		pkg.Path(a.BackupRootDir),
		pkg.SSHPrivateKey(a.SSHPrivateKey),
	)

	backupCleaner := factory.CreateBackupCleaner(
		currentTime,
		pkg.Path(a.BackupRootDir),
		a.BackupKeepAmount,
		a.BackupCleanupEnabled,
	)

	trigger := run.NewTrigger()

	return service.Run(
		ctx,
		a.createSetupResourceDefinition(trigger),
		run.Triggered(a.createBackupCron(sentryClient, backupExectuor), trigger.Done()),
		run.Triggered(a.createCleanupCron(sentryClient, backupCleaner), trigger.Done()),
		a.createHttpServer(sentryClient, backupExectuor, backupCleaner),
	)
}

func (a *application) createBackupCron(sentryClient libsentry.Client, backupExectuor pkg.BackupExectuor) run.Func {
	return factory.CreateBackupCron(sentryClient, backupExectuor, a.Kubeconfig, k8s.Namespace(a.Namespace), libcron.Expression(a.CronExpression))
}

func (a *application) createCleanupCron(sentryClient libsentry.Client, backupCleaner pkg.BackupCleaner) run.Func {
	return factory.CreateCleanupCron(sentryClient, backupCleaner, a.Kubeconfig, k8s.Namespace(a.Namespace), libcron.Expression(a.CronExpression))
}

func (a *application) createSetupResourceDefinition(trigger run.Trigger) func(ctx context.Context) error {
	return factory.CreateSetupResourceDefinition(a.Kubeconfig, k8s.Namespace(a.Namespace), trigger)
}

func (a *application) createHttpServer(
	sentryClient libsentry.Client,
	backupExectuor pkg.BackupExectuor,
	backupCleaner pkg.BackupCleaner,
) run.Func {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		router := mux.NewRouter()
		router.Path("/healthz").Handler(libhttp.NewPrintHandler("OK"))
		router.Path("/readiness").Handler(libhttp.NewPrintHandler("OK"))
		router.Path("/metrics").Handler(promhttp.Handler())
		router.Path("/setloglevel/{level}").Handler(log.NewSetLoglevelHandler(ctx, log.NewLogLevelSetter(2, 5*time.Minute)))

		router.Path("/status").Handler(libhttp.NewErrorHandler(
			factory.CreateStatusHandler(a.Kubeconfig, k8s.Namespace(a.Namespace), pkg.Path(a.BackupRootDir)),
		))

		router.Path("/list").Handler(libhttp.NewErrorHandler(
			factory.CreateListHandler(a.Kubeconfig, k8s.Namespace(a.Namespace)),
		))

		router.Path("/backup/all").Handler(libhttp.NewBackgroundRunHandler(ctx,
			factory.CreateBackupAction(
				sentryClient,
				backupExectuor,
				a.Kubeconfig,
				k8s.Namespace(a.Namespace),
			).Run,
		))

		router.Path("/backup/{name}").Handler(libhttp.NewErrorHandler(
			factory.CreateBackupHandler(a.Kubeconfig, k8s.Namespace(a.Namespace), backupExectuor),
		))

		router.Path("/cleanup/all").Handler(libhttp.NewBackgroundRunHandler(ctx,
			factory.CreateCleanAction(
				sentryClient,
				backupCleaner,
				a.Kubeconfig,
				k8s.Namespace(a.Namespace),
			).Run,
		))

		router.Path("/cleanup/{name}").Handler(libhttp.NewErrorHandler(
			factory.CreateCleanupHandler(a.Kubeconfig, k8s.Namespace(a.Namespace), backupCleaner),
		))

		glog.V(2).Infof("starting http server listen on %s", a.Listen)
		return libhttp.NewServer(
			a.Listen,
			router,
		).Run(ctx)
	}
}
