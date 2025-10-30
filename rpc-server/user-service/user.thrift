namespace go user

struct userRegisterRequest { #用户注册请求
  1: required string username
  2: required string password
}

struct userRegisterResponse { #用户注册响应
  1: required i64 id
}

struct userLoginRequest { #用户登录请求
  1: required string username
  2: required string password
}

struct userLoginResponse { #用户登录响应
    1:required i32 id
}

struct userChangePasswordRequest { #用户修改密码请求
  1: required string username
  2: required string old_password
  3: required string new_password
}

struct userChangePasswordResponse { #用户修改密码响应
}

struct userRefreshTokenRequest { #刷新token请求
  1: required string refresh_token
  2: required string username
}

struct userRefreshTokenResponse { #刷新token响应
  1: required string token
}

service UserService {
  userRegisterResponse user_register(1: userRegisterRequest req) #用户注册
  userLoginResponse user_login(1: userLoginRequest req) #用户登录
  userChangePasswordResponse user_change_password(1: userChangePasswordRequest req) #用户修改密码
  userRefreshTokenResponse user_refresh_token(1: userRefreshTokenRequest req) #刷新token
}