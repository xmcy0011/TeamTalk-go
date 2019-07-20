#!/bin/sh
DST_DIR=../server/src/base/improto/
SRC_DIR=./gen

if [[ ! -d "${DST_DIR}" ]];then
    mkdir -p ${DST_DIR}
fi

#C++
cp -rf ${SRC_DIR}/go/* ${DST_DIR}

rm -rf ./gen

## 需要手动进去 DST_DIR 更改package后面的名称为 improto