package main

type host struct {
	Name string `json:"name"`
}

func (h *host) Backup() error {
	return nil
}
