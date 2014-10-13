#!/bin/sh


export DEBFULLNAME="Benjamin Borbe"
export EMAIL=bborbe@rocketnews.de

export DEB_SERVER=misc.rn.benjamin-borbe.de
export TARGET_DIR=opt/backup/bin

export NAME=backup
export BINS="backup_cleanup backup_keep backup_latest backup_list backup_old backup_resume"
export INSTALLS="github.com/bborbe/backup/bin/backup_cleanup github.com/bborbe/backup/bin/backup_keep github.com/bborbe/backup/bin/backup_latest github.com/bborbe/backup/bin/backup_list github.com/bborbe/backup/bin/backup_old github.com/bborbe/backup/bin/backup_resume"
export SOURCEDIRECTORY="github.com/bborbe/backup"

export MAJOR=0
export MINOR=1
export BUGFIX=0

# exec
sh src/github.com/bborbe/jenkins/jenkins.sh