package api

import (
	"context"
	"strconv"
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
