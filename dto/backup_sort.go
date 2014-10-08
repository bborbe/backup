package dto

import (
	"github.com/bborbe/stringutil"
)

type BackupByName []Backup

func (v BackupByName) Len() int      { return len(v) }
func (v BackupByName) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v BackupByName) Less(i, j int) bool {
	return stringutil.StringLess(v[i].GetName(), v[j].GetName())
}
