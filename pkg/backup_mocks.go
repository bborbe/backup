// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import (
	backupv1 "github.com/bborbe/backup/k8s/client/clientset/versioned/typed/backup.benjamin-borbe.de/v1"
)

//counterfeiter:generate -o ../mocks/backup-v1-interface.go --fake-name BackupV1Interface . BackupV1Interface
type BackupV1Interface backupv1.BackupV1Interface
