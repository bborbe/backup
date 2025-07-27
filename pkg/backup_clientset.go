// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import (
	"context"

	"github.com/bborbe/errors"
	libk8s "github.com/bborbe/k8s"

	clientset "github.com/bborbe/backup/k8s/client/clientset/versioned"
)

//counterfeiter:generate -o ../mocks/backup-clientset.go --fake-name BackupClientset . BackupClientset
type BackupClientset clientset.Interface

func CreateBackupClientset(ctx context.Context, kubeconfig string) (BackupClientset, error) {
	config, err := libk8s.CreateConfig(kubeconfig)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "build k8s config failed")
	}
	backupClientset, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "build backup clientset failed")
	}
	return backupClientset, nil
}
