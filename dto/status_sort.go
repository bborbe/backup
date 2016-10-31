package dto

import (
	"github.com/bborbe/stringutil"
)

type StatusByName []Status

func (v StatusByName) Len() int {
	return len(v)
}

func (v StatusByName) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v StatusByName) Less(i, j int) bool {
	return stringutil.StringLess(v[i].Host, v[j].Host)
}

type StatusByDate []Status

func (v StatusByDate) Len() int {
	return len(v)
}

func (v StatusByDate) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v StatusByDate) Less(i, j int) bool {
	return v[i].LatestBackup.Less(v[j].LatestBackup)
}
