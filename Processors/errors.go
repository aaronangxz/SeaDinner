package Processors

import "log"

const (
	CONST_ERROR_OK        = 0
	CONST_ERROR_PARAM     = 1
	CONST_ERROR_INVALID   = 2
	CONST_ERROR_NOT_FOUND = 3
)

type ErrorResp struct {
	DebugMsg  string
	ErrorCode int32
}

func RespSuccess() ErrorResp {
	return ErrorResp{
		DebugMsg:  "",
		ErrorCode: CONST_ERROR_OK,
	}
}

func RespParamError(s string) ErrorResp {
	log.Println(s)
	return ErrorResp{
		DebugMsg:  s,
		ErrorCode: CONST_ERROR_PARAM,
	}
}
