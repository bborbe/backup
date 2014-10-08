package testutil

import (
	"fmt"
	"os"

	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

const BACKUP_ROOT_DIR = "/tmp/backuproot"

func ClearBackupRootDir(root string) error {
	return os.RemoveAll(root)
}

func CreateRootDir(root string) error {
	var fileMode os.FileMode
	fileMode = 0777
	return os.Mkdir(root, fileMode)
}

func CreateHostDir(root string, host string) error {
	var fileMode os.FileMode
	fileMode = 0777
	dir := fmt.Sprintf("%s%c%s", root, os.PathSeparator, host)
	logger.Debugf("create hostdir %s", dir)
	return os.Mkdir(dir, fileMode)
}

func CreateBackupDir(root string, host string, backup string) error {
	var fileMode os.FileMode
	fileMode = 0777
	dir := fmt.Sprintf("%s%c%s%c%s", root, os.PathSeparator, host, os.PathSeparator, backup)
	logger.Debugf("create backupdir %s", dir)
	return os.Mkdir(dir, fileMode)
}
