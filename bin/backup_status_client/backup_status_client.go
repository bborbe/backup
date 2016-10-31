package main

import (
	"fmt"
	"net/http"

	flag "github.com/bborbe/flagenv"

	"runtime"

	"time"

	"github.com/bborbe/backup/model"
	backup_status_fetcher "github.com/bborbe/backup/status/client/fetcher"
	"github.com/bborbe/backup/status/client/fetcher/cache"
	backup_status_handler "github.com/bborbe/backup/status/client/handler"
	http_client_builder "github.com/bborbe/http/client_builder"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/golang/glog"
)

const (
	defaultPort                  int = 8080
	defaultStatusServerAddress       = "http://backup.pn.benjamin-borbe.de:1080"
	parameterPort                    = "port"
	parameterStatusServerAddress     = "server"
	parameterCacheTTL                = "cache-ttl"
)

var (
	statusServerAddressPtr = flag.String(parameterStatusServerAddress, defaultStatusServerAddress, "backup status server address")
	portnumberPtr          = flag.Int(parameterPort, defaultPort, "server port")
	cacheTTLPtr            = flag.Duration(parameterCacheTTL, 30*time.Second, "cache ttl")
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
	statusServerAddress := *statusServerAddressPtr
	if len(statusServerAddress) == 0 {
		return nil, fmt.Errorf("parameter %s missing", parameterStatusServerAddress)
	}
	cacheTTL := model.CacheTTL(*cacheTTLPtr)
	if cacheTTL.IsEmpty() {
		return nil, fmt.Errorf("parameter %s missing", parameterCacheTTL)
	}

	httpClient := http_client_builder.New().WithoutProxy().Build()
	statusFetcher := backup_status_fetcher.New(httpClient.Get, statusServerAddress)
	cachedStatusFetcher := cache.New(statusFetcher, cacheTTL)
	handler := backup_status_handler.New(cachedStatusFetcher.StatusList)
	glog.Infof("start server on port: %v with status api: %v", port, statusServerAddress)
	return &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}, nil
}
