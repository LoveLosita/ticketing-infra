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
	userGroup.POST("/change-pwd", api.UserChangePasswordHandler)
	userGroup.POST("/refresh-token", api.UserRefreshTokenHandler)
	h.Spin()
}
