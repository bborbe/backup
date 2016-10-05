package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"runtime"

	backup_config "github.com/bborbe/backup/constants"
	backup_service "github.com/bborbe/backup/service"
	"github.com/golang/glog"
)

const (
	NO_HOST         = "-"
	parameterTarget = "target"
	parameterHost   = "host"
)

var (
	rootdirPtr = flag.String(parameterTarget, backup_config.DEFAULT_ROOT_DIR, "string")
	hostPtr    = flag.String(parameterHost, NO_HOST, "string")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	writer := os.Stdout
	glog.V(2).Infof("use backup dir %s", *rootdirPtr)
	backupService := backup_service.NewBackupService(*rootdirPtr)
	err := do(writer, backupService, *hostPtr)
	if err != nil {
		glog.Exit(err)
	}
}

func do(writer io.Writer, backupService backup_service.BackupService, hostname string) error {
	glog.V(2).Info("start")
	if hostname == NO_HOST {
		return fmt.Errorf("parameter host missing")
	}
	host, err := backupService.GetHost(hostname)
	if err != nil {
		fmt.Fprintf(writer, "host %s not found", hostname)
		return err
	}
	err = backupService.Resume(host)
	if err != nil {
		fmt.Fprintf(writer, "resume backup for host %s failed\n", hostname)
		glog.Warning(err)
	} else {
		fmt.Fprintf(writer, "resume backup for host %s success\n", hostname)
	}
	glog.V(2).Info("done")
	return nil
}
