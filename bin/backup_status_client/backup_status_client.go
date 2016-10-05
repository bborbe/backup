package main

import (
	flag "github.com/bborbe/flagenv"

	"net/http"
	"runtime"

	backup_status_client "github.com/bborbe/backup/status_client"
	http_client_builder "github.com/bborbe/http/client_builder"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/golang/glog"
)

const (
	defaultPort     int = 8080
	defaultServer       = "http://backup.pn.benjamin-borbe.de:7777"
	parameterPort       = "port"
	parameterServer     = "server"
)

var (
	serverPtr     = flag.String(parameterServer, defaultServer, "backup status server address")
	portnumberPtr = flag.Int(parameterPort, defaultPort, "server port")
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
