package http



type JSONResponse struct {
	Code  	int64      `json:"code"`
	Msg 	string      `json:"msg"`
	Desc    string      `json:"desc"`
	Data    interface{} `json:"data"`
}

func NewJSONResponse(code int64, msg string, desc string, Data interface{}) *JSONResponse {
	return &JSONResponse{
		Code: code,
		Msg:  msg,
		Desc: desc,
		Data: Data,
	}
}

func NewJSONResponseErr(code int64, msg string, desc string, Data interface{}) *JSONResponse {
	return &JSONResponse{
		Code: code,
		Msg:  msg,
		Desc: desc,
		Data: Data,
	}
}

func NewJSONResponseSuccess(code int64, msg string, desc string, Data interface{}) *JSONResponse {
	return &JSONResponse{
		Code: code,
		Msg:  msg,
		Desc: desc,
		Data: Data,
	}
}



