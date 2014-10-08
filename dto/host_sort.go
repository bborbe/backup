package dto

import (
	"github.com/bborbe/backup/util"
)

type HostByDate []Host

func (v HostByDate) Len() int           { return len(v) }
func (v HostByDate) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v HostByDate) Less(i, j int) bool { return util.StringLess(v[i].GetName(), v[j].GetName()) }
