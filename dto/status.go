package dto

type Status interface {
	GetHost() string
	SetHost(string)
	GetStatus() bool
	SetStatus(bool)
}

type status struct {
	Host         string `json:"host"`
	Status       bool   `json:"status"`
	LatestBackup string `json:"latestBackup"`
}

func NewStatus() *status {
	h := new(status)
	return h
}

func (h *status) GetHost() string {
	return h.Host
}

func (h *status) SetHost(host string) {
	h.Host = host
}

func (h *status) GetStatus() bool {
	return h.Status
}

func (h *status) SetStatus(status bool) {
	h.Status = status
}

func (h *status) GetLatestBackup() string {
	return h.LatestBackup
}

func (h *status) SetLatestBackup(latestBackup string) {
	h.LatestBackup = latestBackup
}
