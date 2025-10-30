package conv

import (
	"ticketing-infra/rpc-server/user-service/kitex_gen/user"
	"ticketing-infra/rpc-server/user-service/model"
)

func ToModelRegisterUser(registerUser user.UserRegisterRequest) model.User {
	var newUser model.User
	newUser.Username = registerUser.Username
	newUser.Password = registerUser.Password
	return newUser
}

func ToProtoRegisterUser(respID int) user.UserRegisterResponse {
	var registerResp user.UserRegisterResponse
	registerResp.Id = int64(respID)
	return registerResp
}

func ToModelLoginUser(registerUser user.UserLoginRequest) model.User {
	var newUser model.User
	newUser.Username = registerUser.Username
	newUser.Password = registerUser.Password
	return newUser
}
