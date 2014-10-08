package fileutil

import (
	"os"

	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

func IsDir(dir string) (bool, error) {
	file, err := os.Open(dir)
	if err != nil {
		logger.Debugf("open dir %s failed: %v", dir, err)
		return false, nil
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		return false, err
	}
	return fileinfo.IsDir(), nil
}
