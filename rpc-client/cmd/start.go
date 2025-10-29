package cmd

import (
	"log"
	init_client "ticketing-infra/rpc-client/inits"
	"ticketing-infra/rpc-client/routers"
)

func Start() {
	//1.启动kitex客户端
	err := init_client.InitUserSvClient()
	if err != nil {
		log.Fatal(err)
	}
	//2.启动路由
	routers.RegisterRouters()
}
