// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import (
	"context"

	"github.com/bborbe/errors"

	backupv1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
)

func NewTargetFinderByHostname(k8sConnector K8sConnector) TargetFinder {
	return &targetFinderByHostname{
		k8sConnector: k8sConnector,
	}
}

type targetFinderByHostname struct {
	k8sConnector K8sConnector
}

func (t *targetFinderByHostname) Target(ctx context.Context, hostname string) (*backupv1.Target, error) {
	targets, err := t.k8sConnector.Targets(ctx)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "get targets failed")
	}
	for _, target := range targets {
		if target.Spec.Host.String() == hostname {
			return &target, nil
		}
	}
	return nil, errors.Wrapf(ctx, TargetNotFoundError, "target with hostname '%s' not found", hostname)
}
