#!/bin/sh

rm -rf $GOPATH/package
mkdir -p $GOPATH/package/DEBIAN
mkdir -p $GOPATH/package/opt/backup/bin

cp $GOPATH/src/github.com/bborbe/backup/control $GOPATH/package/DEBIAN/control
cp $GOPATH/bin/backup_latest $GOPATH/package/opt/backup/bin/
cp $GOPATH/bin/backup_list $GOPATH/package/opt/backup/bin/
cp $GOPATH/bin/backup_old $GOPATH/package/opt/backup/bin/
cp $GOPATH/bin/backup_cleanup $GOPATH/package/opt/backup/bin/

chmod 555 $GOPATH/package/opt/backup/bin/*
dpkg -b $GOPATH/package $GOPATH/backup.deb
