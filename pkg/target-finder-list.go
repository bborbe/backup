// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import (
	"context"

	"github.com/bborbe/errors"

	backupv1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
)

type TargetFinderList []TargetFinder

func (t TargetFinderList) Target(ctx context.Context, name string) (*backupv1.Target, error) {
	for _, finder := range t {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			target, err := finder.Target(ctx, name)
			if err != nil {
				continue
			}
			return target, nil
		}
	}
	return nil, errors.Wrapf(ctx, TargetNotFoundError, "target '%s' not found in any finder", name)
}
