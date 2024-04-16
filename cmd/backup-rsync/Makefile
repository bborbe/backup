export SENTRY_DSN_KEY=yqZjWq

run:
	@go run -mod=vendor main.go \
	-sentry-dsn="$(shell teamvault-url --teamvault-config ~/.teamvault.json --teamvault-key=${SENTRY_DSN_KEY})" \
	-backup-root-dir="/tmp/backup" \
	-backup-host="hetzner-1.benjamin-borbe.de" \
	-backup-user="bborbe" \
	-backup-port="22" \
	-backup-dirs="/" \
	-backup-excludes="/root,/sys,/proc,/dev,/etc,/var,/data,/run,/tmp,/opt,/boot,/core,/usr,/home/webdav,/lost+found,/swapfile" \
	-ssh-key="/Users/bborbe/.ssh/id_ed25519_personal" \
	-v=4
