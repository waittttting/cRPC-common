package common

const (
	CX_SUCCESS = "success"
	CX_FAIL = "fail"
)

const CX_SUCCESS_INT = 200

type CXError struct {
	ErrCode int64
	ErrMsg  string
}

func (ce *CXError) Error() string {
	return ce.ErrMsg
}

func NewError(errCode int64, errMsg string) error {
	return &CXError{
		ErrCode: errCode,
		ErrMsg: errMsg,
	}
}

var (
	ErrDB           = NewError(3001,"database error")
	ErrInvalidParam = NewError(3002,"invalid param")
	ErrTimeOut      = NewError(3003,"timeout")
	ErrAuth         = NewError(3004,"permission denied")
)

