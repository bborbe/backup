package cache

import (
	backup_dto "github.com/bborbe/backup/dto"
	"github.com/bborbe/backup/model"
	"github.com/bborbe/backup/status/client/fetcher"
	"github.com/wunderlist/ttlcache"
)

type cache struct {
	backupStatusFetcher fetcher.BackupStatusFetcher
	cache               *ttlcache.Cache
}

func New(
	backupStatusFetcher fetcher.BackupStatusFetcher,
	ttl model.CacheTTL,
) *cache {
	c := new(cache)
	c.backupStatusFetcher = backupStatusFetcher
	c.cache = ttlcache.NewCache(ttl.Duration())
	return c
}

func (c *cache) StatusList() ([]backup_dto.Status, error) {
	return c.backupStatusFetcher.StatusList()
}
