package main

import (
	"flag"

	"net/http"
	"os"

	"runtime"

	backup_config "github.com/bborbe/backup/config"
	backup_status_client "github.com/bborbe/backup/status_client"
	http_client_builder "github.com/bborbe/http/client_builder"
	"github.com/bborbe/log"
	"github.com/facebookgo/grace/gracehttp"
)

const (
	DEFAULT_PORT       int = 8080
	DEFAULT_SERVER         = "http://backup.pn.benjamin-borbe.de:7777"
	PARAMETER_LOGLEVEL     = "loglevel"
	PARAMETER_PORT         = "port"
	PARAMETER_SERVER       = "server"
)

var (
	logger        = log.DefaultLogger
	logLevelPtr   = flag.String(PARAMETER_LOGLEVEL, log.LogLevelToString(backup_config.DEFAULT_LOG_LEVEL), log.FLAG_USAGE)
	serverPtr     = flag.String(PARAMETER_SERVER, DEFAULT_SERVER, "backup status server address")
	portnumberPtr = flag.Int(PARAMETER_PORT, DEFAULT_PORT, "server port")
)

func main() {
	defer logger.Close()
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	err := do(
		*portnumberPtr,
		*serverPtr,
	)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(
	port int,
	serverAddr string,
) error {
	server, err := createServer(
		port,
		serverAddr,
	)
	if err != nil {
		return err
	}
	logger.Debugf("start server")
	return gracehttp.Serve(server)
}

func createServer(
	port int,
	server string,
) (*http.Server, error) {
	logger.Tracef("server %s", *serverPtr)
	logger.Tracef("portnumberPtr %d", *portnumberPtr)
	logger.Debugf("backup status server started at port %d", *portnumberPtr)

	httpClient := http_client_builder.New().WithoutProxy().Build()

	return backup_status_client.NewServer(httpClient.Get, port, server), nil
}
