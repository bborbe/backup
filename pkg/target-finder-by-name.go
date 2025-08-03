// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import (
	"context"

	"github.com/bborbe/errors"

	backupv1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
)

func NewTargetFinder(k8sConnector K8sConnector) TargetFinder {
	return &targetFinder{
		k8sConnector: k8sConnector,
	}
}

type targetFinder struct {
	k8sConnector K8sConnector
}

func (t *targetFinder) Target(ctx context.Context, name string) (*backupv1.Target, error) {
	target, err := t.k8sConnector.Target(ctx, name)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "get target failed")
	}
	return target, nil
}
