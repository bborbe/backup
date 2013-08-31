#!/bin/sh

MAJOR=0
MINOR=1
BUGFIX=0
BUILD=${BUILD_NUMBER}
NAME=backup
VERSION=$MAJOR.$MINOR.$BUGFIX.$BUILD

export VERSION
export NAME
export GOPATH=${WORKSPACE}
export REPORT_DIR=${WORKSPACE}/test-reports

rm -rf $REPORT_DIR bin pkg package *.deb backup*
mkdir -p $REPORT_DIR
PACKAGES=`cd src && find github.com/bborbe/backup -name "*_test.go" | /usr/lib/go/bin/dirof | /usr/lib/go/bin/unique`
for PACKAGE in $PACKAGES
do
        XML=$REPORT_DIR/`/usr/lib/go/bin/pkg2xmlname $PACKAGE`
        OUT=$XML.out
        go test -i $PACKAGE
        go test -v $PACKAGE | tee $OUT
        cat $OUT
        /usr/lib/go/bin/go2xunit -fail -input $OUT -output $XML
done

go install github.com/bborbe/backup/bin/backup_cleanup github.com/bborbe/backup/bin/backup_latest github.com/bborbe/backup/bin/backup_list github.com/bborbe/backup/bin/backup_old github.com/bborbe/backup/bin/backup_keep
sh src/github.com/bborbe/backup/build_deb.sh
dput -u misc.rn.benjamin-borbe.de ${NAME}_${VERSION}-1_amd64.changes
