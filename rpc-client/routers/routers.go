package routers

import (
	"ticketing-infra/rpc-client/api"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRouters() {
	h := server.Default()
	userGroup := h.Group("/user")

	userGroup.POST("/register", api.UserRegisterHandler)
	userGroup.POST("/login", api.UserLoginHandler)
	h.Spin()
}
