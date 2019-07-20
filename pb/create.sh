#!/bin/sh
#export PATH=$PATH:../server/src/protobuf/bin

SRC_DIR=./
DST_DIR=./gen

#go
mkdir -p $DST_DIR/go
protoc --go_out=${DST_DIR}/go/ $SRC_DIR/IM.BaseDefine.proto
protoc --go_out=${DST_DIR}/go/ $SRC_DIR/IM.Buddy.proto
protoc --go_out=${DST_DIR}/go/ $SRC_DIR/IM.Group.proto
protoc --go_out=${DST_DIR}/go/ $SRC_DIR/IM.Login.proto
protoc --go_out=${DST_DIR}/go/ $SRC_DIR/IM.Message.proto
protoc --go_out=${DST_DIR}/go/ $SRC_DIR/IM.Other.proto
protoc --go_out=${DST_DIR}/go/ $SRC_DIR/IM.Server.proto

#C++
#mkdir -p $DST_DIR/cpp
#protoc -I=$SRC_DIR --cpp_out=$DST_DIR/cpp/ $SRC_DIR/*.proto

#JAVA
#mkdir -p $DST_DIR/java
#protoc -I=$SRC_DIR --java_out=$DST_DIR/java/ $SRC_DIR/*.proto

#OBJC
#mkdir -p $DST_DIR/objc
#protoc -I=$SRC_DIR --objc_out=$DST_DIR/objc/ $SRC_DIR/*.proto

#PYTHON
#mkdir -p $DST_DIR/python
#protoc -I=$SRC_DIR --python_out=$DST_DIR/python/ $SRC_DIR/*.proto
