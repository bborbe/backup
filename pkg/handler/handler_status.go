package handler

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/bborbe/errors"
	libhttp "github.com/bborbe/http"
	glog "github.com/golang/glog"

	v1 "github.com/bborbe/backup/k8s/apis/backup.benjamin-borbe.de/v1"
	"github.com/bborbe/backup/pkg"
)

func NewStatusHandler(k8sConnector pkg.K8sConnector, backupRootDir pkg.Path) libhttp.WithError {
	return libhttp.NewJsonHandler(
		libhttp.JsonHandlerFunc(func(ctx context.Context, req *http.Request) (interface{}, error) {

			targets, err := k8sConnector.Targets(ctx)
			if err != nil {
				return nil, errors.Wrapf(ctx, err, "list targets failed")
			}
			result := map[v1.BackupHost]string{}
			for _, target := range targets {
				host := target.Spec.Host
				backupDir := backupRootDir.Join(host.String())
				glog.V(4).Infof("search for backups in %s", backupDir)
				entries, err := os.ReadDir(backupDir.String())
				if err != nil {
					return nil, errors.Wrapf(ctx, err, "list failed")
				}
				glog.V(4).Infof("found %d entries in %s", len(entries), backupDir)
				var latestBackup *time.Time
				for _, entry := range entries {
					backupTime, err := time.Parse(time.DateOnly, entry.Name())
					if err != nil {
						glog.V(4).Infof("name(%s) is not valid  => skip", entry.Name())
						continue
					}
					if latestBackup == nil || backupTime.After(*latestBackup) {
						latestBackup = &backupTime
						result[host] = backupTime.Format(time.DateOnly)
					}
				}
			}
			return result, nil
		}),
	)
}
