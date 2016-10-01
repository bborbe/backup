package main

import (
	"os"
)

type targetDirectory string

func (t targetDirectory) String() string {
	return string(t)
}

func (t targetDirectory) IsValid() error {
	_, err := os.Stat(t.String())
	return err
}
