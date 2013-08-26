package util

import "github.com/bborbe/backup/dto"

type BackupByDate []dto.Backup

func (v BackupByDate) Len() int           { return len(v) }
func (v BackupByDate) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v BackupByDate) Less(i, j int) bool { return StringLess(v[i].GetName(), v[j].GetName()) }

func StringLess(stringA, stringB string) bool {
	a := []byte(stringA)
	b := []byte(stringB)
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] < b[i] {
			return true
		}
	}
	return len(a) < len(b)
}
