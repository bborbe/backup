package testutil

import (
	"fmt"
	"os"

	"github.com/golang/glog"
)

const BACKUP_ROOT_DIR = "/tmp/backuproot"

func ClearRootDir(root string) error {
	glog.V(2).Infof("ClearRootDir - root: %s", root)
	return os.RemoveAll(root)
}

func CreateRootDir(root string) error {
	glog.V(2).Infof("CreateRootDir - root: %s", root)
	var fileMode os.FileMode
	fileMode = 0777
	return os.Mkdir(root, fileMode)
}

func CreateHostDir(root string, host string) error {
	glog.V(2).Infof("CreateHostDir root: %s host: %s", root, host)
	var fileMode os.FileMode
	fileMode = 0777
	dir := fmt.Sprintf("%s%c%s", root, os.PathSeparator, host)
	glog.V(2).Infof("create hostdir %s", dir)
	return os.Mkdir(dir, fileMode)
}

func CreateBackupDir(root string, host string, backup string) error {
	glog.V(2).Infof("CreateBackupDir root: %s host: %s backup: %s", root, host, backup)
	var fileMode os.FileMode
	fileMode = 0777
	dir := fmt.Sprintf("%s%c%s%c%s", root, os.PathSeparator, host, os.PathSeparator, backup)
	glog.V(2).Infof("create backupdir %s", dir)
	return os.Mkdir(dir, fileMode)
}

func CreateBackupCurrentSymlink(root string, host string, backup string) error {
	glog.V(2).Infof("CreateBackupCurrentSymlink root: %s host: %s backup: %s", root, host, backup)
	return os.Symlink(fmt.Sprintf("%s%c%s%c%s", root, os.PathSeparator, host, os.PathSeparator, backup), fmt.Sprintf("%s%c%s%c%s", root, os.PathSeparator, host, os.PathSeparator, "current"))
}

func CreateFile(path string) error {
	glog.V(2).Infof("CreateFile file: %s", path)
	var fileMode os.FileMode
	fileMode = 0666
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileMode)
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.WriteString("hello world")
	if err != nil {
		return err
	}
	return nil
}
