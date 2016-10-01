package main

import (
	"net/http"
	"runtime"

	backup_config "github.com/bborbe/backup/constants"
	backup_status_server "github.com/bborbe/backup/status_server"
	flag "github.com/bborbe/flagenv"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/golang/glog"
)

const (
	DEFAULT_PORT   int = 8002
	PARAMETER_ROOT     = "rootdir"
	PARAMETER_PORT     = "port"
)

var (
	rootdirPtr    = flag.String(PARAMETER_ROOT, backup_config.DEFAULT_ROOT_DIR, "root directory for backups")
	portnumberPtr = flag.Int(PARAMETER_PORT, DEFAULT_PORT, "server port")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := do(
		*portnumberPtr,
		*rootdirPtr,
	)
	if err != nil {
		glog.Exit(err)
	}

}

func do(
	port int,
	rootdir string,
) error {
	server, err := createServer(
		port,
		rootdir,
	)
	if err != nil {
		return err
	}
	glog.V(2).Infof("start server")
	return gracehttp.Serve(server)
}

func createServer(
	port int,
	rootdir string,
) (*http.Server, error) {
	return backup_status_server.NewServer(port, rootdir), nil
}
