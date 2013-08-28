#!/bin/sh

NAME=backup
VERSION=0.1
BINS="backup_cleanup backup_keep backup_latest backup_list backup_old"


#########################################################################

# Create scripts source dir
DIR=$NAME-$VERSION
rm -rf $DIR
mkdir $DIR

# Copy bins
for BIN in $BINS
do
cp bin/$BIN $DIR/$BIN
done
cd $DIR

# Create skeleton
echo foo | DEBFULLNAME="Benjamin Borbe" dh_make --single --indep --createorig --copyright bsd --email bborbe@rocketnews.de

# Remove make calls
grep -v makefile debian/rules > debian/rules.new
mv debian/rules.new debian/rules

# Add copyright
cp ../src/github.com/bborbe/backup/LICENSE debian/copyright

# Add to install
for BIN in $BINS
do
echo $BIN opt/backup/bin | tee -a debian/install
done

# We don't want a quilt based package
echo "1.0" > debian/source/format

# Remove the example files
rm debian/*.ex
rm debian/README*

# Build package
debuild -us -uc
