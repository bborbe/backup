package dto

type Status struct {
	Host         string `json:"host"`
	Status       bool   `json:"status"`
	LatestBackup string `json:"latestBackup"`
}
