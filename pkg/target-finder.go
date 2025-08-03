// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import (
	"context"
	stderrors "errors"

	backupv1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
)

var (
	TargetNotFoundError = stderrors.New("target not found")
)

//counterfeiter:generate -o ../mocks/target-finder.go --fake-name TargetFinder . TargetFinder
type TargetFinder interface {
	Target(ctx context.Context, name string) (*backupv1.Target, error)
}
