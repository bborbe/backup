# Backup

Tools for backup via rsync.

## Install

`go get github.com/bborbe/backup/bin/backup_cleanup`

`go get github.com/bborbe/backup/bin/backup_keep`

`go get github.com/bborbe/backup/bin/backup_latest`

`go get github.com/bborbe/backup/bin/backup_list`

`go get github.com/bborbe/backup/bin/backup_old`

`go get github.com/bborbe/backup/bin/backup_resume`

`go get github.com/bborbe/backup/bin/backup_status_server`

## Backup Rsync Client

`go get github.com/bborbe/backup/bin/backup_rsync_client`

via config

```
backup_rsync_client \
-logtostderr \
-v=2 \
-lock=/var/run/backup_rsync_client.lock \
-target=/backup/ \
-config=backup_rsync_client_sample.json \
-one-time
```

via args

```
backup_rsync_client \
-logtostderr \
-v=2 \
-lock=/var/run/backup_rsync_client.lock \
-target=/tmp \
-user=bborbe \
-host=localhost \
-port=22 \
-dir=/backup/ \
-exclude_from=/tmp/excludes \
-one-time
```
