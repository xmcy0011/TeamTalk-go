#!/bin/sh
DST_DIR=../server/src/base/protocol/
SRC_DIR=./gen

if [[ ! -d "${DST_DIR}" ]];then
    mkdir -p ${DST_DIR}
fi

#C++
cp -rf ${SRC_DIR}/go/* ${DST_DIR}

rm -rf ./gen
