package main

import (
	"container/list"
	log "github.com/golang/glog"
	"net"
	"strconv"
	"sync"
)

var msgServerList *list.List
var msgServerListMutex sync.Mutex

type MsgServer struct {
	remoteIp   string
	remotePort int
	conn       net.Conn

	lastHeartBeat uint32
}

func init()  {
	msgServerList = list.New()
}

func ListenMsgServerConn(listenIp string, port int) {
	ln, err := net.Listen("tcp", listenIp+":"+strconv.Itoa(port))
	if err != nil {
		log.Error("listen on ", listenIp, " failed:", err.Error())
	}
	for {
		con, _ := ln.Accept()
		//go newConnect(con)
	}
}
