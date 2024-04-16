// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/bborbe/errors"
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
)

func main() {
	app := &application{}
	os.Exit(service.Main(context.Background(), app, &app.SentryDSN))
}

type application struct {
	SentryDSN      string `required:"true" arg:"sentry-dsn" env:"SENTRY_DSN" usage:"SentryDSN" display:"length"`
	Listen         string `required:"true" arg:"listen" env:"LISTEN" usage:"address to listen to"`
	Kubeconfig     string `required:"false" arg:"kubeconfig" env:"KUBECONFIG" usage:"Path to k8s config"`
	CronExpression string `required:"true" arg:"cron-expression" env:"CRON_EXPRESSION" usage:"Cron expression to determine when service is run" default:"0 0 0 * * ?"`
	BackupRootDir  string `required:"true" arg:"backup-root-dir" env:"BACKUP_ROOT_DIR" usage:"Directory all backups are stored"`
	SSHPrivateKey  string `required:"true" arg:"ssh-key" env:"SSH_KEY" usage:"path to ssh private key"`
	Namespace      string `required:"true" arg:"namespace" env:"NAMESPACE" usage:"kubernetes namespace"`
}

func (a *application) Run(ctx context.Context, sentryClient libsentry.Client) error {
	currentTime := libtime.NewCurrentTime()

	trigger := run.NewTrigger()

	return service.Run(
		ctx,
		a.createSetupResourceDefinition(trigger),
		a.createHttpServer(currentTime),
	)
}

func (a *application) createSetupResourceDefinition(trigger run.Trigger) func(ctx context.Context) error {
	return pkg.CreateSetupResourceDefinition(a.Kubeconfig, k8s.Namespace(a.Namespace), trigger)
}

func (a *application) createHttpServer(
	currentTimeGetter libtime.CurrentTimeGetter,
) run.Func {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		router := mux.NewRouter()
		router.Path("/healthz").Handler(libhttp.NewPrintHandler("OK"))
		router.Path("/readiness").Handler(libhttp.NewPrintHandler("OK"))
		router.Path("/metrics").Handler(promhttp.Handler())
		router.Path("/setloglevel/{level}").Handler(log.NewSetLoglevelHandler(ctx, log.NewLogLevelSetter(2, 5*time.Minute)))

		router.Path("/list").Handler(libhttp.NewErrorHandler(
			libhttp.NewJsonHandler(
				libhttp.JsonHandlerFunc(func(ctx context.Context, req *http.Request) (interface{}, error) {
					k8sConnector := pkg.NewK8sConnector(
						a.Kubeconfig,
						k8s.Namespace(a.Namespace),
					)
					targets, err := k8sConnector.Targets(ctx)
					if err != nil {
						return nil, errors.Wrapf(ctx, err, "list targets failed")
					}
					return targets.Specs(), nil
				}),
			),
		))

		router.Path("/backup/{name}").Handler(libhttp.NewErrorHandler(
			libhttp.WithErrorFunc(func(ctx context.Context, resp http.ResponseWriter, req *http.Request) error {
				vars := mux.Vars(req)
				k8sConnector := pkg.NewK8sConnector(
					a.Kubeconfig,
					k8s.Namespace(a.Namespace),
				)
				target, err := k8sConnector.Target(ctx, vars["name"])
				if err != nil {
					return errors.Wrapf(ctx, err, "get target failed")
				}
				backupExectuor := pkg.CreateBackupExectuor(
					currentTimeGetter,
					pkg.BackupRootDirectory(a.BackupRootDir),
					pkg.SSHPrivateKey(a.SSHPrivateKey),
				)
				if err := backupExectuor.Backup(ctx, target.Spec); err != nil {
					return errors.Wrapf(ctx, err, "backup %s failed", target.Name)
				}
				libhttp.WriteAndGlog(resp, "backup %s completed", target.Name)
				return nil
			}),
		))

		router.Path("/trigger").Handler(libhttp.NewBackgroundRunHandler(ctx,
			pkg.CreateBackupCron(
				currentTimeGetter,
				a.Kubeconfig,
				pkg.BackupRootDirectory(a.BackupRootDir),
				pkg.SSHPrivateKey(a.SSHPrivateKey),
				k8s.Namespace(a.Namespace),
			).Run,
		))

		glog.V(2).Infof("starting http server listen on %s", a.Listen)
		return libhttp.NewServer(
			a.Listen,
			router,
		).Run(ctx)
	}
}
