package dto

import (
	"fmt"
	"time"
)

type BackupDate string

func (b BackupDate) String() string {
	return string(b)
}
func (b BackupDate) Time() (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05", b.String())
}

func (b BackupDate) Less(date BackupDate) bool {
	t1, err := b.Time()
	if err != nil {
		return false
	}
	t2, err := date.Time()
	if err != nil {
		return false
	}
	return t1.Before(t2)
}

func (b BackupDate) Age(now time.Time) string {
	t, err := b.Time()
	if err != nil {
		return "-"
	}
	return FormatDuration(now.Sub(t))
}

type Status struct {
	Host         string     `json:"host"`
	Status       bool       `json:"status"`
	LatestBackup BackupDate `json:"latestBackup"`
}

func FormatDuration(duration time.Duration) string {
	if hours := int(duration.Hours()); hours > 0 {
		if hours > 24 {
			return fmt.Sprintf("%dd", hours/24)
		}
		return fmt.Sprintf("%dh", hours)
	}
	if minutes := int(duration.Minutes()); minutes > 0 {
		return fmt.Sprintf("%dm", minutes)
	}
	if seconds := int(duration.Seconds()); seconds > 0 {
		return fmt.Sprintf("%ds", seconds)
	}
	return duration.String()
}

type StatusByBackupDate []*Status

func (v StatusByBackupDate) Len() int           { return len(v) }
func (v StatusByBackupDate) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v StatusByBackupDate) Less(i, j int) bool { return v[i].LatestBackup.Less(v[j].LatestBackup) }
