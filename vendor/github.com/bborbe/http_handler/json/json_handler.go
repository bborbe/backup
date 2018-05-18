package json

import (
	"net/http"

	"encoding/json"
	"reflect"

	error_handler "github.com/bborbe/http_handler/error"
	"github.com/golang/glog"
)

type handler struct {
	m interface{}
}

func New(m interface{}) *handler {
	h := new(handler)
	h.m = m
	return h
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	glog.V(4).Info("write json")
	glog.V(4).Infof("object to convert %v", h.m)
	b, err := json.Marshal(h.m)
	if err != nil {
		glog.V(2).Infof("Marshal json failed: %v", err)
		e := error_handler.NewMessage(http.StatusInternalServerError, err.Error())
		e.ServeHTTP(responseWriter, request)
		return
	}
	if glog.V(4) {
		glog.Infof("json string %s", string(b))
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)

	glog.V(4).Infof("object type %v", reflect.TypeOf(h.m).Kind())
	if reflect.TypeOf(h.m).Kind() == reflect.Slice && string(b) == "null" {
		responseWriter.Write([]byte("[]"))
	} else {
		responseWriter.Write(b)
	}

}
