package errs

import (
	"github.com/jiang-ruoxi/gopkg/server/api"
	"net/http"
)

const Success = 10000

type Err struct {
	code     int
	httpCode int
	message  string
}

func (err *Err) Code() int {
	return err.code
}

func (err *Err) HttpCode() int {
	return err.httpCode
}

func (err *Err) Error() string {
	return err.message
}

func ErrResp(err api.Error) (httpCode int, rsp Response) {
	httpCode = err.HttpCode()
	rsp = Response{
		Code: err.Code(),
		Msg:  err.Error(),
	}
	return
}

type Response struct {
	Code int         `json:"code" validate:"required"` // 响应码
	Msg  string      `json:"msg" validate:"required"`  // 响应消息
	Data interface{} `json:"data"`                     // 响应数据
}

func SucResp(data interface{}) (resCode int, res Response) {
	resCode = 200
	res = Response{
		Code: Success,
		Msg:  "success",
		Data: data,
	}
	return
}

func NewError(message string, args ...int) *Err {
	var (
		badRequestCode = http.StatusBadRequest
		httpCode       = http.StatusOK
	)

	for key, arg := range args {
		if key == 0 {
			badRequestCode = arg
		} else if key == 1 {
			httpCode = arg
		}
	}

	return &Err{
		code:     badRequestCode,
		httpCode: httpCode,
		message:  message,
	}
}
