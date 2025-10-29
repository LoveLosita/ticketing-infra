package response

type Response struct { //响应结构体
	Status string `json:"status"`
	Info   string `json:"info"`
}

type FinalResponse struct { //最终响应结构体
	Status string      `json:"status"`
	Info   string      `json:"info"`
	Data   interface{} `json:"data"`
}

// 实现error接口
func (r Response) Error() string { // 实现 error 接口
	return r.Info
}

func InternalError(err error) Response { //服务器错误
	return Response{
		Status: "500",
		Info:   err.Error(),
	}
}

func Respond(response Response, data interface{}) FinalResponse { //传入一个响应结构体和数据，返回一个最终响应结构体
	var finalResponse FinalResponse
	finalResponse.Status = response.Status
	finalResponse.Info = response.Info
	finalResponse.Data = data
	return finalResponse
}

var ( //请求相关的响应
	Ok = Response{ //正常
		Status: "10000",
		Info:   "success",
	}
	WrongParamType = Response{ //参数错误
		Status: "40002",
		Info:   "wrong param type",
	}
)
