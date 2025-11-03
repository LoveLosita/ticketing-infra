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
	//1.先将请求转换为model.User
	loginUser := conv.ToModelLoginUser(*req)
	//2.再看看用户是否存在
	result, err := dao.IfUsernameExists(loginUser.Username)
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, "Database error when checking username existence")
	}
	if result == false { //用户不存在
		return nil, kerrors.NewBizStatusError(40003, "Username does not exist")
	}
	//3.获取用户的加密密码
	hashedPwd, err := dao.GetUserHashedPassword(loginUser.Username)
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, "Database error when getting user hashed password")
	}
	//4.对比密码
	result, err = utils.CompareHashPwdAndPwd(hashedPwd, loginUser.Password)
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, "Error when comparing hashed password and password")
	}
	if result == false { //密码错误
		return nil, kerrors.NewBizStatusError(40004, "Wrong password")
	}
	//5.获取用户ID
	userId, err := dao.GetUserIDByUsername(loginUser.Username)
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, "Database error when getting user ID by username")
	}
	//6.返回验证成功响应
	return &user.UserLoginResponse{Id: int32(userId)}, nil
}

// UserChangePassword implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserChangePassword(ctx context.Context, req *user.UserChangePasswordRequest) (resp *user.UserChangePasswordResponse, err error) {
	//1.先将请求转换为model.User
	changePwdUser := conv.ToModelChangePwdUser(*req)
	//2。查看该用户是否存在
	result, err := dao.IfUsernameExists(changePwdUser.Username)
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, "Database error when checking username existence")
	}
	if result == false { //用户不存在
		return nil, kerrors.NewBizStatusError(40003, "Username does not exist")
	}
	//3.获取用户的加密密码
	hashedPwd, err := dao.GetUserHashedPassword(changePwdUser.Username)
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, "Database error when getting user hashed password")
	}
	//4.对比旧密码
	result, err = utils.CompareHashPwdAndPwd(hashedPwd, changePwdUser.Password)
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, "Error when comparing hashed password and password")
	}
	if result == false { //旧密码错误
		return nil, kerrors.NewBizStatusError(40009, "Wrong old password")
	}
	//5.修改密码
	//5.1.加密新密码
	newHashedPwd, err := utils.HashPassword(req.NewPassword_)
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, "Error when hashing new password")
	}
	//5.2.更新密码到数据库
	err = dao.ChangeUserPassword(changePwdUser.Username, newHashedPwd)
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, "Database error when changing user password")
	}
	//6.返回修改成功响应
	return &user.UserChangePasswordResponse{}, nil
}

func (s *UserServiceImpl) UserSetAdmin(ctx context.Context, req *user.UserSetAdminRequest) (resp *user.UserSetAdminResponse, err error) {
	//1.查看该用户是否存在
	result, err := dao.IfUserIDExists(int(req.TargetId))
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, "Database error when checking user ID existence")
	}
	if result == false { //用户不存在
		return nil, kerrors.NewBizStatusError(40003, "User ID does not exist")
	}
	//2.检查操作者是否是Owner角色
	role, err := dao.GetUserRoleByID(int(req.OperatorId))
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, "Database error when getting user role by ID")
	}
	if role != "owner" {
		return nil, kerrors.NewBizStatusError(40012, "does not have permission")
	}
	//3.检查目标用户是否已经是管理员或者更高角色
	targetRole, err := dao.GetUserRoleByID(int(req.TargetId))
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, "Database error when getting target user role by ID")
	}
	if targetRole == "admin" || targetRole == "owner" {
		return nil, kerrors.NewBizStatusError(40013, "Target user is already admin or higher")
	}
	//4.设置用户为管理员
	err = dao.SetUserRoleToAdmin(int(req.TargetId))
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, "Database error when setting user role to admin")
	}
	//5.返回设置成功响应
	return &user.UserSetAdminResponse{}, nil
}
