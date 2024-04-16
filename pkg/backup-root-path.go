package pkg

type BackupRootDirectory string

func (f BackupRootDirectory) String() string {
	return string(f)
}
