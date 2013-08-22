package object

type Host interface {
	GetName() string
}

type host struct {
	name string
}

func NewHost(name string) *host {
	h := new(host)
	h.name = name
	return h
}

func (h *host) GetName() string {
	return h.name
}
