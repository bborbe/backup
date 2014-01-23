package dto

type Backup interface {
	GetName() string
	SetName(string)
}

type backup struct {
	Name string `json:"name"`
}

func NewBackup() *backup {
	h := new(backup)
	return h
}

func (h *backup) GetName() string {
	return h.Name
}

func (h *backup) SetName(name string) {
	h.Name = name
}
