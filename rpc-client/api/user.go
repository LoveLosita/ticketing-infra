package api

import (
	"context"
	"errors"
	"strconv"
	"ticketing-infra/rpc-client/auth"
	"ticketing-infra/rpc-client/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/pkg/kerrors"

	init_client "ticketing-infra/rpc-client/inits"
	"ticketing-infra/rpc-client/kitex-gens/user/kitex_gen/user"
	"ticketing-infra/rpc-client/response"
)

func UserRegisterHandler(ctx context.Context, c *app.RequestContext) {
	//1.先获取前端传来的参数
	var registerUser user.UserRegisterRequest
	err := c.BindJSON(&registerUser)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.WrongParamType)
		return
	}
	//2.调用rpc服务端的用户注册接口
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) //设置超时时间
	defer cancel()
	resp, err := init_client.NewUserClient.UserRegister(ctx, &registerUser)
	if err != nil {
		if bizErr, isBizErr := kerrors.FromBizStatusError(err); isBizErr {
			res := strconv.Itoa(int(bizErr.BizStatusCode()))
			c.JSON(consts.StatusBadRequest, response.Response{Status: res, Info: bizErr.BizMessage()})
			return
		} else {
			c.JSON(consts.StatusInternalServerError, response.InternalError(err))
			return
		}
	}
	//3.返回结果给前端
	c.JSON(consts.StatusOK, response.Respond(response.Ok, resp))
}

func UserLoginHandler(ctx context.Context, c *app.RequestContext) {
	//1.先获取前端传来的参数
	var loginUser user.UserLoginRequest
	err := c.BindJSON(&loginUser)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.WrongParamType)
		return
	}
	//2.调用rpc服务端的用户登录接口
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) //设置超时时间
	defer cancel()
	resp, err := init_client.NewUserClient.UserLogin(ctx, &loginUser)
	if err != nil {
		if bizErr, isBizErr := kerrors.FromBizStatusError(err); isBizErr {
			res := strconv.Itoa(int(bizErr.BizStatusCode()))
			c.JSON(consts.StatusBadRequest, response.Response{Status: res, Info: bizErr.BizMessage()})
			return
		} else {
			c.JSON(consts.StatusInternalServerError, response.InternalError(err))
			return
		}
	}
	//3.生成token并返回结果给前端
	accessToken, refreshToken, err := auth.GenerateTokens(int(resp.Id))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, response.InternalError(err))
		return
	}
	loginResp := response.Respond(response.Ok, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
	c.JSON(consts.StatusOK, loginResp)
}

func UserChangePasswordHandler(ctx context.Context, c *app.RequestContext) {
	//1.先获取前端传来的参数
	var changePwdReq user.UserChangePasswordRequest
	err := c.BindJSON(&changePwdReq)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.WrongParamType)
		return
	}
	//2.调用rpc服务端的修改密码接口
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) //设置超时时间
	defer cancel()
	_, err = init_client.NewUserClient.UserChangePassword(ctx, &changePwdReq)
	if err != nil {
		if bizErr, isBizErr := kerrors.FromBizStatusError(err); isBizErr {
			res := strconv.Itoa(int(bizErr.BizStatusCode()))
			c.JSON(consts.StatusBadRequest, response.Response{Status: res, Info: bizErr.BizMessage()})
			return
		} else {
			c.JSON(consts.StatusInternalServerError, response.InternalError(err))
			return
		}
	}
	//3.返回结果给前端
	c.JSON(consts.StatusOK, response.Ok)
}

func UserRefreshTokenHandler(ctx context.Context, c *app.RequestContext) {
	//1.先获取前端传来的参数
	var refreshReq model.Tokens
	err := c.BindJSON(&refreshReq)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.WrongParamType)
		return
	}
	//2.直接刷新token
	tokens, err := auth.RefreshTokenHandler(refreshReq.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, response.InvalidRefreshToken), errors.Is(err, response.InvalidClaims),
			errors.Is(err, response.InvalidTokenSingingMethod): //如果是无效刷新令牌或者无效claims或者无效签名方法
			c.JSON(consts.StatusBadRequest, err)
			return
		default:
			c.JSON(consts.StatusInternalServerError, response.InternalError(err))
		}
	}
	//3.返回结果给前端
	c.JSON(consts.StatusOK, response.Respond(response.Ok, tokens))
}

func UserSetAdminHandler(ctx context.Context, c *app.RequestContext) {
	//1.先获取前端传来的参数
	var setAdminReq user.UserSetAdminRequest
	err := c.BindJSON(&setAdminReq)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.WrongParamType)
		return
	}
	//再从上下文中获取操作者id
	handlerID := c.GetFloat64("user_id")
	setAdminReq.OperatorId = int32(handlerID)
	//2.调用rpc服务端的设置管理员接口
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) //设置超时时间
	defer cancel()
	_, err = init_client.NewUserClient.UserSetAdmin(ctx, &setAdminReq)
	if err != nil {
		if bizErr, isBizErr := kerrors.FromBizStatusError(err); isBizErr {
			res := strconv.Itoa(int(bizErr.BizStatusCode()))
			c.JSON(consts.StatusBadRequest, response.Response{Status: res, Info: bizErr.BizMessage()})
			return
		} else {
			c.JSON(consts.StatusInternalServerError, response.InternalError(err))
			return
		}
	}
	//3.返回结果给前端
	c.JSON(consts.StatusOK, response.Ok)
}
