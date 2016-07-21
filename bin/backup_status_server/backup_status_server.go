package main

import (
	"net/http"
	"os"

	"runtime"

	backup_config "github.com/bborbe/backup/config"
	backup_status_server "github.com/bborbe/backup/status_server"
	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/log"
	"github.com/facebookgo/grace/gracehttp"
)

const (
	DEFAULT_PORT       int = 8002
	PARAMETER_LOGLEVEL     = "loglevel"
	PARAMETER_ROOT         = "rootdir"
	PARAMETER_PORT         = "port"
)

var (
	logger        = log.DefaultLogger
	logLevelPtr   = flag.String(PARAMETER_LOGLEVEL, log.LogLevelToString(backup_config.DEFAULT_LOG_LEVEL), log.FLAG_USAGE)
	rootdirPtr    = flag.String(PARAMETER_ROOT, backup_config.DEFAULT_ROOT_DIR, "root directory for backups")
	portnumberPtr = flag.Int(PARAMETER_PORT, DEFAULT_PORT, "server port")
)

func main() {
	defer logger.Close()
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	server, err := createServer(*portnumberPtr, *rootdirPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
	logger.Debugf("start server")
	gracehttp.Serve(server)
}

func createServer(port int, rootdir string) (*http.Server, error) {
	return backup_status_server.NewServer(*portnumberPtr, *rootdirPtr), nil
}
