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
	"github.com/bborbe/log"
	"github.com/bborbe/run"
	libsentry "github.com/bborbe/sentry"
	"github.com/bborbe/service"
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
	CronExpression string `required:"false" arg:"cron-schedule-expression" env:"CRON_SCHEDULE_EXPRESSION" usage:"Cron schedule expression to determine when service is run" default:"@every 1m"`
}

func (a *application) Run(ctx context.Context, sentryClient libsentry.Client) error {

	trigger := run.NewTrigger()

	return service.Run(
		ctx,
		pkg.CreateSetupResourceDefinition(a.Kubeconfig, trigger),
		a.createHttpServer(),
	)
}

func (a *application) createHttpServer() run.Func {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		router := mux.NewRouter()
		router.Path("/healthz").Handler(libhttp.NewPrintHandler("OK"))
		router.Path("/readiness").Handler(libhttp.NewPrintHandler("OK"))
		router.Path("/metrics").Handler(promhttp.Handler())
		router.Path("/setloglevel/{level}").Handler(log.NewSetLoglevelHandler(ctx, log.NewLogLevelSetter(2, 5*time.Minute)))

		router.Path("/list").Handler(libhttp.NewErrorHandler(
			libhttp.WithErrorFunc(func(ctx context.Context, resp http.ResponseWriter, req *http.Request) error {
				targets, err := pkg.NewK8sConnector(a.Kubeconfig).Targets(ctx)
				if err != nil {
					return errors.Wrapf(ctx, err, "list targets failed")
				}
				libhttp.WriteAndGlog(resp, "targets %+v", targets)
				return nil
			}),
		))

		glog.V(2).Infof("starting http server listen on %s", a.Listen)
		return libhttp.NewServer(
			a.Listen,
			router,
		).Run(ctx)
	}
}
