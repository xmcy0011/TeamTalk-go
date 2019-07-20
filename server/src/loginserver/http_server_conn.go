package main

import (
	"encoding/json"
	"fmt"
	log "github.com/golang/glog"
	"net/http"
	"strconv"
)

// 我们模拟一个json数据返回
type ServerInfo struct {
	BackupIP   string `json:"backupIP"` // 注意，首字母需要大写，否则json序列化会失败！
	Code       int    `json:"code"`
	// Discovery  string `json:"discovery"`
	// MsfsBackup string `json:"msfsBackup"`
	// MsfsPrior  string `json:"msfsPrior"`
	Msg        string `json:"msg"`
	Port       string `json:"port"`
	PriorIP    string `json:"priorIP"`
}

func ListenHttpServerConn(listenIp string, listenPort int) {
	log.Info("login_server start listen.")

	http.HandleFunc("/msg_server", func(w http.ResponseWriter, r *http.Request) {
		log.Info("remote http request,host:%s \n", r.Host)

		res := ServerInfo{
			Code:       0,
			Msg:        "",
			PriorIP:    "106.14.172.35",
			BackupIP:   "106.14.172.35",
			Port:       "9090",
			// MsfsPrior:  "http://106.14.172.35:8700/",
			// MsfsBackup: "http://106.14.172.35:8700/",
			// Discovery:  "http://106.14.172.35/api/discovery",
		}
		buf, _ := json.Marshal(res)
		_, e := fmt.Fprint(w, string(buf))
		if e != nil {
			log.Error("http response has error:%s \n", e.Error())
		}
	})

	log.Fatal(http.ListenAndServe(listenIp+":"+strconv.Itoa(listenPort), nil))
}
