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

type Config struct {
	LoginServerIp   string
	LoginServerPort int
}

func ReadConfig(fileName string) (Config, error) {
	var config = Config{}

	cfg, err := goconfig.LoadConfigFile(fileName)
	if err != nil {
		return config, err
	}

	config.LoginServerIp, err = cfg.GetValue("server", "login_server_ip")
	if err != nil {
		log.Fatal("config login_server_ip miss...")
	}

	config.LoginServerPort, err = cfg.Int("server", "login_server_port")
	if err != nil {
		log.Fatal("config login_server_port miss...")
	}

	return config, nil
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
		fileName = "msgserver.conf"
	}
	cfg, err := ReadConfig(fileName)
	if err != nil {
		log.Info("read config file error:", err.Error())
		os.Exit(-1)
	}

	// write pid to server.pid file
	if base.WritePid() {


		// 优雅退出
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		WaitExit(c)
	}
}
