package main

import (
	"fmt"
	"net/http"
	"runtime"

	backup_config "github.com/bborbe/backup/constants"
	backup_service "github.com/bborbe/backup/service"
	backup_status_checker "github.com/bborbe/backup/status/server/checker"
	backup_status_handler "github.com/bborbe/backup/status/server/handler"
	flag "github.com/bborbe/flagenv"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/golang/glog"
	"github.com/bborbe/backup/model"
)

const (
	defaultPort int = 8002
	parameterTarget = "target"
	parameterPort = "port"
)

var (
	rootdirPtr = flag.String(parameterTarget, backup_config.DEFAULT_ROOT_DIR, "root directory for backups")
	portnumberPtr = flag.Int(parameterPort, defaultPort, "server port")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := do(); err != nil {
		glog.Exit(err)
	}

}

func do() error {
	server, err := createServer()
	if err != nil {
		return err
	}
	glog.V(2).Infof("start server")
	return gracehttp.Serve(server)
}

func createServer() (*http.Server, error) {
	port := model.Port(*portnumberPtr)
	rootdir := *rootdirPtr
	backupService := backup_service.NewBackupService(rootdir)
	statusChecker := backup_status_checker.New(backupService)
	handler := backup_status_handler.New(statusChecker)
	glog.V(2).Infof("create http server on %s", port.Address())
	return &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}, nil
}
