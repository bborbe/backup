package handler

import (
	"context"
	"net/http"

	"github.com/bborbe/backup/pkg"
	"github.com/bborbe/errors"
	libhttp "github.com/bborbe/http"
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
