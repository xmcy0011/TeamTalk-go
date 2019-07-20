package base

import (
	"TeamTalk-go-flutter/server/src/base/improto"
	"github.com/golang/protobuf/proto"
	"net"
)

type ResultCallback func(res Response)

//type Request struct {
//	Header  *ImPduHeader
//	Message *proto.Message
//	Callback     ResultCallback
//	NeedResponse bool
//}

type Response struct {
	Header *ImPduHeader
	Data   []byte
}

type ImConn struct {
	remoteEndPoint string
	network        string
	conn           net.Conn

	onRead ResultCallback
}

func NewClient(conn net.Conn, onReadCallback ResultCallback) *ImConn {
	im := &ImConn{
		network:        conn.RemoteAddr().Network(),
		remoteEndPoint: conn.RemoteAddr().String(),
		conn:           conn,
		onRead:         onReadCallback,
	}
	return im
}

func (im *ImConn) Run() {
	go im.write()
	go im.read()
}

func (im *ImConn) Send(message proto.Message, serviceId improto.ServiceID, cmdId int32) {
	
}

func (im *ImConn) write() {

}

func (im *ImConn) read() {

}
