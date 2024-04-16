export SENTRY_DSN_KEY=yqZjWq

run:
	@go run -mod=vendor main.go \
	-sentry-dsn="$(shell teamvault-url --teamvault-config ~/.teamvault.json --teamvault-key=${SENTRY_DSN_KEY})" \
	-backup-root-dir="/tmp/backup" \
	-backup-host="hetzner-1.benjamin-borbe.de" \
	-backup-user="root" \
	-backup-port="22" \
	-backup-dirs="/mnt/a,/mnt/b" \
	-v=4
