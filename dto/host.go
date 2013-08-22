package dto

type Host interface {
	GetName() string
	SetName(name string)
}

type host struct {
	name string
}

func NewHost() *host {
	h := new(host)
	return h
}

func (h *host) GetName() string {
	return h.name
}

func (h *host) SetName(name string) {
	h.name = name
}
