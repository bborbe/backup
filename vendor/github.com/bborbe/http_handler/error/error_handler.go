package error

import (
	"net/http"

	"encoding/json"

	"github.com/golang/glog"
)

type handler struct {
	status  int
	message string
}

func New(status int) *handler {
	return NewMessage(status, http.StatusText(status))
}

func NewMessage(status int, message string) *handler {
	h := new(handler)
	h.status = status
	h.message = message
	return h
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	glog.V(4).Info("handle error")

	var data struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}
	data.Message = h.message
	data.Status = h.status
	glog.V(4).Infof("set status: %d", h.status)
	responseWriter.WriteHeader(h.status)
	responseWriter.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(responseWriter).Encode(&data); err != nil {
		glog.Warningf("render failureRenderer failed! %v", err)
	}
}
