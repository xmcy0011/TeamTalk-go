package main

import (
	"TeamTalk-go-flutter/server/src/base"
	"TeamTalk-go-flutter/server/src/base/improto"
	"container/list"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"time"
)

type MsgServer struct {
	Conn *base.ImConn

	PriorIp  string
	BackupIp string
	Port     int

	UserCount       int
	MaxConnectLimit int
}

var msgInfoTicker *time.Ticker
var msgServerList *list.List

func init() {
	msgServerList = list.New()

	msgInfoTicker = time.NewTicker(time.Duration(1) * time.Second)
	defer msgInfoTicker.Stop()
	//go handleRequetMsgInfo()
}

func NewMsgServer(con net.Conn) *MsgServer {
	msg := &MsgServer{
		Conn:      base.NewConn(con),
		UserCount: 0,
	}
	return msg
}

func ListenMsgServerConn(listenIp string, port int) {
	glog.Info("start listen msg_server on:", listenIp+":"+strconv.Itoa(port))

	ln, err := net.Listen("tcp", listenIp+":"+strconv.Itoa(port))
	if err != nil {
		glog.Error("listen on ", listenIp, " failed:", err.Error())
	}
	for {
		con, _ := ln.Accept()
		glog.Info("msg_server_conn connect from", con.RemoteAddr().String())

		msgServer := NewMsgServer(con)
		go msgServer.Conn.RunSync(HandlePdu)
		msgServerList.PushBack(msgServer)
	}
}

func HandlePdu(im *base.ImConn, data *base.Response) {
	switch data.Header.CommandId {
	case uint16(improto.OtherCmdID_CID_OTHER_USER_CNT_UPDATE):
		{
			hanldeMsgServerUserCount(im, data)
		}
		break
	case uint16(improto.OtherCmdID_CID_OTHER_MSG_SERV_INFO):
		{
			handleMsgServerInfo(im, data)
		}
		break
	}
}

func hanldeMsgServerUserCount(im *base.ImConn, data *base.Response) {
	req := improto.IMUserCntUpdate{}
	err := proto.Unmarshal(data.Data, &req)
	if err != nil {
		glog.Error("hanldeMsgServerUserCount() parse proto error:", err.Error())
	}

	glog.Info("hanldeMsgServerUserCount,user_count=", req.UserCount)

	var isFind = false
	for i := msgServerList.Front(); i != nil; i = i.Next() {
		item := i.Value.(*MsgServer)
		if item.Conn == im {
			item.UserCount = int(req.UserCount)
			isFind = true
			break
		}
	}

	if !isFind {
		glog.Error("hanldeMsgServerUserCount,cant find msgServerConn from msgServerList")
	}
}

func handleMsgServerInfo(im *base.ImConn, data *base.Response) {
	req := improto.IMMsgServInfo{}
	err := proto.Unmarshal(data.Data, &req)
	if err != nil {
		glog.Error("handleMsgServerInfo() parse proto error:", err.Error())
	}

	glog.Info("handleMsgServerInfo,ip1=", req.Ip1, ",ip2=", req.Ip2, ",port=", req.Port, ",host=", req.HostName,
		",maxConnect=", req.MaxConnCnt, ",curConnect=", req.CurConnCnt)

	var isFind = false
	for i := msgServerList.Front(); i != nil; i = i.Next() {
		item := i.Value.(*MsgServer)
		if item.Conn == im {
			item.PriorIp = req.Ip1
			item.BackupIp = req.Ip2
			item.Port = int(req.Port)
			item.MaxConnectLimit = int(req.MaxConnCnt)
			isFind = true
			break
		}
	}

	if !isFind {
		glog.Error("handleMsgServerInfo,cant find msgServerConn from msgServerList")
	}
}

// 定时请求msg_server上的用户信息
//func handleRequetMsgInfo() {
//	for {
//		select {
//		case <-msgInfoTicker.C:
//			{
//				for i := msgServerList.Front(); i != nil; i = i.Next() {
//
//				}
//				break
//			}
//		}
//	}
//}
