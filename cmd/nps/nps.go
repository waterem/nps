package main

import (
	"flag"
	"github.com/cnlh/nps/lib/common"
	"github.com/cnlh/nps/lib/daemon"
	"github.com/cnlh/nps/lib/file"
	"github.com/cnlh/nps/lib/install"
	"github.com/cnlh/nps/server"
	"github.com/cnlh/nps/server/test"
	"github.com/cnlh/nps/vender/github.com/astaxie/beego"
	"github.com/cnlh/nps/vender/github.com/astaxie/beego/logs"
	_ "github.com/cnlh/nps/web/routers"
	"log"
	"os"
	"path/filepath"
)

var (
	level   string
	logType = flag.String("log", "stdout", "Log output mode（stdout|file）")
)

func main() {
	flag.Parse()
	beego.LoadAppConfig("ini", filepath.Join(common.GetRunPath(), "conf", "nps.conf"))
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "test":
			test.TestServerConfig()
			log.Println("test ok, no error")
			return
		case "start", "restart", "stop", "status":
			daemon.InitDaemon("nps", common.GetRunPath(), common.GetTmpPath())
		case "install":
			install.InstallNps()
			return
		}
	}
	if level = beego.AppConfig.String("logLevel"); level == "" {
		level = "7"
	}
	logs.Reset()
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	if *logType == "stdout" {
		logs.SetLogger(logs.AdapterConsole, `{"level":`+level+`,"color":true}`)
	} else {
		logs.SetLogger(logs.AdapterFile, `{"level":`+level+`,"filename":"nps_log.log"}`)
	}
	task := &file.Tunnel{
		Mode: "webServer",
	}
	bridgePort, err := beego.AppConfig.Int("bridgePort")
	if err != nil {
		logs.Error("Getting bridgePort error", err)
		os.Exit(0)
	}
	server.StartNewServer(bridgePort, task, beego.AppConfig.String("bridgeType"))
}
