package main

import (
	"fmt"
	"log"
	"minerdb/min"
	"minerdb/raft/api/jrpc/jrpc_server"
	"minerdb/raft/api/rest/middle_ware"
	"minerdb/raft/api/rest/route"
	"minerdb/raft/discover/zk_discover"
	"minerdb/raft/starter/app"
	"minerdb/raft/starter/config"
	"runtime"
)

func init() {
	if runtime.GOOS == "windows" {
		log.Fatalln("MinerDB is only compatible with Mac and Linux")
	}
}

func main() {
	fmt.Println("开始启动！！！")
	cfg, err := config.New()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("startApiJrpc结束")
	a := app.NewApp(cfg)
	fmt.Println("NewApp结束")
	startApiJrpc(a)
	// time.Sleep(10 * time.Second)
	err = a.Node.SetRaft()
	if err != nil {
		log.Fatalln(err)
	}
	startApiRest(a)
	select {}
}

func startApiJrpc(a *app.App) {
	fmt.Println("[proto] Starting proto jrpc_server...")
	err := zk_discover.RegisterNode(a.Config.CurrentNode.ID)
	if err != nil {
		fmt.Println("a.Config.CurrentNode.ID无法被注册，进程将要退出")
		log.Fatalln("a.Config.CurrentNode.ID无法被注册", err)
	}
	jrpc_server.Start(a)
}

func startApiRest(a *app.App) {
	fmt.Println("！！！！！a.Config.CurrentNode.ApiAddress: ", a.Config.CurrentNode.ApiAddress)
	errListen := newApiRest(a).Run(a.Config.CurrentNode.ApiAddress)
	if errListen != nil {
		log.Fatalln("api can't be started:", errListen)
	}
}

func newApiRest(a *app.App) *min.Engine {
	// 注册fiber的中间件和路由
	middle_ware.InitMiddlewares(a.HttpGroup)
	route.Register(a)
	return a.HttpEngine
}
