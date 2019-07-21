package base

import (
	"TeamTalk-go-flutter/server/src/base/improto"
	"bytes"
	"container/list"
	"encoding/binary"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"math"
	"net"
	"sync"
	"time"
)

type Request struct {
	Header  *ImPduHeader
	Message *proto.Message
	//Callback     ResultCallback
	//NeedResponse bool
}

type Response struct {
	Header *ImPduHeader
	Data   []byte
}

type ResultCallback func(im *ImConn, res *Response)

type IImConn interface {
	RunSync()
	Send(message proto.Message, serviceId improto.ServiceID, cmdId int32) error
}

type ImConn struct {
	RemoteEndPoint string
	Network        string

	IsOpen bool

	conn              net.Conn
	writeChan         chan []byte
	lastRecvHeartTime int // second

	onRead ResultCallback
}

var connList *list.List
var connListMutex sync.Mutex
var heartBeatTicker *time.Ticker

const maxHeartInterval = 60 * 1000 // 心跳超时时间

// public

func init() {
	connList = list.New()

	// 心跳线程
	heartBeatTicker = time.NewTicker(1 * time.Millisecond)
	defer heartBeatTicker.Stop()
	go heartBeatHandle()
}

func NewConn(conn net.Conn) *ImConn {
	im := &ImConn{
		Network:        conn.RemoteAddr().Network(),
		RemoteEndPoint: conn.RemoteAddr().String(),
		conn:           conn,
		writeChan:      make(chan []byte),
		IsOpen:         true,
	}
	connList.PushBack(im)
	return im
}

func (im *ImConn) RunSync(onReadCallback ResultCallback) {
	im.onRead = onReadCallback

	go im.write()
	im.read()
}

func (im *ImConn) Send(message proto.Message, serviceId improto.ServiceID, cmdId int32) error {
	buff, err := proto.Marshal(message)
	if err != nil {
		glog.Error("parse proto error", err.Error())
		return err
	}
	im.writeChan <- buff
	return nil
}

func (im *ImConn) Close() {
	defer connListMutex.Unlock()
	connListMutex.Lock()

	if im.conn != nil {
		close(im.writeChan)

		err := im.conn.Close()
		if err != nil {
			glog.Error("close conn err:", err.Error())
		}

		im.IsOpen = false
	}

	// remove conn
	for i := connList.Front(); i != nil; i = i.Next() {
		item := i.Value.(ImConn)
		if &item == im {
			connList.Remove(i)
			break
		}
	}
}

func (im *ImConn) OnTimer() {
	if math.Abs(float64(time.Now().Nanosecond()-im.lastRecvHeartTime)) > maxHeartInterval {
		glog.Error("connect", im.RemoteEndPoint, "time out")
		im.Close()
	}
}

// private

func (im *ImConn) write() {
	for {
		buff, ok := <-im.writeChan
		if !ok {
			break
		}
		sendLen, err := im.conn.Write(buff)
		if err != nil {
			glog.Error("socket send error:", err.Error())
			break
		}
		if sendLen <= 0 {
			glog.Error("socket send failed")
			break
		}
	}
}

func (im *ImConn) read() {
	var buffer = make([]byte, 10*1024) // 10 KB
	var writeOffset = 0
	for {
		dataLen, err := im.conn.Read(buffer[writeOffset:])
		if err != nil {
			glog.Error("socket read error:", err.Error())
			break
		}
		if dataLen <= 0 {
			glog.Error("socket read error:", err.Error())
			break
		}
		writeOffset += dataLen

		// tcp 粘包
		for {
			if !isPduAvailable(buffer, dataLen) {
				break
			}

			head := &ImPduHeader{}
			head.ReadHeader(buffer, dataLen)
			tempBuff := buffer[pduHeaderLen:head.Length]

			// 心跳
			if head.CommandId == uint16(improto.OtherCmdID_CID_OTHER_HEARTBEAT) {
				im.lastRecvHeartTime = time.Now().Nanosecond()

				// receive
				req := improto.IMHeartBeat{}
				_ = im.Send(&req, improto.ServiceID_SID_OTHER, int32(improto.OtherCmdID_CID_OTHER_HEARTBEAT))

			} else {
				res := &Response{
					Header: head,
					Data:   tempBuff,
				}
				// callback
				if im.onRead != nil {
					im.onRead(im, res)
				}
			}

			writeOffset -= int(head.Length)
			if uint32(writeOffset) > head.Length {
				glog.Warning("conn read reset buffer")
				resetBuf := buffer[head.Length:writeOffset]
				copy(buffer, resetBuf)
			} else {
				break
			}
		}
	}
}

func isPduAvailable(data []byte, len int) bool {
	if len < pduHeaderLen {
		return false
	}

	tempBuf := bytes.NewBuffer(data)
	var packetLen uint32
	err := binary.Read(tempBuf, binary.BigEndian, &packetLen)
	if err != nil {
		glog.Error("binary.Read error", err.Error())
		return false
	}

	if packetLen > uint32(len) {
		return false
	}
	if packetLen == 0 {
		glog.Error("pdu len is 0")
		return false
	}
	return true
}

func heartBeatHandle() {
	// need fixed ? never exit
	for {
		select {
		case <-heartBeatTicker.C:
			for i := connList.Front(); i != nil; i = i.Next() {
				item := i.Value.(ImConn)
				item.OnTimer()
			}
			break
		}
	}
}
