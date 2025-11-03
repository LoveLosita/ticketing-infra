package routers

import (
	"ticketing-infra/rpc-client/api"
	"ticketing-infra/rpc-client/middleware"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRouters() {
	h := server.Default()
	userGroup := h.Group("/user")

	userGroup.POST("/register", api.UserRegisterHandler)
	userGroup.POST("/login", api.UserLoginHandler)
	userGroup.POST("/change-pwd", api.UserChangePasswordHandler)
	userGroup.POST("/refresh-token", api.UserRefreshTokenHandler)
	userGroup.POST("/set-admin", middleware.JWTTokenAuth(), api.UserSetAdminHandler)
	h.Spin()
}
