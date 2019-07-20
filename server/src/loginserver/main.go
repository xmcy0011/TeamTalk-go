package main

import (
	"TeamTalk-go-flutter/server/src/base"
	"flag"
	"github.com/Unknwon/goconfig"
	log "github.com/golang/glog"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	VERSION     string
	BuildTime   string
	GoVersion   string
	GitCommitId string
	GitBranch   string
)

type Config struct {
	HttpListenIp string
	HttpPort     int
	MsgListIp    string
	MsgPort      int
}

func ReadConfig(fileName string) (error, Config) {
	var config = Config{}

	cfg, err := goconfig.LoadConfigFile(fileName)
	if err != nil {
		return err, config
	}

	config.HttpListenIp, err = cfg.GetValue("server", "http_list_ip")
	if err != nil {
		log.Fatal("config http_list_ip miss...")
	}

	config.HttpPort, err = cfg.Int("server", "http_port")
	if err != nil {
		log.Fatal("config http_port miss...")
	}

	config.MsgListIp, err = cfg.GetValue("server", "msg_server_listen_ip")
	if err != nil {
		log.Fatal("config msg_server_listen_ip miss...")
	}

	config.MsgPort, err = cfg.Int("server", "msg_server_port")
	if err != nil {
		log.Fatal("config msg_server_port miss...")
	}

	return nil, config
}

// 优雅退出
func WaitExit(c chan os.Signal) {
	for i := range c {
		switch i {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			// clean up here
			os.Exit(0)
		}
	}
}

func main() {
	flag.Parse()
	defer log.Flush()

	log.MaxSize = 50 * 1024 * 1024
	log.Infof("\nVersion:     %s\nBuilt:       %s\nGo version:  %s\nGit branch:  %s\nGit commit:  %s\n",
		VERSION, BuildTime, GoVersion, GitBranch, GitCommitId)
	rand.Seed(time.Now().UnixNano())
	//runtime.GOMAXPROCS(runtime.NumCPU())

	// read config
	var fileName = ""
	if len(flag.Args()) == 1 {
		fileName = flag.Arg(0)
		log.Info("config_file = ", fileName)
	} else {
		fileName = "loginserver.conf"
	}
	err, cfg := ReadConfig(fileName)
	if err != nil {
		log.Info("read config file error:", err.Error())
		os.Exit(-1)
	}

	// write pid to server.pid file
	if base.WritePid() {
		go ListenMsgServerConn(cfg.MsgListIp, cfg.MsgPort)
		go ListenHttpServerConn(cfg.HttpListenIp, cfg.HttpPort)

		// 优雅退出
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		WaitExit(c)
	}
}
