#!/bin/sh


DEBFULLNAME="Benjamin Borbe"
EMAIL=bborbe@rocketnews.de

DEB_SERVER=misc.rn.benjamin-borbe.de
TARGET_DIR=opt/backup/bin

NAME=backup
BINS="backup_cleanup backup_keep backup_latest backup_list backup_old"
INSTALLS="github.com/bborbe/backup/bin/backup_cleanup github.com/bborbe/backup/bin/backup_latest github.com/bborbe/backup/bin/backup_list github.com/bborbe/backup/bin/backup_old github.com/bborbe/backup/bin/backup_keep"
SOURCEDIRECTORY="github.com/bborbe/backup"

MAJOR=0
MINOR=1
BUGFIX=0

#########################################################################

BUILD=${BUILD_NUMBER}
VERSION=$MAJOR.$MINOR.$BUGFIX.$BUILD

export VERSION
export NAME
export GOPATH=${WORKSPACE}
export REPORT_DIR=${WORKSPACE}/test-reports

rm -rf $REPORT_DIR bin pkg package *.deb backup*
mkdir -p $REPORT_DIR
PACKAGES=`cd src && find $SOURCEDIRECTORY -name "*_test.go" | /usr/lib/go/bin/dirof | /usr/lib/go/bin/unique`
for PACKAGE in $PACKAGES
do
        XML=$REPORT_DIR/`/usr/lib/go/bin/pkg2xmlname $PACKAGE`
        OUT=$XML.out
        go test -i $PACKAGE
        go test -v $PACKAGE | tee $OUT
        cat $OUT
        /usr/lib/go/bin/go2xunit -fail -input $OUT -output $XML
done

go install $INSTALLS

# Create scripts source dir
DIR=$NAME-$VERSION
echo $DIR
rm -rf $DIR
mkdir $DIR

# Copy bins
for BIN in $BINS
do
cp bin/$BIN $DIR/$BIN
done
cd $DIR

# Create skeleton
echo foo | dh_make --single --indep --createorig --copyright bsd --email $EMAIL

# Remove make calls
grep -v makefile debian/rules > debian/rules.new
mv debian/rules.new debian/rules

# Set distrubtion
sed -i.bak 's/unstable/bborbe-unstable/g' debian/changelog

# Add copyright
cp ../src/${SOURCEDIRECTORY}/LICENSE debian/copyright

# Add to install
for BIN in $BINS
do
echo $BIN $TARGET_DIR | tee -a debian/install
done

# We don't want a quilt based package
echo "1.0" > debian/source/format

# Remove the example files
rm debian/*.ex
rm debian/README*

# Build package
debuild -us -uc

dput -u $DEB_SERVER ${NAME}_${VERSION}-1_amd64.changes
