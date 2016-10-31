package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	backup_dto "github.com/bborbe/backup/dto"
	"github.com/golang/glog"
)

type BackupStatusFetcher interface {
	StatusList() ([]backup_dto.Status, error)
}

type getUrl func(url string) (resp *http.Response, err error)

type fetcher struct {
	getUrl  getUrl
	address string
}

func New(getUrl getUrl, address string) *fetcher {
	f := new(fetcher)
	f.getUrl = getUrl
	f.address = address
	return f
}

func (f *fetcher) StatusList() ([]backup_dto.Status, error) {
	resp, err := f.getUrl(f.address)
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
