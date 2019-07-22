package main

import (
	"TeamTalk-go-flutter/server/src/base"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"net/http"
	"strconv"
)

// http json数据 输出
type ServerInfo struct {
	BackupIP string `json:"backupIP"` // 注意，首字母需要大写，否则json序列化会失败！
	Code     int    `json:"code"`
	// Discovery  string `json:"discovery"`
	// MsfsBackup string `json:"msfsBackup"`
	// MsfsPrior  string `json:"msfsPrior"`
	Msg     string `json:"msg"`
	Port    string `json:"port"`
	PriorIP string `json:"priorIP"`
}

func ListenHttpServerConn(listenIp string, port int) {
	glog.Info("start listen login_server on:", listenIp+":"+strconv.Itoa(port))

	http.HandleFunc("/msg_server", func(w http.ResponseWriter, r *http.Request) {
		var res ServerInfo
		if msgServerList.Len() < 0 {
			res = ServerInfo{
				Code:     base.LoginNoMsgServer,
				Msg:      base.LoginNoMsgServerDesc,
				PriorIP:  "106.14.172.35",
				BackupIP: "106.14.172.35",
				Port:     "9090",
			}
			glog.Error("not msgServer connect")
		} else {
			// 负载均衡算法：最小连接法
			var minServer = msgServerList.Front().Value.(*MsgServer)
			for i := msgServerList.Front().Next(); i != nil; i = i.Next() {
				cur := i.Value.(*MsgServer)
				if cur.UserCount < minServer.UserCount {
					minServer = cur
				}
			}

			if minServer.UserCount >= minServer.MaxConnectLimit {
				res = ServerInfo{
					Code:     base.LoginNoFreeMsgServer,
					Msg:      base.LoginNoFreeMsgServerDesc,
					PriorIP:  "",
					BackupIP: "",
					Port:     "",
				}
				glog.Error("msgServer:user is full")
			} else {
				res = ServerInfo{
					Code:     0,
					Msg:      "success",
					PriorIP:  minServer.PriorIp,
					BackupIP: minServer.BackupIp,
					Port:     strconv.Itoa(minServer.Port),
				}
			}
		}

		buf, _ := json.Marshal(res)
		_, e := fmt.Fprint(w, string(buf))
		glog.Info("remote http request,host:", r.Host, ",msg_info:", string(buf))
		if e != nil {
			glog.Error("http response has error:", e.Error())
		}
	})

	glog.Fatal(http.ListenAndServe(listenIp+":"+strconv.Itoa(port), nil))
}
