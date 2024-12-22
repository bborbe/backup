package pkg

import (
	v1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
	libtime "github.com/bborbe/time"
)

type Backups []Backup

type Backup struct {
	Date libtime.Date  `json:"date"`
	Path Path          `json:"path"`
	Host v1.BackupHost `json:"host"`
}
