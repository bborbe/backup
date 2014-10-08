package host

import (
	"github.com/bborbe/backup/util"
)

type HostByName []Host

func (v HostByName) Len() int           { return len(v) }
func (v HostByName) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v HostByName) Less(i, j int) bool { return util.StringLess(v[i].Name(), v[j].Name()) }
