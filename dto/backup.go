package dto

type Backup interface {
	GetName() string
	SetName(string)
}

type backup struct {
	name string
}

func NewBackup() *backup {
	h := new(backup)
	return h
}

func (h *backup) GetName() string {
	return h.name
}

func (h *backup) SetName(name string) {
	h.name = name
}
