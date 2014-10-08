package dto

import (
	"github.com/bborbe/backup/util"
)

type BackupByDate []Backup

func (v BackupByDate) Len() int           { return len(v) }
func (v BackupByDate) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v BackupByDate) Less(i, j int) bool { return util.StringLess(v[i].GetName(), v[j].GetName()) }
