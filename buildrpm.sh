#!/bin/bash -x
WORK_DIR=`pwd`
BUILD_NAME=goqos
BUILD_SPEC=rpm/${BUILD_NAME}.spec
BUILD_TAR=${BUILD_NAME}.tar.gz
BUILD_DIR=${WORK_DIR}/rpmbuild
BUILD_NUMBER=${BUILD_NUMBER:=0}

ARCHIVE_DIR=${WORK_DIR}/${BUILD_NAME}

/bin/rm   -rf ${BUILD_DIR}
/bin/mkdir -p ${BUILD_DIR}/{BUILD,RPMS,SOURCES,SPECS,SRPMS}

/usr/bin/git archive -v --format=tar.gz --prefix=${BUILD_NAME}/ HEAD -o ${BUILD_TAR}
/bin/mv -f ${BUILD_TAR} ${WORK_DIR}/rpmbuild/SOURCES/

/usr/bin/rpmbuild --define "release ${BUILD_NUMBER}%{?dist}" --define "_topdir ${BUILD_DIR}" -bb ${BUILD_SPEC}
