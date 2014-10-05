package rootdir

type Rootdir interface {
	Path() string
}

type rootdir string

func New(dir string) Rootdir {
	d := rootdir(dir)
	return &d
}

func (r *rootdir) Path() string {
	return string(*r)
}
