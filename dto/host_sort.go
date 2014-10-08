package dto

import (
	"github.com/bborbe/stringutil"
)

type HostByName []Host

func (v HostByName) Len() int           { return len(v) }
func (v HostByName) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v HostByName) Less(i, j int) bool { return stringutil.StringLess(v[i].GetName(), v[j].GetName()) }
