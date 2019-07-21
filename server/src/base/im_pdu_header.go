package base

import (
	"bytes"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"log"
	"sync/atomic"
)

const pduHeaderLen = 16
const pduVersion = 1
const UINT16_MAX = ^uint16(0)

var pduSeq uint32 = 1 // max uint16

type ImPduHeader struct {
	Length    uint32 // the whole pdu length
	Version   uint16 // pdu version number
	Flag      uint16 // not used
	ServiceId uint16 //
	CommandId uint16 //
	SeqNum    uint16 // 包序号
	Reversed  uint16 // 保留

	pbMessage proto.Message // 消息体
}

// public

func (it *ImPduHeader) SetPduMsg(message proto.Message) {
	it.pbMessage = message
}

func (it *ImPduHeader) ReadHeader(data []byte, len int) {
	if len >= pduHeaderLen {
		buffer := bytes.NewBuffer(data)
		_ = binary.Read(buffer, binary.BigEndian, &it.Length)
		_ = binary.Read(buffer, binary.BigEndian, &it.Version)
		_ = binary.Read(buffer, binary.BigEndian, &it.Flag)
		_ = binary.Read(buffer, binary.BigEndian, &it.ServiceId)
		_ = binary.Read(buffer, binary.BigEndian, &it.CommandId)
		_ = binary.Read(buffer, binary.BigEndian, &it.SeqNum)
		_ = binary.Read(buffer, binary.BigEndian, &it.Reversed)
	}
}

func (it *ImPduHeader) GetBuffer() ([]byte, error) {
	// write header
	tempSlice := make([]byte, 0)
	buffer := bytes.NewBuffer(tempSlice)

	data, err := proto.Marshal(it.pbMessage)
	if err != nil {
		log.Println("parse pb error:", err)
		return nil, err
	}

	// 设置头信息
	it.Length = uint32(len(data)) + pduHeaderLen
	it.Version = uint16(pduVersion)

	headerData := it.getHeaderBuffer()

	_ = binary.Write(buffer, binary.BigEndian, headerData)
	_ = binary.Write(buffer, binary.BigEndian, data)

	return buffer.Bytes(), nil
}

func (it *ImPduHeader) IncreSeq() {
	// 序号全局唯一
	it.SeqNum = getSeq()
}

func (it *ImPduHeader) GetBodyBuffer() []byte {
	data, _ := proto.Marshal(it.pbMessage)
	return data
}

// private

func (it *ImPduHeader) getHeaderBuffer() []byte {
	tempSlice := make([]byte, 0)
	buffer := bytes.NewBuffer(tempSlice)
	_ = binary.Write(buffer, binary.BigEndian, it.Length)
	_ = binary.Write(buffer, binary.BigEndian, it.Version)
	_ = binary.Write(buffer, binary.BigEndian, it.Flag)
	_ = binary.Write(buffer, binary.BigEndian, it.ServiceId)
	_ = binary.Write(buffer, binary.BigEndian, it.CommandId)
	_ = binary.Write(buffer, binary.BigEndian, it.SeqNum)
	_ = binary.Write(buffer, binary.BigEndian, it.Reversed)

	return buffer.Bytes()
}

// 获取递增唯一序号
func getSeq() uint16 {
	// 原子操作
	atomic.AddUint32(&pduSeq, 1)
	// 溢出
	if pduSeq > uint32(UINT16_MAX) {
		atomic.StoreUint32(&pduSeq, 1) // 原子操作
	}
	return uint16(pduSeq)
}