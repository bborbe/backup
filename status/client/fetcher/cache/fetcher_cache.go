package cache

import (
	backup_dto "github.com/bborbe/backup/dto"
	"github.com/bborbe/backup/model"
	"github.com/bborbe/backup/status/client/fetcher"
	"github.com/golang/glog"
	"sync"
	"time"
)

type cache struct {
	backupStatusFetcher fetcher.BackupStatusFetcher
	ttl                 model.CacheTTL

	mutex sync.Mutex
	list  []backup_dto.Status
	time  time.Time
}

func New(
	backupStatusFetcher fetcher.BackupStatusFetcher,
	ttl model.CacheTTL,
) *cache {
	c := new(cache)
	c.backupStatusFetcher = backupStatusFetcher
	c.ttl = ttl
	return c
}

func (c *cache) StatusList() ([]backup_dto.Status, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.time.IsZero() || time.Now().Sub(c.time) > c.ttl.Duration() {
		glog.V(2).Infof("cache expired, fetch status list")
		list, err := c.backupStatusFetcher.StatusList()
		if err != nil {
			return nil, err
		}
		c.time = time.Now()
		c.list = list
	} else {
		glog.V(2).Infof("use cached status list")
	}
	return c.list, nil
}
