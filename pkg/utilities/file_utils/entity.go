package file_utils

import (
	"net/http"
)

type StatusCode struct {
	Code     string `json:"code"`
	Msg      string `json:"msg"`
	HttpCode int    `json:"http_code"`
}

var (
	WRE = StatusCode{Code: "WRE-500", Msg: "write %s file error", HttpCode: http.StatusInternalServerError}
	RRE = StatusCode{Code: "RRE-500", Msg: "read %s file error", HttpCode: http.StatusInternalServerError}
	CFE = StatusCode{Code: "CFE-500", Msg: "create %s file error", HttpCode: http.StatusInternalServerError}
)
