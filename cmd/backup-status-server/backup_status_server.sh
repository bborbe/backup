#!/bin/sh
# backup_status_server
#
# copy script to location /etc/init.d/backup_status_server
#

case "$1" in
	start)
		echo "Starting backup_status_server"
		/opt/backup/bin/backup_status_server -logtostderr -v=2 -port=7777 -rootdir=/rsync > /var/log/backup_status_server.log &
	;;
	stop)
		echo "Stopping backup_status_server"
		pid=`ps ax|grep backup_status_server | grep -v init.d |awk '{ print $1 }'`
		kill $pid  > /dev/null 2>&1
	;;
	restart)
		$0 stop
		sleep 2
		$0 start
	;;
	*)
		echo "Usage: /etc/init.d/backup_status_server {start|stop|restart}"
		exit 1
	;;
esac

exit 0
