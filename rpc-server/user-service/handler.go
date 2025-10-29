package user_service

import (
	"context"
	"ticketing-infra/rpc-server/user-service/conv"
	"ticketing-infra/rpc-server/user-service/dao"
	user "ticketing-infra/rpc-server/user-service/kitex_gen/user"
	"ticketing-infra/rpc-server/utils"

	"github.com/cloudwego/kitex/pkg/kerrors"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// UserRegister implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserRegister(ctx context.Context, req *user.UserRegisterRequest) (resp *user.UserRegisterResponse, err error) {
	//1.先将请求转换为model.User
	registerUser := conv.ToModelRegisterUser(*req)
	//2.检查用户名是否已经存在
	result, err := dao.IfUsernameExists(registerUser.Username)
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, "Database error when checking username existence")
	}
	if result == true { //已经存在
		return nil, kerrors.NewBizStatusError(40001, "Username already exists")
	}
	//3.插入新用户信息
	//3.1.加密密码
	hashedPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, "Error when hashing password")
	}
	registerUser.Password = hashedPwd
	userId, err := dao.InsertNewUserInfo(registerUser)
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, "Database error when inserting new registerUser")
	}
	//4.将插入结果转换为响应并返回
	resp1 := conv.ToProtoRegisterUser(userId)
	return &resp1, nil
}

// UserLogin implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserLogin(ctx context.Context, req *user.UserLoginRequest) (resp *user.UserLoginResponse, err error) {
	// TODO: Your code here...
	return
}

// UserChangePassword implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserChangePassword(ctx context.Context, req *user.UserChangePasswordRequest) (resp *user.UserChangePasswordResponse, err error) {
	// TODO: Your code here...
	return
}

// UserRefreshToken implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserRefreshToken(ctx context.Context, req *user.UserRefreshTokenRequest) (resp *user.UserRefreshTokenResponse, err error) {
	// TODO: Your code here...
	return
}
