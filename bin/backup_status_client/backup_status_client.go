package main

import (
	"flag"

	"net/http"
	"runtime"

	backup_status_client "github.com/bborbe/backup/status_client"
	http_client_builder "github.com/bborbe/http/client_builder"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/golang/glog"
)

const (
	DEFAULT_PORT     int = 8080
	DEFAULT_SERVER       = "http://backup.pn.benjamin-borbe.de:7777"
	PARAMETER_PORT       = "port"
	PARAMETER_SERVER     = "server"
)

var (
	serverPtr     = flag.String(PARAMETER_SERVER, DEFAULT_SERVER, "backup status server address")
	portnumberPtr = flag.Int(PARAMETER_PORT, DEFAULT_PORT, "server port")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := do(
		*portnumberPtr,
		*serverPtr,
	)
	if err != nil {
		glog.Exit(err)
	}
}

func do(
	port int,
	serverAddr string,
) error {
	glog.Infof("port: %v server: %v", port, serverAddr)
	server, err := createServer(
		port,
		serverAddr,
	)
	if err != nil {
		return err
	}
	glog.V(2).Infof("start server")
	return gracehttp.Serve(server)
}

func createServer(
	port int,
	server string,
) (*http.Server, error) {

	httpClient := http_client_builder.New().WithoutProxy().Build()

	return backup_status_client.NewServer(httpClient.Get, port, server), nil
}
