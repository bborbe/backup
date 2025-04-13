// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import (
	libtime "github.com/bborbe/time"

	v1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
)

type Backups []Backup

type Backup struct {
	Date libtime.Date  `json:"date"`
	Path Path          `json:"path"`
	Host v1.BackupHost `json:"host"`
}
