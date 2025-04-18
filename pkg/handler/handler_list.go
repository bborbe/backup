// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"context"
	"net/http"

	"github.com/bborbe/errors"
	libhttp "github.com/bborbe/http"

	"github.com/bborbe/backup/pkg"
)

func NewListHandler(k8sConnector pkg.K8sConnector) libhttp.WithError {
	return libhttp.NewJsonHandler(
		libhttp.JsonHandlerFunc(func(ctx context.Context, req *http.Request) (interface{}, error) {
			targets, err := k8sConnector.Targets(ctx)
			if err != nil {
				return nil, errors.Wrapf(ctx, err, "list targets failed")
			}
			return targets.Specs(), nil
		}),
	)
}
