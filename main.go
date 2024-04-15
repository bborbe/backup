// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"os"
	"time"

	libhttp "github.com/bborbe/http"
	"github.com/bborbe/log"
	"github.com/bborbe/run"
	libsentry "github.com/bborbe/sentry"
	"github.com/bborbe/service"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	app := &application{}
	os.Exit(service.Main(context.Background(), app, &app.SentryDSN))
}

type application struct {
	SentryDSN string `required:"true" arg:"sentry-dsn" env:"SENTRY_DSN" usage:"SentryDSN" display:"length"`
	Listen    string `required:"true" arg:"listen" env:"LISTEN" usage:"address to listen to"`
}

func (a *application) Run(ctx context.Context, sentryClient libsentry.Client) error {
	return service.Run(
		ctx,
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

		glog.V(2).Infof("starting http server listen on %s", a.Listen)
		return libhttp.NewServer(
			a.Listen,
			router,
		).Run(ctx)
	}
}
