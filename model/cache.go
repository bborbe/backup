package model

import "time"

type CacheTTL time.Duration

func (c CacheTTL) IsEmpty() bool {
	return int64(c) == 0
}

func (c CacheTTL) Duration() time.Duration {
	return time.Duration(c)
}
