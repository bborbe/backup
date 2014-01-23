package dto

type Status interface {
	GetHost() string
	SetHost(string)
	GetStatus() bool
	SetStatus(bool)
}

type status struct {
	host   string
	status bool
}

func NewStatus() *status {
	h := new(status)
	return h
}

func (h *status) GetHost() string {
	return h.host
}

func (h *status) SetHost(host string) {
	h.host = host
}

func (h *status) GetStatus() bool {
	return h.status
}

func (h *status) SetStatus(status bool) {
	h.status = status
}
