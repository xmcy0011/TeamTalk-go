package main

import (
	"TeamTalk-go-flutter/server/src/base"
	"TeamTalk-go-flutter/server/src/base/improto"
	"flag"
	"github.com/golang/glog"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestListenMsgServerConn(t *testing.T) {
	flag.Parse()
	defer glog.Flush()

	go ListenMsgServerConn("127.0.0.1", 8001)

	// 发起1个连接
	time.Sleep(time.Duration(1000) * time.Millisecond)
	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		t.Fatal("connect error", err.Error())
	}

	// 1.向login_server上报自己的信息
	req := improto.IMMsgServInfo{
		Port:       8000,
		Ip1:        "192.168.1.100", // 电信
		Ip2:        "192.168.1.100", // 网通
		MaxConnCnt: 50000,
		CurConnCnt: 0,
		HostName:   "192.168.1.100",
	}

	head := base.ImPduHeader{}
	head.SetPduMsg(&req)
	head.IncreSeq()
	head.CommandId = uint16(improto.OtherCmdID_CID_OTHER_MSG_SERV_INFO)
	head.ServiceId = uint16(improto.ServiceID_SID_OTHER)

	buffer, _ := head.GetBuffer()
	_, err = conn.Write(buffer)
	if err != nil {
		t.Error(err.Error())
	}

	// 2.向login_server定时更新用户数量
	//var buffer = make([]byte, 10*1024)
	for i := 0; i < 3; i++ {
		// body
		req := improto.IMUserCntUpdate{
			UserCount: uint32(i * 3),
		}

		// head
		head := base.ImPduHeader{}
		head.SetPduMsg(&req)
		head.IncreSeq()
		head.CommandId = uint16(improto.OtherCmdID_CID_OTHER_USER_CNT_UPDATE)
		head.ServiceId = uint16(improto.ServiceID_SID_OTHER)

		buffer, _ := head.GetBuffer()
		sendLen, err := conn.Write(buffer)

		if err != nil {
			t.Error(err.Error())
		}
		if sendLen < 0 {
			t.Error("send failed")
		}

		time.Sleep(time.Duration(5) * time.Second)
	}
}

func TestListenHttpServerConn(t *testing.T) {
	// msg_server
	go TestListenMsgServerConn(t)

	go ListenHttpServerConn("127.0.0.1", 8099)

	time.Sleep(time.Duration(5) * time.Second)
	for i := 0; i < 3; i++ {
		res, err := http.Get("http://127.0.0.1:8099/msg_server")
		if err != nil {
			t.Error(err.Error())
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err.Error())
		}
		_ = res.Body.Close()

		glog.Info("http://127.0.01:8099/msg_server res:", string(body))

		time.Sleep(time.Duration(6) * time.Second)
	}

}
