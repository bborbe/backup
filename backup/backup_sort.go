package backup

import (
	"github.com/bborbe/backup/util"
)

type BackupByName []Backup

func (v BackupByName) Len() int           { return len(v) }
func (v BackupByName) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v BackupByName) Less(i, j int) bool { return util.StringLess(v[i].Name(), v[j].Name()) }
