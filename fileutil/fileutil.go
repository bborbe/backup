package fileutil

import (
	"os"

	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func IsDir(dir string) (bool, error) {
	logger.Debugf("IsDir %s", dir)
	file, err := os.Open(dir)
	if err != nil {
		logger.Debugf("IsDir - open dir %s failed: %v", dir, err)
		return false, nil
	}
	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		logger.Debugf("IsDir get state for dir %s failed: %v", dir, err)
		return false, err
	}
	return fileinfo.IsDir(), nil
}
