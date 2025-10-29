package user_service

import (
	"log"
	"net"
	"ticketing-infra/rpc-server/inits"
	user "ticketing-infra/rpc-server/user-service/kitex_gen/user/userservice"
	"ticketing-infra/rpc-server/user-service/model"

	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
)

func Start() {
	//1.连接mysql
	err := inits.ConnectDB()
	if err != nil {
		log.Fatalf("init.ConnectDB error: %v", err)
	}
	//自动迁移
	err = inits.Db.AutoMigrate(model.User{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	//2.连接redis
	err = inits.InitRedis()
	if err != nil {
		log.Fatalf("init.InitRedis error: %v", err)
	}
	//3.启动服务
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8889")
	if err != nil {
		log.Fatalf("net.ResolveTCPAddr error: %v", err)
	}
	svr := user.NewServer(new(UserServiceImpl),
		server.WithServiceAddr(addr),                            // ← 明确 TTHeader
		server.WithMetaHandler(transmeta.ServerTTHeaderHandler)) // ← 必须)
	err = svr.Run()
	if err != nil {
		log.Println(err)
	}
}
