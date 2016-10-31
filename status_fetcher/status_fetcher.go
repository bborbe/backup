package status_fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	backup_dto "github.com/bborbe/backup/dto"
	"github.com/golang/glog"
)

type getUrl func(url string) (resp *http.Response, err error)

type fetcher struct {
	getUrl getUrl
}

func New(getUrl getUrl) *fetcher {
	f := new(fetcher)
	f.getUrl = getUrl
	return f
}

func (f *fetcher) StatusList(address string) ([]backup_dto.Status, error) {
	resp, err := f.getUrl(address)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("request failed: %s", (content))
	}
	glog.V(4).Infof(string(content))
	var statusList []backup_dto.Status
	err = json.Unmarshal(content, &statusList)
	if err != nil {
		glog.V(1).Infof("unmarshal jsoni failed: %v", err)
		return nil, err
	}
	return statusList, nil
}
