package dto

type Host interface {
	GetName() string
	SetName(name string)
}

type host struct {
	Name string `json:"name"`
}

func NewHost() *host {
	h := new(host)
	return h
}

func (h *host) GetName() string {
	return h.Name
}

func (h *host) SetName(name string) {
	h.Name = name
}
