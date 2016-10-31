package main

import (
	"fmt"
	"net/http"

	flag "github.com/bborbe/flagenv"

	"runtime"

	backup_dto "github.com/bborbe/backup/dto"
	backup_status_handler "github.com/bborbe/backup/status_client_handler"
	"github.com/bborbe/backup/status_fetcher"
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
	port := *portnumberPtr
	if port <= 0 {
		return nil, fmt.Errorf("parameter %s missing", parameterPort)
	}
	server := *serverPtr
	if len(server) == 0 {
		return nil, fmt.Errorf("parameter %s missing", parameterServer)
	}
	glog.Infof("port: %v server: %v", port, server)

	httpClient := http_client_builder.New().WithoutProxy().Build()
	statusFetcher := status_fetcher.New(httpClient.Get)
	handler := backup_status_handler.New(func() ([]backup_dto.Status, error) {
		return statusFetcher.StatusList(server)
	})
	return &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}, nil
}
