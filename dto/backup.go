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

func CreateBackups(backups []string) []Backup {
	result := make([]Backup, len(backups))
	for i, backup := range backups {
		result[i] = CreateBackup(backup)
	}
	return result
}

func CreateBackup(backup string) Backup {
	h := NewBackup()
	h.SetName(backup)
	return h
}
