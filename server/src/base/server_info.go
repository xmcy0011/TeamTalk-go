package base

import (
	log "github.com/golang/glog"
	"net"
)

type ServerInfo struct {
	Ip   string
	Port int

	Conn *ImConn
}

func ConnectServer(list []string) []ServerInfo {
	srvArr := make([]ServerInfo, len(list))

	for i := 0; i < len(list); i++ {
		serverEndpoint := list[i]
		srv := srvArr[i]

		conn, err := net.Dial("tcp", serverEndpoint)
		if err != nil {
			log.Error("connect", serverEndpoint, "error:", err.Error())
			srv.Conn = nil
		} else {
			srv.Conn = NewConn(conn)
		}
	}

	return nil
}

func ConnectServerCheck(list []ServerInfo) {

}
