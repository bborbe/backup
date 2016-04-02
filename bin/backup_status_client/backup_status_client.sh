#!/bin/sh
# backup_status_client
#
# copy script to location /etc/init.d/backup_status_client
#

case "$1" in
	start)
		echo "Starting backup_status_client"
		/opt/backup/bin/backup_status_client -loglevel=ERROR -port=7777 -addr=/rsync > /var/log/backup_status_client.log &
	;;
	stop)
		echo "Stopping backup_status_client"
		pid=`ps ax|grep backup_status_client | grep -v init.d |awk '{ print $1 }'`
		kill $pid  > /dev/null 2>&1
	;;
	restart)
		$0 stop
		sleep 2
		$0 start
	;;
	*)
		echo "Usage: /etc/init.d/backup_status_client {start|stop|restart}"
		exit 1
	;;
esac

exit 0