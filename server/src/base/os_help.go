package base

import (
	log "github.com/golang/glog"
	"os"
	"strconv"
)

func IsExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func WritePid() bool {
	// 记录进程id
	var pid = os.Getpid()
	var err error

	const pidFile = "server.pid"
	var fileHandle *os.File
	if !IsExist(pidFile) {
		fileHandle, err = os.Create(pidFile)
	} else {
		fileHandle, err = os.OpenFile(pidFile, os.O_RDWR, os.ModePerm)
	}

	if err != nil {
		log.Fatal("write pid file error:", err.Error())
		return false
	}

	_, err = fileHandle.Write([]byte(strconv.Itoa(pid)))
	if err != nil {
		log.Fatal("write pid file error:", err.Error())
		return false
	}

	err = fileHandle.Close()
	if err != nil {
		log.Fatal("write pid file error:", err.Error())
		return false
	}
	return true
}
