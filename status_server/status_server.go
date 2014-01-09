package status_server

import (
	"github.com/bborbe/server"
	"github.com/bborbe/server/handler/static"
	"strconv"
)

func NewServer(port int, rootdir string) server.Server {
	addr := "0.0.0.0:" + strconv.Itoa(port)
	handler := static.NewHandlerStaticContent("test")
	return server.NewServer(addr, handler)
}
